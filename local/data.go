package local

import (
	"os"
	"path"
	"path/filepath"

	"github.com/juniorsundar/neorg_roamio/logger"
)

var DataHome string
var DatabasePath string

func GetLocalDir() {
    // Get the XDG_DATA_HOME value
    DataHome = os.Getenv("XDG_DATA_HOME")
    if DataHome == "" {
        DataHome = filepath.Join(os.Getenv("HOME"), ".local", "share", "roamio")
    }

	err := os.MkdirAll(DataHome, 0755)
	if err != nil {
		logger.LogErr.Println(err)
	}
}

func GetDatabase() {
    if DataHome == "" {
        GetLocalDir()
    }

    if ConfigData.Workspace.Name == "" {
        ParseConfig()
    }

    DatabasePath = path.Join(DataHome, ConfigData.Workspace.Name+".db3")
}
