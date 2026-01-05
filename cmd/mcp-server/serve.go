package main

import (
	"fmt"
	"os"
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
	flagBindings := map[string]string{
		"mcp.transport.type":            "transport",
		"mcp.transport.port":            "port",
		"mcp.transport.host":            "host",
		"mcp.transport.timeout":         "transport-timeout",
		"openstack.auth_url":            "os-auth-url",
		"openstack.username":            "os-username",
		"openstack.password":            "os-password",
		"openstack.project_name":        "os-project-name",
		"openstack.project_id":          "os-project-id",
		"openstack.region":              "os-region",
		"openstack.endpoint_type":       "os-endpoint-type",
		"openstack.user_domain_name":    "os-user-domain",
		"openstack.project_domain_name": "os-project-domain",
		"openstack.verify_ssl":          "os-verify-ssl",
		"openstack.ca_cert_file":        "os-cacert",
		"openstack.timeout":             "os-timeout",
		"openstack.max_retries":         "os-max-retries",
		"mcp.read_only":                 "read-only",
	}
	for key, flag := range flagBindings {
		if err := viper.BindPFlag(key, cmd.Flags().Lookup(flag)); err != nil {
			fmt.Fprintf(os.Stderr, "Error binding flag %s: %v\n", flag, err)
			os.Exit(1)
		}
	}

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
