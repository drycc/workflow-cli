package cmd

import (
	"fmt"
	"os"
	"strings"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/certs"
	dtime "github.com/drycc/controller-sdk-go/pkg/time"
	"github.com/drycc/workflow-cli/settings"
)

const dateFormat = "2 Jan 2006"

// CertsList lists certs registered with the controller.
func (d *DryccCmd) CertsList(results int) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = s.Limit
	}

	certList, _, err := certs.List(s.Client, results)
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
func (d *DryccCmd) CertAdd(cert string, key string, name string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	d.Print("Adding SSL endpoint... ")
	quit := progress(d.WOut)
	err = d.doCertAdd(s.Client, cert, key, name)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

func (d *DryccCmd) doCertAdd(c *drycc.Client, cert string, key string, name string) error {
	certFile, err := os.ReadFile(cert)
	if err != nil {
		return err
	}

	keyFile, err := os.ReadFile(key)
	if err != nil {
		return err
	}

	_, err = certs.New(c, string(certFile), string(keyFile), name)
	return d.checkAPICompatibility(c, err)
}

// CertRemove deletes a cert from the controller.
func (d *DryccCmd) CertRemove(name string) error {
	s, err := settings.Load(d.ConfigFile)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("Removing %s... ", name)
	quit := progress(d.WOut)

	err = certs.Delete(s.Client, name)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// CertInfo gets info about certficiate
func (d *DryccCmd) CertInfo(name string) error {
	s, err := settings.Load(d.ConfigFile)
	if err != nil {
		return err
	}
	cert, err := certs.Get(s.Client, name)
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
func (d *DryccCmd) CertAttach(name, domain string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	d.Printf("Attaching certificate %s to domain %s... ", name, domain)
	quit := progress(d.WOut)

	err = certs.Attach(s.Client, name, domain)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) == nil {
		d.Println("done")
	}

	return err
}

// CertDetach detaches a certificate from a domain
func (d *DryccCmd) CertDetach(name, domain string) error {
	s, err := settings.Load(d.ConfigFile)

	if err != nil {
		return err
	}

	d.Printf("Detaching certificate %s from domain %s... ", name, domain)
	quit := progress(d.WOut)

	err = certs.Detach(s.Client, name, domain)
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
