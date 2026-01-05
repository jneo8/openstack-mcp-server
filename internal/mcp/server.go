package mcp

import (
	"context"
	"fmt"

	"github.com/jneo8/openstack-mcp-server/internal/config"
	"github.com/jneo8/openstack-mcp-server/internal/mcp/handlers"
	"github.com/jneo8/openstack-mcp-server/internal/o7k"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

// Server represents the MCP server with all its dependencies
type Server struct {
	config     *config.MCPConfig
	osClient   *o7k.Client
	mcpServer  *server.MCPServer
	handlers   []handlers.Handler
	httpServer *server.StreamableHTTPServer
}

// NewServer creates a new MCP server instance
func NewServer(cfg *config.MCPConfig, osClient *o7k.Client) (*Server, error) {
	log.Info().
		Str("server_name", cfg.ServerName).
		Str("server_version", cfg.ServerVersion).
		Str("transport", cfg.Transport.Type).
		Bool("read_only", cfg.ReadOnly).
		Msg("Creating MCP server")

	// Create MCP server with tool capabilities
	mcpServer := server.NewMCPServer(
		cfg.ServerName,
		cfg.ServerVersion,
		server.WithToolCapabilities(true),
	)

	// Create handlers
	volumeHandler := handlers.NewVolumeHandler(osClient)
	handlerList := []handlers.Handler{
		volumeHandler,
		// Add more handlers here (NetworkHandler, ComputeHandler, etc.)
	}

	// Create server instance
	s := &Server{
		config:     cfg,
		osClient:   osClient,
		mcpServer:  mcpServer,
		handlers:   handlerList,
	}

	// Register tools from all handlers
	for _, handler := range handlerList {
		if err := handler.RegisterTools(mcpServer, cfg.ReadOnly); err != nil {
			return nil, fmt.Errorf("registering tools: %w", err)
		}
	}

	// For HTTP transport, create the HTTP server
	if cfg.Transport.Type == "http-streaming" {
		s.httpServer = server.NewStreamableHTTPServer(mcpServer)
	}

	log.Info().
		Int("tools_count", len(mcpServer.ListTools())).
		Msg("MCP server created successfully")

	return s, nil
}

// Start starts the MCP server with the configured transport
func (s *Server) Start(ctx context.Context) error {
	log.Info().
		Str("transport", s.config.Transport.Type).
		Msg("Starting MCP server")

	switch s.config.Transport.Type {
	case "stdio":
		return s.startStdio(ctx)
	case "http-streaming":
		return s.startHTTP(ctx)
	default:
		return fmt.Errorf("unsupported transport type: %s", s.config.Transport.Type)
	}
}

// Shutdown gracefully shuts down the MCP server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down MCP server")

	// For HTTP transport, explicitly shut down the HTTP server
	if s.httpServer != nil {
		log.Info().Msg("Shutting down HTTP server")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Error shutting down HTTP server")
			return err
		}
		log.Info().Msg("HTTP server shutdown complete")
	}

	log.Info().Msg("MCP server shutdown complete")
	return nil
}

// startStdio starts the MCP server using stdio transport
func (s *Server) startStdio(ctx context.Context) error {
	log.Info().Msg("Starting stdio transport")

	// Use mark3labs/mcp-go's built-in stdio server
	// This is a blocking call that reads from stdin and writes to stdout
	if err := server.ServeStdio(s.mcpServer); err != nil {
		log.Error().Err(err).Msg("Stdio server error")
		return err
	}

	return nil
}

// startHTTP starts the MCP server using HTTP streamable transport
func (s *Server) startHTTP(ctx context.Context) error {
	log.Info().
		Str("host", s.config.Transport.Host).
		Int("port", s.config.Transport.Port).
		Msg("Starting HTTP streamable transport")

	// Build address
	addr := fmt.Sprintf("%s:%d", s.config.Transport.Host, s.config.Transport.Port)

	// Start HTTP server in goroutine
	errChan := make(chan error, 1)
	go func() {
		log.Info().Str("address", addr).Msg("HTTP server listening")
		if err := s.httpServer.Start(addr); err != nil {
			log.Error().Err(err).Msg("HTTP server error")
			errChan <- err
		}
	}()

	// Wait for either context cancellation or server error
	select {
	case <-ctx.Done():
		log.Info().Msg("Context cancelled, shutting down HTTP server")
		// The SDK's server will handle shutdown internally when context is done
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
