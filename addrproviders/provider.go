package addrproviders

import (
	"go-dyndns/util"
	"net"
)

type ProviderType string

const (
	Web      ProviderType = "web"
	FritzBox ProviderType = "fritzbox"
)

type AddressProvider interface {
	GetIP() (net.IP, error)
}

func CreateProvider(provider ProviderType, httpClient util.HttpClient) (AddressProvider, error) {
	switch provider {
	case Web:
		return createWebProvider(httpClient)

	case FritzBox:
		return createFritzBoxProvider(httpClient)
	}

	return nil, UnknownProvider(string(provider))
}
