package cmd

import (
	"strings"

	"github.com/drycc/controller-sdk-go/allowlist"
)

// Allowlistlist lists the addresses allowlisted for app
func (d *DryccCmd) AllowlistList(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	allowlist, err := allowlist.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== %s Allowlisted Addresses\n", appID)

	for _, ip := range allowlist.Addresses {
		d.Println(ip)
	}
	return nil
}

// AllowlistAdd adds the addresses to the app's allowlist.
func (d *DryccCmd) AllowlistAdd(appID, IPs string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Adding %s to %s allowlist...\n", IPs, appID)

	quit := progress(d.WOut)
	_, err = allowlist.Add(s.Client, appID, strings.Split(IPs, ","))
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// AllowlistRemove deletes the addresses from the app's Allowlist.
func (d *DryccCmd) AllowlistRemove(appID, IPs string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Removing %s from %s allowlist...\n", IPs, appID)

	quit := progress(d.WOut)
	err = allowlist.Delete(s.Client, appID, strings.Split(IPs, ","))
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}
