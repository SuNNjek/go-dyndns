package addrproviders

import (
	"context"
	"go-dyndns/util"
	"net"
)

type ProviderType string

const (
	Web      ProviderType = "web"
	FritzBox ProviderType = "fritzbox"
)

type AddressProvider interface {
	GetIPv4(ctx context.Context) (net.IP, error)
}

type AddressV6Provider interface {
	AddressProvider

	GetIPv6Prefix(ctx context.Context) (*util.IPv6Prefix, error)
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
