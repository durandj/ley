package configuration

import (
	"fmt"

	"github.com/durandj/ley/internal/common/configuration"
)

// ServiceConfiguration contains general configuration required for
// just running the service.
type ServiceConfiguration struct {
	EnvironmentType configuration.EnvironmentType `envconfig:"environment_type"`
	Host            string                        `default:"127.0.0.1"`
	Port            int                           `default:"8080"`
}

// Address gets the host and port combination that the service should
// run on.
func (config ServiceConfiguration) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
