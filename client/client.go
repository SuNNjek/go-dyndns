package client

import (
	"context"
	"errors"
	"github.com/google/wire"
	"go-dyndns/addrproviders"
	"go-dyndns/cache"
	"go-dyndns/log"
	"go-dyndns/updater"
	"time"
)

var clientSet = wire.NewSet(loadConfig, newDynDnsClient, wire.FieldsOf(new(*clientConfig), "IpProvider"))

type DynDnsClient struct {
	config   *clientConfig
	logger   log.Logger
	cache    cache.Cache
	provider addrproviders.AddressProvider
	updater  updater.Updater
}

func newDynDnsClient(
	config *clientConfig,
	logger log.Logger,
	cache cache.Cache,
	provider addrproviders.AddressProvider,
	updater updater.Updater,
) (*DynDnsClient, error) {
	return &DynDnsClient{
		config:   config,
		logger:   logger,
		cache:    cache,
		updater:  updater,
		provider: provider,
	}, nil
}

func (c *DynDnsClient) Run(ctx context.Context) error {
	if err := ctx.Err(); errors.Is(context.Canceled, err) {
		return err
	}

	ticker := time.NewTicker(c.config.Delay)
	defer ticker.Stop()

	for {
		if err := c.doUpdate(ctx); err != nil {
			c.logger.Warn("Error while updating IP address: %v", err)
		}

		select {
		// Continue the loop when the ticker has sent a signal
		case <-ticker.C:
			continue

		// Abort when the context is done
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *DynDnsClient) doUpdate(ctx context.Context) error {
	lastIp, err := c.cache.GetLastIp()
	if err != nil {
		c.logger.Warn("Failed to get last IP from cache: %v", err)
	}

	ip, err := c.provider.GetIP(ctx)
	if err != nil {
		return err
	}

	// Don't send an update request if the IP matches that from cache.
	// DynDNS providers don't like it when you send too many requests for the same IP ;)
	if ip.Equal(lastIp) {
		c.logger.Info("IP is already up to date")
		return nil
	}

	if err = c.updater.UpdateIP(ctx, ip); err != nil {
		return err
	} else {
		c.logger.Info("IP successfully updated")
	}

	return c.cache.SetLastIp(ip)
}
