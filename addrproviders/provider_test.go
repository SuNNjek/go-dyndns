package addrproviders

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateWebProvider(t *testing.T) {
	t.Setenv("WEB_URL", "example.com")

	providerWeb, err := CreateProvider(Web)

	assert.Nil(t, err)
	assert.IsType(t, new(webProvider), providerWeb)
}

func TestCreateFritzboxProvider(t *testing.T) {
	providerFritz, err := CreateProvider(FritzBox)

	assert.Nil(t, err)
	assert.IsType(t, new(fritzBoxProvider), providerFritz)
}

func TestUnknownProvider(t *testing.T) {
	provider, err := CreateProvider("asdf")

	assert.Nil(t, provider)
	assert.IsType(t, new(UnknownProviderError), err)
}
