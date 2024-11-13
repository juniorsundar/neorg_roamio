package local

import (
	"os"
	"path/filepath"
)

var dataHome string

func GetLocalDir() {
    // Get the XDG_DATA_HOME value
    dataHome := os.Getenv("XDG_DATA_HOME")
    if dataHome == "" {
        dataHome = filepath.Join(os.Getenv("HOME"), ".local", "share", "roamio")
    }
}
