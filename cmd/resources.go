package cmd

import (
	"fmt"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/resources"
	"github.com/drycc/pkg/prettyprint"
	"io"
	"regexp"
	"strings"
)

// ResourcesCreate create a resource for the application
func (d *DryccCmd) ResourcesCreate(appID, plan string, name string, params []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Creating %s to %s... ", name, appID)

	paramsMap, err := parseParams(params)
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
		d.Println("Could not find any resources")
	} else {
		printResources(d, appID, resources, d.WOut)
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
	// todo format data json to yaml
	printResourceDetail(d, appID, resource, d.WOut)
	//d.Println(resource)
	return nil
}

// ResourceDelete delete a resource from the application
func (d *DryccCmd) ResourceDelete(appID, name string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Deleting %s from %s... ", name, appID)

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
func (d *DryccCmd) ResourcePut(appID, plan string, name string, params []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Updating %s to %s... ", name, appID)

	paramsMap, err := parseParams(params)
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
	bindAction := api.Binding{BindAction: "bind"}
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
	bindAction := api.Binding{BindAction: "unbind"}
	_, err = resources.Binding(s.Client, appID, name, bindAction)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n")

	return nil
}

// printResources format Resources data
func printResources(d *DryccCmd, appID string, resources api.Resources, wOut io.Writer) {

	fmt.Fprintf(wOut, "=== %s resources\n", appID)
	resourceNames := make([]string, len(resources))

	for _, resource := range resources {
		resourceNames = append(resourceNames, resource.Name)
	}
	lenResourceNames := sliceMaxLen(resourceNames) + 5

	for _, resource := range resources {
		spaces := strings.Repeat(" ", lenResourceNames-len(resource.Name))
		fmt.Fprintf(wOut, "%s%s%s\n", resource.Name, spaces, resource.Plan)
	}
}

func printResourceDetail(d *DryccCmd, appID string, resource api.Resource, wOut io.Writer) {
	d.Printf("=== %s resource %s\n", appID, resource.Name)

	dataMap := make(map[string]string)
	for key, value := range resource.Data {
		dataMap[key+":"] = fmt.Sprintf("%v", value)
	}
	optionsMap := make(map[string]string)
	for key, value := range resource.Options {
		optionsMap[key+":"] = fmt.Sprintf("%v", value)
	}
	lenDataMap := mapMaxLen(dataMap)
	lenOptionsMap := mapMaxLen(optionsMap)
	tempArray := []int{lenDataMap, lenOptionsMap, len("binding:")}
	max := maxNum(tempArray...) + 5
	d.Print(fmt.Sprintf("plan:%s%s\n", strings.Repeat(" ", max-len("plan:")), resource.Plan))
	d.Print(fmt.Sprintf("status:%s%s\n", strings.Repeat(" ", max-len("status:")), resource.Status))
	d.Print(fmt.Sprintf("binding:%s%s\n", strings.Repeat(" ", max-len("binding:")), resource.Binding))

	if lenDataMap != 0 {
		d.Println()
		d.Print(prettyprint.PrettyTabs(dataMap, max-lenDataMap))
	}
	if lenOptionsMap != 0 {
		d.Println()
		d.Print(prettyprint.PrettyTabs(optionsMap, max-lenOptionsMap))
	}
}

func mapMaxLen(msg map[string]string) int {
	// find the longest key so we know how much padding to use
	max := 0
	for key := range msg {
		if len(key) > max {
			max = len(key)
		}
	}
	return max
}

func sliceMaxLen(msgs []string) int {
	// find the longest member so we know how much padding to use
	max := 0
	for _, msg := range msgs {
		if len(msg) > max {
			max = len(msg)
		}
	}
	return max
}

func maxNum(tempArray ...int) int {
	// find the longest num so we know how much padding to use
	max := 0
	for _, temp := range tempArray {
		if max < temp {
			max = temp
		}
	}
	return max
}

// parseParams transfer params to map
func parseParams(params []string) (map[string]interface{}, error) {
	paramsMap := make(map[string]interface{})
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
