package cmd

import (
  "fmt"

  "github.com/drycc/pkg/prettyprint"

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
  servicesMap := make(map[string]string)
  if len(services) > 0 {
    for _, service := range services {
      servicesMap[service.ProcfileType] = fmt.Sprintf("%v", service.PathPattern)
    }
    d.Print(prettyprint.PrettyTabs(servicesMap, 5))
  }
  return nil
}

// ServicesAdd adds a service to an app.
func (d *DryccCmd) ServicesAdd(appID, procfileType string, pathPattern string) error {
  s, appID, err := load(d.ConfigFile, appID)

  if err != nil {
    return err
  }

  d.Printf("Adding %s (%s) to %s... ", procfileType, pathPattern, appID)

  quit := progress(d.WOut)
  _, err = services.New(s.Client, appID, procfileType, pathPattern)
  quit <- true
  <-quit
  if d.checkAPICompatibility(s.Client, err) != nil {
    return err
  }

  d.Println("done")
  return nil
}

// ServicesRemove removes a service for procfileType registered with an app.
func (d *DryccCmd) ServicesRemove(appID, procfileType string) error {
  s, appID, err := load(d.ConfigFile, appID)

  if err != nil {
    return err
  }

  d.Printf("Removing %s from %s... ", procfileType, appID)

  quit := progress(d.WOut)
  err = services.Delete(s.Client, appID, procfileType)
  quit <- true
  <-quit
  if d.checkAPICompatibility(s.Client, err) != nil {
    return err
  }

  d.Println("done")
  return nil
}
