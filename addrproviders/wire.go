//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package addrproviders

import "github.com/google/wire"

func createWebProvider() (*webProvider, error) {
	wire.Build(webSet)
	return &webProvider{}, nil
}

func createFritzBoxProvider() (*fritzBoxProvider, error) {
	wire.Build(fritzBoxSet)
	return &fritzBoxProvider{}, nil
}
