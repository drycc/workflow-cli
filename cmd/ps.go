package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/containerd/console"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/ps"
	"github.com/drycc/workflow-cli/pkg/logging"
	"golang.org/x/net/websocket"
	yaml "gopkg.in/yaml.v3"
)

const (
	stdinChannel  = "\x00"
	stdoutChannel = "\x01"
	stderrChannel = "\x02"
	errorChannel  = "\x03"
	resizeChannel = "\x04"
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

// PodLogs returns the logs from an pod.
func (d *DryccCmd) PsLogs(appID, podID string, lines int, follow bool, container string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	request := api.PodLogsRequest{
		Lines:     lines,
		Follow:    follow,
		Container: container,
	}
	conn, err := ps.Logs(s.Client, appID, podID, request)
	if err != nil {
		return err
	}
	defer conn.Close()
	for {
		var message string
		err := websocket.Message.Receive(conn, &message)
		if err != nil {
			if err != io.EOF {
				log.Printf("error: %v", err)
			}
			break
		}
		logging.PrintLog(os.Stdout, strings.TrimRight(string(message), "\n"))
	}
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
func (d *DryccCmd) PsRestart(appID string, targets []string, confirm string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	if len(targets) == 0 && (confirm == "" || confirm != "yes") {
		d.Printf(` !    WARNING: Potentially Restart Action
 !    This command will restart all processes of the application
 !    To proceed, type "yes" !

> `)

		fmt.Scanln(&confirm)
		if confirm != "yes" {
			return fmt.Errorf("cancel the restart action")
		}
	}
	d.Printf("Restarting processes... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress(d.WOut)
	ptypes := strings.Join(targets, ",")
	targetMap := map[string]string{
		"types": ptypes,
	}
	err = ps.Restart(s.Client, appID, targetMap)
	quit <- true
	<-quit
	if err != nil {
		return err
	}

	d.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))
	return nil
}

// PsDescribe describe an app's processes.
func (d *DryccCmd) PsDescribe(appID, podID string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	// The 1000 is fake for now until API understands limits
	podState, _, err := ps.Describe(s.Client, appID, podID, 1000)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{})
	for _, containerState := range podState {
		table.Append([]string{"Container:", containerState.Container})
		table.Append([]string{"Image:", containerState.Image})
		table.Append([]string{"Command:"})
		for _, command := range containerState.Command {
			table.Append([]string{"", fmt.Sprintf("- %v", command)})
		}
		table.Append([]string{"Args:"})
		for _, arg := range containerState.Args {
			table.Append([]string{"", fmt.Sprintf("- %v", arg)})
		}
		// State
		for key := range containerState.State {
			table.Append([]string{"State:", key})
			value := containerState.State[key]
			for innerKey := range value {
				table.Append([]string{fmt.Sprintf("  %s:", innerKey), strconv.Quote(fmt.Sprintf("%v", value[innerKey]))})
			}
		}
		// LastState
		for key := range containerState.LastState {
			table.Append([]string{"Last State:", key})
			value := containerState.LastState[key]
			for innerKey := range value {
				table.Append([]string{fmt.Sprintf("  %s:", innerKey), strconv.Quote(fmt.Sprintf("%v", value[innerKey]))})
			}
		}
		table.Append([]string{"Ready:", fmt.Sprintf("%v", containerState.Ready)})
		table.Append([]string{"Restart Count:", fmt.Sprintf("%v", containerState.RestartCount)})
		table.Append([]string{})
	}
	table.Render()
	return nil
}

func printProcesses(d *DryccCmd, appID string, input []api.Pods) {
	processes := ps.ByType(input)

	if len(processes) == 0 {
		d.Println(fmt.Sprintf("No processes found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"NAME", "RELEASE", "STATE", "PTYPE", "READY", "RESTARTS", "STARTED"})
		for _, process := range processes {
			for _, pod := range process.PodsList {
				table.Append([]string{
					pod.Name,
					pod.Release,
					pod.State,
					pod.Type,
					pod.Ready,
					fmt.Sprintf("%v", pod.Restarts),
					pod.Started.Format("2006-01-02T15:04:05MST"),
				})
			}
		}
		table.Render()
	}
}

func printExec(d *DryccCmd, conn *websocket.Conn) error {
	var data string
	err := websocket.Message.Receive(conn, &data)
	if err != nil {
		if err != io.EOF {
			log.Printf("error: %v", err)
		}
		return nil
	}
	message, err := parseChannelMessage(data)
	if err == nil {
		d.Printf("%s", message)
	}
	return err
}

func runRecvTask(conn *websocket.Conn, c console.Console, recvChan, sendChan chan string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			var data string
			err := websocket.Message.Receive(conn, &data)
			if err != nil {
				cancel()
				break
			}
			message, err := parseChannelMessage(data)
			if err != nil {
				cancel()
				break
			}
			recvChan <- message
		}
	}()
	go func() {
		buf := make([]byte, 1024)
		for {
			size, err := c.Read(buf)
			if err == io.EOF {
				cancel()
				break
			} else if err != nil {
				continue
			}
			sendChan <- string(buf[:size])
		}
	}()
	return ctx, cancel
}

func runResizeTask(conn *websocket.Conn, c console.Console) {
	go func() {
		var size console.WinSize
		for {
			if tmpSize, err := c.Size(); err == nil {
				if size.Height != tmpSize.Height || size.Width != tmpSize.Width {
					size = tmpSize
					message := fmt.Sprintf(`{"Height": %d, "Width": %d}`, size.Height, size.Width)
					if err := websocket.Message.Send(conn, resizeChannel+message); err != nil {
						break
					}
				}
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()
}

func streamExec(conn *websocket.Conn, tty bool) error {
	c := console.Current()
	defer c.Reset()
	if tty {
		if err := c.SetRaw(); err != nil {
			return err
		}
		runResizeTask(conn, c)
	}
	recvChan, sendChan := make(chan string, 10), make(chan string, 10)
	ctx, cancel := runRecvTask(conn, c, recvChan, sendChan)
	defer cancel()
	defer close(recvChan)
	defer close(sendChan)
	for {
		select {
		case <-ctx.Done():
			return nil
		case message := <-sendChan:
			if err := websocket.Message.Send(conn, stdinChannel+message); err != nil {
				return err
			}
		case message := <-recvChan:
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

func parseChannelMessage(data string) (string, error) {
	channel, message := data[0], data[1:]
	if string(channel) == errorChannel {
		data := make(map[string]interface{})
		yaml.Unmarshal([]byte(message), data)
		if value, hasKey := data["message"]; hasKey {
			if message, ok := value.(string); ok {
				return message, nil
			}
			return "", fmt.Errorf("message is not string, type: %T", message)
		}
		return "", nil
	}
	return message, nil
}
