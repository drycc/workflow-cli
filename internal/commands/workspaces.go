package commands

import (
	"fmt"

	"github.com/drycc/controller-sdk-go/workspaces"
	"github.com/drycc/controller-sdk-go/workspaces/invitations"
	"github.com/drycc/controller-sdk-go/workspaces/members"
	"github.com/drycc/workflow-cli/pkg/settings"
)

// WorkspacesList lists workspaces visible to the current user.
func (d *DryccCmd) WorkspacesList(results int) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	wkspaces, count, err := workspaces.List(s.Client, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count > 0 {
		table := d.getDefaultFormatTable([]string{"NAME", "EMAIL", "CREATED", "UPDATED"})
		for _, ws := range wkspaces {
			table.Append([]string{
				ws.Name,
				ws.Email,
				d.formatTime(ws.Created),
				d.formatTime(ws.Updated),
			})
		}
		table.Render()
	} else {
		d.Println("No workspaces found.")
	}
	return nil
}

// WorkspacesCreate creates a new workspace.
func (d *DryccCmd) WorkspacesCreate(name, email string) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	d.Print("Creating Workspace... ")
	quit := progress(d.WOut)
	ws, err := workspaces.Create(s.Client, name, email)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("done, created %s\n", ws.Name)
	return nil
}

// WorkspacesInfo shows detailed information about a workspace.
func (d *DryccCmd) WorkspacesInfo(name string, results int) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	ws, err := workspaces.Get(s.Client, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"Name:", ws.Name})
	table.Append([]string{"Email:", ws.Email})
	table.Append([]string{"Created:", d.formatTime(ws.Created)})
	table.Append([]string{"Updated:", d.formatTime(ws.Updated)})

	if results == defaultLimit {
		results = s.Limit
	}

	// print members
	mems, _, err := members.List(s.Client, name, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if len(mems) > 0 {
		table.Append([]string{"Members:"})
		for index, m := range mems {
			table.Append([]string{"", "User:", m.User})
			table.Append([]string{"", "Email:", m.Email})
			table.Append([]string{"", "Role:", m.Role})
			table.Append([]string{"", "Alerts:", fmt.Sprintf("%v", m.Alerts)})
			if len(mems) > index+1 {
				table.Append([]string{""})
			}
		}
	} else {
		table.Append([]string{"Members:", safeGetString("")})
	}

	table.Render()
	return nil
}

// WorkspacesDelete deletes a workspace.
func (d *DryccCmd) WorkspacesDelete(name, confirm string) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	if confirm == "" {
		d.Printf(` !    WARNING: Potentially Destructive Action
 !    This command will destroy the workspace: %s
 !    To proceed, type "%s" or re-run this command with --confirm=%s

> `, name, name, name)

		fmt.Scanln(&confirm)
	}

	if confirm != name {
		return fmt.Errorf("workspace %s does not match confirm %s, aborting", name, confirm)
	}

	d.Printf("Destroying %s... ", name)
	quit := progress(d.WOut)
	err = workspaces.Delete(s.Client, name)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// WorkspacesInvite invites a user to a workspace by email.
func (d *DryccCmd) WorkspacesInvite(workspace, email string) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	d.Printf("Inviting %s to %s... ", email, workspace)
	quit := progress(d.WOut)
	_, err = invitations.Create(s.Client, workspace, email)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// WorkspacesRemove removes a user from a workspace.
func (d *DryccCmd) WorkspacesRemove(workspace, username string) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	d.Printf("Removing %s from %s... ", username, workspace)
	quit := progress(d.WOut)
	err = members.Delete(s.Client, workspace, username)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// WorkspacesUpdate updates a workspace member's role and/or alerts setting.
func (d *DryccCmd) WorkspacesUpdate(workspace, username, role string, alerts *bool) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}

	d.Printf("Updating %s in %s... ", username, workspace)
	quit := progress(d.WOut)
	member, err := members.Update(s.Client, workspace, username, role, alerts)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("done, %s is now role=%s alerts=%v in %s\n", member.User, member.Role, member.Alerts, workspace)
	return nil
}
