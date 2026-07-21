package commands

import (
	"fmt"
	"os"
	"sort"

	"github.com/drycc/controller-sdk-go/addons"
	"github.com/drycc/workflow-cli/internal/loader"
	"github.com/drycc/workflow-cli/pkg/coder"
	"github.com/drycc/workflow-cli/pkg/settings"
	"sigs.k8s.io/yaml"
)

// AddonsList lists addon instances attached to an app.
func (d *DryccCmd) AddonsList(appID string, results int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	if results == defaultLimit {
		results = s.Limit
	}

	instances, count, err := addons.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count == 0 {
		d.Println(fmt.Sprintf("No addons found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"NAME", "KIND", "PLAN"})
		for _, instance := range instances {
			table.Append([]string{instance.Name, instance.Kind, instance.Plan})
		}
		table.Render()
	}
	return nil
}

// AddonsInfo shows detailed information about an addon instance.
func (d *DryccCmd) AddonsInfo(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	info, err := addons.Get(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	c := &coder.AddonCoder{Info: info}
	yamlBytes, err := c.Encode()
	if err != nil {
		return err
	}
	d.Println(string(yamlBytes))
	return nil
}

// AddonsConnection shows the connection details for an addon instance.
func (d *DryccCmd) AddonsConnection(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	connection, err := addons.Connection(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	keys := make([]string, 0, len(connection))
	for key := range connection {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	table := d.getDefaultFormatTable([]string{"KEY", "VALUE"})
	for _, key := range keys {
		table.Append([]string{key, fmt.Sprintf("%v", connection[key])})
	}
	table.Render()
	return nil
}

// AddonsApply creates or updates an addon instance from a YAML file.
func (d *DryccCmd) AddonsApply(appID, filePath string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return err
	}

	c := &coder.AddonCoder{}
	if err := c.Decode(jsonData); err != nil {
		return fmt.Errorf("invalid addon configuration: %w", err)
	}
	if c.Name == "" {
		return fmt.Errorf("invalid addon configuration: missing metadata.name")
	}
	if c.Request.Kind == "" {
		return fmt.Errorf("invalid addon configuration: missing kind")
	}
	if c.Request.Plan == "" {
		return fmt.Errorf("invalid addon configuration: missing spec.plan")
	}

	d.Printf("Applying addon %s to %s... ", c.Name, appID)

	quit := progress(d.WOut)
	_, err = addons.Upsert(s.Client, appID, c.Name, c.Request)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// AddonsRemove removes an addon instance from an app.
func (d *DryccCmd) AddonsRemove(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	d.Printf("Removing addon %s from %s... ", name, appID)

	quit := progress(d.WOut)
	err = addons.Delete(s.Client, appID, name)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// AddonsClassesList lists all addon classes in the catalog.
func (d *DryccCmd) AddonsClassesList(results int) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}
	if results == defaultLimit {
		results = s.Limit
	}

	classes, count, err := addons.ListClasses(s.Client, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count == 0 {
		d.Println("No addon classes found.")
	} else {
		table := d.getDefaultFormatTable([]string{"NAME", "KIND", "DESCRIPTION"})
		for _, class := range classes {
			table.Append([]string{class.Name, class.Kind, class.Description})
		}
		table.Render()
	}
	return nil
}

// AddonsClassesInfo shows detailed information about an addon class.
func (d *DryccCmd) AddonsClassesInfo(name string) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	class, err := addons.GetClass(s.Client, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	yamlBytes, err := coder.EncodeClass(class)
	if err != nil {
		return err
	}
	d.Println(string(yamlBytes))
	return nil
}
