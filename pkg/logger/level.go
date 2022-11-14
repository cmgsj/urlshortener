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
	resetColor = "\033[0m"
)

var Level = loggingLevels{
	ALL:   LoggingLevel{name: "ALL", color: white, value: 0},
	TRACE: LoggingLevel{name: "TRACE", color: blue, value: 1},
	DEBUG: LoggingLevel{name: "DEBUG", color: cyan, value: 2},
	INFO:  LoggingLevel{name: "INFO", color: green, value: 3},
	WARN:  LoggingLevel{name: "WARN", color: yellow, value: 4},
	ERROR: LoggingLevel{name: "ERROR", color: red, value: 5},
	FATAL: LoggingLevel{name: "FATAL", color: magenta, value: 6},
	NONE:  LoggingLevel{name: "OFF", color: white, value: 7},
}

type LoggingLevel struct {
	name  string
	color string
	value int
}

type loggingLevels struct {
	ALL   LoggingLevel
	TRACE LoggingLevel
	DEBUG LoggingLevel
	INFO  LoggingLevel
	WARN  LoggingLevel
	ERROR LoggingLevel
	FATAL LoggingLevel
	NONE  LoggingLevel
}

func LevelFromString(s string) LoggingLevel {
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
		Warn("Invalid log level:", s)
		return Level.ALL
	}
}
