# OpenStack MCP Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server for OpenStack, enabling AI assistants and MCP clients to interact with OpenStack infrastructure through a standardized interface.

## Overview

This MCP server provides programmatic access to OpenStack resources, allowing AI assistants like Claude to manage and query OpenStack infrastructure using natural language. Built with the [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) SDK, it supports both stdio and HTTP streaming transports.

## Features

Current OpenStack service support:

- [x] **Block Storage (Cinder)** - Volume management
  - List volumes
  - Get volume details
  - Create volumes
  - Update volume metadata
  - Delete volumes
- [ ] **Compute (Nova)** - Virtual machine management (coming soon)
- [ ] **Network (Neutron)** - Network management (coming soon)
- [ ] **Image (Glance)** - Image management (coming soon)
- [ ] **Identity (Keystone)** - User and project management (coming soon)
- [ ] **Object Storage (Swift)** - Object storage operations (coming soon)

## Prerequisites

- **Go 1.24.2+** - Required to build and run the server
- **OpenStack Environment** - Access to an OpenStack cloud with credentials
- **MCP Client** - Such as [MCP Inspector](https://github.com/modelcontextprotocol/inspector) for testing

## Quick Start

### 1. Set OpenStack Credentials

Export your OpenStack authentication credentials as environment variables:

```bash
export OS_AUTH_URL=https://your-openstack.example.com:5000/v3
export OS_USERNAME=your-username
export OS_PASSWORD=your-password
export OS_PROJECT_NAME=your-project
export OS_USER_DOMAIN_NAME=Default
export OS_PROJECT_DOMAIN_NAME=Default
export OS_REGION_NAME=RegionOne

# Optional: Disable SSL verification for development (not recommended for production)
export OSMCP_OS_VERIFY_SSL=false
```

### 2. Start the Server

```bash
# Stdio transport (default)
go run ./cmd/mcp-server serve

# HTTP transport
go run ./cmd/mcp-server serve --transport http --port 8080
```

## Testing with MCP Inspector

The [MCP Inspector](https://github.com/modelcontextprotocol/inspector) is a great tool for testing and debugging your MCP server.

## Available Tools

### Block Storage (Cinder)

| Tool | Description | Read-Only |
|------|-------------|-----------|
| `volumes_list` | List all volumes in the current project | Yes |
| `volume_get` | Get detailed information about a specific volume | Yes |
| `volume_create` | Create a new block storage volume | No |
| `volume_update` | Update volume metadata (name, description) | No |
| `volume_delete` | Delete a volume | No |


### Configuration File

Create `.openstack-mcp-server.yaml` in your home directory or project root:

```yaml
openstack:
  auth_url: https://your-openstack.example.com:5000/v3
  username: your-username
  password: your-password
  project_name: your-project
  user_domain_name: Default
  project_domain_name: Default
  region: RegionOne
  verify_ssl: true

mcp:
  server_name: openstack-mcp-server
  server_version: 0.1.0
  read_only: false
  transport:
    type: http
    host: localhost
    port: 8080

logging:
  level: info
```

### Environment Variables

All configuration can be set via environment variables using the `OSMCP_` prefix:

**OpenStack settings** support both standard OpenStack variables and OSMCP-prefixed variables:
- Standard: `OS_AUTH_URL`, `OS_USERNAME`, `OS_PASSWORD`, `OS_PROJECT_NAME`, `OS_REGION_NAME`, etc.
- OSMCP-prefixed: `OSMCP_OS_AUTH_URL`, `OSMCP_OS_USERNAME`, `OSMCP_OS_PASSWORD`, `OSMCP_OS_VERIFY_SSL`, etc.

**MCP settings** use the `OSMCP_` prefix:
- `OSMCP_READONLY=true`
- `OSMCP_TRANSPORT_TYPE=http`
- `OSMCP_TRANSPORT_HOST=localhost`
- `OSMCP_TRANSPORT_PORT=8080`
- `OSMCP_TRANSPORT_TIMEOUT=30s`

**Logging settings**:
- `OSMCP_LOGGING_LEVEL=debug`

## Read-Only Mode

Run the server in read-only mode to disable all write operations (create, update, delete):

```bash
go run ./cmd/mcp-server serve --read-only
```
