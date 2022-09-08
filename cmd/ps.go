package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/ps"
	"github.com/gorilla/websocket"
	"golang.org/x/term"
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
		streamExec(conn)
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

	psType, psName := "", ""
	if target != "" {
		psType, psName = parseType(target, appID)
	}

	d.Printf("Restarting processes... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress(d.WOut)

	processes, err := ps.Restart(s.Client, appID, psType, psName)
	quit <- true
	<-quit
	if err == drycc.ErrPodNotFound {
		return fmt.Errorf("Could not find process type %s in app %s", psType, appID)
	} else if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if len(processes) == 0 {
		d.Println("Could not find any processes to restart")
	} else {
		d.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))
		printProcesses(appID, processes, d.WOut)
	}

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

func streamExec(conn *websocket.Conn) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	t := term.NewTerminal(os.Stdin, "")

	ctx, cancel := context.WithCancel(context.Background())

	stdin := make(chan string)
	defer close(stdin)
	go func() {
		for {
			line, err := t.ReadLine()
			if err != nil {
				cancel()
				break
			}
			stdin <- line
		}
	}()

	go func() {
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil || messageType == websocket.CloseMessage {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Printf("error: %v", err)
				}
				cancel()
				break
			} else {
				t.Write(message)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case line := <-stdin:
			if err := conn.WriteMessage(websocket.TextMessage, []byte(line+"\n")); err != nil {
				return nil
			}
		}
	}
}

func parseType(target string, appID string) (string, string) {
	var psType, psName string

	if strings.Contains(target, "-") {
		replaced := strings.Replace(target, appID+"-", "", 1)
		parts := strings.Split(replaced, "-")
		// the API requires the type, for now
		// regex matches against how Deployment pod name is constructed
		regex := regexp.MustCompile("[a-z0-9]{8,10}-[a-z0-9]{5}$")
		if regex.MatchString(replaced) || len(parts) == 2 {
			psType = parts[0]
		} else {
			psType = parts[1]
		}
		// process name is the full pod
		psName = target
	} else {
		psType = target
	}

	return psType, psName
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
