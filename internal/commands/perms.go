package commands

import (
	"strings"

	"github.com/drycc/controller-sdk-go/perms"
	"github.com/drycc/workflow-cli/internal/utils"
)

// PermList prints which users have permissions.
func (d *DryccCmd) PermList(appID string, results int) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	perms, _, err := perms.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{"USERNAME", "PERMISSIONS"})
	for _, perm := range perms {
		p := strings.Join(perm.Permissions, ",")
		table.Append([]string{perm.Username, p})
	}
	table.Render()
	return nil
}

// PermCreate create user perms.
func (d *DryccCmd) PermCreate(appID, username, permissions string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	d.Printf("Adding user %s as a collaborator for %s... ", username, permissions)
	quit := progress(d.WOut)
	err = perms.Create(s.Client, appID, username, permissions)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")

	return nil
}

// PermUpdate update user perms.
func (d *DryccCmd) PermUpdate(appID, username, permissions string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	d.Printf("Updating user %s as a collaborator for %s... ", username, permissions)
	quit := progress(d.WOut)
	err = perms.Update(s.Client, appID, username, permissions)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")

	return nil
}

// PermDelete removes a user from an app.
func (d *DryccCmd) PermDelete(appID, username string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	d.Printf("Removing user permission... ")
	quit := progress(d.WOut)
	err = perms.Delete(s.Client, appID, username)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	d.Println("done")
	return nil
}
