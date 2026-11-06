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
	"github.com/olekukonko/tablewriter/pkg/twwarp"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
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

const (
	defaultTablePadding  = "    "
	defaultHeaderPadding = " "
)

type defaultFormatRenderer struct {
	*renderer.Blueprint
	rowLengths []int
	rowIndex   int
	rowLength  int
}

func (r *defaultFormatRenderer) Header(headers [][]string, ctx tw.Formatting) {
	keys := ctx.Row.Widths.SortedKeys()
	if len(keys) > 0 {
		last := keys[len(keys)-1]
		ctx.Row.Widths = ctx.Row.Widths.Clone()
		ctx.Row.Widths.Set(last, max(0, ctx.Row.Widths.Get(last)-len(defaultTablePadding)+len(defaultHeaderPadding)))

		cell := ctx.Row.Current[last]
		cell.Padding.Right = defaultHeaderPadding
		ctx.Row.Current[last] = cell
	}
	r.Blueprint.Header(headers, ctx)
}

func (r *defaultFormatRenderer) Row(row []string, ctx tw.Formatting) {
	if !ctx.IsSubRow && r.rowIndex < len(r.rowLengths) {
		r.rowLength = r.rowLengths[r.rowIndex]
		r.rowIndex++
	}

	keys := ctx.Row.Widths.SortedKeys()
	if r.rowLength < len(keys) {
		cells := ctx.Row.Current
		ctx.Row.Widths = ctx.Row.Widths.Clone()
		ctx.Row.Current = make(map[int]tw.CellContext, r.rowLength)
		for i, key := range keys {
			if i < r.rowLength {
				ctx.Row.Current[key] = cells[key]
			} else {
				ctx.Row.Widths.Delete(key)
			}
		}
	}
	r.Blueprint.Row(row, ctx)
}

type defaultFormatTable struct {
	*tablewriter.Table
	renderer *defaultFormatRenderer
}

func (t *defaultFormatTable) Append(rows ...interface{}) error {
	if err := t.Table.Append(rows...); err != nil {
		return err
	}

	rowLength := len(rows)
	if len(rows) == 1 {
		if row, ok := rows[0].([]string); ok {
			rowLength = len(row)
		}
	}
	t.renderer.rowLengths = append(t.renderer.rowLengths, rowLength)
	return nil
}

// getDefaultFormatTable return default format ascii table
func (d *DryccCmd) getDefaultFormatTable(headers []string) *defaultFormatTable {
	formattedHeaders := make([]string, len(headers))
	for i, header := range headers {
		formattedHeaders[i] = tw.Title(header)
	}

	formatRenderer := &defaultFormatRenderer{Blueprint: renderer.NewBlueprint()}
	table := tablewriter.NewTable(d.WOut,
		tablewriter.WithRenderer(formatRenderer),
		tablewriter.WithHeaderAlignment(tw.AlignLeft),
		tablewriter.WithRowAlignment(tw.AlignLeft),
		tablewriter.WithHeaderAutoWrap(tw.WrapNone),
		tablewriter.WithRowAutoWrap(tw.WrapNone),
		tablewriter.WithHeaderAutoFormat(tw.Off),
		tablewriter.WithTrimSpace(tw.Off),
		tablewriter.WithTrimLine(tw.Off),
		tablewriter.WithPadding(tw.Padding{Right: defaultTablePadding, Overwrite: true}),
		tablewriter.WithRendition(tw.Rendition{
			Borders: tw.BorderNone,
			Symbols: tw.NewSymbols(tw.StyleNone),
			Settings: tw.Settings{
				Lines:      tw.LinesNone,
				Separators: tw.SeparatorsNone,
			},
		}),
		tablewriter.WithHeader(formattedHeaders),
	)
	return &defaultFormatTable{Table: table, renderer: formatRenderer}
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
	sa, _ := twwarp.WrapString(s, defaultLines)
	return strings.Join(sa, "\n")
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
		if ptypes[i] == api.PtypeWeb {
			return true
		}
		if ptypes[j] == api.PtypeWeb {
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
