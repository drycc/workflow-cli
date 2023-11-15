package cmd

import (
	"fmt"

	"github.com/drycc/controller-sdk-go/perms"
	"github.com/drycc/workflow-cli/pkg/git"
	"github.com/drycc/workflow-cli/settings"
)

// PermsList prints which users have permissions.
func (d *DryccCmd) PermsList(appID string, admin bool, results int) error {
	s, appID, err := permsLoad(d.ConfigFile, appID, admin)

	if err != nil {
		return err
	}

	var users []string

	if admin {
		if results == defaultLimit {
			results = s.Limit
		}
		users, _, err = perms.ListAdmins(s.Client, results)
	} else {
		users, err = perms.List(s.Client, appID)
	}

	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	table := d.getDefaultFormatTable([]string{"USERNAME", "ADMIN"})
	for _, user := range users {
		table.Append([]string{user, fmt.Sprintf("%v", admin)})
	}
	table.Render()
	return nil
}

// PermCreate adds a user to an app or makes them an administrator.
func (d *DryccCmd) PermCreate(appID string, username string, admin bool) error {

	s, appID, err := permsLoad(d.ConfigFile, appID, admin)

	if err != nil {
		return err
	}

	if admin {
		d.Printf("Adding %s to system administrators... ", username)
		err = perms.NewAdmin(s.Client, username)
	} else {
		d.Printf("Adding %s to %s collaborators... ", username, appID)
		err = perms.New(s.Client, appID, username)
	}

	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")

	return nil
}

// PermDelete removes a user from an app or revokes admin privileges.
func (d *DryccCmd) PermDelete(appID, username string, admin bool) error {

	s, appID, err := permsLoad(d.ConfigFile, appID, admin)

	if err != nil {
		return err
	}

	if admin {
		d.Printf("Removing %s from system administrators... ", username)
		err = perms.DeleteAdmin(s.Client, username)
	} else {
		d.Printf("Removing %s from %s collaborators... ", username, appID)
		err = perms.Delete(s.Client, appID, username)
	}

	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")

	return nil
}

func permsLoad(cf, appID string, admin bool) (*settings.Settings, string, error) {
	s, err := settings.Load(cf)

	if err != nil {
		return nil, "", err
	}

	if !admin && appID == "" {
		appID, err = git.DetectAppName(git.DefaultCmd, s.Client.ControllerURL.Host)

		if err != nil {
			return nil, "", err
		}
	}

	return s, appID, err
}
