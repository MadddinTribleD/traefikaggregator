package log

import (
	"fmt"
	"os"
)

type LogLevel string

const (
	InfoLevel  LogLevel = "INFO"
	DebugLevel LogLevel = "DEBUG"
	ErrorLevel LogLevel = "ERROR"
)

func Log(level LogLevel, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)

	file := os.Stdout
	if level == ErrorLevel {
		file = os.Stderr
	}

	file.WriteString(fmt.Sprintf("%s: %s\n", level, message))
}

func Info(format string, args ...interface{}) {
	Log(InfoLevel, format, args...)
}

func Debug(format string, args ...interface{}) {
	Log(DebugLevel, format, args...)
}

func Error(format string, args ...interface{}) {
	Log(ErrorLevel, format, args...)
}
