package utils

import (
	"os"
	"path/filepath"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/workflow-cli/pkg/git"
	"github.com/drycc/workflow-cli/pkg/settings"
	"github.com/stretchr/testify/assert"
)

func TestLoadProject(t *testing.T) {
	name, err := os.MkdirTemp("", "drycc-cli-unit-test-load")
	assert.NoError(t, err)
	defer os.RemoveAll(name)

	filename := filepath.Join(name, "test.json")
	host := "drycc.example.com"
	client, err := drycc.New(false, host, "")
	assert.NoError(t, err)

	config := settings.Settings{
		Username: "test",
		Client:   client,
	}

	filename, err = config.Save(filename)
	assert.NoError(t, err)

	appID, _, err := LoadAppSettings(filename, "test")
	assert.NoError(t, err)
	assert.Equal(t, appID, "test", "app")

	assert.NoError(t, os.Chdir(name))

	os.Setenv("DRYCC_APP", "testapp")
	appID, _, err = LoadAppSettings(filename, "")
	assert.NoError(t, err)
	assert.Equal(t, appID, "testapp", "app")
	os.Unsetenv("DRYCC_APP")

	appID, _, err = LoadAppSettings(filename, "")
	assert.NoError(t, err)
	assert.Equal(t, appID, filepath.Base(name), "app")

	assert.NoError(t, git.Init(git.DefaultCmd))
	assert.NoError(t, git.CreateRemote(git.DefaultCmd, host, "drycc", "testing"))

	appID, _, err = LoadAppSettings(filename, "")
	assert.NoError(t, err)
	assert.Equal(t, appID, "testing", "app")
}
