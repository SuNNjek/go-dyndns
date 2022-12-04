package log

import (
	"errors"
	"strings"
)

var (
	ParseLogLevelError = errors.New("failed to parse log level")
	InvalidLogLevel    = errors.New("invalid log level")
)

const (
	Trace LogLevel = iota
	Debug
	Info
	Warn
	Error
)

type LogLevel int

func (l LogLevel) String() string {
	if res, err := l.MarshalText(); err == nil {
		return strings.ToUpper(string(res))
	} else {
		return "UNKNOWN"
	}
}

func (l LogLevel) MarshalText() (text []byte, err error) {
	switch l {
	case Error:
		return []byte("error"), nil
	case Warn:
		return []byte("warn"), nil
	case Info:
		return []byte("info"), nil
	case Debug:
		return []byte("debug"), nil
	case Trace:
		return []byte("trace"), nil
	}

	return nil, InvalidLogLevel
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	if level, err := parseLogLevel(string(text)); err != nil {
		return err
	} else {
		*l = level
		return nil
	}
}

func parseLogLevel(str string) (LogLevel, error) {
	switch strings.ToLower(str) {
	case "error":
		return Error, nil
	case "warn":
		return Warn, nil
	case "info":
		return Info, nil
	case "debug":
		return Debug, nil
	case "trace":
		return Trace, nil
	}

	return -1, ParseLogLevelError
}
