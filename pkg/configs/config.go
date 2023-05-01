package configs

import (
	"github.com/joeshaw/envdecode"

	"github.com/pkg/errors"
)

// Config represents all HTTP server configuration options.
type Config struct {
	ServerName            string `env:"SERVER_NAME"`
	ServerHost            string `env:"SERVER_HOST"`
	ServerPort            int    `env:"SERVER_PORT"`
	LogLevel              string `env:"LOG_LEVEL"`
	LogFormat             string `env:"LOG_FORMAT"`
	CosmosSDKGRPCEndpoint string `env:"COSMOS_SDK_GRPC_ENDPOINT"`
}

// NewConfig constructs a new instance of ServerConfig via decoding
// the mapped env vars with envdecode library.
func NewConfig() (*Config, error) {
	var config Config

	err := envdecode.Decode(&config)

	if err != nil {
		if err != envdecode.ErrNoTargetFieldsAreSet {
			return nil, errors.WithStack(err)
		}
	}

	return &config, nil
}
