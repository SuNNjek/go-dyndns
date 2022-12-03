//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package updater

import (
	"github.com/google/wire"
	"go-dyndns/util"
)

func CreateUpdater() (Updater, error) {
	wire.Build(dynDnsSet, util.DefaultHttpClientValue, wire.Bind(new(Updater), new(*dynDnsUpdater)))
	return &dynDnsUpdater{}, nil
}
