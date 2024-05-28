package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/workflow-cli/pkg/git"
	"github.com/drycc/workflow-cli/settings"
	"github.com/olekukonko/tablewriter"
)

var defaultLimit = -1
var defaultLines = 64

func progress(wOut io.Writer) chan bool {
	frames := []string{"...", "o..", ".o.", "..o"}
	backspaces := strings.Repeat("\b", 3)
	tick := time.NewTicker(400 * time.Millisecond)
	quit := make(chan bool)
	go func() {
		for {
			for _, frame := range frames {
				fmt.Fprint(wOut, frame)
				select {
				case <-quit:
					fmt.Fprint(wOut, backspaces)
					close(quit)
					return
				case <-tick.C:
					fmt.Fprint(wOut, backspaces)
				}
			}
		}
	}()
	return quit
}

// load loads settings file and looks up the app name
func load(cf string, appID string) (*settings.Settings, string, error) {
	s, err := settings.Load(cf)

	if err != nil {
		return nil, "", err
	}

	if appID == "" {
		appID, err = git.DetectAppName(git.DefaultCmd, s.Client.ControllerURL.Host)

		if err != nil {
			return nil, "", err
		}
	}

	return s, appID, nil
}

func drinkOfChoice() string {
	drink := os.Getenv("DRYCC_DRINK_OF_CHOICE")

	if drink == "" {
		drink = "coffee"
	}

	return drink
}

func limitCount(objs, total int) string {
	if objs == total {
		return "\n"
	}

	return fmt.Sprintf(" (%d of %d)\n", objs, total)
}

// checkAPICompatibility handles specific behavior for certain errors,
// such as printing an warning for the API mismatch error
func (d *DryccCmd) checkAPICompatibility(c *drycc.Client, err error) error {
	if err == drycc.ErrAPIMismatch {
		if !d.Warned {
			d.PrintErrf(`!    WARNING: Client and server API versions do not match. Please consider upgrading.
!    Client version: %s
!    Server version: %s
`, drycc.APIVersion, c.ControllerAPIVersion)
			d.Warned = true
		}

		// API mismatch isn't fatal, so after warning continue on.
		return nil
	}

	return err
}

// getDefaultFormatTable return default format ascii table
func (d *DryccCmd) getDefaultFormatTable(headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(d.WOut)
	table.SetHeader(headers)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding(fmt.Sprintf("%4s", " "))
	table.SetNoWhiteSpace(true)
	return table
}

// format time string to local time
func (d *DryccCmd) formatTime(timeStr string) string {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return timeStr
	}
	if d.Location != nil {
		return t.In(d.Location).Format(time.RFC3339)
	}
	return t.In(time.UTC).Format(time.RFC3339)
}

// wrapString wraps s into a paragraph of lines of length lim, with minimal raggedness.
func (d *DryccCmd) wrapString(s string) string {
	sa, _ := tablewriter.WrapString(s, defaultLines)
	return strings.Join(sa, "\r\n")
}

func sortKeys(data map[string]interface{}) *[]string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return &keys
}

func safeGetString(data string) string {
	if data == "" {
		return "<none>"
	}
	return data
}
