package cmd

import (
	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/auth"
	"github.com/drycc/workflow-cli/settings"
)

// Login to a Drycc controller.
func (d *DryccCmd) Login(controller string, sslVerify bool, username, password string) error {
	c, err := drycc.New(sslVerify, controller, "")

	if err != nil {
		return err
	}

	// Set user agent for temporary client.
	c.UserAgent = settings.UserAgent

	if err = c.CheckConnection(); d.checkAPICompatibility(c, err) != nil {
		return err
	}

	token, err := d.TokensAdd(c, username, password, "workflow-cli", "yes", false)
	if err != nil {
		return err
	}
	// save settings
	s := settings.Settings{Client: c}
	s.Client.Token = token.Token
	s.Username = token.Username
	filename, err := s.Save(d.ConfigFile)
	if err != nil {
		return err
	}
	d.Printf("Logged in as %s\n", token.Username)
	d.Printf("Configuration file written to %s\n", filename)
	return nil
}

// Logout from a Drycc controller.
func (d *DryccCmd) Logout() error {
	if err := settings.Delete(d.ConfigFile); err != nil {
		return err
	}

	d.Println("Logged out")
	return nil
}

// Whoami prints the logged in user. If all is true, it fetches info from the controller to know
// more about the user.
func (d *DryccCmd) Whoami(all bool) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if all {
		user, err := auth.Whoami(s.Client)
		if err != nil {
			return err
		}
		d.Println(user)
	} else {
		d.Printf("You are %s at %s\n", s.Username, s.Client.ControllerURL.String())
	}
	return nil
}
