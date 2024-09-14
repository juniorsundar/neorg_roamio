package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/juniorsundar/neorg_roamio/dir_watcher"
	"github.com/juniorsundar/neorg_roamio/logger"
)

func main() {

	dirPtr := flag.String("dir", ".", "Roam directory address.")
	portPtr := flag.String("port", "8080", "Roam server port.")
	verbosePtr := flag.Bool("verbose", false, "Enable verbose logging.")
	noColorPtr := flag.Bool("no-color", false, "Disable colored output during logging.")

	flag.Parse()

	// call initLoggers for use
	logger.InitLoggers(*noColorPtr)

	if *verbosePtr {
		flagString := fmt.Sprintf("\n\tDirectory: %s\n\tPort: %s\n\tVerbose: %t\n\tANSI Colors: %t\n",
			*dirPtr, *portPtr, *verbosePtr, *noColorPtr)
		logger.LogInfo.Println(flagString)
	}

	info, err := os.Stat(*dirPtr)
	if os.IsNotExist(err) {
		logger.LogErr.Fatalf("Directory %s does not exist.", *dirPtr)
	} else if os.IsPermission(err) {
		logger.LogErr.Fatalf("Directory %s does not have edit permission.", *dirPtr)
	}
	if !info.IsDir() {
		logger.LogErr.Fatalf("%s is not a directory", *dirPtr)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.LogErr.Fatalln("Something went wrong in creating watcher.")
	}
	defer watcher.Close()

	err = dir_watcher.AddWatchDirRecursively(watcher, *dirPtr)
	if err != nil {
		logger.LogErr.Fatalln(err)
	}
	logger.LogWarn.Printf("Watching: %s", *dirPtr)

	go dir_watcher.DirWatcher(watcher, *verbosePtr)
	<-make(chan struct{})
}
