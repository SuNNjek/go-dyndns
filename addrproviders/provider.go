package addrproviders

import (
	"fmt"
	"net"
)

type ProviderType string

var (
	Web      ProviderType = "web"
	FritzBox ProviderType = "fritzbox"
)

type AddressProvider interface {
	GetIP() (net.IP, error)
}

func CreateProvider(provider ProviderType) (AddressProvider, error) {
	switch provider {
	case Web:
		return createWebProvider()

	case FritzBox:
		return createFritzBoxProvider()
	}

	return nil, fmt.Errorf("unknown provider type: %s", provider)
}
