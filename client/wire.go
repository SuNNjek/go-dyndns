//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package client

import (
	"github.com/google/wire"
	"go-dyndns/addrproviders"
	"go-dyndns/updater"
)

func CreateClient() (*DynDnsClient, error) {
	wire.Build(clientSet, addrproviders.CreateProvider, updater.CreateUpdater)
	return &DynDnsClient{}, nil
}
