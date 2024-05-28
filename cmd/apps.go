package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/apps"
	"github.com/drycc/controller-sdk-go/appsettings"
	"github.com/drycc/controller-sdk-go/domains"
	"github.com/drycc/controller-sdk-go/ps"
	"github.com/drycc/workflow-cli/pkg/git"
	"github.com/drycc/workflow-cli/pkg/logging"
	"github.com/drycc/workflow-cli/pkg/webbrowser"
	"github.com/drycc/workflow-cli/settings"
	"golang.org/x/net/websocket"
)

// AppCreate creates an app.
func (d *DryccCmd) AppCreate(id, remote string, noRemote bool) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	d.Print("Creating Application... ")
	quit := progress(d.WOut)
	app, err := apps.New(s.Client, id)

	quit <- true
	<-quit

	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("done, created %s\n", app.ID)

	if !noRemote {
		if err = git.CreateRemote(git.DefaultCmd, s.Client.ControllerURL.Host, remote, app.ID); err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("error: remote %s already exists.", remote)) {
				msg := "A git remote with the name %s already exists. To overwrite this remote run:\n"
				msg += "drycc git:remote --force --remote %s --app %s"
				return fmt.Errorf(msg, remote, remote, app.ID)
			}
			return err
		}

		d.Printf(remoteCreationMsg, remote, app.ID)
	}

	if noRemote {
		d.Printf("If you want to add a git remote for this app later, use `drycc git:remote -a %s`\n", app.ID)
	}

	return nil
}

// AppsList lists apps on the Drycc controller.
func (d *DryccCmd) AppsList(results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	apps, count, err := apps.List(s.Client, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count > 0 {
		table := d.getDefaultFormatTable([]string{"ID", "UUID", "OWNER", "CREATED", "UPDATED"})
		for _, app := range apps {
			table.Append([]string{
				app.ID,
				app.UUID,
				app.Owner,
				d.formatTime(app.Created),
				d.formatTime(app.Updated),
			})
		}
		table.Render()
	} else {
		d.Println("No apps found.")
	}
	return nil
}

// AppInfo prints info about app.
func (d *DryccCmd) AppInfo(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	app, err := apps.Get(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	url, err := d.appURL(s, appID)
	if err != nil {
		return err
	}

	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"App:", app.ID})
	table.Append([]string{"URL:", url})
	table.Append([]string{"UUID:", app.UUID})
	table.Append([]string{"Owner:", app.Owner})
	table.Append([]string{"Created:", d.formatTime(app.Created)})
	table.Append([]string{"Updated:", d.formatTime(app.Updated)})

	// print the app processes
	processes, _, err := ps.List(s.Client, appID, defaultLimit)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if len(processes) > 0 {
		table.Append([]string{"Processes:"})
		for index, process := range processes {
			table.Append([]string{"", "Name:", process.Name})
			table.Append([]string{"", "Release:", process.Release})
			table.Append([]string{"", "State:", process.State})
			table.Append([]string{"", "Type:", process.Type})
			table.Append([]string{"", "Started:", process.Started.Format("2006-01-02T15:04:05MST")})
			if len(processes) > index+1 {
				table.Append([]string{""})
			}
		}
	} else {
		table.Append([]string{"Processes:", safeGetString("")})
	}

	domains, _, err := domains.List(s.Client, appID, defaultLimit)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if len(domains) > 0 {
		table.Append([]string{"Domains:"})
		for index, domain := range domains {
			table.Append([]string{"", "Domain:", domain.Domain})
			table.Append([]string{"", "Created:", d.formatTime(domain.Created)})
			table.Append([]string{"", "Updated:", d.formatTime(domain.Updated)})
			if len(domains) > index+1 {
				table.Append([]string{""})
			}
		}
	} else {
		table.Append([]string{"Domains:", safeGetString("")})
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if len(appSettings.Label) > 0 {
		table.Append([]string{"Labels:"})
		for index, label := range *sortKeys(appSettings.Label) {
			table.Append([]string{"", "Key:", label})
			table.Append([]string{"", "Value:", fmt.Sprintf("%v", appSettings.Label[label])})
			if len(appSettings.Label) > index+1 {
				table.Append([]string{""})
			}
		}
	} else {
		table.Append([]string{"Labels:", safeGetString("")})
	}
	table.Render()
	return nil
}

// AppOpen opens an app in the default webbrowser.
func (d *DryccCmd) AppOpen(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	u, err := d.appURL(s, appID)
	if err != nil {
		return err
	}

	if u == "" {
		return fmt.Errorf(noDomainAssignedMsg, appID)
	}

	if !(strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")) {
		u = "http://" + u
	}

	return webbrowser.Webbrowser(u)
}

// AppLogs returns the logs from an app.
func (d *DryccCmd) AppLogs(appID string, lines int, follow bool, timeout int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	request := api.AppLogsRequest{
		Lines:   lines,
		Follow:  follow,
		Timeout: timeout,
	}
	conn, err := apps.Logs(s.Client, appID, request)
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

// AppRun runs a one time command in the app.
func (d *DryccCmd) AppRun(appID, command string, volumeVars []string, timeout, expires uint32) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Running '%s'...\n", command)
	volumeMap, err := parseMount(volumeVars)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if err := apps.Run(s.Client, appID, command, volumeMap, timeout, expires); d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	return nil
}

func parseMount(volumeVars []string) (map[string]interface{}, error) {
	volumeMap := make(map[string]interface{})

	regex := regexp.MustCompile(`^([A-z_]+[A-z0-9_]*):([\s\S]*)$`)
	for _, volume := range volumeVars {
		if regex.MatchString(volume) {
			captures := regex.FindStringSubmatch(volume)
			volumeMap[captures[1]] = captures[2]
		} else {
			return nil, fmt.Errorf("'%s' does not match the pattern 'key:var', ex: MODE:test", volume)
		}
	}
	return volumeMap, nil
}

// AppDestroy destroys an app.
func (d *DryccCmd) AppDestroy(appID, confirm string) error {
	gitSession := false

	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if appID == "" {
		appID, err = git.DetectAppName(git.DefaultCmd, s.Client.ControllerURL.Host)

		if err != nil {
			return err
		}

		gitSession = true
	}

	if confirm == "" {
		d.Printf(` !    WARNING: Potentially Destructive Action
 !    This command will destroy the application: %s
 !    To proceed, type "%s" or re-run this command with --confirm=%s

> `, appID, appID, appID)

		fmt.Scanln(&confirm)
	}

	if confirm != appID {
		return fmt.Errorf("app %s does not match confirm %s, aborting", appID, confirm)
	}

	startTime := time.Now()
	d.Printf("Destroying %s...\n", appID)

	if err = apps.Delete(s.Client, appID); d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))

	if gitSession {
		return d.GitRemove(appID)
	}

	return nil
}

// AppTransfer transfers app ownership to another user.
func (d *DryccCmd) AppTransfer(appID, username string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Transferring %s to %s... ", appID, username)

	err = apps.Transfer(s.Client, appID, username)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")

	return nil
}

const noDomainAssignedMsg = "no domain assigned to %s"

// appURL grabs the first domain an app has and returns this.
func (d *DryccCmd) appURL(s *settings.Settings, appID string) (string, error) {
	domains, _, err := domains.List(s.Client, appID, 1)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return "", err
	}

	if len(domains) == 0 {
		return "", nil
	}

	return expandURL(s.Client.ControllerURL.Host, domains[0].Domain), nil
}

// expandURL expands an app url if necessary.
func expandURL(host, u string) string {
	if strings.Contains(u, ".") {
		// If domain is a full url.
		return u
	}

	// If domain is a subdomain, look up the controller url and replace the subdomain.
	parts := strings.Split(host, ".")
	parts[0] = u
	return strings.Join(parts, ".")
}
