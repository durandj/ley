package logging

import (
	"fmt"

	"github.com/durandj/ley/internal/common/configuration"
	"go.uber.org/zap"
)

// NewZapLogger creates a logger which is backed by Zap.
func NewZapLogger(
	environmentType configuration.EnvironmentType,
	logLevel zap.AtomicLevel,
) (*zap.Logger, error) {
	var loggerConfig zap.Config
	if environmentType == configuration.EnvironmentTypeDev {
		loggerConfig = zap.NewDevelopmentConfig()
	} else {
		loggerConfig = zap.NewProductionConfig()
	}

	loggerConfig.Level = logLevel

	zapLogger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("Unable to create Zap logger: %w", err)
	}

	return zapLogger, nil
}
