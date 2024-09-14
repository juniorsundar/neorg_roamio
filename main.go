package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Creating custom coloured loggers
var (
	logErr  *log.Logger
	logWarn *log.Logger
	logInfo *log.Logger
)

// Create prefixes to log messages with ANSI codes. I chose to do this over 3rd
// party packages to take advantage of Go's rich standard libraries.
//
// @param noColorCoding: Disable ANSI encoding.
func initLoggers(noColorCoding bool) {
	red := "\033[31m"
	green := "\033[32m"
	yellow := "\033[33m"
	reset := "\033[0m"

	var errorPrefix, warnPrefix, infoPrefix string

	if noColorCoding {
		errorPrefix = fmt.Sprintf("ERROR: ")
		warnPrefix = fmt.Sprintf("WARN: ")
		infoPrefix = fmt.Sprintf("INFO: ")
	} else {
		errorPrefix = fmt.Sprintf("%sERROR: %s", red, reset)
		warnPrefix = fmt.Sprintf("%sWARN: %s", yellow, reset)
		infoPrefix = fmt.Sprintf("%sINFO: %s", green, reset)
	}

	logErr = log.New(os.Stderr, errorPrefix, log.LstdFlags|log.Lshortfile)
	logWarn = log.New(os.Stdout, warnPrefix, log.LstdFlags)
	logInfo = log.New(os.Stdout, infoPrefix, log.LstdFlags)
}

// Initialise watcher for subdirectories recursively.
//
// @param watcher: Main `watcher` from main loop.
//
// @param dir: Root directory as source of recursive search.
func addWatchDirRecursively(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name()[0] == '.' {
				logWarn.Println("Skipping hidden directory:", path)
				return filepath.SkipDir
			}

			err = watcher.Add(path)
			if err != nil {
				return err
			}
			logWarn.Printf("Watching: %s", path)
		}
		return nil
	})
}

func dirWatcher(watcher *fsnotify.Watcher, verbose bool) {
	for {
		// Select is like Switch in that it observes what comes out of the
		// channels. In this case, the two channels of interest are
		// watcher.Events/Write. If any one of them gets a message then their
		// respective case is triggered.
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			file, err := os.Stat(event.Name)

			if verbose {
				logInfo.Println("Event:", event.Name)
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				if err == nil && file.IsDir() {
					err = addWatchDirRecursively(watcher, event.Name)
					if err != nil {
						logErr.Fatalln(err)
					}
				}
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove {
				if os.IsNotExist(err) {
					watcher.Remove(event.Name)
					logWarn.Printf("Removed %s from watcher.", event.Name)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logErr.Println(err)
		}
	}
}

func main() {

	dirPtr := flag.String("dir", ".", "Roam directory address.")
	portPtr := flag.String("port", "8080", "Roam server port.")
	verbosePtr := flag.Bool("verbose", false, "Enable verbose logging.")
	noColorPtr := flag.Bool("no-color", false, "Disable colored output during logging.")

	flag.Parse()

	// call initLoggers for use
	initLoggers(*noColorPtr)

	if *verbosePtr {
		flagString := fmt.Sprintf("\n\tDirectory: %s\n\tPort: %s\n\tVerbose: %t\n\tANSI Colors: %t\n",
			*dirPtr, *portPtr, *verbosePtr, *noColorPtr)
		logInfo.Println(flagString)
	}

	info, err := os.Stat(*dirPtr)
	if os.IsNotExist(err) {
		logErr.Fatalf("Directory %s does not exist.", *dirPtr)
	} else if os.IsPermission(err) {
		logErr.Fatalf("Directory %s does not have edit permission.", *dirPtr)
	}
	if !info.IsDir() {
		logErr.Fatalf("%s is not a directory", *dirPtr)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logErr.Fatalln("Something went wrong in creating watcher.")
	}
	defer watcher.Close()

	err = addWatchDirRecursively(watcher, *dirPtr)
	if err != nil {
		logErr.Fatalln(err)
	}
	logWarn.Printf("Watching: %s", *dirPtr)

	go dirWatcher(watcher, *verbosePtr)
	<-make(chan struct{})
}
