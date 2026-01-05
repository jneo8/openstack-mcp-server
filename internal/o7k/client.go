package o7k

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/jneo8/openstack-mcp-server/internal/config"
	"github.com/rs/zerolog/log"
)

// Client represents an OpenStack client with authenticated connections
type Client struct {
	provider        *gophercloud.ProviderClient
	blockStorageV3  *gophercloud.ServiceClient
	config          *config.OpenStackConfig
}

// NewClient creates a new OpenStack client with authentication
func NewClient(cfg *config.OpenStackConfig) (*Client, error) {
	// Create authentication options
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: cfg.AuthURL,
		Username:         cfg.Username,
		Password:         cfg.Password,
		DomainName:       cfg.UserDomain,
		TenantName:       cfg.ProjectName,
		TenantID:         cfg.ProjectID,
	}

	// If project domain is specified, use scoped auth
	if cfg.ProjectDomain != "" {
		authOpts.DomainName = cfg.ProjectDomain
	}

	log.Debug().
		Str("auth_url", cfg.AuthURL).
		Str("username", cfg.Username).
		Str("project_name", cfg.ProjectName).
		Str("region", cfg.Region).
		Bool("verify_ssl", cfg.VerifySSL).
		Msg("Authenticating with OpenStack")

	// Configure HTTP client with TLS settings
	httpClient := http.Client{
		Timeout: cfg.Timeout,
	}

	// Disable SSL verification if configured (for development/self-signed certs)
	if !cfg.VerifySSL {
		log.Warn().Msg("SSL verification disabled - not recommended for production")
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	// Create provider client with custom HTTP client
	ctx := context.Background()
	provider, err := openstack.NewClient(authOpts.IdentityEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider client: %w", err)
	}

	// Set HTTP client and configure provider
	provider.HTTPClient = httpClient
	provider.MaxBackoffRetries = uint(cfg.MaxRetries)

	// Authenticate
	err = openstack.Authenticate(ctx, provider, authOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	log.Info().Msg("Successfully authenticated with OpenStack")

	client := &Client{
		provider: provider,
		config:   cfg,
	}

	// Initialize Block Storage (Cinder) v3 client
	if err := client.initBlockStorage(); err != nil {
		return nil, fmt.Errorf("failed to initialize block storage client: %w", err)
	}

	return client, nil
}

// initBlockStorage initializes the Block Storage (Cinder) v3 service client
func (c *Client) initBlockStorage() error {
	endpointOpts := gophercloud.EndpointOpts{
		Region:       c.config.Region,
		Availability: gophercloud.Availability(c.config.EndpointType),
	}

	client, err := openstack.NewBlockStorageV3(c.provider, endpointOpts)
	if err != nil {
		return fmt.Errorf("creating block storage v3 client: %w", err)
	}

	c.blockStorageV3 = client
	log.Debug().Msg("Initialized Block Storage v3 client")
	return nil
}

// Close closes the OpenStack client connections
func (c *Client) Close() error {
	// Gophercloud doesn't require explicit connection closing
	// but we can clean up resources if needed
	log.Debug().Msg("Closing OpenStack client")
	return nil
}
