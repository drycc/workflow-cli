package cmd

import (
	"fmt"

	"github.com/drycc/controller-sdk-go/domains"
)

// DomainsList lists domains registered with an app.
func (d *DryccCmd) DomainsList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	domains, count, err := domains.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count > 0 {
		table := d.getDefaultFormatTable([]string{"APP", "OWNER", "PTYPE", "CREATED", "UPDATED", "DOMAIN"})
		for _, domain := range domains {
			table.Append([]string{
				domain.App,
				domain.Owner,
				domain.ProcfileType,
				domain.Created,
				domain.Updated,
				domain.Domain,
			})
		}
		table.Render()
	} else {
		d.Println(fmt.Sprintf("No domains found in %s app.", appID))
	}
	return nil
}

// DomainsAdd adds a domain to an app.
func (d *DryccCmd) DomainsAdd(appID, domain, procfileType string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Adding %s to %s... ", domain, appID)

	quit := progress(d.WOut)
	_, err = domains.New(s.Client, appID, domain, procfileType)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// DomainsRemove removes a domain registered with an app.
func (d *DryccCmd) DomainsRemove(appID, domain string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Removing %s from %s... ", domain, appID)

	quit := progress(d.WOut)
	err = domains.Delete(s.Client, appID, domain)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}
