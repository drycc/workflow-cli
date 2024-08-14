package cmd

import (
	"fmt"

	"github.com/drycc/controller-sdk-go/perms"
	"github.com/drycc/workflow-cli/settings"
)

// PermCodes prints all perm codenames.
func (d *DryccCmd) PermCodes(results int) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}
	codenames, _, err := perms.Codes(s.Client, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{"CODENAME", "DESCRIPTION"})
	for _, code := range codenames {
		table.Append([]string{code.Codename, code.Description})
	}
	table.Render()
	return nil
}

// PermList prints which users have permissions.
func (d *DryccCmd) PermList(codename string, results int) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}
	perms, _, err := perms.List(s.Client, codename, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{"ID", "CODENAME", "UNIQUEID", "USERNAME"})
	for _, perm := range perms {
		table.Append([]string{fmt.Sprintf("%d", perm.ID), perm.Codename, perm.Uniqueid, perm.Username})
	}
	table.Render()
	return nil
}

// PermCreate create user perm to user.
func (d *DryccCmd) PermCreate(codename, uniqueid, username string) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}
	d.Printf("Adding %s to %s:%s collaborators... ", username, codename, uniqueid)
	err = perms.Create(s.Client, codename, uniqueid, username)

	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")

	return nil
}

// PermDelete removes a user from an app or revokes admin privileges.
func (d *DryccCmd) PermDelete(userPermID uint64) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}
	d.Printf("Removing user perm with id %d... ", userPermID)
	err = perms.Delete(s.Client, fmt.Sprintf("%d", userPermID))

	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	d.Println("done")
	return nil
}
