package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"sigs.k8s.io/yaml"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/resources"
	"github.com/drycc/workflow-cli/settings"
)

// ResourceServices list resource service
func (d *DryccCmd) ResourcesServices(results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}
	services, count, err := resources.Services(s.Client, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if count == 0 {
		d.Println("Could not find any services")
	} else {
		table := d.getDefaultFormatTable([]string{"ID", "NAME", "UPDATEABLE"})
		for _, service := range services {
			table.Append([]string{
				service.ID,
				service.Name,
				strconv.FormatBool(service.Updateable),
			})
		}
		table.Render()
	}
	return nil
}

// ResourcePlans list resource plans
func (d *DryccCmd) ResourcesPlans(serviceName string, results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}
	plans, count, err := resources.Plans(s.Client, serviceName, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if count == 0 {
		d.Println(fmt.Sprintf("Could not find any plans in %s service.", serviceName))
	} else {
		table := d.getDefaultFormatTable([]string{"ID", "NAME", "DESCRIPTION"})
		for _, plan := range plans {
			table.Append([]string{
				plan.ID,
				plan.Name,
				plan.Description,
			})
		}
		table.Render()
	}
	return nil
}

// ResourcesCreate create a resource for the application
func (d *DryccCmd) ResourcesCreate(appID, plan string, name string, params []string, values string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	paramsMap := make(map[string]interface{})
	if values != "" {
		valueFile, err := os.Stat(values)
		if err != nil {
			return err
		}
		if valueFile.Size() == 0 {
			return fmt.Errorf("%s is empty", values)
		}
		rawValues, err := os.ReadFile(values)
		if err != nil {
			return err
		}
		parsed := make(map[string]interface{})
		err = yaml.Unmarshal(rawValues, &parsed)
		if err != nil {
			return err
		}
		paramsMap["rawValues"] = base64.StdEncoding.EncodeToString([]byte(rawValues))
	}

	d.Printf("Creating %s to %s... ", name, appID)

	paramsMap, err = parseParams(paramsMap, params)
	if err != nil {
		return err
	}

	quit := progress(d.WOut)
	resource := api.Resource{
		Name:    name,
		Plan:    plan,
		Options: paramsMap,
	}
	_, err = resources.Create(s.Client, appID, resource)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// ResourcesList list resources in the application
func (d *DryccCmd) ResourcesList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}
	resources, count, err := resources.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if count == 0 {
		d.Println(fmt.Sprintf("No resources found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"UUID", "NAME", "OWNER", "PLAN", "UPDATED"})
		for _, resource := range resources {
			table.Append([]string{
				resource.UUID,
				resource.Name,
				resource.Owner,
				resource.Plan,
				resource.Updated,
			})
		}
		table.Render()
	}
	return nil
}

// ResourceGet describe a resource from the application
func (d *DryccCmd) ResourceGet(appID, name string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	//d.Printf(" %s from %s... ", name, appID)

	resource, err := resources.Get(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"App:", appID})
	table.Append([]string{"UUID:", resource.UUID})
	table.Append([]string{"Name:", resource.Name})
	table.Append([]string{"Plan:", resource.Plan})
	table.Append([]string{"Owner:", resource.Owner})
	table.Append([]string{"Status:", resource.Status})
	table.Append([]string{"Binding:", resource.Binding})
	table.Append([]string{"Data:"})
	for _, key := range *sortKeys(resource.Data) {
		table.Append([]string{"", fmt.Sprintf("%s:", key), fmt.Sprintf("%s", resource.Data[key])})
	}
	table.Append([]string{"Options:"})
	for _, key := range *sortKeys(resource.Options) {
		table.Append([]string{"", fmt.Sprintf("%s:", key), fmt.Sprintf("%s", resource.Options[key])})
	}
	table.Append([]string{"Message:", safeGetString(resource.Message)})
	table.Append([]string{"Created:", resource.Created})
	table.Append([]string{"Updated:", resource.Updated})
	table.Render()
	return nil
}

// ResourceDelete delete a resource from the application
func (d *DryccCmd) ResourceDelete(appID, name, confirm string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if confirm == "" {
		d.Printf(` !    WARNING: Potentially Destructive Action
 !    This command will destroy the resource: %s
 !    To proceed, type "%s" or re-run this command with --confirm=%s

> `, name, name, name)

		fmt.Scanln(&confirm)
	}

	if confirm != name {
		return fmt.Errorf("resource %s does not match confirm %s, aborting", name, confirm)
	}

	d.Printf("Destroying %s...\n", name)

	quit := progress(d.WOut)
	err = resources.Delete(s.Client, appID, name)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// ResourcePut update a resource for the application
func (d *DryccCmd) ResourcePut(appID, plan string, name string, params []string, values string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	paramsMap := make(map[string]interface{})
	if values != "" {
		valueFile, err := os.Stat(values)
		if err != nil {
			return err
		}
		if valueFile.Size() == 0 {
			return fmt.Errorf("%s is empty", values)
		}
		rawValues, err := os.ReadFile(values)
		if err != nil {
			return err
		}
		parsed := make(map[string]interface{})
		err = yaml.Unmarshal(rawValues, &parsed)
		if err != nil {
			return err
		}
		paramsMap["rawValues"] = base64.StdEncoding.EncodeToString([]byte(rawValues))
	}

	d.Printf("Updating %s to %s... ", name, appID)

	paramsMap, err = parseParams(paramsMap, params)
	if err != nil {
		return err
	}

	quit := progress(d.WOut)
	resource := api.Resource{
		Plan:    plan,
		Options: paramsMap,
	}
	_, err = resources.Put(s.Client, appID, name, resource)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// ResourceBind mount a resource to process of the application
func (d *DryccCmd) ResourceBind(appID string, name string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Binding resource... ")

	quit := progress(d.WOut)
	bindAction := api.ResourceBinding{BindAction: "bind"}
	_, err = resources.Binding(s.Client, appID, name, bindAction)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n")

	return nil
}

// ResourceUnbind resource a resource the application
func (d *DryccCmd) ResourceUnbind(appID string, name string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Unbinding resource... ")

	quit := progress(d.WOut)
	bindAction := api.ResourceBinding{BindAction: "unbind"}
	_, err = resources.Binding(s.Client, appID, name, bindAction)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n")

	return nil
}

// parseParams transfer params to map
func parseParams(paramsMap map[string]interface{}, params []string) (map[string]interface{}, error) {
	regex := regexp.MustCompile(`^([A-z_]+[A-z0-9_]*[\.{1}[A-z0-9_]+]*)=([\s\S]*)$`)
	for _, param := range params {
		if regex.MatchString(param) {
			captures := regex.FindStringSubmatch(param)
			paramsMap[captures[1]] = captures[2]
		} else {
			return nil, fmt.Errorf("'%s' does not match the pattern 'key=var', ex: MODE=test", param)
		}
	}

	return paramsMap, nil
}
