package template

import (
	"bytes"
	"strings"
	"text/template"
)

func CustomExample(example string, kwargs map[string]string) string {
	tpl, err := template.New("template").Funcs(template.FuncMap{
		"rpad": func(str string, length int) string {
			padding := length - len(str)
			if padding > 0 {
				return str + strings.Repeat(" ", padding)
			}
			return str
		},
	}).Parse(`{{.Example}}

Args:
  {{- range $key, $value := .Kwargs}}
  {{rpad $key ($.Padding)}}{{$value}}
  {{- end}}`)
	if err != nil {
		return example
	}
	padding, miniPadding := 0, 3
	for key := range kwargs {
		if padding < len(key)+miniPadding {
			padding = len(key) + miniPadding
		}
	}
	var buf bytes.Buffer
	if err = tpl.Execute(&buf, struct {
		Example string
		Kwargs  map[string]string
		Padding int
	}{example, kwargs, padding}); err != nil {
		return example
	}
	return buf.String()
}
