package client

import (
	"github.com/kelseyhightower/envconfig"
	"go-dyndns/addrproviders"
	"time"
)

type clientConfig struct {
	IpProvider addrproviders.ProviderType `default:"web"`
	EnableIPv6 bool                       `default:"false"`
	Delay      time.Duration              `default:"10m"`
}

func loadConfig() (*clientConfig, error) {
	var config clientConfig
	if err := envconfig.Process("client", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
