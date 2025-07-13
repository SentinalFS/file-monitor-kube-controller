# file-monitor-kube-controller

A Kubernetes controller that monitors file system changes and updates Custom Resource Definitions (CRDs) with file information such as inode, file name, size, and modification time.

## Features

- Queries Kubernetes CRDs using the k8s.io/client-go library
- Monitors file system changes in specified paths
- Updates CRD status with file information (inode, name, size, modification time)
- Supports both in-cluster and out-of-cluster configurations
- Namespace-aware file monitoring

## Prerequisites

- Go 1.24 or later
- Kubernetes cluster (for deployment)
- kubectl configured to access your cluster

## Installation

1. **Install the CRD:**
   ```bash
   kubectl apply -f filemonitor-crd.yaml
   ```

2. **Deploy the controller:**
   ```bash
   kubectl apply -f deployment.yaml
   ```

3. **Create a FileMonitor resource:**
   ```bash
   kubectl apply -f example-filemonitor.yaml
   ```

## Usage

### Running locally (development)

1. **Build the application:**
   ```bash
   go build -o file-monitor-controller
   ```

2. **Run the controller:**
   ```bash
   ./file-monitor-controller
   ```

### Running in Kubernetes

1. **Build and push Docker image:**
   ```bash
   docker build -t your-registry/file-monitor-controller:latest .
   docker push your-registry/file-monitor-controller:latest
   ```

2. **Update the deployment image:**
   ```bash
   # Edit deployment.yaml to use your image
   kubectl apply -f deployment.yaml
   ```

## CRD Structure

The FileMonitor CRD has the following structure:

```yaml
apiVersion: sentinalfs.io/v1
kind: FileMonitor
metadata:
  name: example-file-monitor
  namespace: default
spec:
  path: "/tmp"           # Path to monitor
  namespace: "default"   # Target namespace
status:
  files:                 # Updated by controller
  - name: "example.txt"
    inode: 12345
    size: 1024
    modTime: "2025-01-13T10:00:00Z"
    path: "/tmp/example.txt"
    isDir: false
```

## Controller Functionality

The controller performs the following operations in a continuous loop:

1. **Query CRDs:** Lists all FileMonitor resources across namespaces
2. **Namespace-specific queries:** Queries CRDs in specific namespaces
3. **File information updates:** Updates CRD status with current file information
4. **Error handling:** Gracefully handles missing CRDs and API errors

## Key Functions

- `queryCRDs()`: Queries all CRDs across all namespaces
- `queryCRDsInNamespace()`: Queries CRDs in a specific namespace
- `updateCRDWithFileInfo()`: Updates CRD status with file information
- `initKubernetesClients()`: Initializes both regular and dynamic K8s clients

## Configuration

The controller automatically detects the runtime environment:
- **In-cluster:** Uses service account credentials
- **Local development:** Uses kubeconfig from `~/.kube/config`

## RBAC Permissions

The controller requires the following permissions:
- Get, list, watch, create, update, patch, delete on `filemonitors` resources
- Get, update, patch on `filemonitors/status`
- Create, patch on `events` for logging

## Monitoring

The controller logs its operations including:
- CRD discovery and enumeration
- File system monitoring events
- Error conditions and retries
- Status updates

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request