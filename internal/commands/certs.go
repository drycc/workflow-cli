package commands

import (
	"fmt"
	"os"
	"strings"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/certs"
	dtime "github.com/drycc/controller-sdk-go/pkg/time"
	"github.com/drycc/workflow-cli/internal/utils"
)

const dateFormat = "2 Jan 2006"

// CertsList lists certs registered with the controller.
func (d *DryccCmd) CertsList(appID string, results int) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	certList, _, err := certs.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	if len(certList) == 0 {
		d.Println("No certs")
	} else {
		table := d.getDefaultFormatTable([]string{"NAME", "COMMON-NAME", "EXPIRES", "SAN", "DOMAINS"})
		for _, cert := range certList {
			expires := safeGetString("")
			if cert.Expires.Time != nil {
				expires = cert.Expires.Format(dateFormat)
			}
			san := safeGetString(strings.Join(cert.SubjectAltName[:], ","))
			if len(san) > 32 {
				san = fmt.Sprintf("%s[...]", san[:32])
			}
			domains := safeGetString(strings.Join(cert.Domains[:], ","))
			if len(domains) > 32 {
				domains = fmt.Sprintf("%s[...]", domains[:32])
			}
			table.Append([]string{cert.Name, cert.CommonName, expires, san, domains})
		}
		table.Render()
	}
	return nil
}

// CertAdd adds a cert to the controller.
func (d *DryccCmd) CertAdd(appID string, cert string, key string, name string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Print("Adding SSL endpoint... ")
	quit := progress(d.WOut)
	err = d.doCertAdd(s.Client, appID, cert, key, name)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

func (d *DryccCmd) doCertAdd(c *drycc.Client, appID string, cert string, key string, name string) error {
	certFile, err := os.ReadFile(cert)
	if err != nil {
		return err
	}

	keyFile, err := os.ReadFile(key)
	if err != nil {
		return err
	}

	_, err = certs.New(c, appID, string(certFile), string(keyFile), name)
	return d.checkAPICompatibility(c, err)
}

// CertRemove deletes a cert from the controller.
func (d *DryccCmd) CertRemove(appID string, name string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("Removing %s... ", name)
	quit := progress(d.WOut)

	err = certs.Delete(s.Client, appID, name)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// CertInfo gets info about certficiate
func (d *DryccCmd) CertInfo(appID string, name string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)
	if err != nil {
		return err
	}
	cert, err := certs.Get(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"Name:", safeGetString(cert.Name)})
	table.Append([]string{"Common Name(s):", safeGetString(cert.CommonName)})
	table.Append([]string{"Expires At:", safeGetTime(cert.Expires, dateFormat)})
	table.Append([]string{"Starts At:", safeGetTime(cert.Starts, dateFormat)})
	table.Append([]string{"Fingerprint:", safeGetString(cert.Fingerprint)})
	table.Append([]string{"Subject Alt Name:", safeGetString(strings.Join(cert.SubjectAltName[:], ","))})
	table.Append([]string{"Issuer:", safeGetString(cert.Issuer)})
	table.Append([]string{"Subject:", safeGetString(cert.Subject)})
	table.Append([]string{""})
	table.Append([]string{"Connected Domains:", safeGetString(strings.Join(cert.Domains[:], ","))})
	table.Append([]string{"Owner:", safeGetString(cert.Owner)})
	table.Append([]string{"Created:", d.formatTime(cert.Created)})
	table.Append([]string{"Updated:", d.formatTime(cert.Updated)})
	table.Render()
	return nil
}

// CertAttach attaches a certificate to a domain
func (d *DryccCmd) CertAttach(appID string, name, domain string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Attaching certificate %s to domain %s... ", name, domain)
	quit := progress(d.WOut)

	err = certs.Attach(s.Client, appID, name, domain)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) == nil {
		d.Println("done")
	}

	return err
}

// CertDetach detaches a certificate from a domain
func (d *DryccCmd) CertDetach(appID, name, domain string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Detaching certificate %s from domain %s... ", name, domain)
	quit := progress(d.WOut)

	err = certs.Detach(s.Client, appID, name, domain)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

func safeGetTime(t dtime.Time, format string) string {
	out := ""
	if t.Time != nil {
		out = t.Format(format)
	}

	return safeGetString(out)
}
