# OpenStack MCP Server - Architecture

## Overview

The OpenStack MCP Server is a Model Context Protocol (MCP) server that provides integration with OpenStack infrastructure. It enables AI assistants and other MCP clients to interact with OpenStack resources through a standardized interface.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        MCP Client                            │
│                   (Claude Desktop, etc.)                     │
└──────────────────────────┬──────────────────────────────────┘
                           │ MCP Protocol (stdio/SSE)
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                    MCP Server Layer                          │
│              Transport • Protocol • Handlers                 │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                   Application Layer                          │
│         Lifecycle • Orchestration • DI Container            │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                  OpenStack Client Layer                      │
│         Auth • Service Clients • Resource Abstraction       │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                    OpenStack Cloud                           │
│              (Nova, Cinder, Neutron, Glance, etc.)          │
└─────────────────────────────────────────────────────────────┘
```

## Tech Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Language** | Go 1.24+ | System programming, performance, concurrency |
| **CLI Framework** | Cobra | Command-line interface structure |
| **Configuration** | Viper | Multi-source configuration management |
| **Dependency Injection** | Manual DI | Explicit dependency wiring |
| **MCP SDK** | MCP Go SDK | MCP protocol implementation |
| **OpenStack Client** | Gophercloud | OpenStack API client library |
| **Logging** | Zap/Logrus | Structured logging |

## Package Architecture

### `cmd/` - Command Line Interface

**Purpose**: Application entry point and CLI commands

**Responsibilities**:
- Define CLI commands and flags using Cobra
- Bootstrap application with configuration
- Handle graceful shutdown signals

---

### `internal/config/` - Configuration Management

**Purpose**: Centralized configuration using Cobra and Viper ecosystem

**Responsibilities**:
- Load configuration from multiple sources (files, environment variables, flags)
- Provide type-safe configuration access
- Validate configuration values
- Integrate with Cobra command-line flags

---

### `internal/app/` - Application Lifecycle

**Purpose**: Application lifecycle management and process orchestration

**Responsibilities**:
- Initialize all components in correct order
- Manage service dependencies through manual dependency injection
- Handle graceful startup and shutdown
- Coordinate between different layers (MCP, OpenStack)
- Provide application-level error handling

---

### `internal/o7k/` - OpenStack Integration

**Purpose**: High-level wrapper for connecting to OpenStack

**Responsibilities**:
- Authenticate with OpenStack using Gophercloud
- Provide abstracted interfaces for OpenStack services:
  - **Compute** (Nova): VM instances, flavors, keypairs
  - **Volume** (Cinder): Block storage volumes, snapshots
  - **Network** (Neutron): Networks, subnets, security groups, floating IPs
  - **Image** (Glance): Operating system images
  - **Identity** (Keystone): Authentication and projects
- Handle API retries and error handling
- Convert between OpenStack models and internal representations

---

### `internal/mcp/` - MCP Server Implementation

**Purpose**: Building the MCP server and protocol implementation

**Responsibilities**:
- Implement MCP JSON-RPC 2.0 protocol
- Handle transport layer (stdio/SSE)
- Register and expose MCP primitives:
  - **Resources**: Read-only OpenStack resource data (servers, volumes, networks)
  - **Tools**: Actions to create, modify, delete OpenStack resources
  - **Prompts**: Templates for common OpenStack workflows
- Route incoming requests to appropriate handlers
- Convert between MCP protocol types and internal models

---

## References

- [Model Context Protocol](https://spec.modelcontextprotocol.io/)
- [OpenStack API](https://docs.openstack.org/api-ref/)
- [Gophercloud](https://github.com/gophercloud/gophercloud)
- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
