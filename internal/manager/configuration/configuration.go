package configuration

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Configuration holds service configuration.
type Configuration struct {
	Service ServiceConfiguration
	Logging LoggingConfiguration
	DB      DBConfiguration
}

// NewFromEnvironment loads configuration from the environment.
func NewFromEnvironment() (*Configuration, error) {
	config := Configuration{}

	if err := envconfig.Process("ley_manager", &config); err != nil {
		return nil, fmt.Errorf("Unable to load configuration from environment: %w", err)
	}

	return &config, nil
}
