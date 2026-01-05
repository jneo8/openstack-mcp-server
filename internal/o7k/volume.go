// Package o7k (OpenStack) provides a high-level wrapper connect to openstack.
package o7k

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/rs/zerolog/log"
)

// Volume represents an OpenStack volume with common attributes
type Volume struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Size        int               `json:"size"`         // Size in GB
	Status      string            `json:"status"`       // creating, available, in-use, etc.
	VolumeType  string            `json:"volume_type"`
	Bootable    bool              `json:"bootable"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

// CreateVolumeOpts contains options for creating a volume
type CreateVolumeOpts struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Size        int               `json:"size"`         // Size in GB (required)
	VolumeType  string            `json:"volume_type,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// UpdateVolumeOpts contains options for updating a volume
type UpdateVolumeOpts struct {
	Name        *string            `json:"name,omitempty"`
	Description *string            `json:"description,omitempty"`
	Metadata    *map[string]string `json:"metadata,omitempty"`
}

// CreateVolume creates a new volume in OpenStack
func (c *Client) CreateVolume(ctx context.Context, opts CreateVolumeOpts) (*Volume, error) {
	if c.blockStorageV3 == nil {
		return nil, fmt.Errorf("block storage client not initialized")
	}

	log.Info().
		Str("name", opts.Name).
		Int("size", opts.Size).
		Str("volume_type", opts.VolumeType).
		Msg("Creating volume")

	createOpts := volumes.CreateOpts{
		Size:        opts.Size,
		Name:        opts.Name,
		Description: opts.Description,
		VolumeType:  opts.VolumeType,
		Metadata:    opts.Metadata,
	}

	vol, err := volumes.Create(ctx, c.blockStorageV3, createOpts, nil).Extract()
	if err != nil {
		return nil, fmt.Errorf("creating volume: %w", err)
	}

	log.Info().
		Str("id", vol.ID).
		Str("name", vol.Name).
		Msg("Volume created successfully")

	return convertVolume(vol), nil
}

// GetVolume retrieves a volume by ID
func (c *Client) GetVolume(ctx context.Context, volumeID string) (*Volume, error) {
	if c.blockStorageV3 == nil {
		return nil, fmt.Errorf("block storage client not initialized")
	}

	log.Debug().
		Str("volume_id", volumeID).
		Msg("Getting volume")

	vol, err := volumes.Get(ctx, c.blockStorageV3, volumeID).Extract()
	if err != nil {
		return nil, fmt.Errorf("getting volume %s: %w", volumeID, err)
	}

	return convertVolume(vol), nil
}

// ListVolumes lists all volumes accessible to the current project
func (c *Client) ListVolumes(ctx context.Context) ([]Volume, error) {
	if c.blockStorageV3 == nil {
		return nil, fmt.Errorf("block storage client not initialized")
	}

	log.Debug().Msg("Listing volumes")

	allPages, err := volumes.List(c.blockStorageV3, volumes.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing volumes: %w", err)
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return nil, fmt.Errorf("extracting volumes: %w", err)
	}

	result := make([]Volume, len(allVolumes))
	for i, vol := range allVolumes {
		result[i] = *convertVolume(&vol)
	}

	log.Debug().Int("count", len(result)).Msg("Listed volumes")
	return result, nil
}

// UpdateVolume updates a volume's metadata
func (c *Client) UpdateVolume(ctx context.Context, volumeID string, opts UpdateVolumeOpts) (*Volume, error) {
	if c.blockStorageV3 == nil {
		return nil, fmt.Errorf("block storage client not initialized")
	}

	log.Info().
		Str("volume_id", volumeID).
		Msg("Updating volume")

	updateOpts := volumes.UpdateOpts{}
	if opts.Name != nil {
		updateOpts.Name = opts.Name
	}
	if opts.Description != nil {
		updateOpts.Description = opts.Description
	}
	if opts.Metadata != nil {
		updateOpts.Metadata = *opts.Metadata
	}

	vol, err := volumes.Update(ctx, c.blockStorageV3, volumeID, updateOpts).Extract()
	if err != nil {
		return nil, fmt.Errorf("updating volume %s: %w", volumeID, err)
	}

	log.Info().
		Str("volume_id", volumeID).
		Msg("Volume updated successfully")

	return convertVolume(vol), nil
}

// DeleteVolume deletes a volume by ID
func (c *Client) DeleteVolume(ctx context.Context, volumeID string) error {
	if c.blockStorageV3 == nil {
		return fmt.Errorf("block storage client not initialized")
	}

	log.Info().
		Str("volume_id", volumeID).
		Msg("Deleting volume")

	err := volumes.Delete(ctx, c.blockStorageV3, volumeID, volumes.DeleteOpts{}).ExtractErr()
	if err != nil {
		return fmt.Errorf("deleting volume %s: %w", volumeID, err)
	}

	log.Info().
		Str("volume_id", volumeID).
		Msg("Volume deleted successfully")

	return nil
}

// convertVolume converts a Gophercloud volume to our Volume type
func convertVolume(vol *volumes.Volume) *Volume {
	return &Volume{
		ID:          vol.ID,
		Name:        vol.Name,
		Description: vol.Description,
		Size:        vol.Size,
		Status:      vol.Status,
		VolumeType:  vol.VolumeType,
		Bootable:    vol.Bootable == "true",
		Metadata:    vol.Metadata,
		CreatedAt:   vol.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   vol.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
