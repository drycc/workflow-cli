package cmd

import (
	"fmt"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/config"
)

// TagsList lists an app's tags.
func (d *DryccCmd) TagsList(appID, ptype string, version int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	config, err := config.List(s.Client, appID, version)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if len(config.Tags) == 0 {
		d.Println(fmt.Sprintf("No tags found in %s app.", appID))
	} else {
		if tags, ok := config.Tags[ptype]; ok {
			table := d.getDefaultFormatTable([]string{"KEY", "VALUE"})
			for _, key := range *sortKeys(tags) {
				table.Append([]string{
					key,
					fmt.Sprintf("%v", config.Tags[ptype][key]),
				})
			}
			table.Render()
		}
	}
	return nil
}

// TagsSet sets an app's tags.
func (d *DryccCmd) TagsSet(appID, ptype string, tags []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	tagsMap, err := parseTags(tags)
	if err != nil {
		return err
	}

	d.Print("Applying tags... ")

	quit := progress(d.WOut)
	configObj := api.Config{Tags: make(map[string]api.ConfigTags)}
	configObj.Tags[ptype] = tagsMap

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.TagsList(appID, ptype, -1)
}

// TagsUnset removes an app's tags.
func (d *DryccCmd) TagsUnset(appID, ptype string, tags []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Applying tags... ")

	quit := progress(d.WOut)

	configObj := api.Config{Tags: make(map[string]api.ConfigTags)}
	configTags := make(api.ConfigTags)
	for _, tag := range tags {
		configTags[tag] = nil
	}
	configObj.Tags[ptype] = configTags

	_, err = config.Set(s.Client, appID, configObj)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Print("done\n\n")

	return d.TagsList(appID, ptype, -1)
}

func parseTags(tags []string) (map[string]interface{}, error) {
	tagMap := make(map[string]interface{})

	for _, tag := range tags {
		key, value, err := parseTag(tag)

		if err != nil {
			return nil, err
		}

		tagMap[key] = value
	}

	return tagMap, nil
}

func parseTag(tag string) (string, string, error) {
	parts := strings.Split(tag, "=")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf(`%s is invalid, Must be in format key=value
Examples: rack=1 evironment=production`, tag)
	}

	return parts[0], parts[1], nil
}
