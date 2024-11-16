package local

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/juniorsundar/neorg_roamio/logger"
	"gopkg.in/yaml.v3"
)

var configHome string
var ConfigPath string
var ConfigData Config

type Config struct {
	Workspace struct {
		Name string `yaml:"name"`
		Root string `yaml:"root"`
		Port string `yaml:"port"`
	} `yaml:"workspace"`
	Logging struct {
		Verbosity bool `yaml:"verbosity"`
		Color     bool `yaml:"color"`
	} `yaml:"logging"`
}

func GetConfig(configName string) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(os.Getenv("HOME"), ".config", "roamio")
	}

	err := os.MkdirAll(configHome, 0755)
	if err != nil {
		logger.LogErr.Println(err)
	}

	ConfigPath = filepath.Join(configHome, configName+".yaml")
	_, err = os.Stat(ConfigPath)
	if os.IsNotExist(err) {
		var tempConfig Config
		// Create a temp uuid for name of the workspace
		tempConfig.Workspace.Name = uuid.New().String()
		tempConfig.Workspace.Port = "8080"

		data, _ := yaml.Marshal(&tempConfig)
		err = os.WriteFile(ConfigPath, []byte(data), 0666)
	}
}

func ParseConfig() error {
	data, err := os.ReadFile(ConfigPath)

	err = yaml.Unmarshal(data, &ConfigData)
	if err != nil {
		logger.LogErr.Fatalf("Error parsing config file %s: %v", ConfigPath, err)
		return err
	}

	// Check if the config file has all the data that you need
	if ConfigData.Workspace.Root == "" {
		return errors.New(fmt.Sprintf("'%s' missing Workspace:Root.", ConfigPath))
	}
	if ConfigData.Workspace.Name == "" {
		return errors.New(fmt.Sprintf("'%s' missing Workspace:Name.", ConfigPath))
	}
	if ConfigData.Workspace.Port == "" {
		return errors.New(fmt.Sprintf("'%s' missing Workspace:Port.", ConfigPath))
	}
	return nil
}
