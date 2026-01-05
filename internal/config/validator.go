package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// Validate performs comprehensive configuration validation
func (c *Config) Validate() error {
	var errors ValidationErrors

	// Validate OpenStack configuration
	errors = append(errors, c.validateOpenStack()...)

	// Validate MCP configuration
	errors = append(errors, c.validateMCP()...)

	// Validate Logging configuration
	errors = append(errors, c.validateLogging()...)

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// validateOpenStack validates OpenStack-specific configuration
func (c *Config) validateOpenStack() []ValidationError {
	var errors []ValidationError

	// Required fields
	if c.OpenStack.AuthURL == "" {
		errors = append(errors, ValidationError{
			Field:   "openstack.auth_url",
			Message: "authentication URL is required",
		})
	} else {
		// Validate URL format
		if _, err := url.Parse(c.OpenStack.AuthURL); err != nil {
			errors = append(errors, ValidationError{
				Field:   "openstack.auth_url",
				Message: fmt.Sprintf("invalid URL format: %v", err),
			})
		}
	}

	if c.OpenStack.Username == "" {
		errors = append(errors, ValidationError{
			Field:   "openstack.username",
			Message: "username is required",
		})
	}

	if c.OpenStack.Password == "" {
		errors = append(errors, ValidationError{
			Field:   "openstack.password",
			Message: "password is required",
		})
	}

	// Project scope validation (need at least one)
	if c.OpenStack.ProjectName == "" && c.OpenStack.ProjectID == "" {
		errors = append(errors, ValidationError{
			Field:   "openstack.project",
			Message: "either project_name or project_id is required",
		})
	}

	// Validate endpoint type
	validEndpoints := map[string]bool{"public": true, "internal": true, "admin": true}
	if c.OpenStack.EndpointType != "" && !validEndpoints[c.OpenStack.EndpointType] {
		errors = append(errors, ValidationError{
			Field:   "openstack.endpoint_type",
			Message: "must be one of: public, internal, admin",
		})
	}

	// Validate CA cert file if specified
	if c.OpenStack.CACertFile != "" {
		if _, err := os.Stat(c.OpenStack.CACertFile); os.IsNotExist(err) {
			errors = append(errors, ValidationError{
				Field:   "openstack.ca_cert_file",
				Message: fmt.Sprintf("file does not exist: %s", c.OpenStack.CACertFile),
			})
		}
	}

	// Validate timeout and retries
	if c.OpenStack.Timeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "openstack.timeout",
			Message: "must be greater than 0",
		})
	}

	if c.OpenStack.MaxRetries < 0 {
		errors = append(errors, ValidationError{
			Field:   "openstack.max_retries",
			Message: "must be non-negative",
		})
	}

	return errors
}

// validateMCP validates MCP-specific configuration
func (c *Config) validateMCP() []ValidationError {
	var errors []ValidationError

	// Validate transport type
	validTransports := map[string]bool{"stdio": true, "http": true}
	if !validTransports[c.MCP.Transport.Type] {
		errors = append(errors, ValidationError{
			Field:   "mcp.transport.type",
			Message: "must be either 'stdio' or 'http'",
		})
	}

	// HTTP-specific validation
	if c.MCP.Transport.Type == "http" {
		if c.MCP.Transport.Port <= 0 || c.MCP.Transport.Port > 65535 {
			errors = append(errors, ValidationError{
				Field:   "mcp.transport.port",
				Message: "must be between 1 and 65535",
			})
		}

		if c.MCP.Transport.Host == "" {
			errors = append(errors, ValidationError{
				Field:   "mcp.transport.host",
				Message: "host is required for http transport",
			})
		}
	}

	// Validate timeout
	if c.MCP.Transport.Timeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "mcp.transport.timeout",
			Message: "must be greater than 0",
		})
	}

	return errors
}

// validateLogging validates logging configuration
func (c *Config) validateLogging() []ValidationError {
	var errors []ValidationError

	// Validate log level
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[c.Logging.Level] {
		errors = append(errors, ValidationError{
			Field:   "logging.level",
			Message: "must be one of: debug, info, warn, error",
		})
	}

	return errors
}
