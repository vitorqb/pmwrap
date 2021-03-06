package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/vitorqb/pmwrap/package/testUtils"
)


func TestLoadConfigReturnsProperConfig(t *testing.T) {
	// ARRANGE
	cmd := cobra.Command{}
	configPath, err := testUtils.GetTestDataDirectory()
	assert.Nil(t, err)
	aViper, err := loadViper(&cmd, configPath, "config")
	assert.Nil(t, err)

	// ACT
	aConfig, err := loadConfig(*aViper)

	// ASSERT
	assert.Nil(t, err)
	expectedConfig := config{
		DmenuCommand: []string{"dmenu", "--command"},
		PinEntryCommand: []string{"pinentry-qt"},
		NotifySendCommand: "custom-notify-send",
	}
	assert.Equal(t, aConfig, expectedConfig)
}

func TestLoadConfigReturnsDefaultValues(t *testing.T) {
	cmd := cobra.Command{}
	aViper, err := loadViper(&cmd, "/not/exists", "not-exists")
	assert.Nil(t, err)
	aConfig, err := loadConfig(*aViper)
	assert.Nil(t, err)
	expectedConfig := config{
		DmenuCommand: []string{"dmenu"},
		PinEntryCommand: []string{"pinentry"},
		NotifySendCommand: "notify-send",
	}
	assert.Equal(t, aConfig, expectedConfig)	
}
