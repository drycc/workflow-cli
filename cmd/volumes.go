package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/volumes"
	"github.com/schollz/progressbar/v3"
	"sigs.k8s.io/yaml"
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
	s, appID, err := load(d.ConfigFile, appID)

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
func (d *DryccCmd) VolumesCreate(appID, name, vType, size string, parameters map[string]interface{}) error {
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

// VolumesClient a client for manage volume file
func (d *DryccCmd) VolumesClient(appID, cmd string, args ...string) error {
	switch cmd {
	case "ls":
		return d.volumesClientLs(appID, args[0])
	case "cp":
		return d.volumesClientCp(appID, args[0], args[1])
	case "rm":
		return d.volumesClientRm(appID, args[0])
	default:
		return fmt.Errorf("unknown command %s", cmd)
	}
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

// volumesClientLs get all directory entries sorted by filename.
func (d *DryccCmd) volumesClientLs(appID, vol string) error {

	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	name, path, err := parseVol(vol)
	if err != nil {
		return err
	}
	dirs, _, err := volumes.ListDir(s.Client, appID, name, path, 3000)
	if err != nil {
		return err
	}

	table := d.getDefaultFormatTable([]string{})
	for _, dir := range dirs {
		if dir.Type == "dir" {
			dir.Name = fmt.Sprintf("%s/", dir.Name)
		}
		table.Append([]string{fmt.Sprintf("[%s]", d.formatTime(dir.Timestamp)), dir.Size, dir.Name})
	}
	table.Render()
	return nil
}

func (d *DryccCmd) volumesClientGetAll(client *drycc.Client, appID, volumeID, volumePath, localPath string) error {
	if _, err := os.Stat(localPath); err != nil && os.IsNotExist(err) {
		os.MkdirAll(localPath, os.ModePerm)
	}
	dirs, _, err := volumes.ListDir(client, appID, volumeID, volumePath, 3000)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		_, subpath := path.Split(dir.Path)
		filepath := path.Join(localPath, subpath)
		if dir.Type == "file" {
			res, err := volumes.GetFile(client, appID, volumeID, dir.Path)
			if err != nil {
				return err
			}
			w, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			bar := d.newProgressbar(res.ContentLength, "↓", filepath)
			defer w.Close()
			if _, err = io.Copy(io.MultiWriter(w, bar), res.Body); err != nil {
				return err
			}
		} else {
			os.MkdirAll(filepath, os.ModePerm)
			if err := d.volumesClientGetAll(client, appID, volumeID, dir.Path, filepath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DryccCmd) volumesClientPostAll(client *drycc.Client, appID, volumeID, volumePath, localPath string) error {
	if file, err := os.Stat(localPath); err != nil {
		return err
	} else if !file.IsDir() {
		file, err := os.Open(localPath)
		if err != nil {
			return err
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return err
		}
		reader := progressbar.NewReader(file, d.newProgressbar(stat.Size(), "↑", localPath))
		if _, err := volumes.PostFile(client, appID, volumeID, volumePath, file.Name(), stat.Size(), &reader); err != nil {
			return err
		}
		return nil
	}
	if entries, err := os.ReadDir(localPath); err == nil {
		for _, entry := range entries {
			var dstFilepath string
			if entry.IsDir() {
				dstFilepath = path.Join(volumePath, entry.Name())
			} else {
				dstFilepath = volumePath
			}
			if err := d.volumesClientPostAll(client, appID, volumeID, dstFilepath, path.Join(localPath, entry.Name())); err != nil {
				return err
			}
		}
	} else {
		return err
	}
	return nil
}

// volumesClientCp copy files between volume and local file
func (d *DryccCmd) volumesClientCp(appID, src, dst string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	if strings.HasPrefix(src, "vol://") {
		f, err := os.Stat(dst)
		if err != nil {
			return err
		}
		if !f.IsDir() {
			return fmt.Errorf("the local path must be an existing dir")
		}
		volumeID, volumePath, err := parseVol(src)
		if err != nil {
			return err
		}
		if dirs, _, err := volumes.ListDir(s.Client, appID, volumeID, volumePath, 3000); err == nil && (len(dirs) != 1 || dirs[0].Type != "file") {
			dst = mergeDestDir(dst, volumePath)
		}
		return d.volumesClientGetAll(s.Client, appID, volumeID, volumePath, dst)
	} else if strings.HasPrefix(dst, "vol://") {
		volumeID, volumePath, err := parseVol(dst)
		if err != nil {
			return err
		}
		if dirs, _, err := volumes.ListDir(s.Client, appID, volumeID, volumePath, 3000); err == nil {
			names := strings.Split(strings.Trim(src, "/"), "/")
			if len(dirs) == 1 && dirs[0].Type == "file" && strings.HasSuffix(strings.Trim(volumePath, "/"), names[len(names)-1]) {
				return fmt.Errorf("the volume path cannot be an existing file")
			}
		}
		if file, err := os.Stat(src); err == nil && file.IsDir() {
			volumePath = mergeDestDir(volumePath, src)
		}
		return d.volumesClientPostAll(s.Client, appID, volumeID, volumePath, src)
	}
	return nil
}

// volumesClientRm delete a file from volume
func (d *DryccCmd) volumesClientRm(appID, vol string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	host, path, err := parseVol(vol)
	if err != nil {
		return err
	}
	res, err := volumes.DeleteFile(s.Client, appID, host, path)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("incorrect http status code %d", res.StatusCode)
	}

	return nil
}

func parseVolume(volumeVars []string) (map[string]interface{}, error) {
	volumeMap := make(map[string]interface{})
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
	description = fmt.Sprintf("%-32s", description)
	if len(description) > 32 {
		description = fmt.Sprintf("...%s", description[len(description)-29:])
	}
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
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][%s][reset] %s", icon, description)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}
