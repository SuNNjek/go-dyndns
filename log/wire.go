//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package log

import (
	"github.com/google/wire"
)

func CreateLogger() (Logger, error) {
	wire.Build(loadConfig, writerLoggerSet)
	return &writerLogger{}, nil
}
