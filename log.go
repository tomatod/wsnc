package main

import (
	"fmt"
	"log"
	"os"
)

const (
	LEVEL_INFO  int = 0
	LEVEL_WARN      = 1
	LEVEL_ERROR     = 2
	LEVEL_DEBUG     = 3
)

var logLevels = []string{
	"[info] ",
	"[warn] ",
	"[error]",
	"[debug]",
}

func logf(level int, format string, value ...interface{}) {
	if appConfig.IsLevel {
		format = logLevels[level] + " " + format
	}
	if appConfig.IsTime {
		log.Printf(format, value...)
		return
	}
	fmt.Fprintf(os.Stderr, format+"\n", value...)
}

func infoLogf(format string, value ...interface{}) {
	logf(LEVEL_INFO, format, value...)
}

func warnLogf(format string, value ...interface{}) {
	logf(LEVEL_WARN, format, value...)
}

func errLogf(format string, value ...interface{}) error {
	logf(LEVEL_ERROR, format, value...)
	return fmt.Errorf(format, value...)
}

func debugLogf(format string, value ...interface{}) {
	if !appConfig.IsDebug {
		return
	}
	logf(LEVEL_DEBUG, format, value...)
}

// Right log function.
func rlogf(format string, value ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", value...)
}

func rdebugf(format string, value ...interface{}) {
	if !appConfig.IsDebug {
		return
	}
	rlogf(format, value)
}
