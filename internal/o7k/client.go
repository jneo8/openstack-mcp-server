package o7k

import (
	"context"
	"fmt"

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
		Msg("Authenticating with OpenStack")

	// Authenticate and get provider client
	ctx := context.Background()
	provider, err := openstack.AuthenticatedClient(ctx, authOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set HTTP client options
	provider.HTTPClient.Timeout = cfg.Timeout
	provider.MaxBackoffRetries = uint(cfg.MaxRetries)

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
