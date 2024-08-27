package cmd

import (
	"fmt"
	"regexp"
	"strconv"

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
	if len(services) > 0 {
		table := d.getDefaultFormatTable([]string{"PTYPE", "PORT", "PROTOCOL", "TARGET-PORT", "DOMAIN"})
		for _, service := range services {
			for _, port := range service.Ports {
				table.Append([]string{
					service.Ptype,
					fmt.Sprint(port.Port),
					port.Protocol,
					fmt.Sprint(port.TargetPort),
					service.Domain,
				})
			}
		}
		table.Render()
	} else {
		d.Println(fmt.Sprintf("No services found in %s app.", appID))
	}
	return nil
}

// ServicesAdd adds a service to an app.
func (d *DryccCmd) ServicesAdd(appID, ptype string, ports string, protocol string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	portArray, err := parsePorts(ports)
	if err != nil {
		return err
	}
	d.Printf("Adding %s (%d) to %s... ", ptype, portArray[0], appID)

	quit := progress(d.WOut)
	err = services.New(s.Client, appID, ptype, portArray[0], protocol, portArray[1])
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// ServicesRemove removes a service for Ptype registered with an app.
func (d *DryccCmd) ServicesRemove(appID, ptype string, protocol string, port int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Removing %s from %s... ", ptype, appID)

	quit := progress(d.WOut)
	err = services.Delete(s.Client, appID, ptype, protocol, port)
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
