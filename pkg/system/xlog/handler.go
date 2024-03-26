package xlog

//
//import (
//	"bytes"
//	"context"
//	"fmt"
//	"io"
//	"log/slog"
//	"os"
//
//	"github.com/gookit/color"
//)
//
//type Handler struct {
//	level slog.Level
//	path  string
//}
//
//func (h Handler) Enabled(_ context.Context, level slog.Level) bool {
//	return level >= h.level
//}
//
//func (h Handler) Handle(_ context.Context, r slog.Record) error {
//	var buf bytes.Buffer
//	buf.WriteString(h.color(r.Level, fmt.Sprintf("[%s]", r.Level.String())))
//	buf.WriteByte(' ')
//	buf.WriteString(h.color(r.Level, r.Message))
//	return h.print(&buf)
//}
//
//func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (h Handler) WithGroup(name string) slog.Handler {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (h Handler) color(level slog.Level, msg string) string {
//	switch level {
//	case slog.LevelDebug:
//		return color.Gray.Sprintf(msg)
//	case slog.LevelInfo:
//		return color.Cyan.Sprintf(msg)
//	case slog.LevelWarn:
//		return color.Yellow.Sprintf(msg)
//	case slog.LevelError:
//		return color.Red.Sprintf(msg)
//	}
//	return ""
//}
//
//func (h Handler) print(reader io.Reader) error {
//	if h.path != "" {
//		file, err := os.OpenFile(h.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
//		if err != nil {
//			return err
//		}
//		writer := io.MultiWriter(os.Stderr, file)
//		_, err = io.Copy(writer, reader)
//		return err
//	}
//	_, err := io.Copy(os.Stderr, reader)
//	return err
//}
