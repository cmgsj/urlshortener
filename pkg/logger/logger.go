package logger

import (
	"fmt"
	"log"
	"os"
)

var logLevel = Level.ALL

func SetLogLevel(level LoggingLevel) {
	logLevel = level
	Info("Log level set to", level.name)
}

func Trace(args ...interface{}) {
	logWithLevel(Level.TRACE, args...)
}

func Debug(args ...interface{}) {
	logWithLevel(Level.DEBUG, args...)
}

func Info(args ...interface{}) {
	logWithLevel(Level.INFO, args...)
}

func Warn(args ...interface{}) {
	logWithLevel(Level.WARN, args...)
}

func Error(args ...interface{}) {
	logWithLevel(Level.ERROR, args...)
}

func Fatal(args ...interface{}) {
	logWithLevel(Level.FATAL, args...)
	os.Exit(1)
}

func logWithLevel(level LoggingLevel, args ...interface{}) {
	if logLevel.value <= level.value {
		log.Printf("%s%s:%s %s", level.color, level.name, resetColor, fmt.Sprintln(args...))
	}
}
