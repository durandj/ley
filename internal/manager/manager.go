package manager

import (
	"context"
	"database/sql"
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
	db         *sql.DB
}

// New creates a service instance from the configuration.
func New(config *configuration.Configuration) (*Server, error) {
	logger, err := logging.NewZapLogger(
		config.Service.EnvironmentType,
		config.Logging.Level.AsAtomicLevel(),
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to setup logger: %w", err)
	}

	dbConnectionString, err := config.DB.ConnectionString()
	if err != nil {
		return nil, fmt.Errorf("Unable to create database connection string: %w", err)
	}

	db, err := sql.Open(string(config.DB.Type), dbConnectionString)
	if err != nil {
		return nil, fmt.Errorf("Invalid database connection string: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	return &Server{
		logger: logger,
		httpServer: http.Server{
			Addr:    config.Service.Address(),
			Handler: NewController(db),
		},
		db: db,
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

// CleanUp is called when the server needs resources freed to terminate
// cleanly.
func (server *Server) CleanUp() {
	_ = server.db.Close()
}
