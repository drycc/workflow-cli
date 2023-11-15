package cmd

import (
	"fmt"

	"github.com/drycc/controller-sdk-go/users"
	"github.com/drycc/workflow-cli/settings"
)

// UsersList lists users registered with the controller.
func (d *DryccCmd) UsersList(results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	users, _, err := users.List(s.Client, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{"USERNAME", "EMAIL", "ADMIN", "STAFF", "ACTIVE", "DATE-JOIN"})
	for _, user := range users {
		table.Append([]string{
			user.Username,
			user.Email,
			fmt.Sprintf("%v", user.IsSuperuser),
			fmt.Sprintf("%v", user.IsStaff),
			fmt.Sprintf("%v", user.IsActive),
			user.DateJoined,
		})
	}
	table.Render()
	return nil
}

// UsersEnable enable user with the controller.
func (d *DryccCmd) UsersEnable(username string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}
	d.Printf("Enabling user %s... ", username)
	err = users.Enable(s.Client, username)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	d.Println("This modification is only temporary and will be reverted when the user login again.")
	return nil
}

// UsersDisable disable user with the controller.
func (d *DryccCmd) UsersDisable(username string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	d.Printf("Disabling user %s... ", username)
	err = users.Disable(s.Client, username)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	d.Println("This modification is only temporary and will be reverted when the user login again.")
	return nil
}
