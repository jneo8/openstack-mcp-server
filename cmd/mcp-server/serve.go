package main

import (
	"fmt"
	"time"

	"github.com/jneo8/openstack-mcp-server/internal/app"
	"github.com/jneo8/openstack-mcp-server/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server",
		Long: `Start the OpenStack MCP server and begin listening for MCP client connections.

The server will use the configured transport (stdio or http-streaming) to communicate
with MCP clients and provide access to OpenStack resources, tools, and prompts.`,
		RunE: runServe,
	}

	// Serve-specific flags
	cmd.Flags().String("transport", "stdio", "transport type (stdio, http-streaming)")
	cmd.Flags().Int("port", 8080, "port for http-streaming transport")
	cmd.Flags().String("host", "localhost", "host for http-streaming transport")
	cmd.Flags().Duration("transport-timeout", 30*time.Second, "transport timeout")

	// OpenStack auth flags (can override config file)
	cmd.Flags().String("os-auth-url", "", "OpenStack authentication URL")
	cmd.Flags().String("os-username", "", "OpenStack username")
	cmd.Flags().String("os-password", "", "OpenStack password")
	cmd.Flags().String("os-project-name", "", "OpenStack project name")
	cmd.Flags().String("os-project-id", "", "OpenStack project ID")
	cmd.Flags().String("os-region", "", "OpenStack region")
	cmd.Flags().String("os-endpoint-type", "public", "OpenStack endpoint type (public, internal, admin)")
	cmd.Flags().String("os-user-domain", "Default", "OpenStack user domain name")
	cmd.Flags().String("os-project-domain", "Default", "OpenStack project domain name")
	cmd.Flags().Bool("os-verify-ssl", true, "verify SSL certificates")
	cmd.Flags().String("os-cacert", "", "path to CA certificate file")
	cmd.Flags().Duration("os-timeout", 30*time.Second, "OpenStack API timeout")
	cmd.Flags().Int("os-max-retries", 3, "maximum number of retries for OpenStack API calls")

	// MCP server flags
	cmd.Flags().Bool("read-only", false, "run in read-only mode (disable tools)")

	// Bind flags to viper
	viper.BindPFlag("mcp.transport.type", cmd.Flags().Lookup("transport"))
	viper.BindPFlag("mcp.transport.port", cmd.Flags().Lookup("port"))
	viper.BindPFlag("mcp.transport.host", cmd.Flags().Lookup("host"))
	viper.BindPFlag("mcp.transport.timeout", cmd.Flags().Lookup("transport-timeout"))

	viper.BindPFlag("openstack.auth_url", cmd.Flags().Lookup("os-auth-url"))
	viper.BindPFlag("openstack.username", cmd.Flags().Lookup("os-username"))
	viper.BindPFlag("openstack.password", cmd.Flags().Lookup("os-password"))
	viper.BindPFlag("openstack.project_name", cmd.Flags().Lookup("os-project-name"))
	viper.BindPFlag("openstack.project_id", cmd.Flags().Lookup("os-project-id"))
	viper.BindPFlag("openstack.region", cmd.Flags().Lookup("os-region"))
	viper.BindPFlag("openstack.endpoint_type", cmd.Flags().Lookup("os-endpoint-type"))
	viper.BindPFlag("openstack.user_domain_name", cmd.Flags().Lookup("os-user-domain"))
	viper.BindPFlag("openstack.project_domain_name", cmd.Flags().Lookup("os-project-domain"))
	viper.BindPFlag("openstack.verify_ssl", cmd.Flags().Lookup("os-verify-ssl"))
	viper.BindPFlag("openstack.ca_cert_file", cmd.Flags().Lookup("os-cacert"))
	viper.BindPFlag("openstack.timeout", cmd.Flags().Lookup("os-timeout"))
	viper.BindPFlag("openstack.max_retries", cmd.Flags().Lookup("os-max-retries"))

	viper.BindPFlag("mcp.read_only", cmd.Flags().Lookup("read-only"))

	return cmd
}

func runServe(cmd *cobra.Command, args []string) error {
	// Setup context with signal handling
	ctx := setupContext()

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// Setup logger based on configuration
	config.SetupLogger(cfg.Logging.Level)

	log.Info().Msg("Starting OpenStack MCP Server...")

	// Create application instance with manual DI
	application, err := app.NewApp(cfg)
	if err != nil {
		return fmt.Errorf("creating application: %w", err)
	}

	// Run the application (blocking)
	if err := application.Run(ctx); err != nil {
		return fmt.Errorf("running application: %w", err)
	}

	// Graceful shutdown
	if err := application.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down application: %w", err)
	}

	log.Info().Msg("Server stopped gracefully")
	return nil
}
