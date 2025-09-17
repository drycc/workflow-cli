package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/olekukonko/tablewriter"
	yaml "gopkg.in/yaml.v3"
)

var (
	defaultLimit = -1
	defaultLines = 64
)

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
	return strings.Join(sa, "\n")
}

// fixateString  fix s width.
func (d *DryccCmd) fixateString(s string, width int) string {
	switch {
	case len(s) > width:
		trimSize := len(s) - width + 3
		if trimSize < len(s) {
			s = "..." + s[trimSize:]
		}
	case len(s) < width:
		s += strings.Repeat(" ", width-len(s))
	}
	return s
}

// indentString indent s into a paragraph of lines of length lim, with minimal raggedness.
func (d *DryccCmd) indentString(s string, indent int) string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimRight(line, "\r")
		padding := indent + len(line)
		lines = append(lines, fmt.Sprintf("% "+strconv.Itoa(padding)+"s", line))
	}
	return strings.Join(lines, "\n")
}

// toYamlString convert object to yaml string
func (d *DryccCmd) toYamlString(v any, indent int) string {
	buf := bytes.Buffer{}
	encode := yaml.NewEncoder(&buf)
	encode.SetIndent(indent)
	encode.Encode(v)
	return buf.String()
}

func sortKeys(data map[string]any) *[]string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return &keys
}

func sortConfigValues(data []api.ConfigValue) []api.ConfigValue {
	sort.Slice(data, func(i, j int) bool {
		return strings.ToLower(data[i].Name) < strings.ToLower(data[j].Name)
	})
	return data
}

func sortPtypes(ptypes []string) []string {
	sort.Slice(ptypes, func(i, j int) bool {
		if ptypes[i] == "web" {
			return true
		}
		if ptypes[j] == "web" {
			return false
		}
		return ptypes[i] < ptypes[j]
	})
	return ptypes
}

func safeGetString(data string) string {
	if data == "" {
		return "<none>"
	}
	return data
}

// ResponseLimit converts a limit value to the format expected by the API.
// If limit is 0, it returns -1 to indicate no limit.
func ResponseLimit(limit int) (int, error) {
	if limit == 0 {
		return -1, nil
	}
	return limit, nil
}
