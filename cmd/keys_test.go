package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/workflow-cli/pkg/testutil"
	"github.com/drycc/workflow-cli/settings"
	"github.com/stretchr/testify/assert"
)

func TestGetKey(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp("", "drycc-key")
	assert.NoError(t, err)

	toWrite := []byte("ssh-rsa abc test@example.com")

	expected := api.KeyCreateRequest{
		ID:     "test@example.com",
		Public: string(toWrite),
		Name:   file.Name(),
	}

	_, err = file.Write(toWrite)
	assert.NoError(t, err)
	file.Close()

	key, err := getKey(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, key, expected, "key")

	_, err = getKey("notarealkey")
	assert.NotEqual(t, err, nil, "file error")
}

func TestGetKeyNoComment(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp("", "drycc-key")
	assert.NoError(t, err)

	toWrite := []byte("ssh-rsa abc")

	expected := api.KeyCreateRequest{
		ID:     filepath.Base(file.Name()),
		Public: string(toWrite),
		Name:   file.Name(),
	}

	_, err = file.Write(toWrite)
	assert.NoError(t, err)

	key, err := getKey(file.Name())
	assert.NoError(t, err)

	assert.Equal(t, key, expected, "key")
}

func TestGetInvalidKey(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp("", "drycc-key")
	assert.NoError(t, err)

	toWrite := []byte("not a key")
	_, err = file.Write(toWrite)
	assert.NoError(t, err)

	expected := fmt.Sprintf("%s is not a valid ssh key", file.Name())

	_, err = getKey(file.Name())
	assert.Equal(t, err.Error(), expected, "error")
}

func TestListKeys(t *testing.T) {
	name, err := os.MkdirTemp("", "drycc-key")
	assert.NoError(t, err)
	settings.SetHome(name)

	folder := filepath.Join(name, ".ssh")

	err = os.Mkdir(folder, 0755)
	assert.NoError(t, err)

	toWrite := []byte("ssh-rsa abc test@example.com")
	fileNames := []string{"test1.pub", "test2.pub"}

	expected := []api.KeyCreateRequest{
		{
			ID:     "test@example.com",
			Public: string(toWrite),
			Name:   filepath.Join(folder, fileNames[0]),
		},
		{
			ID:     "test@example.com",
			Public: string(toWrite),
			Name:   filepath.Join(folder, fileNames[1]),
		},
	}

	for _, file := range fileNames {
		os.WriteFile(filepath.Join(folder, file), toWrite, 0775)
		assert.NoError(t, err)
	}

	keys, err := listKeys(io.Discard)
	assert.NoError(t, err)

	assert.Equal(t, keys, expected, "key")

	var b bytes.Buffer
	// Write bad ssh key
	filename := filepath.Join(folder, "test3.pub")
	os.WriteFile(filename, []byte("ssh-rsa"), 0775)
	_, err = listKeys(&b)
	assert.Equal(t, b.String(), filename+" is not a valid ssh key\n", "output")
	assert.NoError(t, err)

}

type chooseKeyCases struct {
	Reader      io.Reader
	Err         bool
	ExpectedErr string
	ExpectedKey *api.KeyCreateRequest
	LoadKey     bool
}

func TestChooseKey(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp("", "drycc-key")
	assert.NoError(t, err)
	toWrite := []byte("ssh-rsa abc test@example.com")
	_, err = file.Write(toWrite)
	assert.NoError(t, err)
	file.Close()

	testKeys := []api.KeyCreateRequest{
		{
			ID:     "test@example.com",
			Public: "ssh-rsa 123 abc@example.com",
			Name:   ".ssh/public/id_rsa.pub",
		},
		{
			ID:     "example@example.com",
			Public: "ssh-rsa abc123 example@example.com",
			Name:   ".ssh/public/id_rsa.pub",
		},
	}

	expectedWrittenKey := api.KeyCreateRequest{
		ID:     "test@example.com",
		Public: string(toWrite),
		Name:   file.Name(),
	}

	checks := []chooseKeyCases{
		{strings.NewReader("-1"), true, "-1 is not a valid option", nil, false},
		{strings.NewReader("3"), true, "3 is not a valid option", nil, false},
		{strings.NewReader("a"), true, "a is not a valid integer", nil, false},
		{strings.NewReader("1"), false, "", &testKeys[0], false},
		{strings.NewReader("0\n" + file.Name()), false, "", &expectedWrittenKey, true},
	}

	var b bytes.Buffer
	for _, check := range checks {
		b.Reset()
		key, err := chooseKey(testKeys, check.Reader, &b)
		expectedOut := `Found the following SSH public keys:
1) id_rsa.pub test@example.com
2) id_rsa.pub example@example.com
0) Enter path to pubfile (or use keys:add <key_path>)
Which would you like to use with Drycc? `

		if check.LoadKey {
			expectedOut += "Enter the path to the pubkey file: "
		}
		assert.Equal(t, b.String(), expectedOut, "output")

		if check.Err {
			assert.Equal(t, err.Error(), check.ExpectedErr, "error")
		} else {
			assert.Equal(t, key, *check.ExpectedKey, "key")
		}
	}
}

func TestKeysList(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/keys/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
				{
					"created": "2014-01-01T00:00:00UTC",
					"id": "cpike@starfleet.ufp",
					"owner": "cpike",
					"public": "ssh-rsa abc cpike@starfleet.ufp",
					"updated": "2014-01-01T00:00:00UTC",
					"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
				},
				{
					"created": "2014-01-01T00:00:00UTC",
					"id": "cpike@1701.ncc.starfleet.ufp",
					"owner": "cpike",
					"public": "ssh-rsa 123 cpike@1701.ncc.starfleet.ufp",
					"updated": "2014-01-01T00:00:00UTC",
					"uuid": "le19f5b5-4a72-4f94-a10c-d2a374jcd075"
				}
			]
		}`)
	})

	err = cmdr.KeysList(-1)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `ID                              OWNER    KEY                           
cpike@starfleet.ufp             cpike    ssh-rsa abc cpik...rfleet.ufp    
cpike@1701.ncc.starfleet.ufp    cpike    ssh-rsa 123 cpik...rfleet.ufp    
`, "output")
}

func TestKeysListLimit(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/keys/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		fmt.Fprintf(w, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
				{
					"created": "2014-01-01T00:00:00UTC",
					"id": "cpike@starfleet.ufp",
					"owner": "cpike",
					"public": "ssh-rsa abc cpike@starfleet.ufp",
					"updated": "2014-01-01T00:00:00UTC",
					"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
				}
			]
		}`)
	})

	err = cmdr.KeysList(1)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `ID                     OWNER    KEY                           
cpike@starfleet.ufp    cpike    ssh-rsa abc cpik...rfleet.ufp    
`, "output")
}

func TestKeyRemove(t *testing.T) {
	t.Parallel()
	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	server.Mux.HandleFunc("/v2/keys/cpike@starfleet.ufp", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})

	err = cmdr.KeyRemove("cpike@starfleet.ufp")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), "Removing cpike@starfleet.ufp SSH Key... done\n", "output")
}

func TestKeyAdd(t *testing.T) {
	// Set temp home dir so no unknown files are listed.
	name, err := os.MkdirTemp("", "drycc-key")
	assert.NoError(t, err)
	settings.SetHome(name)
	folder := filepath.Join(name, ".ssh")
	err = os.Mkdir(folder, 0755)
	assert.NoError(t, err)

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	keyFile, err := os.CreateTemp("", "drycc-cli-unit-test-ssh-key")
	assert.NoError(t, err)
	toWrite := []byte("ssh-rsa abc test@example.com")
	_, err = keyFile.Write(toWrite)
	assert.NoError(t, err)
	keyFile.Close()

	server.Mux.HandleFunc("/v2/keys/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.KeyCreateRequest{ID: "test@example.com", Public: string(toWrite)}, r)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})

	out := fmt.Sprintf("Uploading %s to drycc... done\n", filepath.Base(keyFile.Name()))

	err = cmdr.KeyAdd("", keyFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), out, "output")

	b.Reset()
	cmdr.WIn = strings.NewReader("0\n" + keyFile.Name())
	err = cmdr.KeyAdd("", "")
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), `Found the following SSH public keys:
0) Enter path to pubfile (or use keys:add <key_path>)
Which would you like to use with Drycc? Enter the path to the pubkey file: `+out, "output")
}

func TestKeyAddName(t *testing.T) {
	// Set temp home dir so no unknown files are listed.
	name, err := os.MkdirTemp("", "drycc-key")
	assert.NoError(t, err)
	settings.SetHome(name)
	folder := filepath.Join(name, ".ssh")
	err = os.Mkdir(folder, 0755)
	assert.NoError(t, err)

	cf, server, err := testutil.NewTestServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	var b bytes.Buffer
	cmdr := DryccCmd{WOut: &b, ConfigFile: cf}

	keyFile, err := os.CreateTemp("", "drycc-cli-unit-test-ssh-key")
	assert.NoError(t, err)
	// generate with one name but used another in the add
	toWrite := []byte("ssh-rsa abc test@example.com")
	_, err = keyFile.Write(toWrite)
	assert.NoError(t, err)
	keyFile.Close()

	server.Mux.HandleFunc("/v2/keys/", func(w http.ResponseWriter, r *http.Request) {
		testutil.SetHeaders(w)
		testutil.AssertBody(t, api.KeyCreateRequest{ID: "drycc-test-key", Public: string(toWrite)}, r)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "{}")
	})

	out := fmt.Sprintf("Uploading %s to drycc... done\n", filepath.Base(keyFile.Name()))

	err = cmdr.KeyAdd("drycc-test-key", keyFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, testutil.StripProgress(b.String()), out, "output")
}
