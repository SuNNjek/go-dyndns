//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package client

import (
	"github.com/google/wire"
	"go-dyndns/addrproviders"
	"go-dyndns/cache"
	"go-dyndns/log"
	"go-dyndns/updater"
	"go-dyndns/util"
)

func CreateClient(logger log.Logger) (*DynDnsClient, error) {
	wire.Build(
		clientSet,
		util.DefaultHttpClientValue,
		cache.CreateCache,
		addrproviders.CreateProvider,
		updater.CreateUpdater,
	)

	return &DynDnsClient{}, nil
}
