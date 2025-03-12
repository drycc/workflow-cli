package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

const keyFixture = "fdbf3b34742e4ed2be4dfa848af13007"

func TestLogin(t *testing.T) {
	t.Skip("Skip long running tests")
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{}`)
	})

	server.Mux.HandleFunc("/v2/auth/login/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		}
		w.WriteHeader(http.StatusFound)
		if len(body) == 0 {
			testutil.SetHeaders(w)
			w.Header().Add("Location", fmt.Sprintf("/v2/login/drycc/?key=%s/", keyFixture))
			w.WriteHeader(http.StatusOK)
			w.Write(nil)
		} else {
			w.Write([]byte(fmt.Sprintf(`{"key": "%s"}`, keyFixture)))
		}
	})

	server.Mux.HandleFunc(fmt.Sprintf("/v2/auth/token/%s/", keyFixture), func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"username":"test-user","token":"eaf2d1d85f6b410b81d94bfec159019b"}`))
		w.Write(nil)
	})
	err = cmdr.Login(server.Server.URL, false, "", "")
	assert.NoError(t, err)
}

func TestLogout(t *testing.T) {
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	err = cmdr.Logout()
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "Logged out\n", "output")
}

func TestWhoami(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/auth/whoami/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "email": "test@example.com",
  "username": "test",
  "first_name": "",
  "last_name": "",
  "is_superuser": true,
  "is_staff": true,
  "groups": [],
  "user_permissions": [],
  "last_login": "2016-09-12T22:15:26Z",
  "date_joined": "2015-09-12T22:15:26Z",
  "is_active": true
}`)
	})

	err = cmdr.Whoami(false)
	assert.NoError(t, err)
	expected := fmt.Sprintf("You are test at %s\n", server.Server.URL)
	assert.Equal(t, b.String(), expected, "output")
}

func TestWhoamiAll(t *testing.T) {
	t.Parallel()

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/auth/whoami/", func(w http.ResponseWriter, _ *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
  "email": "test@example.com",
  "username": "test",
  "first_name": "test",
  "last_name": "test",
  "is_superuser": true,
  "is_staff": true,
  "groups": [],
  "user_permissions": [],
  "last_login": "2016-09-12T22:15:26Z",
  "date_joined": "2015-09-12T22:15:26Z",
  "is_active": true
}`)
	})

	err = cmdr.Whoami(true)
	assert.NoError(t, err)
	expected := `ID: 0
Username: test
Email: test@example.com
First Name: test
Last Name: test
Last Login: 2016-09-12T22:15:26Z
Is Superuser: true
Is Staff: true
Is Active: true
Date Joined: 2015-09-12T22:15:26Z
`
	assert.Equal(t, b.String(), expected, "output")
}
