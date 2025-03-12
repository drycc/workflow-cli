package commands

import (
	"fmt"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/appsettings"
	"github.com/drycc/workflow-cli/internal/utils"
)

// LabelsList list app's labels
func (d *DryccCmd) LabelsList(appID string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if len(appSettings.Label) == 0 {
		d.Println(fmt.Sprintf("No labels found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"OWNER", "KEY", "VALUE"})
		for _, key := range *sortKeys(appSettings.Label) {
			table.Append([]string{
				appSettings.Owner,
				key,
				fmt.Sprintf("%v", appSettings.Label[key]),
			})
		}
		table.Render()
	}

	return nil
}

// LabelsSet sets labels for app
func (d *DryccCmd) LabelsSet(appID string, labels []string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	labelsMap, err := parseLabels(labels)
	if err != nil {
		return err
	}

	d.Printf("Applying labels on %s... ", appID)

	quit := progress(d.WOut)

	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Label: labelsMap})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

// LabelsUnset removes labels for the app.
func (d *DryccCmd) LabelsUnset(appID string, labels []string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	labelsMap := make(map[string]interface{})

	for _, label := range labels {
		labelsMap[label] = nil
	}

	d.Printf("Removing labels on %s... ", appID)

	quit := progress(d.WOut)

	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Label: labelsMap})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

func parseLabels(labels []string) (map[string]interface{}, error) {
	labelsMap := make(map[string]interface{})

	for _, label := range labels {
		key, value, err := parseLabel(label)

		if err != nil {
			return nil, err
		}

		labelsMap[key] = value
	}

	return labelsMap, nil
}

func parseLabel(label string) (string, string, error) {
	parts := strings.Split(label, "=")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf(`%s is invalid, Must be in format key=value
Examples: git_repo=https://github.com/drycc/workflow team=frontend`, label)
	}

	return parts[0], parts[1], nil
}
