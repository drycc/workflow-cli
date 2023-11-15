package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUsersList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/users/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
				{
					"id": 2,
					"last_login": "2014-10-19T22:01:00.601Z",
					"is_superuser": false,
					"username": "test",
					"first_name": "test",
					"last_name": "testerson",
					"email": "test@example.com",
					"is_staff": false,
					"is_active": true,
					"date_joined": "2014-10-19T22:01:00.601Z",
					"groups": [],
					"user_permissions": []
				},
				{
					"id": 1,
					"last_login": "2014-10-19T22:01:00.601Z",
					"is_superuser": true,
					"username": "jkirk",
					"first_name": "james",
					"last_name": "kirk",
					"email": "jkrik@starfleet.ufp.gov",
					"is_staff": true,
					"is_active": true,
					"date_joined": "2014-10-19T22:01:00.601Z",
					"groups": [],
					"user_permissions": []
				}
			]
		}`)
	})

	err = cmdr.UsersList(-1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `USERNAME    EMAIL                      ADMIN    STAFF    ACTIVE    DATE-JOIN                
test        test@example.com           false    false    true      2014-10-19T22:01:00.601Z    
jkirk       jkrik@starfleet.ufp.gov    true     true     true      2014-10-19T22:01:00.601Z    
`, "output")
}

func TestUsersListLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/users/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
				{
					"id": 2,
					"last_login": "2014-10-19T22:01:00.601Z",
					"is_superuser": false,
					"username": "test",
					"first_name": "test",
					"last_name": "testerson",
					"email": "test@example.com",
					"is_staff": false,
					"is_active": true,
					"date_joined": "2014-10-19T22:01:00.601Z",
					"groups": [],
					"user_permissions": []
				}
			]
		}`)
	})

	err = cmdr.UsersList(1)
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `USERNAME    EMAIL               ADMIN    STAFF    ACTIVE    DATE-JOIN                
test        test@example.com    false    false    true      2014-10-19T22:01:00.601Z    
`, "output")
}

func TestUsersEnable(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/users/test/enable/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, ``)
	})

	err = cmdr.UsersEnable("test")
	assert.NoError(t, err)

	assert.Equal(t, b.String(), `Enabling user test... done
This modification is only temporary and will be reverted when the user login again.
`, "output")
}

func TestUsersDisable(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/users/test/disable/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, ``)
	})

	err = cmdr.UsersDisable("test")
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `Disabling user test... done
This modification is only temporary and will be reverted when the user login again.
`, "output")
}
