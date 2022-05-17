package manager

import (
	"context"
	"fmt"
	"net/http"

	"github.com/durandj/ley/internal/common/logging"
	"github.com/durandj/ley/internal/manager/configuration"
	"go.uber.org/zap"
)

// Server runs the management part of Ley.
type Server struct {
	logger     *zap.Logger
	httpServer http.Server
}

// New creates a service instance from the configuration.
func New(config *configuration.Configuration) (*Server, error) {
	fmt.Println(config)

	logger, err := logging.NewZapLogger(
		config.Service.EnvironmentType,
		config.Logging.Level.AsAtomicLevel(),
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to setup logger: %w", err)
	}

	return &Server{
		logger: logger,
		httpServer: http.Server{
			Addr:    config.Service.Address(),
			Handler: nil,
		},
	}, nil
}

// Run starts the service.
func (server *Server) Run(ctx context.Context) error {
	server.logger.Info(fmt.Sprintf("Starting HTTP server '%s'", server.httpServer.Addr))

	errChannel := make(chan error, 1)

	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil {
			errChannel <- fmt.Errorf("Server stopped: %w", err)
		}

		errChannel <- nil
	}()

	select {
	case err := <-errChannel:
		return err

	case <-ctx.Done():
		return fmt.Errorf("Server stopped: %w", ctx.Err())
	}
}
