package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/containerd/console"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/ps"
	"github.com/gorilla/websocket"
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

	printProcesses(appID, processes, d.WOut)

	return nil
}

// PsList lists an app's processes.
func (d *DryccCmd) PsExec(appID, podID string, tty, stdin bool, command []string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	conn, err := ps.Exec(s.Client, appID, podID, tty, stdin, command)
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

	d.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))

	processes, _, err := ps.List(s.Client, appID, s.Limit)
	if err != nil {
		return err
	}

	printProcesses(appID, processes, d.WOut)
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

func printProcesses(appID string, input []api.Pods, wOut io.Writer) {
	processes := ps.ByType(input)

	fmt.Fprintf(wOut, "=== %s Processes\n", appID)

	for _, process := range processes {
		fmt.Fprintf(wOut, "--- %s:\n", process.Type)

		for _, pod := range process.PodsList {
			fmt.Fprintf(wOut, "%s %s (%s)\n", pod.Name, pod.State, pod.Release)
		}
	}
}

func printExec(d *DryccCmd, conn *websocket.Conn) error {
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
			log.Printf("error: %v", err)
		}
		return nil
	}
	if messageType == websocket.TextMessage {
		d.Printf("%s", string(message))
	} else {
		d.Printf(base64.StdEncoding.EncodeToString(message))
	}
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

	recvQueue := make(chan []byte)
	defer close(recvQueue)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil || messageType == websocket.CloseMessage {
				cancel()
				break
			} else {
				recvQueue <- message
			}
		}
	}()

	sendQueue := make(chan []byte)
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
				sendQueue <- buf[:size]
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case message := <-sendQueue:
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return err
			}
		case message := <-recvQueue:
			c.Write(message)
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
