package cmd

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/drycc/controller-sdk-go/services"
)

// ServicesList lists extra services for the app
func (d *DryccCmd) ServicesList(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	services, err := services.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== %s Services\n", appID)
	if len(services) > 0 {
		table := tablewriter.NewWriter(d.WOut)
		table.SetHeader([]string{"Type", "Name", "Port", "Protocol", "TargetPort"})
		for _, service := range services {
			for _, port := range service.Ports {
				table.Append([]string{service.ProcfileType, port.Name, fmt.Sprint(port.Port), port.Protocol, fmt.Sprint(port.TargetPort)})
			}
		}
		table.SetAutoMergeCellsByColumnIndex([]int{0})
		table.SetRowLine(true)
		table.Render()
	}
	return nil
}

// ServicesAdd adds a service to an app.
func (d *DryccCmd) ServicesAdd(appID, procfileType string, ports string, protocol string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	portArray, err := parsePorts(ports)
	if err != nil {
		return err
	}
	d.Printf("Adding %s (%d) to %s... ", procfileType, portArray[0], appID)

	quit := progress(d.WOut)
	err = services.New(s.Client, appID, procfileType, portArray[0], protocol, portArray[1])
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// ServicesRemove removes a service for procfileType registered with an app.
func (d *DryccCmd) ServicesRemove(appID, procfileType string, protocol string, port int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Removing %s from %s... ", procfileType, appID)

	quit := progress(d.WOut)
	err = services.Delete(s.Client, appID, procfileType, protocol, port)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// parsePorts transfer ports to [2]int
func parsePorts(param string) ([2]int, error) {
	var ports [2]int
	var err error
	regex := regexp.MustCompile(`(^[1-9]+[0-9_]+):([1-9]+[0-9_]+)$`)

	if regex.MatchString(param) {
		captures := regex.FindStringSubmatch(param)
		ports[0], _ = strconv.Atoi(captures[1])
		ports[1], _ = strconv.Atoi(captures[2])
	} else {
		err = fmt.Errorf("'%s' does not match the pattern 'port:targatPort', ex: 80:8000", param)
	}

	return ports, err
}
