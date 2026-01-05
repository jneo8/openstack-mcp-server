package app

import (
	"context"
	"fmt"

	"github.com/jneo8/openstack-mcp-server/internal/config"
	"github.com/jneo8/openstack-mcp-server/internal/o7k"
	"github.com/rs/zerolog/log"
)

// App represents the application with all dependencies
type App struct {
	config   *config.Config
	osClient *o7k.Client
	// TODO: Add MCP server when implemented
}

// NewApp creates a new application instance with manual DI
func NewApp(cfg *config.Config) (*App, error) {
	log.Info().Msg("Initializing application")

	// Create OpenStack client
	osClient, err := o7k.NewClient(&cfg.OpenStack)
	if err != nil {
		return nil, fmt.Errorf("creating OpenStack client: %w", err)
	}

	log.Info().Msg("OpenStack client initialized successfully")

	// TODO: Create MCP server
	// TODO: Wire MCP server with OpenStack client

	return &App{
		config:   cfg,
		osClient: osClient,
	}, nil
}

// Run starts the application (blocking)
func (a *App) Run(ctx context.Context) error {
	log.Info().Msg("Starting application")

	// TODO: Start MCP server
	// For now, just demonstrate that OpenStack client is ready
	log.Info().
		Str("transport", a.config.MCP.Transport.Type).
		Bool("read_only", a.config.MCP.ReadOnly).
		Msg("MCP server configuration loaded (not yet started)")

	// Block until context is cancelled
	<-ctx.Done()
	log.Info().Msg("Application context cancelled")
	return ctx.Err()
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down application")

	// TODO: Stop MCP server gracefully

	// Close OpenStack client
	if a.osClient != nil {
		if err := a.osClient.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing OpenStack client")
			return fmt.Errorf("closing OpenStack client: %w", err)
		}
	}

	log.Info().Msg("Application shutdown complete")
	return nil
}
