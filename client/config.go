package client

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type ClientConfig struct {
	IpProvider string        `default:"web"`
	Delay      time.Duration `default:"10m"`
}

func LoadConfig() (*ClientConfig, error) {
	var config ClientConfig
	if err := envconfig.Process("client", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
