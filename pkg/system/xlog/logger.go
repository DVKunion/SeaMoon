package xlog

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gookit/color"
)

type Logger struct {
	*slog.Logger
}

func (l Logger) Write(p []byte) (n int, err error) {
	s := strings.Replace(string(p), "\n", "", -1)
	group := "[GIN] "
	if strings.Contains(s, "[GIN") {
		group = ""
	}
	if strings.Contains(s, "debug") {
		l.Debug(color.Gray.Sprintf("%s%s", group, s))
	}
	if strings.Contains(s, "WARNING") || strings.Contains(s, "| 403 |") ||
		strings.Contains(s, "| 401 |") || strings.Contains(s, "| 400 |") {
		l.Warn(color.Yellow.Sprintf("%s%s", group, s))
	}

	if strings.Contains(s, "| 500 |") {
		l.Error(color.Red.Sprintf("%s%s", group, s))
	} else {
		l.Info(color.Cyan.Sprintf("%s%s", group, s))
	}
	return len(p), nil
}

var defaultLog = &Logger{slog.Default()}

var index = "UNEXPECT"

func GetLogger() *Logger {
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
	defaultLog.Info(color.Cyan.Sprintf("[%s] %s", index, msg), args...)
}

// Warn logs at LevelWarn.
func Warn(msg string, args ...any) {
	if strings.Contains(msg, " ") {
		index = strings.ToUpper(strings.Split(msg, " ")[0])
	}
	defaultLog.Info(color.Yellow.Sprintf("[%s] %s", index, msg), args...)
}

// Error logs at LevelError.
func Error(msg string, args ...any) {
	if strings.Contains(msg, " ") {
		index = strings.ToUpper(strings.Split(msg, " ")[0])
	}
	defaultLog.Error(color.Red.Sprintf("[%s] %s", index, msg), args...)
}
