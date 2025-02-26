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
var infoFunc = func(s string) string { return color.New(color.FgHiWhite).Sprint(s) }
var debugFunc = func(s string) string { return color.New(color.FgWhite).Sprint(s) }

// SetupLog sets up the logger
func SetupLog(dbg, noColor bool) {
	colorized = !noColor
	debug = dbg

	logger = log.New(log.Out(color.Output), log.Err(color.Error), log.Format("{{.Message}}"))
}

// Errorf logs error message
func Errorf(format string, args ...interface{}) {
	if colorized {
		logger.Logf(errorFunc(format), args...)
		return
	}
	logger.Logf(format, args...)
}

// Warnf logs warning message
func Warnf(format string, args ...interface{}) {
	if colorized {
		logger.Logf(warnFunc(format), args...)
		return
	}
	logger.Logf(format, args...)
}

// Infof logs info message
func Infof(format string, args ...interface{}) {
	if colorized {
		logger.Logf(infoFunc(format), args...)
		return
	}
	logger.Logf(format, args...)
}

// Debugf logs debug message if debug mode is on
func Debugf(format string, args ...interface{}) {
	if !debug {
		return
	}

	if colorized {
		logger.Logf(debugFunc(format), args...)
		return
	}
	logger.Logf(format, args...)
}
