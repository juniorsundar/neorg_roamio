package logger

import (
	"fmt"
	"log"
	"os"
)

// Creating custom coloured loggers
var (
	LogErr  *log.Logger
	LogWarn *log.Logger
	LogInfo *log.Logger
)

// Create prefixes to log messages with ANSI codes. I chose to do this over 3rd
// party packages to take advantage of Go's rich standard libraries.
//
// @param noColorCoding: Disable ANSI encoding.
func InitLoggers(noColorCoding bool) {
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

	LogErr = log.New(os.Stderr, errorPrefix, log.LstdFlags|log.Lshortfile)
	LogWarn = log.New(os.Stdout, warnPrefix, log.LstdFlags)
	LogInfo = log.New(os.Stdout, infoPrefix, log.LstdFlags)
}
