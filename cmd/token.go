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
	"github.com/drycc/controller-sdk-go/tokens"
	"github.com/drycc/workflow-cli/settings"
)

func (d *DryccCmd) TokensList(results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	tokens, _, err := tokens.List(s.Client, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{"UUID", "OWNER", "ALIAS", "KEY", "CREATE", "UPDATED"})
	for _, token := range tokens {
		table.Append([]string{
			token.UUID,
			token.Owner,
			safeGetString(token.Alias),
			token.Key,
			token.Created,
			token.Updated,
		})
	}
	table.Render()
	return nil
}

func (d *DryccCmd) TokensAdd(c *drycc.Client, username, password, alias, confirm string, render bool) (*api.AuthTokenResponse, error) {
	if c == nil {
		s, err := settings.Load(d.ConfigFile)

		if err != nil {
			return nil, err
		}
		c = s.Client
	}

	if confirm == "" {
		d.Printf(` !    WARNING: Make sure to copy your token now.
 !    You won't be able to see it again, please confirm whether to continue.
 !    To proceed, type "yes" !

> `)

		fmt.Scanln(&confirm)
	}

	if confirm != "yes" {
		return nil, fmt.Errorf("cancel the creation of %s's token, aborting", username)
	}

	key, err := auth.Login(c, username, password)
	if d.checkAPICompatibility(c, err) != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
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
			return nil, err
		}
		key = u.Query()["key"][0]
	}
	quit := progress(d.WOut)
	token, err := d.doToken(c, key, alias)
	quit <- true
	<-quit
	if render {
		table := d.getDefaultFormatTable([]string{"USERNAME", "TOKEN"})
		table.Append([]string{token.Username, token.Token})
		table.Render()
	}
	return token, err
}

func (d *DryccCmd) TokensRemove(id, confirm string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if confirm == "" {
		d.Printf(` !    WARNING: You cannot undo this action.
 !    Any using this token will no longer be able to access the Controller API.
 !    To proceed, type "yes" !

> `)

		fmt.Scanln(&confirm)
	}

	if confirm != "yes" {
		d.Println("skip")
		return nil
	}

	err = tokens.Delete(s.Client, id)
	if err != nil {
		return err
	}
	d.Println("done")
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

func (d *DryccCmd) doToken(c *drycc.Client, key, alias string) (*api.AuthTokenResponse, error) {
	var token api.AuthTokenResponse
	for i := 0; i <= 120; i++ {
		token, _ = auth.Token(c, key, alias)
		time.Sleep(time.Duration(5) * time.Second)
		if token.Token != "" && token.Username != "" {
			break
		}
	}
	if token.Token == "" || token.Token == "fail" {
		return nil, errors.New("logged fail")
	}
	return &token, nil
}
