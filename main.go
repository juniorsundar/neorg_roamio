package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/juniorsundar/neorg_roamio/local"
	"github.com/juniorsundar/neorg_roamio/logger"
)

func main() {

	// --------------------

	configFilePtr := flag.String("config", "config", "Get configuration file.")

	flag.Parse()

	// --------------------

	configFile := *configFilePtr
	// call initLoggers for use
	logger.InitLoggers(false, false)

	local.GetConfig(configFile)
	err := local.ParseConfig()

	if err != nil {
		logger.LogErr.Fatalln("Config file missing root directory.")
	}

	// call initLoggers for use
	logger.InitLoggers(local.ConfigData.Logging.Color, local.ConfigData.Logging.Verbosity)

	// validate roamio initialisation
	flagString := fmt.Sprintf("\n\tDirectory: %s\n\tPort: %s\n\tANSI Colors: %t\n\tVerbosity: %t\n",
		local.ConfigData.Workspace.Root,
		local.ConfigData.Workspace.Port,
		local.ConfigData.Logging.Color,
		local.ConfigData.Logging.Verbosity)
	logger.LogInfo.Println(flagString)

	// Check if workspace exists
	info, err := os.Stat(local.ConfigData.Workspace.Root)
	if os.IsNotExist(err) {
		logger.LogErr.Fatalf("Directory %s does not exist.", local.ConfigData.Workspace.Root)
	} else if os.IsPermission(err) {
		logger.LogErr.Fatalf("Directory %s does not have edit permission.", local.ConfigData.Workspace.Root)
	}
	if !info.IsDir() {
		logger.LogErr.Fatalf("%s is not a directory", local.ConfigData.Workspace.Root)
	}

	// initialise dir_watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.LogErr.Fatalln("Something went wrong in creating watcher.")
	}
	defer watcher.Close()

	err = addWatchDirRecursively(watcher, local.ConfigData.Workspace.Root)
	if err != nil {
		logger.LogErr.Fatalln(err)
	}
	logger.LogWarn.Printf("Watching: %s", local.ConfigData.Workspace.Root)
	go dirWatcher(watcher)
	relativeFileList := listFilesRecursively(watcher, ".norg")

	// Verify if Cache exists
	if cacheExists() {
		// invalidate if exists
		err := invalidateCache(relativeFileList)
		if err != nil {
			logger.LogErr.Fatalln(err)
		}
	} else {
		// build cache if not
		err := buildCache()
		if err != nil {
			logger.LogErr.Fatalln(err)
		}
	}

	<-make(chan struct{})
}
