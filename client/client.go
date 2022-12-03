package client

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"go-dyndns/addrproviders"
	"go-dyndns/updater"
	"time"
)

var Set = wire.NewSet(LoadConfig, NewDynDnsClient)

type DynDnsClient struct {
	config   *ClientConfig
	provider addrproviders.AddressProvider
	updater  updater.Updater
}

func NewDynDnsClient(config *ClientConfig, updater updater.Updater) (*DynDnsClient, error) {
	provider, err := addrproviders.CreateProvider(config.IpProvider)
	if err != nil {
		return nil, err
	}

	return &DynDnsClient{
		config:   config,
		updater:  updater,
		provider: provider,
	}, nil
}

func (c *DynDnsClient) Run(ctx context.Context) error {
	ticker := time.NewTicker(c.config.Delay)
	defer ticker.Stop()

	for {
		if err := c.doUpdate(); err != nil {
			fmt.Println(err)
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

func (c *DynDnsClient) doUpdate() error {
	ip, err := c.provider.GetIP()
	if err != nil {
		return err
	}

	return c.updater.UpdateIP(ip)
}
