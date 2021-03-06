package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Default configuration path passed to Viper
var CONFIG_PATH = "$HOME/.config/"
var CONFIG_NAME = "pmwrap"

// A struct containing all global configuration
type config struct {
	DmenuCommand []string `mapstructure:"DmenuCommand"`
	PinEntryCommand []string `mapstructure:"PinEntryCommand"`
	NotifySendCommand string `mapstructure:"NotifySendCommand"`
}

// Global instances, populated by the load methods
var globalViper *viper.Viper
var globalConfig *config

// Creates and configures a new Viper instance
func loadViper(cmd *cobra.Command, configPath string, configName string) (*viper.Viper, error) {
	aViper := viper.New()
	aViper.SetConfigName(configName)
	aViper.SetConfigType("yaml")
	aViper.AddConfigPath(configPath)

	// Dmenu config
	cmd.PersistentFlags().StringArrayP("dmenu-command", "", []string{}, "dmenu command to use for querying user to select a value.")
	err := aViper.BindPFlag("DmenuCommand", cmd.PersistentFlags().Lookup("dmenu-command"))
	if err != nil {
		return nil, err
	}
	aViper.SetDefault("DmenuCommand", []string{"dmenu"})

	// Pinenry config
	cmd.PersistentFlags().StringArrayP("pinentry", "", []string{}, "pinentry command to use to query user for pin.")
	err = aViper.BindPFlag("PinEntryCommand", cmd.PersistentFlags().Lookup("pinentry"))
	if err != nil {
		return nil, err
	}
	aViper.SetDefault("PinEntryCommand", "pinentry")

	// notify-send config
	cmd.PersistentFlags().StringP("notify-send-command", "", "", "notify-send command to send user notifications.")
	err = aViper.BindPFlag("NotifySendCommand", cmd.PersistentFlags().Lookup("notify-send-command"))
	if err != nil {
		return nil, err
	}
	aViper.SetDefault("NotifySendCommand", "notify-send")
	

	// Return
	globalViper = aViper
	return aViper, nil
}

// Loads viper by populating the globalViper variable with a configured viper instance.
func LoadViper(cmd *cobra.Command) (error) {
	aViper, err := loadViper(cmd, CONFIG_PATH, CONFIG_NAME)
	globalViper = aViper
	return err
}

// Loads the configuration using a viper instance.
func loadConfig(aViper viper.Viper) (config, error) {
	var aConfig config
	err := aViper.ReadInConfig()
	if err != nil {
		_, isNotFound := err.(viper.ConfigFileNotFoundError)
		if ! isNotFound {
			return aConfig, err
		}
	}
	err = aViper.Unmarshal(&aConfig)
	return aConfig, err
}

// Loads the configuration by populating the globalConfig variable
func LoadConfig() error {
	aConfig, err := loadConfig(*globalViper)
	globalConfig = &aConfig
	return err
}

func GetConfig() config {
	return *globalConfig
}
