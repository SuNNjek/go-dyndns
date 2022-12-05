package cache

import (
	"github.com/kelseyhightower/envconfig"
	"net"
)

type Cache interface {
	GetLastIp() (net.IP, error)
	SetLastIp(ip net.IP) error
}

type cacheConfig struct {
	File string
}

func CreateCache() (Cache, error) {
	var config cacheConfig
	if err := envconfig.Process("cache", &config); err != nil {
		return nil, err
	}

	if config.File != "" {
		return newFileCache(config.File), nil
	}

	return &memoryCache{}, nil
}
