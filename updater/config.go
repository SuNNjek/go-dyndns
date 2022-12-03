package updater

import (
	"github.com/kelseyhightower/envconfig"
)

type UpdaterConfig struct {
	Host         string   `required:"true"`
	User         string   `required:"true"`
	PasswordFile string   `required:"true"`
	Domains      []string `required:"true"`
}

func LoadConfig() (*UpdaterConfig, error) {
	var config UpdaterConfig
	if err := envconfig.Process("updater", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
