package logger

import (
	"fmt"
	"log"
	"os"
)

var logLevel = Level.ALL

func SetLogLevel(level loggerLevel) {
	logLevel = level
}

func Trace(args ...interface{}) {
	logWithLevel(Level.TRACE, args...)
}

func Tracef(format string, args ...interface{}) {
	logfWithLevel(Level.TRACE, format, args...)
}

func Debug(args ...interface{}) {
	logWithLevel(Level.DEBUG, args...)
}

func Debugf(format string, args ...interface{}) {
	logfWithLevel(Level.DEBUG, format, args...)
}

func Info(args ...interface{}) {
	logWithLevel(Level.INFO, args...)
}

func Infof(format string, args ...interface{}) {
	logfWithLevel(Level.INFO, format, args...)
}

func Warn(args ...interface{}) {
	logWithLevel(Level.WARN, args...)
}

func Warnf(format string, args ...interface{}) {
	logfWithLevel(Level.WARN, format, args...)
}

func Error(args ...interface{}) {
	logWithLevel(Level.ERROR, args...)
}

func Errorf(format string, args ...interface{}) {
	logfWithLevel(Level.ERROR, format, args...)
}

func Fatal(args ...interface{}) {
	logWithLevel(Level.FATAL, args...)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	logfWithLevel(Level.FATAL, format, args...)
	os.Exit(1)
}

func logWithLevel(level loggerLevel, args ...interface{}) {
	if logLevel.value <= level.value {
		log.Printf("%s%s:%s %s", level.color, level.name, resetcolor, fmt.Sprintln(args...))
	}
}

func logfWithLevel(level loggerLevel, format string, args ...interface{}) {
	if logLevel.value <= level.value {
		log.Printf(fmt.Sprintf("%s%s:%s %s\n", level.color, level.name, resetcolor, format), args...)
	}
}
