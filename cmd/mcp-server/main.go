package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jneo8/openstack-mcp-server/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Build-time version information (set via -ldflags)
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// newRootCommand creates the root cobra command
func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp-server",
		Short: "OpenStack MCP Server - Model Context Protocol server for OpenStack",
		Long: `OpenStack MCP Server provides a Model Context Protocol (MCP) interface
to OpenStack infrastructure, enabling AI assistants and other MCP clients
to interact with OpenStack resources through a standardized interface.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	// Global/Persistent flags (available to all subcommands)
	cmd.PersistentFlags().String("config", "", "config file path (default searches ~/.openstack-mcp-server.yaml, /etc/openstack-mcp-server/config.yaml)")
	cmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")

	// Bind persistent flags to viper
	if err := viper.BindPFlag("logging.level", cmd.PersistentFlags().Lookup("log-level")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding log-level flag: %v\n", err)
		os.Exit(1)
	}

	// Add subcommands
	cmd.AddCommand(newServeCommand())
	cmd.AddCommand(newVersionCommand())

	return cmd
}

// initConfig initializes viper configuration
func initConfig() error {
	// Set config file name and type
	viper.SetConfigName(".openstack-mcp-server")
	viper.SetConfigType("yaml")

	// Add config search paths
	viper.AddConfigPath(".") // Current directory
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)                            // Home directory
		viper.AddConfigPath(fmt.Sprintf("%s/.config", home)) // XDG config
	}
	viper.AddConfigPath("/etc/openstack-mcp-server") // System-wide config

	// Enable environment variable support
	viper.SetEnvPrefix("OSMCP") // OSMCP_LOG_LEVEL, etc.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// Manually bind OpenStack environment variables (OS_ prefix is standard)
	envBindings := map[string]string{
		"openstack.auth_url":            "OS_AUTH_URL",
		"openstack.username":            "OS_USERNAME",
		"openstack.password":            "OS_PASSWORD",
		"openstack.project_name":        "OS_PROJECT_NAME",
		"openstack.project_id":          "OS_PROJECT_ID",
		"openstack.project_domain_name": "OS_PROJECT_DOMAIN_NAME",
		"openstack.user_domain_name":    "OS_USER_DOMAIN_NAME",
		"openstack.region":              "OS_REGION_NAME",
		"openstack.endpoint_type":       "OS_ENDPOINT_TYPE",
		"openstack.ca_cert_file":        "OS_CACERT",
	}
	for key, env := range envBindings {
		if err := viper.BindEnv(key, env); err != nil {
			return fmt.Errorf("binding env %s: %w", key, err)
		}
	}

	// Set defaults using viper
	setDefaults()

	// Read config file if exists (not finding it is not an error)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("reading config file: %w", err)
		}
	}

	return nil
}

// setDefaults sets default values in viper
func setDefaults() {
	// OpenStack defaults
	viper.SetDefault("openstack.endpoint_type", "public")
	viper.SetDefault("openstack.timeout", 30*time.Second)
	viper.SetDefault("openstack.max_retries", 3)
	viper.SetDefault("openstack.verify_ssl", true)
	viper.SetDefault("openstack.user_domain_name", "Default")
	viper.SetDefault("openstack.project_domain_name", "Default")

	// MCP defaults
	viper.SetDefault("mcp.transport.type", "stdio")
	viper.SetDefault("mcp.transport.port", 8080)
	viper.SetDefault("mcp.transport.host", "localhost")
	viper.SetDefault("mcp.transport.timeout", 30*time.Second)
	viper.SetDefault("mcp.server_name", "openstack-mcp-server")
	viper.SetDefault("mcp.server_version", "0.1.0")
	viper.SetDefault("mcp.read_only", false)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
}

// loadConfig loads and validates configuration from all sources
func loadConfig() (*config.Config, error) {
	// Initialize viper configuration
	if err := initConfig(); err != nil {
		return nil, err
	}

	// Unmarshal into config struct
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}

// setupContext creates a context that cancels on SIGINT/SIGTERM
func setupContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	// Handle graceful shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		fmt.Fprintf(os.Stderr, "\nReceived signal %v, shutting down gracefully...\n", sig)
		cancel()
	}()

	return ctx
}
