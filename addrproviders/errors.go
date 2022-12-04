package addrproviders

import (
	"errors"
	"fmt"
)

var (
	InvalidResponseError = errors.New("invalid response received")
	ParseIpError         = errors.New("failed to parse IP")
)

type UnknownProviderError struct {
	providerType string
}

func UnknownProvider(providerType string) *UnknownProviderError {
	return &UnknownProviderError{providerType: providerType}
}

func (u *UnknownProviderError) Error() string {
	return fmt.Sprintf("unknown provider type: %s", u.providerType)
}
