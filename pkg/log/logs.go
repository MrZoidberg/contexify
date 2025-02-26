package log

import (
	"github.com/fatih/color"
	log "github.com/go-pkgz/lgr"
)

var logger = log.Default()
var colorized = true
var debug = false

var errorFunc = func(s string) string { return color.New(color.FgHiRed).Sprint(s) }
var warnFunc = func(s string) string { return color.New(color.FgHiYellow).Sprint(s) }
var InfoFunc = func(s string) string { return color.New(color.FgHiWhite).Sprint(s) }
var debugFunc = func(s string) string { return color.New(color.FgWhite).Sprint(s) }

func SetupLog(dbg bool, noColor bool) {
	colorized = !noColor
	debug = dbg

	logger = log.New(log.Out(color.Output), log.Err(color.Error), log.Format("{{.Message}}"))
}

func Errorf(format string, args ...interface{}) {
	// format = "[ERROR] " + format
	if colorized {
		logger.Logf(errorFunc(format), args...)
		return
	}
	logger.Logf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	// format = "[WARN] " + format
	if colorized {
		logger.Logf(warnFunc(format), args...)
		return
	}
	logger.Logf(format, args...)
}

func Infof(format string, args ...interface{}) {
	// format = "[INFO] " + format
	if colorized {
		logger.Logf(InfoFunc(format), args...)
		return
	}
	logger.Logf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	if !debug {
		return
	}
	// format = "[DEBUG] " + format
	if colorized {
		logger.Logf(debugFunc(format), args...)
		return
	}
	logger.Logf(format, args...)
}
