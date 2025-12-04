package commands

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/volumes"
	"github.com/drycc/workflow-cli/internal/loader"
	"github.com/schollz/progressbar/v3"
	"sigs.k8s.io/yaml"
)

// VolumesList list volumes in the application
func (d *DryccCmd) VolumesList(appID string, results int) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}
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

// VolumesInfo get volume in the application
func (d *DryccCmd) VolumesInfo(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	volume, err := volumes.Get(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"UUID:", volume.UUID})
	table.Append([]string{"Name:", volume.Name})
	table.Append([]string{"Owner:", volume.Owner})
	table.Append([]string{"Type:", volume.Type})
	// table append path
	table.Append([]string{"Path:"})
	path, err := yaml.Marshal(volume.Path)
	if err != nil {
		return err
	}
	table.Append([]string{"", string(path)})
	// table append parameters
	table.Append([]string{"Parameters:"})
	parameters, err := yaml.Marshal(volume.Parameters)
	if err != nil {
		return err
	}
	table.Append([]string{"", string(parameters)})
	table.Append([]string{"Created: ", d.formatTime(volume.Created)})
	table.Append([]string{"Updated: ", d.formatTime(volume.Updated)})
	table.Render()
	return nil
}

// VolumesCreate create a volume for the application
func (d *DryccCmd) VolumesCreate(appID, name, vType, size string, parameters map[string]any) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
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
		Name:       name,
		Size:       size,
		Type:       vType,
		Parameters: parameters,
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
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
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
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
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

// VolumesServe Serve serves an app's volume.
func (d *DryccCmd) VolumesServe(appID, name string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()

	d.Print("Starting WebDAV file access service ")
	quit := progress(d.WOut)
	ctx, filer, err := volumes.Serve(parent, s.Client, appID, name)
	quit <- true
	<-quit
	if err != nil {
		return err
	}

	d.Print("\n\n")
	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"Endpoint:", filer["endpoint"]})
	table.Append([]string{"Username:", filer["username"]})
	table.Append([]string{"Password:", filer["password"]})
	table.Render()
	d.Print("\n")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)
	d.Printf("WebDAV service for volume %s is running. Press Ctrl+C to stop.\n", name)
	for {
		select {
		case <-signalChan:
			return nil
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(2 * time.Second)
		}
	}
}

// VolumesMount mount a volume to process of the application
func (d *DryccCmd) VolumesMount(appID string, name string, volumeVars []string) error {
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
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
	appID, s, err := loader.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	valuesMap := make(map[string]any)
	for _, volumeVar := range volumeVars {
		valuesMap[volumeVar] = nil
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

func parseVolume(volumeVars []string) (map[string]any, error) {
	volumeMap := make(map[string]any)
	regex := regexp.MustCompile(`^([a-z0-9]+(?:-[a-z0-9]+)*)=(\/([\w]+[\w-]*\/?)+)$`)
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

// parseVol format volume url
func parseVol(vol string) (string, string, error) {
	u, err := url.Parse(vol)
	if err != nil {
		return "", "", err
	}
	if u.Scheme != "vol" || u.Host == "" {
		return "", "", fmt.Errorf("vol %s format err", vol)
	}
	return u.Host, strings.TrimPrefix(u.Path, "/"), nil
}

// mergeDestDir merge dest dir
func mergeDestDir(prefix, dir string) string {
	if !strings.HasSuffix(dir, "/") {
		names := strings.Split(dir, "/")
		return strings.Join([]string{prefix, names[len(names)-1]}, "/")
	}
	return prefix
}

// printVolumes format volume data
func printVolumes(d *DryccCmd, volumes api.Volumes) {
	table := d.getDefaultFormatTable([]string{"NAME", "OWNER", "TYPE", "PTYPE", "PATH", "SIZE"})
	for _, volume := range volumes {
		if len(volume.Path) > 0 {
			for _, key := range *sortKeys(volume.Path) {
				table.Append([]string{volume.Name, volume.Owner, volume.Type, key, fmt.Sprintf("%v", volume.Path[key]), volume.Size})
			}
		} else {
			table.Append([]string{volume.Name, volume.Owner, volume.Type, "", "", volume.Size})
		}
	}
	table.Render()
}

func (d *DryccCmd) newProgressbar(maxBytes int64, icon, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		maxBytes,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() { fmt.Fprint(os.Stderr, "\n") }),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][%s][reset] %s", icon, d.fixateString(description, 32))),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}
