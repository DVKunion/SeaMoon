package xlog

import (
	"fmt"
	"log/slog"
	"strings"
)

var defaultLog = slog.Default()

var index = "UNEXPECT"

func Logger() *slog.Logger {
	return defaultLog
}

// Debug logs at LevelDebug.
func Debug(msg string, args ...any) {
	if strings.Contains(msg, " ") {
		index = strings.ToUpper(strings.Split(msg, " ")[0])
	}
	defaultLog.Debug(fmt.Sprintf("[%s] %s", index, msg), args...)
}

// Info logs at LevelInfo.
func Info(msg string, args ...any) {
	if strings.Contains(msg, " ") {
		index = strings.ToUpper(strings.Split(msg, " ")[0])
	}
	defaultLog.Info(fmt.Sprintf("[%s] %s", index, msg), args...)
}

// Warn logs at LevelWarn.
func Warn(msg string, args ...any) {
	if strings.Contains(msg, " ") {
		index = strings.ToUpper(strings.Split(msg, " ")[0])
	}
	defaultLog.Info(fmt.Sprintf("[%s] %s", index, msg), args...)
}

// Error logs at LevelError.
func Error(msg string, args ...any) {
	if strings.Contains(msg, " ") {
		index = strings.ToUpper(strings.Split(msg, " ")[0])
	}
	defaultLog.Error(fmt.Sprintf("[%s] %s", index, msg), args...)
}
