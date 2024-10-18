package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/juniorsundar/neorg_roamio/logger"
)

func main() {

	// --------------------

	dirPtr := flag.String("dir", ".", "Roam directory address.")
	portPtr := flag.String("port", "8080", "Roam server port.")
	noColorPtr := flag.Bool("no-color", false, "Disable colored output during logging.")

	flag.Parse()

	// --------------------

	workspaceRoot = *dirPtr

	// call initLoggers for use
	logger.InitLoggers(*noColorPtr)

	// validate roamio initialisation
	flagString := fmt.Sprintf("\n\tDirectory: %s\n\tPort: %s\n\tANSI Colors: %t\n",
		*dirPtr, *portPtr, *noColorPtr)
	logger.LogInfo.Println(flagString)

	// Check if workspace exists
	info, err := os.Stat(workspaceRoot)
	if os.IsNotExist(err) {
		logger.LogErr.Fatalf("Directory %s does not exist.", *dirPtr)
	} else if os.IsPermission(err) {
		logger.LogErr.Fatalf("Directory %s does not have edit permission.", *dirPtr)
	}
	if !info.IsDir() {
		logger.LogErr.Fatalf("%s is not a directory", *dirPtr)
	}

	// initialise dir_watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.LogErr.Fatalln("Something went wrong in creating watcher.")
	}
	defer watcher.Close()

	err = addWatchDirRecursively(watcher, workspaceRoot)
	if err != nil {
		logger.LogErr.Fatalln(err)
	}
	logger.LogWarn.Printf("Watching: %s", workspaceRoot)
	go dirWatcher(watcher)
	relativeFileList := listFilesRecursively(watcher, ".norg")

	// Verify if Cache exists
	if cacheExists(workspaceRoot) {
		// invalidate if exists
		err := invalidateCache(workspaceRoot, relativeFileList)
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
