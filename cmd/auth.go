package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/auth"
	"github.com/drycc/workflow-cli/settings"
)

func (d *DryccCmd) doLogin(s settings.Settings, username, password string) error {
	key, err := auth.Login(s.Client, username, password)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if err != nil {
		return nil
	}
	if username == "" || password == "" {
		fmt.Printf("Opening browser to %s\n", key)
		d.Print("Waiting for login... ")
		err = d.openBrower(key)
		if err != nil {
			d.Print("Cannot open browser, please visit the website in yourself")
		}
		u, err := url.Parse(key)
		if err != nil {
			return err
		}
		key = u.Query()["key"][0]
	}
	quit := progress(d.WOut)
	d.doToken(s, key)
	quit <- true
	<-quit
	return nil
}

func (d *DryccCmd) openBrower(URL string) error {
	var commands = map[string]string{
		"windows": "start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}
	run, ok := commands[runtime.GOOS]
	if !ok {
		return errors.New("warning: Cannot open browser")
	}
	cmd := exec.Command(run, URL)
	err := cmd.Start()
	if err != nil {
		return errors.New("warning: Cannot open browser")
	}

	return nil
}

func (d *DryccCmd) doToken(s settings.Settings, key string) error {
	var token api.AuthTokenResponse
	for i := 0; i <= 120; i++ {
		token, _ = auth.Token(s.Client, key)
		time.Sleep(time.Duration(5) * time.Second)
		if token.Token != "" && token.Username != "" {
			break
		}
	}
	if token.Token == "" || token.Token == "fail" {
		d.Printf("Logged fail")
	} else {
		s.Client.Token = token.Token
		s.Username = token.Username
		filename, err := s.Save(d.ConfigFile)
		if err != nil {
			return nil
		}
		d.Printf("Logged in as %s\n", token.Username)
		d.Printf("Configuration file written to %s\n", filename)
	}
	return nil
}

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

	s := settings.Settings{Client: c}
	return d.doLogin(s, username, password)
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
