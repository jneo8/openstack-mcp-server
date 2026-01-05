package app

import (
	"context"
	"fmt"

	"github.com/jneo8/openstack-mcp-server/internal/config"
	"github.com/jneo8/openstack-mcp-server/internal/mcp"
	"github.com/jneo8/openstack-mcp-server/internal/o7k"
	"github.com/rs/zerolog/log"
)

// App represents the application with all dependencies
type App struct {
	config    *config.Config
	osClient  *o7k.Client
	mcpServer *mcp.Server
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

	// Create MCP server
	mcpServer, err := mcp.NewServer(&cfg.MCP, osClient)
	if err != nil {
		return nil, fmt.Errorf("creating MCP server: %w", err)
	}

	log.Info().Msg("MCP server initialized successfully")

	return &App{
		config:    cfg,
		osClient:  osClient,
		mcpServer: mcpServer,
	}, nil
}

// Run starts the application (blocking)
func (a *App) Run(ctx context.Context) error {
	log.Info().Msg("Starting application")

	// Start MCP server (blocking)
	if err := a.mcpServer.Start(ctx); err != nil {
		return fmt.Errorf("starting MCP server: %w", err)
	}

	log.Info().Msg("Application context cancelled")
	return ctx.Err()
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down application")

	// Stop MCP server
	if a.mcpServer != nil {
		if err := a.mcpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Error stopping MCP server")
		}
	}

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
