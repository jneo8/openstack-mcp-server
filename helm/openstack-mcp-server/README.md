# OpenStack MCP Server Helm Chart

A Helm chart for deploying the OpenStack MCP Server on Kubernetes.

## Overview

This chart deploys the OpenStack MCP Server, which provides a Model Context Protocol interface for managing OpenStack volume resources (Cinder).

## Prerequisites

- Kubernetes 1.23+
- Helm 3.0+
- OpenStack credentials with appropriate permissions for volume management

## Installation

### Quick Start (Development)

For development/testing with inline secrets (NOT recommended for production):

```bash
helm install my-mcp-server ./helm/openstack-mcp-server \
  --set openstack.authUrl="https://your-openstack.example.com:5000/v3" \
  --set openstack.username="your-username" \
  --set openstack.projectName="your-project" \
  --set openstack.regionName="RegionOne" \
  --set secrets.password="your-password"
```

### Production Installation (Recommended)

1. Create a secret with your OpenStack password:

```bash
kubectl create secret generic openstack-credentials \
  --from-literal=password='your-openstack-password'
```

2. Create a values file (`my-values.yaml`):

```yaml
openstack:
  authUrl: "https://your-openstack.example.com:5000/v3"
  username: "your-username"
  projectName: "your-project"
  regionName: "RegionOne"
  userDomainName: "Default"
  projectDomainName: "Default"
  verifySSL: "true"

secrets:
  useExistingSecret: true
  existingSecretName: "openstack-credentials"

resources:
  limits:
    cpu: 500m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

3. Install the chart:

```bash
helm install my-mcp-server ./helm/openstack-mcp-server \
  -f my-values.yaml
```

### For Self-Signed Certificates

If your OpenStack installation uses self-signed certificates:

```bash
helm install my-mcp-server ./helm/openstack-mcp-server \
  --set openstack.verifySSL="false" \
  # ... other settings
```

## Configuration

### Key Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Docker image repository | `ghcr.io/jneo8/openstack-mcp-server` |
| `image.tag` | Docker image tag | `Chart.appVersion` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `replicaCount` | Number of replicas | `1` |
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `8080` |
| `openstack.authUrl` | OpenStack auth URL | `""` |
| `openstack.username` | OpenStack username | `""` |
| `openstack.projectName` | OpenStack project name | `""` |
| `openstack.regionName` | OpenStack region | `RegionOne` |
| `openstack.userDomainName` | User domain name | `Default` |
| `openstack.projectDomainName` | Project domain name | `Default` |
| `openstack.verifySSL` | Verify SSL certificates | `"true"` |
| `secrets.useExistingSecret` | Use existing secret | `false` |
| `secrets.existingSecretName` | Name of existing secret | `""` |
| `secrets.password` | OpenStack password (inline) | `""` |
| `logLevel` | Logging level | `info` |
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `256Mi` |
| `resources.requests.cpu` | CPU request | `100m` |
| `resources.requests.memory` | Memory request | `128Mi` |

## Accessing the MCP Server

The MCP server is deployed as a ClusterIP service (internal only). To access it:

### Port Forward

```bash
kubectl port-forward svc/my-mcp-server-openstack-mcp-server 8080:8080
```

Then connect to `http://localhost:8080`

### From Within the Cluster

Other pods can access the server at:

```
http://my-mcp-server-openstack-mcp-server.default.svc.cluster.local:8080
```

## Upgrading

```bash
helm upgrade my-mcp-server ./helm/openstack-mcp-server \
  -f my-values.yaml
```

## Uninstalling

```bash
helm uninstall my-mcp-server
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -l app.kubernetes.io/name=openstack-mcp-server
```

### View Logs

```bash
kubectl logs -l app.kubernetes.io/name=openstack-mcp-server -f
```

### Common Issues

1. **Connection refused to OpenStack**: Verify `openstack.authUrl` is correct and accessible from the cluster
2. **TLS certificate errors**: Set `openstack.verifySSL="false"` for self-signed certificates
3. **Authentication failed**: Verify credentials in the secret/values
4. **Pod not starting**: Check logs for detailed error messages

## Security Notes

- The server runs as a non-root user (uid=65532)
- Root filesystem is read-only
- All Linux capabilities are dropped
- For production, always use `secrets.useExistingSecret` instead of inline passwords
- Consider using a secret management solution like External Secrets Operator for sensitive data

## Current Features

This MCP server currently supports:
- OpenStack Volume (Cinder) CRUD operations
  - List volumes
  - Get volume details
  - Create volumes
  - Update volume metadata
  - Delete volumes

## License

See the main project repository for license information.
