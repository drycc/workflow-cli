package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/builds"
)

// BuildsInfo lists an app's builds.
func (d *DryccCmd) BuildsInfo(appID string, version int) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	build, err := builds.Get(s.Client, appID, version)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"App:", build.App})
	table.Append([]string{"Sha:", build.Sha})
	table.Append([]string{"UUID:", build.UUID})
	table.Append([]string{"Owner:", build.Owner})
	table.Append([]string{"Image:", build.Image})
	table.Append([]string{"Stack:", build.Stack})
	table.Append([]string{"Created:", build.Created})
	table.Append([]string{"Updated:", build.Updated})
	table.Render()

	return nil
}

// BuildsCreate creates a build for an app.
func (d *DryccCmd) BuildsCreate(appID, image, stack, procfile, dryccpath, confirm string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	procfileMap := make(map[string]string)
	if _, err := os.Stat(procfile); err == nil {
		contents, err := os.ReadFile(procfile)
		if err != nil {
			return err
		}

		if procfileMap, err = parseProcfile(contents); err != nil {
			return err
		}
	}

	dryccfileMap := make(map[string]interface{})
	if info, err := os.Stat(dryccpath); err == nil && info.IsDir() {
		if dryccfileMap, err = drycc.ParseDryccfile(dryccpath); err != nil {
			return err
		}
	}
	// check procfileMap dryccfileMap stack is exist
	err = buildConfirmAction(s.Client, appID, procfileMap, dryccfileMap, confirm)
	if err != nil {
		return err
	}
	d.Print("Creating build... ")
	quit := progress(d.WOut)
	_, err = builds.New(s.Client, appID, image, stack, procfileMap, dryccfileMap)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")

	return nil
}

func (d *DryccCmd) BuildsFetch(appID string, version int, procfile, dryccpath, confirm string, save bool) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	build, err := builds.Get(s.Client, appID, version)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	// Confirm again
	err = buildFetchConfirmAction(confirm, procfile, dryccpath, save)
	if err != nil {
		return err
	}

	if len(build.Procfile) != 0 {
		err := writeProcfileToPath(d, procfile, build.Procfile, save)
		if err != nil {
			return fmt.Errorf("failed to write Procfile: %w", err)
		}
	}

	if len(build.Dryccfile) != 0 {
		err := writeDryccfileToPath(d, dryccpath, build.Dryccfile, save)
		if err != nil {
			return fmt.Errorf("failed to write Dryccfile: %w", err)
		}
	}
	if save {
		d.Println("done")
	}
	return nil
}

func writeProcfileToPath(d *DryccCmd, procfile string, Procinfo map[string]string, save bool) error {
	if save {
		os.Remove(procfile)
		err := os.WriteFile(procfile, []byte(d.toYamlString(Procinfo, 2)), 0664)
		return err
	}
	d.Println("---\n# Source:", procfile)
	d.Println(d.toYamlString(Procinfo, 2))
	return nil
}

func writeDryccfileToPath(d *DryccCmd, dryccpath string, dryccfile map[string]interface{}, save bool) error {
	// Create the directory if it doesn't exist
	if save {
		os.Remove(dryccpath)
	}
	os.ReadDir(dryccpath)
	err := os.MkdirAll(dryccpath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write config section
	if config, ok := dryccfile["config"].(map[string]interface{}); ok {
		configDir := filepath.Join(dryccpath, "config")
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		for key, values := range config {
			envFilePath := filepath.Join(configDir, key)
			var content string
			for k, v := range values.(map[string]interface{}) {
				// Append each key-value pair with newline
				content += fmt.Sprintf("%s=%v\n", k, v)
			}
			// Write accumulated content once
			if save {
				err := os.WriteFile(envFilePath, []byte(content), 0664)
				if err != nil {
					return fmt.Errorf("failed to write env file: %w", err)
				}
			} else {
				d.Println("---\n# Source:", envFilePath)
				d.Println(content)
			}
		}
	}

	// Write pipeline section
	if pipeline, ok := dryccfile["pipeline"].(map[string]interface{}); ok {
		for fileName, pipelineConfig := range pipeline {
			filePath := filepath.Join(dryccpath, fileName)
			var buf bytes.Buffer
			encoder := yaml.NewEncoder(&buf)
			encoder.SetIndent(2)
			if err := encoder.Encode(pipelineConfig); err != nil {
				return fmt.Errorf("failed to marshal pipeline config: %w", err)
			}
			yamlContent := buf.Bytes()
			if save {
				err = os.WriteFile(filePath, yamlContent, 0664)
				if err != nil {
					return fmt.Errorf("failed to write pipeline file: %w", err)
				}
			} else {
				d.Println("---\n# Source:", filePath)
				d.Println(string(yamlContent))
			}
		}
	}

	return nil
}

func parseProcfile(procfile []byte) (map[string]string, error) {
	procfileMap := make(map[string]string)
	return procfileMap, yaml.Unmarshal(procfile, &procfileMap)
}

func buildConfirmAction(c *drycc.Client, appID string, procfileMap map[string]string,
	dryccfileMap map[string]interface{}, confirm string) error {

	build, _ := builds.Get(c, appID, 0)

	if ((len(build.Procfile) != 0 && len(procfileMap) == 0) || (len(build.Dryccfile) != 0 && len(dryccfileMap) == 0)) && (confirm == "" || confirm != "yes") {
		// hint
		fmt.Printf(` !    WARNING: Potentially Build Create Action
 !    The Procfile or drycc file is empty, not last time
 !    To proceed, type "yes" !

> `)

		fmt.Scanln(&confirm)
		if confirm != "yes" {
			return fmt.Errorf("cancel the build create action")
		}
	}
	return nil
}

func buildFetchConfirmAction(confirm, procfile, dryccpath string, save bool) error {

	if save && (confirm == "" || confirm != "yes") {
		// hint
		msg := fmt.Sprintf(" !    WARNING: Potentially Build Fetch Action\n"+
			" !    This operation will overwrite the current \x1b[1m%s\x1b[0m or \x1b[1m%s\x1b[0m locally\n"+
			" !    To proceed, type \"yes\" !\n\n> ", procfile, dryccpath)

		fmt.Print(msg)
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			return fmt.Errorf("cancel the build fetch action")
		}
	}
	return nil
}
