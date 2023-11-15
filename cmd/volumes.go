package cmd

import (
	"fmt"
	"regexp"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/volumes"
)

// VolumesList list volumes in the application
func (d *DryccCmd) VolumesList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}
	//fmt.Println("-----app-----", appID)  # debug
	volumes, count, err := volumes.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if count == 0 {
		d.Println("Could not find any volume.")
	} else {
		printVolumes(d, volumes)
	}
	return nil
}

// printVolumes format volume data
func printVolumes(d *DryccCmd, volumes api.Volumes) {
	table := d.getDefaultFormatTable([]string{"NAME", "OWNER", "TYPE", "PATH", "SIZE"})
	for _, volume := range volumes {
		if len(volume.Path) > 0 {
			for _, key := range *sortKeys(volume.Path) {
				table.Append([]string{volume.Name, volume.Owner, key, fmt.Sprintf("%v", volume.Path[key]), volume.Size})
			}
		} else {
			table.Append([]string{volume.Name, volume.Owner, "", "", volume.Size})
		}
	}
	table.Render()
}

// VolumesCreate create a volume for the application
func (d *DryccCmd) VolumesCreate(appID, name, size string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	regex := regexp.MustCompile("^([1-9][0-9]*[gG])$")
	if !regex.MatchString(size) {
		return fmt.Errorf(`%s doesn't fit format #unit
Examples: 2G 2g`, size)
	}

	d.Printf("Creating %s to %s... ", name, appID)

	quit := progress(d.WOut)
	volume := api.Volume{
		Name: name,
		Size: size,
	}
	_, err = volumes.Create(s.Client, appID, volume)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// VolumesExpand create a volume for the application
func (d *DryccCmd) VolumesExpand(appID, name, size string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	regex := regexp.MustCompile("^([1-9][0-9]*[gG])$")
	if !regex.MatchString(size) {
		return fmt.Errorf(`%s doesn't fit format #unit
Examples: 2G 2g`, size)
	}

	d.Printf("Expand %s to %s... ", name, appID)

	quit := progress(d.WOut)
	volume := api.Volume{
		Name: name,
		Size: size,
	}
	_, err = volumes.Expand(s.Client, appID, volume)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// VolumesDelete delete a volume from the application
func (d *DryccCmd) VolumesDelete(appID, name string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Deleting %s from %s... ", name, appID)

	quit := progress(d.WOut)
	err = volumes.Delete(s.Client, appID, name)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// VolumesMount mount a volume to process of the application
func (d *DryccCmd) VolumesMount(appID string, name string, volumeVars []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	volumeMap, err := parseVolume(volumeVars)
	if err != nil {
		return err
	}

	d.Print("Mounting volume... ")

	quit := progress(d.WOut)
	volumeObj := api.Volume{Path: volumeMap}
	_, err = volumes.Mount(s.Client, appID, name, volumeObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n")
	d.Print("The pods should be restart, please check the pods up or not.\n")

	return nil
}

// VolumesUnmount unmount a volume from process of the application
func (d *DryccCmd) VolumesUnmount(appID string, name string, volumeVars []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	valuesMap := make(map[string]interface{})
	for _, volumeVar := range volumeVars {
		valuesMap[volumeVar] = nil
	}
	if err != nil {
		return err
	}

	d.Print("Unmounting volume... ")

	quit := progress(d.WOut)
	volumeObj := api.Volume{Path: valuesMap}
	_, err = volumes.Mount(s.Client, appID, name, volumeObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n")
	d.Print("The pods should be restart, please check the pods up or not.\n")

	return nil
}

func parseVolume(volumeVars []string) (map[string]interface{}, error) {
	volumeMap := make(map[string]interface{})

	regex := regexp.MustCompile(`^([A-z_]+[A-z0-9_]*)=([\s\S]*)$`)
	for _, volume := range volumeVars {
		if regex.MatchString(volume) {
			captures := regex.FindStringSubmatch(volume)
			volumeMap[captures[1]] = captures[2]
		} else {
			return nil, fmt.Errorf("'%s' does not match the pattern 'key=var', ex: MODE=test", volume)
		}
	}

	return volumeMap, nil
}
