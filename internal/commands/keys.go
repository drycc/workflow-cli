package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/keys"
	"github.com/drycc/workflow-cli/pkg/settings"
	"github.com/drycc/workflow-cli/pkg/ssh"
)

// KeysList lists a user's keys.
func (d *DryccCmd) KeysList(results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	keys, _, err := keys.List(s.Client, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if len(keys) > 0 {
		table := d.getDefaultFormatTable([]string{"ID", "OWNER", "KEY"})
		for _, key := range keys {
			table.Append([]string{
				key.ID,
				key.Owner,
				fmt.Sprintf("%s...%s", key.Public[:16], key.Public[len(key.Public)-10:]),
			})
		}
		table.Render()
	} else {
		d.Println("No any key found.")
	}
	return nil
}

// KeyRemove removes keys.
func (d *DryccCmd) KeyRemove(keyID string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	d.Printf("Removing %s SSH Key...", keyID)

	if err = keys.Delete(s.Client, keyID); d.checkAPICompatibility(s.Client, err) != nil {
		d.Println()
		return err
	}

	d.Println(" done")
	return nil
}

// KeyAdd adds keys.
func (d *DryccCmd) KeyAdd(name string, keyLocation string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	var key api.KeyCreateRequest

	// check if name is the key
	if name != "" && keyLocation == "" {
		// detect of name is a file
		_, err := os.Stat(name)
		if err == nil {
			keyLocation = name
			name = ""
		}
	}

	if keyLocation == "" {
		ks, err := listKeys(d.WOut)
		if err != nil {
			return err
		}
		key, err = chooseKey(ks, d.WIn, d.WOut)
		if err != nil {
			return err
		}
	} else {
		key, err = getKey(keyLocation)
		if err != nil {
			return err
		}
	}

	// if name is provided by user then overwrite that in the key object
	if name != "" {
		key.ID = name
	}

	d.Printf("Uploading %s to drycc...", filepath.Base(key.Name))

	if _, err = keys.New(s.Client, key.ID, key.Public); d.checkAPICompatibility(s.Client, err) != nil {
		d.Println()
		return err
	}

	d.Println(" done")
	return nil
}

func chooseKey(keys []api.KeyCreateRequest, input io.Reader,
	wOut io.Writer) (api.KeyCreateRequest, error) {
	fmt.Fprintln(wOut, "Found the following SSH public keys:")

	for i, key := range keys {
		fmt.Fprintf(wOut, "%d) %s %s\n", i+1, filepath.Base(key.Name), key.ID)
	}

	fmt.Fprintln(wOut, "0) Enter path to pubfile (or use keys:add <key_path>)")

	var selected string

	fmt.Fprint(wOut, "Which would you like to use with Drycc? ")
	fmt.Fscanln(input, &selected)

	numSelected, err := strconv.Atoi(selected)

	if err != nil {
		return api.KeyCreateRequest{}, fmt.Errorf("%s is not a valid integer", selected)
	}

	if numSelected < 0 || numSelected > len(keys) {
		return api.KeyCreateRequest{}, fmt.Errorf("%d is not a valid option", numSelected)
	}

	if numSelected == 0 {
		var filename string

		fmt.Fprint(wOut, "Enter the path to the pubkey file: ")
		fmt.Fscanln(input, &filename)

		return getKey(filename)
	}

	return keys[numSelected-1], nil
}

func listKeys(wOut io.Writer) ([]api.KeyCreateRequest, error) {
	folder := filepath.Join(settings.FindHome(), ".ssh")
	files, err := os.ReadDir(folder)

	if err != nil {
		return nil, err
	}

	var keys []api.KeyCreateRequest

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".pub" {
			key, err := getKey(filepath.Join(folder, file.Name()))

			if err == nil {
				keys = append(keys, key)
			} else {
				fmt.Fprintln(wOut, err)
			}
		}
	}

	return keys, nil
}

func getKey(filename string) (api.KeyCreateRequest, error) {
	keyContents, err := os.ReadFile(filename)

	if err != nil {
		return api.KeyCreateRequest{}, err
	}

	backupID := strings.Split(filepath.Base(filename), ".")[0]
	keyInfo, err := ssh.ParsePubKey(backupID, keyContents)
	if err != nil {
		return api.KeyCreateRequest{}, fmt.Errorf("%s is not a valid ssh key", filename)
	}
	return api.KeyCreateRequest{ID: keyInfo.ID, Public: keyInfo.Public, Name: filename}, nil
}
