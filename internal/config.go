package internal

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const fileName = "config.yaml"

var dbContainerNames = map[string]string{
	"ce":       "db-ce",
	"commerce": "db-commerce",
	"dxp":      "db-dxp",
	"nightly":  "db-nightly",
	"release":  "db-release",
}
var portalContainerNames = map[string]string{
	"ce":       "lpn-ce",
	"commerce": "lpn-commerce",
	"dxp":      "lpn-dxp",
	"nightly":  "lpn-nightly",
	"release":  "lpn-release",
}

// LPNConfig tool configuration
type LPNConfig struct {
	Container ContainerConfig `mapstructure:"container"`
}

// ContainerConfig container configuration
type ContainerConfig struct {
	Names ContainerNameConfig `mapstructure:"names"`
}

// ContainerNameConfig container configuration
type ContainerNameConfig struct {
	Db     map[string]string `mapstructure:"db"`
	Portal map[string]string `mapstructure:"portal"`
}

func initConfigFile(workspace string, configFile string, defaults map[string]interface{}) *os.File {
	log.Printf("Creating %s with default values in %s.", configFile, workspace)

	configFilePath := filepath.Join(workspace, configFile)

	f, _ := os.Create(configFilePath)

	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(workspace)

	err := v.WriteConfig()
	if err != nil {
		log.Fatalf(`Cannot save default configuration file at %s: %v`, configFilePath, err)
	}

	return f
}

// NewConfig returns a new configuration
func NewConfig(workspace string) *LPNConfig {
	lpnConfig, err := readConfig(workspace, fileName, map[string]interface{}{
		"container": map[string]interface{}{
			"names": map[string]interface{}{
				"db":     dbContainerNames,
				"portal": portalContainerNames,
			},
		},
	})
	if err != nil {
		log.Fatalf("Error when reading config: %v\n", err)
	}

	return &lpnConfig
}

func readConfig(
	workspace string, configFile string, defaults map[string]interface{}) (LPNConfig, error) {

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(workspace)

	err := viper.ReadInConfig()
	if err != nil {
		initConfigFile(workspace, configFile, defaults)
		viper.ReadInConfig()
	}

	var lpnConfig LPNConfig
	err = viper.Unmarshal(&lpnConfig)
	if err != nil {
		log.Fatalf("Unable to decode configuration into struct, %v", err)
	}

	return lpnConfig, err
}
