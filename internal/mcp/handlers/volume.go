package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jneo8/openstack-mcp-server/internal/o7k"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

// VolumeHandler handles volume-related MCP tool execution requests and delegates to OpenStack client
type VolumeHandler struct {
	osClient *o7k.Client
}

// NewVolumeHandler creates a new volume handler
func NewVolumeHandler(osClient *o7k.Client) *VolumeHandler {
	return &VolumeHandler{
		osClient: osClient,
	}
}

// HandleListVolumes handles the volumes_list tool
func (h *VolumeHandler) HandleListVolumes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debug().Msg("Executing volumes_list tool")

	volumes, err := h.osClient.ListVolumes(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list volumes")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list volumes: %v", err)), nil
	}

	// Convert to JSON
	data, err := json.Marshal(volumes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal volumes")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal volumes: %v", err)), nil
	}

	log.Debug().
		Int("count", len(volumes)).
		Msg("Volumes listed successfully")

	return mcp.NewToolResultText(string(data)), nil
}

// HandleGetVolume handles the volume_get tool
func (h *VolumeHandler) HandleGetVolume(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debug().Msg("Executing volume_get tool")

	// Extract volume_id parameter using SDK helper
	volumeID := request.GetString("volume_id", "")
	if volumeID == "" {
		return mcp.NewToolResultError("Missing or invalid 'volume_id' parameter"), nil
	}

	log.Debug().Str("volume_id", volumeID).Msg("Getting volume")

	volume, err := h.osClient.GetVolume(ctx, volumeID)
	if err != nil {
		log.Error().
			Err(err).
			Str("volume_id", volumeID).
			Msg("Failed to get volume")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get volume: %v", err)), nil
	}

	// Convert to JSON
	data, err := json.Marshal(volume)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal volume")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal volume: %v", err)), nil
	}

	log.Debug().
		Str("volume_id", volumeID).
		Str("volume_name", volume.Name).
		Msg("Volume retrieved successfully")

	return mcp.NewToolResultText(string(data)), nil
}

// VolumeCreateArgs defines the arguments for creating a volume
type VolumeCreateArgs struct {
	Name        string `json:"name"`
	Size        int    `json:"size"`
	Description string `json:"description,omitempty"`
	VolumeType  string `json:"volume_type,omitempty"`
}

// HandleCreateVolume handles the volume_create tool
func (h *VolumeHandler) HandleCreateVolume(ctx context.Context, request mcp.CallToolRequest, args VolumeCreateArgs) (*mcp.CallToolResult, error) {
	log.Debug().Msg("Executing volume_create tool")

	// Validate required parameters
	if args.Name == "" {
		return mcp.NewToolResultError("Missing or invalid 'name' parameter"), nil
	}
	if args.Size <= 0 {
		return mcp.NewToolResultError("Size must be a positive number"), nil
	}

	log.Debug().
		Str("name", args.Name).
		Int("size", args.Size).
		Str("description", args.Description).
		Str("volume_type", args.VolumeType).
		Msg("Creating volume")

	// Create volume
	opts := o7k.CreateVolumeOpts{
		Name:        args.Name,
		Size:        args.Size,
		Description: args.Description,
		VolumeType:  args.VolumeType,
	}

	volume, err := h.osClient.CreateVolume(ctx, opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", args.Name).
			Msg("Failed to create volume")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create volume: %v", err)), nil
	}

	// Convert to JSON
	data, err := json.Marshal(volume)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal volume")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal volume: %v", err)), nil
	}

	log.Info().
		Str("volume_id", volume.ID).
		Str("volume_name", volume.Name).
		Int("size", volume.Size).
		Msg("Volume created successfully")

	return mcp.NewToolResultText(string(data)), nil
}

// VolumeUpdateArgs defines the arguments for updating a volume
type VolumeUpdateArgs struct {
	VolumeID    string  `json:"volume_id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// HandleUpdateVolume handles the volume_update tool
func (h *VolumeHandler) HandleUpdateVolume(ctx context.Context, request mcp.CallToolRequest, args VolumeUpdateArgs) (*mcp.CallToolResult, error) {
	log.Debug().Msg("Executing volume_update tool")

	// Validate volume_id
	if args.VolumeID == "" {
		return mcp.NewToolResultError("Missing or invalid 'volume_id' parameter"), nil
	}

	// Check if at least one field is provided
	if args.Name == nil && args.Description == nil {
		return mcp.NewToolResultError("At least one of 'name' or 'description' must be provided"), nil
	}

	log.Debug().
		Str("volume_id", args.VolumeID).
		Msg("Updating volume")

	opts := o7k.UpdateVolumeOpts{
		Name:        args.Name,
		Description: args.Description,
	}

	volume, err := h.osClient.UpdateVolume(ctx, args.VolumeID, opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("volume_id", args.VolumeID).
			Msg("Failed to update volume")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update volume: %v", err)), nil
	}

	// Convert to JSON
	data, err := json.Marshal(volume)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal volume")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal volume: %v", err)), nil
	}

	log.Info().
		Str("volume_id", args.VolumeID).
		Str("volume_name", volume.Name).
		Msg("Volume updated successfully")

	return mcp.NewToolResultText(string(data)), nil
}

// HandleDeleteVolume handles the volume_delete tool
func (h *VolumeHandler) HandleDeleteVolume(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debug().Msg("Executing volume_delete tool")

	// Extract volume_id parameter using SDK helper
	volumeID := request.GetString("volume_id", "")
	if volumeID == "" {
		return mcp.NewToolResultError("Missing or invalid 'volume_id' parameter"), nil
	}

	log.Debug().Str("volume_id", volumeID).Msg("Deleting volume")

	err := h.osClient.DeleteVolume(ctx, volumeID)
	if err != nil {
		log.Error().
			Err(err).
			Str("volume_id", volumeID).
			Msg("Failed to delete volume")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete volume: %v", err)), nil
	}

	log.Info().
		Str("volume_id", volumeID).
		Msg("Volume deleted successfully")

	result := map[string]interface{}{
		"success":   true,
		"volume_id": volumeID,
		"message":   "Volume deleted successfully",
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal result")
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// RegisterTools registers all volume-related tools with the MCP server
func (h *VolumeHandler) RegisterTools(mcpServer *server.MCPServer, readOnly bool) error {
	log.Debug().
		Bool("read_only", readOnly).
		Msg("Registering volume tools")

	tools := h.getToolDefinitions()

	registeredCount := 0
	skippedCount := 0

	for _, toolDef := range tools {
		// Skip write tools if in read-only mode
		if readOnly && !toolDef.ReadOnly {
			log.Debug().
				Str("tool", toolDef.Name).
				Msg("Skipping tool (read-only mode enabled)")
			skippedCount++
			continue
		}

		// Build and register the tool
		tool := toolDef.BuildTool()
		mcpServer.AddTool(tool, toolDef.Handler)

		log.Debug().
			Str("tool", toolDef.Name).
			Bool("read_only", toolDef.ReadOnly).
			Msg("Tool registered")
		registeredCount++
	}

	log.Info().
		Int("registered", registeredCount).
		Int("skipped", skippedCount).
		Msg("Volume tools registration complete")

	return nil
}

// getToolDefinitions returns all volume tool definitions
func (h *VolumeHandler) getToolDefinitions() []ToolDefinition {
	return []ToolDefinition{
		{
			Name:        "volumes_list",
			Description: "List all volumes in the current OpenStack project",
			ReadOnly:    true,
			BuildTool: func() mcp.Tool {
				return mcp.NewTool("volumes_list",
					mcp.WithDescription("List all volumes in the current OpenStack project. Returns an array of volume objects with details like ID, name, size, status, and creation time."),
				)
			},
			Handler: h.HandleListVolumes,
		},
		{
			Name:        "volume_get",
			Description: "Get details of a specific volume by ID",
			ReadOnly:    true,
			BuildTool: func() mcp.Tool {
				return mcp.NewTool("volume_get",
					mcp.WithDescription("Get detailed information about a specific volume by its ID. Returns volume metadata including name, size, status, type, and timestamps."),
					mcp.WithString("volume_id",
						mcp.Required(),
						mcp.Description("The UUID of the volume to retrieve"),
					),
				)
			},
			Handler: h.HandleGetVolume,
		},
		{
			Name:        "volume_create",
			Description: "Create a new volume in OpenStack",
			ReadOnly:    false,
			BuildTool: func() mcp.Tool {
				return mcp.NewTool("volume_create",
					mcp.WithDescription("Create a new block storage volume in OpenStack. The volume will be created in the 'creating' state and transition to 'available' when ready."),
					mcp.WithString("name",
						mcp.Required(),
						mcp.Description("Name of the volume"),
					),
					mcp.WithNumber("size",
						mcp.Required(),
						mcp.Description("Size of the volume in gigabytes (GB). Must be a positive integer."),
					),
					mcp.WithString("description",
						mcp.Description("Optional description of the volume"),
					),
					mcp.WithString("volume_type",
						mcp.Description("Optional volume type (e.g., 'lvm', 'ssd', 'hdd'). Defaults to the configured default volume type."),
					),
				)
			},
			Handler: mcp.NewTypedToolHandler(h.HandleCreateVolume),
		},
		{
			Name:        "volume_update",
			Description: "Update a volume's metadata",
			ReadOnly:    false,
			BuildTool: func() mcp.Tool {
				return mcp.NewTool("volume_update",
					mcp.WithDescription("Update a volume's metadata such as name and description. Note: Cannot change volume size or type after creation."),
					mcp.WithString("volume_id",
						mcp.Required(),
						mcp.Description("The UUID of the volume to update"),
					),
					mcp.WithString("name",
						mcp.Description("New name for the volume (optional)"),
					),
					mcp.WithString("description",
						mcp.Description("New description for the volume (optional)"),
					),
				)
			},
			Handler: mcp.NewTypedToolHandler(h.HandleUpdateVolume),
		},
		{
			Name:        "volume_delete",
			Description: "Delete a volume from OpenStack",
			ReadOnly:    false,
			BuildTool: func() mcp.Tool {
				return mcp.NewTool("volume_delete",
					mcp.WithDescription("Delete a volume from OpenStack. The volume must be in 'available' or 'error' state and not attached to any instance. This operation cannot be undone."),
					mcp.WithString("volume_id",
						mcp.Required(),
						mcp.Description("The UUID of the volume to delete"),
					),
				)
			},
			Handler: h.HandleDeleteVolume,
		},
	}
}
