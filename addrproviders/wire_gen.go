// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package addrproviders

import (
	"go-dyndns/util"
)

// Injectors from wire.go:

func createWebProvider(httpClient util.HttpClient) (*webProvider, error) {
	webProviderConfig, err := loadWebProviderConfig()
	if err != nil {
		return nil, err
	}
	addrprovidersWebProvider := newWebProvider(webProviderConfig, httpClient)
	return addrprovidersWebProvider, nil
}

func createFritzBoxProvider(httpClient util.HttpClient) (*fritzBoxProvider, error) {
	addrprovidersFritzBoxConfig, err := loadFritzBoxConfig()
	if err != nil {
		return nil, err
	}
	addrprovidersFritzBoxProvider := newFritzBoxProvider(addrprovidersFritzBoxConfig, httpClient)
	return addrprovidersFritzBoxProvider, nil
}
