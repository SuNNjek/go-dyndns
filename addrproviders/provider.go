package addrproviders

import (
	"fmt"
	"net"
	"strings"
)

type AddressProvider interface {
	GetIP() (net.IP, error)
}

func CreateProvider(provider string) (AddressProvider, error) {
	switch strings.ToLower(provider) {
	case "web":
		return createWebProvider()

	case "fritzbox":
		return createFritzBoxProvider()
	}

	return nil, fmt.Errorf("unknown provider type: %s", provider)
}
