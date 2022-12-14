package cache

import (
	"github.com/kelseyhightower/envconfig"
	"go-dyndns/updater"
)

type Cache interface {
	GetLastRequest() (*updater.UpdateRequest, error)
	SetLastRequest(req *updater.UpdateRequest) error
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
