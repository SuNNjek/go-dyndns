package client

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-dyndns/addrproviders"
	"go-dyndns/cache"
	"go-dyndns/log"
	"go-dyndns/updater"
	"net"
	"testing"
	"time"
)

func TestDynDnsClient_doUpdate_ShouldUpdateWhenCacheIsEmpty(t *testing.T) {
	cacheMock := &cache.MockCache{}
	providerMock := &addrproviders.MockProvider{}
	updaterMock := &updater.MockUpdater{}

	client := createClient(10*time.Minute, cacheMock, providerMock, updaterMock)
	ip := net.ParseIP("127.0.0.1")

	cacheMock.On("GetLastIp").
		Return(nil, nil)

	providerMock.On("GetIP", mock.Anything).
		Return(ip, nil)

	updaterMock.On("UpdateIP", mock.Anything, ip).
		Return(nil)

	cacheMock.On("SetLastIp", ip).
		Return(nil)

	err := client.doUpdate(context.Background())

	assert.Nil(t, err)

	cacheMock.AssertExpectations(t)
	providerMock.AssertExpectations(t)
	updaterMock.AssertExpectations(t)
}
func TestDynDnsClient_doUpdate_ShouldUpdateWhenCacheReturnsError(t *testing.T) {
	cacheMock := &cache.MockCache{}
	providerMock := &addrproviders.MockProvider{}
	updaterMock := &updater.MockUpdater{}

	client := createClient(10*time.Minute, cacheMock, providerMock, updaterMock)
	ip := net.ParseIP("127.0.0.1")

	cacheMock.On("GetLastIp").
		Return(nil, errors.New("failed to get IP address"))

	providerMock.On("GetIP", mock.Anything).
		Return(ip, nil)

	updaterMock.On("UpdateIP", mock.Anything, ip).
		Return(nil)

	cacheMock.On("SetLastIp", ip).
		Return(nil)

	err := client.doUpdate(context.Background())

	assert.Nil(t, err)

	cacheMock.AssertExpectations(t)
	providerMock.AssertExpectations(t)
	updaterMock.AssertExpectations(t)
}

func TestDynDnsClient_doUpdate_ShouldUpdateWhenIpHasChanged(t *testing.T) {
	cacheMock := &cache.MockCache{}
	providerMock := &addrproviders.MockProvider{}
	updaterMock := &updater.MockUpdater{}

	client := createClient(10*time.Minute, cacheMock, providerMock, updaterMock)
	lastIp := net.ParseIP("127.0.0.1")
	newIp := net.ParseIP("127.0.0.2")

	cacheMock.On("GetLastIp").
		Return(lastIp, nil)

	providerMock.On("GetIP", mock.Anything).
		Return(newIp, nil)

	updaterMock.On("UpdateIP", mock.Anything, newIp).
		Return(nil)

	cacheMock.On("SetLastIp", newIp).
		Return(nil)

	err := client.doUpdate(context.Background())

	assert.Nil(t, err)

	cacheMock.AssertExpectations(t)
	providerMock.AssertExpectations(t)
	updaterMock.AssertExpectations(t)
}

func TestDynDnsClient_doUpdate_ShouldNotUpdateWhenIpHasntChanged(t *testing.T) {
	cacheMock := &cache.MockCache{}
	providerMock := &addrproviders.MockProvider{}
	updaterMock := &updater.MockUpdater{}

	client := createClient(10*time.Minute, cacheMock, providerMock, updaterMock)
	ip := net.ParseIP("127.0.0.1")

	cacheMock.On("GetLastIp").
		Return(ip, nil)

	providerMock.On("GetIP", mock.Anything).
		Return(ip, nil)

	err := client.doUpdate(context.Background())

	assert.Nil(t, err)

	cacheMock.AssertExpectations(t)
	providerMock.AssertExpectations(t)
	updaterMock.AssertExpectations(t)
}

func TestDynDnsClient_doUpdate_ShouldFailIfIpCouldNotBeRetrieved(t *testing.T) {
	cacheMock := &cache.MockCache{}
	providerMock := &addrproviders.MockProvider{}
	updaterMock := &updater.MockUpdater{}

	client := createClient(10*time.Minute, cacheMock, providerMock, updaterMock)
	ip := net.ParseIP("127.0.0.1")

	cacheMock.On("GetLastIp").
		Return(ip, nil)

	providerMock.On("GetIP", mock.Anything).
		Return(nil, errors.New("no IP for you lmao"))

	err := client.doUpdate(context.Background())

	assert.Error(t, err)

	cacheMock.AssertExpectations(t)
	providerMock.AssertExpectations(t)
	updaterMock.AssertExpectations(t)
}

func TestDynDnsClient_Run(t *testing.T) {
	cacheMock := &cache.MockCache{}
	providerMock := &addrproviders.MockProvider{}
	updaterMock := &updater.MockUpdater{}

	client := createClient(2*time.Second, cacheMock, providerMock, updaterMock)
	ip := net.ParseIP("127.0.0.1")

	// doUpdate should be called twice, so at least the cacheMock should be checked twice as well
	cacheMock.On("GetLastIp").
		Twice().
		Return(nil, nil)

	providerMock.On("GetIP", mock.Anything).
		Return(ip, nil)

	updaterMock.On("UpdateIP", mock.Anything, ip).
		Return(nil)

	cacheMock.On("SetLastIp", ip).
		Return(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := client.Run(ctx)

	// Check if error is passed on correctly
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	cacheMock.AssertExpectations(t)
	providerMock.AssertExpectations(t)
	updaterMock.AssertExpectations(t)
}

func createClient(delay time.Duration, cache cache.Cache, provider addrproviders.AddressProvider, updater updater.Updater) *DynDnsClient {
	return newDynDnsClient(
		&clientConfig{
			Delay: delay,
		},
		log.CreateTestLogger(),
		cache,
		provider,
		updater,
	)
}
