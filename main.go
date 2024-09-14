package main

import (
	"flag"
	"fmt"
	"log"

// Creating custom coloured loggers
var (
	errorLogger *log.Logger
	warnLogger  *log.Logger
	infoLogger  *log.Logger
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

	errorLogger = log.New(os.Stderr, errorPrefix, log.LstdFlags|log.Lshortfile)
	warnLogger = log.New(os.Stdout, warnPrefix, log.LstdFlags)
	infoLogger = log.New(os.Stdout, infoPrefix, log.LstdFlags)
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
		infoLogger.Println(flagString)
	}

	info, err := os.Stat(*dirPtr)
	if os.IsNotExist(err) {
		errorLogger.Fatalf("Directory %s does not exist.", *dirPtr)
	} else if os.IsPermission(err) {
		errorLogger.Fatalf("Directory %s does not have edit permission.", *dirPtr)
	}
	if !info.IsDir() {
		errorLogger.Fatalf("%s is not a directory", *dirPtr)
	}
}
