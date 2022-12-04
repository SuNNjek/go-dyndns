package log

import "github.com/kelseyhightower/envconfig"

type loggerConfig struct {
	level LogLevel `default:"info"`
}

func loadConfig() (*loggerConfig, error) {
	var config loggerConfig
	if err := envconfig.Process("log", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
