package cmd

import (
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

	users, count, err := users.List(s.Client, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== Users (*=admin)%s", limitCount(len(users), count))

	for _, user := range users {
		if user.IsSuperuser {
			d.Print("*")
		}
		d.Println(user.Username)
	}
	return nil
}
