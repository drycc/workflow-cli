package commands

import (
	"bytes"
	"os"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestDrinkOfChoice(t *testing.T) {
	os.Setenv("DRYCC_DRINK_OF_CHOICE", "test")
	assert.Equal(t, drinkOfChoice(), "test", "output")
	os.Unsetenv("DRYCC_DRINK_OF_CHOICE")
	assert.Equal(t, drinkOfChoice(), "coffee", "output")
}

func TestLimitsCount(t *testing.T) {
	t.Parallel()
	assert.Equal(t, limitCount(1, 1), "\n", "output")
	assert.Equal(t, limitCount(1, 2), " (1 of 2)\n", "output")
}

func TestAPICompatibility(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	cmdr := DryccCmd{WErr: &b, ConfigFile: ""}
	client := drycc.Client{ControllerAPIVersion: "v1.0"}

	err := cmdr.checkAPICompatibility(&client, drycc.ErrAPIMismatch)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), `!    WARNING: Client and server API versions do not match. Please consider upgrading.
!    Client version: 2.3
!    Server version: v1.0
`, "output")

	// After being warned once, the warning should not be printed again.
	b.Reset()
	err = cmdr.checkAPICompatibility(&client, drycc.ErrAPIMismatch)
	assert.NoError(t, err)
	assert.Equal(t, b.String(), "", "output")

	b.Reset()
	err = cmdr.checkAPICompatibility(&client, drycc.ErrConflict)
	assert.Error(t, drycc.ErrConflict, err)
	assert.Equal(t, b.String(), "", "output")
}
