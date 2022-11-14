package logger

import "strings"

const (
	white      = "\033[1;37m"
	blue       = "\033[1;34m"
	cyan       = "\033[1;36m"
	green      = "\033[1;32m"
	yellow     = "\033[1;33m"
	red        = "\033[1;31m"
	magenta    = "\033[1;35m"
	resetcolor = "\033[0m"
)

var Level = loggerLevels{
	ALL:   loggerLevel{name: "ALL", color: white, value: 0},
	TRACE: loggerLevel{name: "TRACE", color: blue, value: 1},
	DEBUG: loggerLevel{name: "DEBUG", color: cyan, value: 2},
	INFO:  loggerLevel{name: "INFO", color: green, value: 3},
	WARN:  loggerLevel{name: "WARN", color: yellow, value: 4},
	ERROR: loggerLevel{name: "ERROR", color: red, value: 5},
	FATAL: loggerLevel{name: "FATAL", color: magenta, value: 6},
	NONE:  loggerLevel{name: "OFF", color: white, value: 7},
}

type loggerLevel struct {
	name  string
	color string
	value int
}

type loggerLevels struct {
	ALL   loggerLevel
	TRACE loggerLevel
	DEBUG loggerLevel
	INFO  loggerLevel
	WARN  loggerLevel
	ERROR loggerLevel
	FATAL loggerLevel
	NONE  loggerLevel
}

func LevelFromString(s string) loggerLevel {
	switch strings.ToUpper(s) {
	case Level.ALL.name:
		return Level.ALL
	case Level.TRACE.name:
		return Level.TRACE
	case Level.DEBUG.name:
		return Level.DEBUG
	case Level.INFO.name:
		return Level.INFO
	case Level.WARN.name:
		return Level.WARN
	case Level.ERROR.name:
		return Level.ERROR
	case Level.FATAL.name:
		return Level.FATAL
	case Level.NONE.name:
		return Level.NONE
	default:
		Warnf("Invalid log level: %s\n", s)
		return Level.ALL
	}
}
