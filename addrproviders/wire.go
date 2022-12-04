//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package addrproviders

import (
	"github.com/google/wire"
	"go-dyndns/util"
)

func createWebProvider(httpClient util.HttpClient) (*webProvider, error) {
	wire.Build(webSet)
	return &webProvider{}, nil
}

func createFritzBoxProvider(httpClient util.HttpClient) (*fritzBoxProvider, error) {
	wire.Build(fritzBoxSet)
	return &fritzBoxProvider{}, nil
}
