package cmd

import "github.com/drycc/controller-sdk-go/tls"

// TLSInfo prints info about the TLS settings for the given app.
func (d *DryccCmd) TLSInfo(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	tls, err := tls.Info(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== %s TLS\n", appID)
	d.Println(tls)

	return nil
}

// TLSForceEnable enables the router to enforce https-only requests to the application.
func (d *DryccCmd) TLSForceEnable(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

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
	s, appID, err := load(d.ConfigFile, appID)

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
	s, appID, err := load(d.ConfigFile, appID)

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
	s, appID, err := load(d.ConfigFile, appID)

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
