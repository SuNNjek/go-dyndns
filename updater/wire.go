//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package updater

import (
	"github.com/google/wire"
	"go-dyndns/log"
	"go-dyndns/util"
)

func CreateUpdater(logger log.Logger, httpClient util.HttpClient) (Updater, error) {
	wire.Build(dynDnsSet, util.NewFilePasswordProvider, wire.Bind(new(Updater), new(*dynDnsUpdater)))
	return &dynDnsUpdater{}, nil
}
