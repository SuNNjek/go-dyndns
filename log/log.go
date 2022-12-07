package log

import (
	"fmt"
	"github.com/google/wire"
	"io"
	"os"
	"time"
)

var writerLoggerSet = wire.NewSet(createDefaultWriterLogger, wire.Bind(new(Logger), new(*writerLogger)))

type Logger interface {
	Log(level LogLevel, format string, args ...any)

	Fatal(format string, args ...any)
	Error(format string, args ...any)
	Warn(format string, args ...any)
	Info(format string, args ...any)
	Debug(format string, args ...any)
	Trace(format string, args ...any)
}

type writerLogger struct {
	normalOut io.Writer
	errorOut  io.Writer

	level LogLevel
}

func CreateTestLogger() Logger {
	return &writerLogger{
		normalOut: os.Stdout,
		errorOut:  os.Stderr,
		level:     Debug,
	}
}

func createDefaultWriterLogger(config *loggerConfig) *writerLogger {
	return &writerLogger{
		normalOut: os.Stdout,
		errorOut:  os.Stderr,
		level:     config.Level,
	}
}

func (w *writerLogger) Log(level LogLevel, format string, args ...any) {
	if level < w.level {
		return
	}

	var writer io.Writer
	switch level {
	case Error:
		writer = w.errorOut
		break

	default:
		writer = w.normalOut
		break
	}

	message := fmt.Sprintf(format, args...)
	formattedTime := time.Now().Format("2006-01-02 15:04:05.000")
	_, _ = fmt.Fprintf(writer, "[%s][%v]: %s\n", formattedTime, level, message)
}

func (w *writerLogger) Fatal(format string, args ...any) {
	w.Error(format, args...)
	os.Exit(1)
}

func (w *writerLogger) Error(format string, args ...any) {
	w.Log(Error, format, args...)
}

func (w *writerLogger) Warn(format string, args ...any) {
	w.Log(Warn, format, args...)
}

func (w *writerLogger) Info(format string, args ...any) {
	w.Log(Info, format, args...)
}

func (w *writerLogger) Debug(format string, args ...any) {
	w.Log(Debug, format, args...)
}

func (w *writerLogger) Trace(format string, args ...any) {
	w.Log(Trace, format, args...)
}
