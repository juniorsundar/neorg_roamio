package logger

import (
	"fmt"
	"io"
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
// @param colorCoding: Enable ANSI encoding.  
// @param verbosity: Enable verbose output stream.
func InitLoggers(colorCoding bool, verbosity bool) {
	red := "\033[31m"
	green := "\033[32m"
	yellow := "\033[33m"
	reset := "\033[0m"

	var errorPrefix, warnPrefix, infoPrefix string

	if !colorCoding {
		errorPrefix = "ERROR: "
		warnPrefix = "WARN: "
		infoPrefix = "INFO: "
	} else {
		errorPrefix = fmt.Sprintf("%sERROR: %s", red, reset)
		warnPrefix = fmt.Sprintf("%sWARN: %s", yellow, reset)
		infoPrefix = fmt.Sprintf("%sINFO: %s", green, reset)
	}

	// Always output errors to STDERR
	LogErr = log.New(os.Stderr, errorPrefix, log.LstdFlags|log.Lshortfile)

	// Conditionally output WARN and INFO logs based on outputWarnInfo flag
	if verbosity {
		LogWarn = log.New(os.Stdout, warnPrefix, log.LstdFlags)
		LogInfo = log.New(os.Stdout, infoPrefix, log.LstdFlags)
	} else {
		LogWarn = log.New(io.Discard, warnPrefix, log.LstdFlags)
		LogInfo = log.New(io.Discard, infoPrefix, log.LstdFlags)
	}
}
