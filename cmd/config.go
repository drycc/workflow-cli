package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/appsettings"
	"github.com/drycc/controller-sdk-go/config"
	"sigs.k8s.io/yaml"
)

// ConfigInfo for an app
func (d *DryccCmd) ConfigInfo(appID string, ptype string, group string, version int) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID, version)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	// init output struct
	cv := api.ConfigInfo{
		Group: make(map[string][]api.KV),
		Ptype: make(map[string]api.PtypeValue),
	}
	for _, value := range sortConfigValues(config.Values) {
		// display the selected or all
		if (ptype != "" && value.Ptype == ptype) ||
			(group != "" && value.Group == group) ||
			(ptype == "" && group == "") {
			if value.Group != "" {
				cv.Group[value.Group] = append(cv.Group[value.Group], api.KV{Name: value.Name, Value: value.Value})
			} else if value.Ptype != "" {
				temp := cv.Ptype[value.Ptype]
				temp.Env = append(temp.Env, api.KV{Name: value.Name, Value: value.Value})
				cv.Ptype[value.Ptype] = temp

				if len(config.ValuesRefs[value.Ptype]) != 0 {
					temp.Ref = config.ValuesRefs[value.Ptype]
					cv.Ptype[value.Ptype] = temp
				}
			}
		}
	}
	if len(cv.Ptype) == 0 && len(cv.Group) == 0 {
		d.Println()
		return nil
	}

	c, err := yaml.Marshal(cv)
	if err != nil {
		return err
	}
	d.Println(string(c))
	return nil
}

// ConfigSet sets an app's config variables.
func (d *DryccCmd) ConfigSet(appID string, ptype string, group string, configVars []string, confirm string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	if ptype == "" && group == "" {
		group = "global"
	}
	err = configConfirmAction(s.Client, appID, ptype, group, confirm)
	if err != nil {
		return err
	}

	configMap, err := parseConfig(ptype, group, configVars)
	if err != nil {
		return err
	}

	d.Print("Creating config... ")

	quit := progress(d.WOut)
	configObj := api.Config{Values: configMap}
	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	d.Print("done\n\n")
	return nil
}

// ConfigUnset removes a config variable from an app.
func (d *DryccCmd) ConfigUnset(appID string, ptype string, group string, configVars []string, confirm string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if ptype == "" && group == "" {
		group = "global"
	}
	err = configConfirmAction(s.Client, appID, ptype, group, confirm)
	if err != nil {
		return err
	}

	d.Print("Removing config... ")

	quit := progress(d.WOut)
	valuesMaps := []api.ConfigValue{}
	for _, configVar := range configVars {
		valuesMap := api.ConfigValue{
			Ptype: ptype,
			Group: group,
			KV: api.KV{
				Name:  configVar,
				Value: "",
			},
		}
		valuesMaps = append(valuesMaps, valuesMap)
	}

	configObj := api.Config{Values: valuesMaps}
	_, err = config.Set(s.Client, appID, configObj)

	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")
	return nil
}

// ConfigPull pulls an app's config to a file.
func (d *DryccCmd) ConfigPull(appID, ptype, group, fileName string, interactive bool, overwrite bool) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if ptype == "" && group == "" {
		group = "global"
	}
	if ptype != "" && group != "" {
		d.Println("Only one of ptype and group can be selected.")
		return nil
	}
	configVars, err := config.List(s.Client, appID, -1)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	stat, err := os.Stdout.Stat()

	if err != nil {
		return err
	}
	configValues := []api.ConfigValue{}

	for _, value := range configVars.Values {
		if (ptype != "" && value.Ptype == ptype) ||
			(group != "" && value.Group == group) {
			configValues = append(configValues, value)
		}
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		d.Print(formatConfig(configValues))
		return nil
	}

	if !overwrite {
		if _, err := os.Stat(fileName); err == nil {
			return fmt.Errorf("%s already exists, pass -o to overwrite", fileName)
		}
	}

	if interactive {
		contents, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		localConfigVars := strings.Split(string(contents), "\n")
		configMap, err := parseEnv(localConfigVars[:len(localConfigVars)-1])
		if err != nil {
			return err
		}
		for _, value := range configValues {
			localValue, ok := configMap[value.Name]
			if ok {
				if value.Value != localValue {
					var confirm string
					d.Printf("%s: overwrite %s with %s? (y/N) ", value.Name, localValue, value)

					fmt.Scanln(&confirm)

					if strings.ToLower(confirm) == "y" {
						configMap[value.Name] = value.Value
					}
				}
			} else {
				configMap[value.Name] = value.Value
			}
		}
		return os.WriteFile(fileName, []byte(formatEnv(configMap)), 0755)
	}
	return os.WriteFile(fileName, []byte(formatConfig(configValues)), 0755)
}

// ConfigPush pushes an app's config from a file.
func (d *DryccCmd) ConfigPush(appID, ptype string, group string, fileName string, confirm string) error {
	stat, err := os.Stdin.Stat()

	if err != nil {
		return err
	}
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	if ptype == "" && group == "" {
		group = "global"
	}
	var contents []byte

	if (stat.Mode() & os.ModeCharDevice) == 0 {

		err = configConfirmActionStdin(s.Client, appID, ptype, group, confirm)
		if err != nil {
			return err
		}
		buffer := new(bytes.Buffer)
		buffer.ReadFrom(os.Stdin)
		contents = buffer.Bytes()
	} else {
		err = configConfirmAction(s.Client, appID, ptype, group, confirm)
		if err != nil {
			return err
		}

		contents, err = os.ReadFile(fileName)

		if err != nil {
			return err
		}
	}

	file := strings.Split(string(contents), "\n")
	config := []string{}

	for _, configVar := range file {
		// If file has CRLF encoding, the default on windows, strip the CR
		configVar = strings.Trim(configVar, "\r")
		if len(configVar) > 0 {
			config = append(config, configVar)
		}
	}

	return d.ConfigSet(appID, ptype, group, config, "yes")
}

func (d *DryccCmd) ConfigAttach(appID string, ptype string, groups string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Attach config... ")

	quit := progress(d.WOut)
	gs := strings.Split(groups, ",")
	refs := api.ValuesRefs{
		ptype: gs,
	}
	configObj := api.Config{ValuesRefs: refs}
	_, err = config.Set(s.Client, appID, configObj)

	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")
	return nil
}

func (d *DryccCmd) ConfigDetach(appID string, ptype string, groups string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Detach config... ")

	quit := progress(d.WOut)
	gs := strings.Split(groups, ",")
	refs := api.ValuesRefs{
		ptype: gs,
	}
	configObj := api.Config{ValuesRefs: refs}
	err = config.Detach(s.Client, appID, configObj)

	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")
	return nil
}

func parseConfig(ptype, group string, configVars []string) ([]api.ConfigValue, error) {
	configMap := []api.ConfigValue{}

	regex := regexp.MustCompile(`^([A-z0-9_\-\.]+)=([\s\S]*)$`)
	for _, config := range configVars {
		// Skip config that starts with an comment
		if config[0] == '#' {
			continue
		}

		if regex.MatchString(config) {
			captures := regex.FindStringSubmatch(config)
			value := api.ConfigValue{
				Ptype: ptype,
				Group: group,
				KV: api.KV{
					Name:  captures[1],
					Value: captures[2],
				},
			}
			configMap = append(configMap, value)

		} else {
			return nil, fmt.Errorf("'%s' does not match the pattern 'key=var', ex: MODE=test", config)
		}
	}

	return configMap, nil
}

func parseEnv(configVars []string) (map[string]interface{}, error) {
	configMap := make(map[string]interface{})

	regex := regexp.MustCompile(`^([A-z0-9_\-\.]+)=([\s\S]*)$`)
	for _, config := range configVars {
		// Skip config that starts with an comment
		if config[0] == '#' {
			continue
		}

		if regex.MatchString(config) {
			captures := regex.FindStringSubmatch(config)
			configMap[captures[1]] = captures[2]
		} else {
			return nil, fmt.Errorf("'%s' does not match the pattern 'key=var', ex: MODE=test", config)
		}
	}

	return configMap, nil
}

func formatEnv(configVars map[string]interface{}) string {
	var formattedConfig string

	keys := *sortKeys(configVars)
	for _, key := range keys {
		formattedConfig += fmt.Sprintf("%s=%v\n", key, configVars[key])
	}

	return formattedConfig
}

func formatConfig(configVars []api.ConfigValue) string {
	var formattedConfig string

	values := sortConfigValues(configVars)
	for _, value := range values {
		formattedConfig += fmt.Sprintf("%s=%v\n", value.Name, value.Value)
	}

	return formattedConfig
}

func configConfirmAction(s *drycc.Client, appID string, ptype string, group string, confirm string) error {

	if ptype != "" && group != "" {
		fmt.Println("Only one of ptype and group can be selected.")
		return nil
	} else if ptype == "" && group == "" {
		group = "global"
	}

	appSettings, _ := appsettings.List(s, appID)
	autodeploy := true
	if appSettings.Autodeploy != nil && !*appSettings.Autodeploy {
		autodeploy = false
	}
	if ptype == "" && group == "" && (confirm == "" || confirm != "yes") && autodeploy {
		fmt.Printf(` !    WARNING: Potentially Config Action
 !    This command will deploy all processes of the application
 !    To proceed, type "yes" !

> `)

		fmt.Scanln(&confirm)
		if confirm != "yes" {
			return fmt.Errorf("cancel the config action")
		}
	}
	return nil
}

func configConfirmActionStdin(s *drycc.Client, appID string, ptype string, group string, confirm string) error {
	var reader *bufio.Reader
	if runtime.GOOS == "windows" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open("/dev/tty")
		if err != nil {
			return err
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	}

	if ptype != "" && group != "" {
		fmt.Println("Only one of ptype and group can be selected.")
		return nil
	}

	appSettings, _ := appsettings.List(s, appID)
	autodeploy := true
	if appSettings.Autodeploy != nil && !*appSettings.Autodeploy {
		autodeploy = false
	}

	if ptype == "" && group == "" && (confirm == "" || confirm != "yes") && autodeploy {
		fmt.Printf(` !    WARNING: Potentially Config Action
 !    This command will deploy all processes of the application
 !    To proceed, type "yes" !

> `)

		confirm, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		confirm = strings.TrimSpace(confirm)
		if confirm != "yes" {
			return fmt.Errorf("cancel the config action")
		}
	}
	return nil
}
