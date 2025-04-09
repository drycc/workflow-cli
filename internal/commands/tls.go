package commands

import (
	"fmt"

	"github.com/drycc/controller-sdk-go/tls"
	"github.com/drycc/workflow-cli/internal/utils"
)

// TLSInfo prints info about the TLS settings for the given app.
func (d *DryccCmd) TLSInfo(appID string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	tls, err := tls.Info(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	table := d.getDefaultFormatTable([]string{})
	table.Append([]string{"UUID:", tls.UUID})
	table.Append([]string{"Owner:", tls.Owner})
	table.Append([]string{"CertsAuto:", fmt.Sprintf("%v", tls.CertsAutoEnabled != nil && *(tls.CertsAutoEnabled))})
	table.Append([]string{"HTTPSEnforced:", fmt.Sprintf("%v", tls.HTTPSEnforced != nil && *(tls.HTTPSEnforced))})
	// Issuer
	table.Append([]string{"Issuer:", ""})
	if tls.Issuer != nil {
		table.Append([]string{"", "Email:", safeGetString(tls.Issuer.Email)})
		table.Append([]string{"", "Server:", safeGetString(tls.Issuer.Server)})
	}
	// Events
	table.Append([]string{"Events:", ""})
	for _, event := range tls.Events {
		table.Append([]string{"", "Name:", event["name"]})
		table.Append([]string{"", "Kind:", event["kind"]})
		table.Append([]string{"", "Time:", event["time"]})
		table.Append([]string{"", "Type:", event["type"]})
		table.Append([]string{"", "Status:", event["status"]})
		table.Append([]string{"", "Message:", d.wrapString(event["message"])})
		table.Append([]string{""})
	}
	table.Render()
	return nil
}

// TLSForceEnable enables the router to enforce https-only requests to the application.
func (d *DryccCmd) TLSForceEnable(appID string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Enabling https-only requests for %s... ", appID)

	quit := progress(d.WOut)
	_, err = tls.EnableHTTPSEnforced(s.Client, appID)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// TLSForceDisable disables the router to enforce https-only requests to the application.
func (d *DryccCmd) TLSForceDisable(appID string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Disabling https-only requests for %s... ", appID)

	quit := progress(d.WOut)
	_, err = tls.DisableHTTPSEnforced(s.Client, appID)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// TLSAutoEnable enables certs-auto requests to the application.
func (d *DryccCmd) TLSAutoEnable(appID string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Enabling certs-auto requests for %s... ", appID)

	quit := progress(d.WOut)
	_, err = tls.EnableCertsAutoEnabled(s.Client, appID)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// TLSAutoDisable disables certs-auto requests to the application.
func (d *DryccCmd) TLSAutoDisable(appID string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Disabling certs-auto requests for %s... ", appID)

	quit := progress(d.WOut)
	_, err = tls.DisableCertsAutoEnabled(s.Client, appID)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// TLSAutoIssuer add issuer requests to the application.
func (d *DryccCmd) TLSAutoIssuer(appID string, email string, server string, keyID string, keySecret string) error {
	appID, s, err := utils.LoadAppSettings(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	d.Printf("Adding issuer requests for %s... ", appID)

	quit := progress(d.WOut)
	_, err = tls.AddCertsIssuer(s.Client, appID, email, server, keyID, keySecret)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}
