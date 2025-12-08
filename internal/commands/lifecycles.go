package commands

import (
	"fmt"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/config"
	"github.com/drycc/workflow-cli/internal/loader"
)

func getLifecycleHandlerString(ptype, handler, signal string, lifecycleHandler *api.LifecycleHandler) string {
	if lifecycleHandler.Exec != nil {
		return fmt.Sprintf("%s %s exec %v %s", handler, ptype, lifecycleHandler.Exec.Command, signal)
	} else if lifecycleHandler.Sleep != nil {
		return fmt.Sprintf("%s %s sleep %v %s", handler, ptype, lifecycleHandler.Sleep, signal)
	} else if lifecycleHandler.TCPSocket != nil {
		return fmt.Sprintf("%s %s tcp-socket port=%v %s", handler, ptype, lifecycleHandler.TCPSocket.Port, signal)
	} else if lifecycleHandler.HTTPGet != nil {
		return fmt.Sprintf(
			"%s %s http-get headers=%v path=%s port=%d %s",
			handler,
			ptype,
			lifecycleHandler.HTTPGet.HTTPHeaders,
			lifecycleHandler.HTTPGet.Path,
			lifecycleHandler.HTTPGet.Port,
			signal,
		)
	}
	return ""
}

// LifecyclesList lists an app's lifecycles.
func (d *DryccCmd) LifecyclesList(appID, ptype string, version int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID, version)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if len(config.Lifecycle) == 0 {
		d.Println(fmt.Sprintf("No lifecycle found in %s app.", appID))
		return nil
	}

	ptypes := []string{}
	if ptype != "" {
		ptypes = append(ptypes, ptype)
	} else {
		for ptype := range config.Lifecycle {
			ptypes = append(ptypes, ptype)
		}
	}

	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"App:", config.App})
	table.Append([]string{"UUID:", config.UUID})
	table.Append([]string{"Owner:", config.Owner})
	table.Append([]string{"Created:", d.formatTime(config.Created)})
	table.Append([]string{"Updated:", d.formatTime(config.Updated)})
	table.Append([]string{"Lifecycle:"})
	for _, ptype := range sortPtypes(ptypes) {
		if lifecycle, ok := config.Lifecycle[ptype]; ok {
			table.Append([]string{"", fmt.Sprintf("stopSignal=%s", lifecycle.StopSignal)})
			table.Append([]string{"", getLifecycleHandlerString(ptype, "postStart", lifecycle.StopSignal, *lifecycle.PostStart)})
			table.Append([]string{"", getLifecycleHandlerString(ptype, "preStop", lifecycle.StopSignal, *lifecycle.PreStop)})
		}
	}
	table.Render()

	return nil
}

// LifecyclesSet sets an app's lifecycle.
func (d *DryccCmd) LifecyclesSet(appID, ptype string, lifecycle *api.Lifecycle) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Print("Applying lifecycle... ")

	quit := progress(d.WOut)
	configObj := api.Config{Lifecycle: make(map[string]*api.Lifecycle)}
	configObj.Lifecycle[ptype] = lifecycle

	_, err = config.Set(s.Client, appID, configObj, true)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.LifecyclesList(appID, ptype, -1)
}

// LifecyclesUnset removes an app's lifecycle.
func (d *DryccCmd) LifecyclesUnset(appID, ptype string, handlers []string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	d.Print("Applying lifecycle... ")

	quit := progress(d.WOut)

	configObj := api.Config{Lifecycle: make(map[string]*api.Lifecycle)}
	lifecycle := &api.Lifecycle{}
	var nullLifecycleHandler *api.LifecycleHandler = nil
	for _, handler := range handlers {
		switch handler {
		case "postStart":
			lifecycle.PostStart = &nullLifecycleHandler
		case "preStop":
			lifecycle.PreStop = &nullLifecycleHandler
		default:
			return fmt.Errorf("unknown lifecycle handler: %s", handler)
		}
	}
	configObj.Lifecycle[ptype] = lifecycle

	_, err = config.Set(s.Client, appID, configObj, true)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.LifecyclesList(appID, ptype, -1)
}
