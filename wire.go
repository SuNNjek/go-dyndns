//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"go-dyndns/client"
	"go-dyndns/updater"
)

func Init() (*client.DynDnsClient, error) {
	wire.Build(
		updater.Set,
		client.Set,
	)

	return &client.DynDnsClient{}, nil
}
