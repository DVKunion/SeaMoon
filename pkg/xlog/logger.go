package xlog

import (
	"fmt"
	"log/slog"
)

var defaultLog = slog.Default()

// Debug logs at LevelDebug.
func Debug(index string, msg string, args ...any) {
	defaultLog.Debug(fmt.Sprintf("[%s] %s", index, msg), args...)
}

// Info logs at LevelInfo.
func Info(index string, msg string, args ...any) {
	defaultLog.Info(fmt.Sprintf("[%s] %s", index, msg), args...)
}

// Warn logs at LevelWarn.
func Warn(index string, msg string, args ...any) {
	defaultLog.Info(fmt.Sprintf("[%s] %s", index, msg), args...)
}

// Error logs at LevelError.
func Error(index string, msg string, args ...any) {
	defaultLog.Error(fmt.Sprintf("[%s] %s", index, msg), args...)
}
