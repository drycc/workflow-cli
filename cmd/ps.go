package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/containerd/console"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/ps"
	"golang.org/x/net/websocket"
)

// PsList lists an app's processes.
func (d *DryccCmd) PsList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	processes, _, err := ps.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	printProcesses(d, appID, processes)

	return nil
}

// PsList lists an app's processes.
func (d *DryccCmd) PsExec(appID, podID string, tty, stdin bool, command []string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	request := api.Command{
		Tty:     tty,
		Stdin:   stdin,
		Command: command,
	}
	conn, err := ps.Exec(s.Client, appID, podID, request)
	if err != nil {
		return err
	}
	defer conn.Close()
	if stdin {
		streamExec(conn, tty)
	} else {
		printExec(d, conn)
	}
	return nil
}

// PsScale scales an app's processes.
func (d *DryccCmd) PsScale(appID string, targets []string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	targetMap, err := parsePsTargets(targets)
	if err != nil {
		return err
	}

	d.Printf("Scaling processes... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress(d.WOut)

	err = ps.Scale(s.Client, appID, targetMap)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("done in %ds\n\n", int(time.Since(startTime).Seconds()))

	processes, _, err := ps.List(s.Client, appID, s.Limit)
	if err != nil {
		return err
	}

	printProcesses(d, appID, processes)
	return nil
}

// PsRestart restarts an app's processes.
func (d *DryccCmd) PsRestart(appID, target string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Printf("Restarting processes... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress(d.WOut)

	err = ps.Restart(s.Client, appID, target)
	quit <- true
	<-quit
	if err != nil {
		return err
	}

	d.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))
	return nil
}

func printProcesses(d *DryccCmd, appID string, input []api.Pods) {
	processes := ps.ByType(input)

	if len(processes) == 0 {
		d.Println(fmt.Sprintf("No processes found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"NAME", "RELEASE", "STATE", "TYPE", "STARTED"})
		for _, process := range processes {
			for _, pod := range process.PodsList {
				table.Append([]string{
					pod.Name,
					pod.Release,
					pod.State,
					pod.Type,
					pod.Started.Format("2006-01-02T15:04:05MST"),
				})
			}
		}
		table.Render()
	}
}

func printExec(d *DryccCmd, conn *websocket.Conn) error {
	var message string
	err := websocket.Message.Receive(conn, &message)
	if err != nil {
		if err != io.EOF {
			log.Printf("error: %v", err)
		}
		return nil
	}
	d.Printf("%s", message)
	return nil
}

func streamExec(conn *websocket.Conn, tty bool) error {
	c := console.Current()
	defer c.Reset()
	if tty {
		if err := c.SetRaw(); err != nil {
			return err
		}
	}

	recvQueue := make(chan string)
	defer close(recvQueue)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			var message string
			err := websocket.Message.Receive(conn, &message)
			if err != nil {
				cancel()
				break
			}
			recvQueue <- message
		}
	}()

	sendQueue := make(chan string)
	defer close(sendQueue)
	go func() {
		buf := make([]byte, 1024)
		for {
			size, err := c.Read(buf)
			if err == io.EOF {
				cancel()
				break
			} else if err != nil {
				continue
			} else {
				sendQueue <- string(buf[:size])
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case message := <-sendQueue:
			if err := websocket.Message.Send(conn, message); err != nil {
				return err
			}
		case message := <-recvQueue:
			c.Write([]byte(message))
		}
	}
}

func parsePsTargets(targets []string) (map[string]int, error) {
	targetMap := make(map[string]int)
	regex := regexp.MustCompile(`^([a-z0-9]+(?:-[a-z0-9]+)*)=([0-9]+)$`)
	var err error

	for _, target := range targets {
		if regex.MatchString(target) {
			captures := regex.FindStringSubmatch(target)
			targetMap[captures[1]], err = strconv.Atoi(captures[2])

			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("'%s' does not match the pattern 'type=num', ex: web=2", target)
		}
	}

	return targetMap, nil
}
