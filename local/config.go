package local

import (
	"errors"
	"os"
	"path/filepath"

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

    if ConfigData.Workspace.Root == "" {
        return errors.New("Workspace root missing in configuration files.")
    }
    return nil
}
