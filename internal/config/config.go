package config

import (
	"time"
)

// Config represents the complete application configuration
type Config struct {
	OpenStack OpenStackConfig `mapstructure:"openstack"`
	MCP       MCPConfig       `mapstructure:"mcp"`
	Logging   LoggingConfig   `mapstructure:"logging"`
}

// OpenStackConfig contains OpenStack authentication and connection settings
type OpenStackConfig struct {
	// Authentication
	AuthURL  string `mapstructure:"auth_url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`

	// Project/Tenant scope
	ProjectName   string `mapstructure:"project_name"`
	ProjectID     string `mapstructure:"project_id"`
	ProjectDomain string `mapstructure:"project_domain_name"`

	// User domain
	UserDomain string `mapstructure:"user_domain_name"`

	// Region and endpoint settings
	Region       string `mapstructure:"region"`
	EndpointType string `mapstructure:"endpoint_type"` // public, internal, admin

	// Connection settings
	Timeout    time.Duration `mapstructure:"timeout"`
	MaxRetries int           `mapstructure:"max_retries"`
	VerifySSL  bool          `mapstructure:"verify_ssl"`
	CACertFile string        `mapstructure:"ca_cert_file"`
}

// MCPConfig contains MCP server settings
type MCPConfig struct {
	// Transport configuration
	Transport TransportConfig `mapstructure:"transport"`

	// Server behavior
	ServerName    string `mapstructure:"server_name"`
	ServerVersion string `mapstructure:"server_version"`

	// Feature flag
	ReadOnly bool `mapstructure:"read_only"`
}

// TransportConfig defines how MCP communicates
type TransportConfig struct {
	Type string `mapstructure:"type"` // "stdio" or "http-streaming"

	// HTTP streaming-specific settings
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`

	// Connection timeout
	Timeout time.Duration `mapstructure:"timeout"`
}

// LoggingConfig controls application logging
type LoggingConfig struct {
	Level string `mapstructure:"level"` // debug, info, warn, error
}
