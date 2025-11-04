# Kubernetes Config Generator

Generate Kubernetes manifests directly from Go code with a CLI-based tool. Supports both interactive and non-interactive workflows for CI/CD pipelines.

## Features

- **Direct YAML Generation**: Generates standard Kubernetes manifests without external dependencies
- **Interactive Mode**: Prompts for input when no flags are provided
- **CLI Flags**: Non-interactive, suitable for automation
- **Environment Support**: Staging and production configurations
- **Multi-Environment**: Generate both staging and production manifests in a single run
- **Configurable**: All resources, ingress, VPA, and quotas are configurable
- **Secure Defaults**: Security contexts and best practices built-in
- **Minimal Labels**: Only adds labels where functionally required

## Prerequisites

- Go 1.17 or later (for building from source)

## Installation

### Build from Source

```bash
git clone https://github.com/pravinbanjade/k8s-config-generator.git
cd k8s-config-generator
go build -o k8s-config-generator main.go
```

### Download Pre-built Binaries

Download from [releases](https://github.com/pravinbanjade/k8s-config-generator/releases).

## Usage

### Interactive Mode

Run the tool without flags to be prompted for input:

```bash
./k8s-config-generator
```

This will prompt you for:
- Application name
- Docker image repository
- Environment (staging/production/both)
- Docker image tag(s)
- Namespace
- Ingress configuration
- And other optional settings

### Basic Example - Generate Staging Environment

Generate Kubernetes manifests for a staging environment:

```bash
./k8s-config-generator \
  --app-name myapp \
  --image-repo registry.example.com/myapp \
  --image-tag staging-123 \
  --namespace myapp-staging \
  --env staging \
  --container-port 3000 \
  --ingress-enabled \
  --ingress-host stage.example.com \
  --ingress-class nginx \
  --output-dir ./manifests
```

### Production with Ingress TLS

```bash
./k8s-config-generator \
  --app-name myapp \
  --image-repo registry.example.com/myapp \
  --image-tag production-abc123 \
  --namespace myapp-production \
  --env production \
  --container-port 3000 \
  --ingress-enabled \
  --ingress-host prod.example.com \
  --ingress-class nginx \
  --ingress-tls-secret k8s-tls-secret-replica \
  --output-dir ./manifests
```

### Generate Both Environments at Once

Generate manifests for both staging and production in a single run:

```bash
./k8s-config-generator \
  --app-name myapp \
  --image-repo registry.example.com/myapp \
  --all-environments \
  --image-tag-stage staging-123 \
  --image-tag-prod production-abc123 \
  --namespace myapp-staging \
  --container-port 3000 \
  --ingress-enabled \
  --ingress-host-stage stage.example.com \
  --ingress-host-prod prod.example.com \
  --ingress-class nginx \
  --ingress-tls-secret-prod k8s-tls-secret-replica \
  --output-dir ./manifests
```

This will create separate directories for each environment:
- `./manifests/staging/`
- `./manifests/production/`

### Render to stdout

Output manifests to stdout instead of files:

```bash
./k8s-config-generator \
  --app-name myapp \
  --image-repo registry.example.com/myapp \
  --image-tag staging-123 \
  --namespace myapp-staging \
  --env staging \
  --render
```

### With Resource Limits and VPA

```bash
./k8s-config-generator \
  --app-name myapp \
  --image-repo registry.example.com/myapp \
  --image-tag staging-123 \
  --namespace myapp-staging \
  --env staging \
  --container-port 3000 \
  --replicas 2 \
  --resources-requests-cpu 100m \
  --resources-requests-memory 256Mi \
  --resources-limits-cpu 500m \
  --resources-limits-memory 512Mi \
  --vpa-enabled \
  --resource-quota-enabled \
  --output-dir ./manifests
```

### With Image Pull Secrets

```bash
./k8s-config-generator \
  --app-name myapp \
  --image-repo registry.example.com/myapp \
  --image-tag staging-123 \
  --namespace myapp-staging \
  --env staging \
  --image-pull-secret gitlab-credentials \
  --image-pull-secret docker-registry-secret \
  --output-dir ./manifests
```

## CLI Flags

### Required Flags (or provide via interactive prompts)

- `--app-name`: Application name
- `--image-repo`: Docker image repository
- `--image-tag`: Docker image tag (not required if using `--all-environments`)

### Optional Flags

#### Basic Configuration

- `--namespace`: Kubernetes namespace
- `--container-port`: Container port (default: 3000)
- `--env`: Environment (staging|production)
- `--all-environments`: Generate manifests for both staging and production
- `--replicas`: Number of replicas (default: 1)

#### Multi-Environment Configuration (used with `--all-environments`)

- `--image-tag-stage`: Docker image tag for staging
- `--image-tag-prod`: Docker image tag for production
- `--ingress-host-stage`: Ingress host for staging
- `--ingress-host-prod`: Ingress host for production
- `--ingress-tls-secret-stage`: Ingress TLS secret for staging
- `--ingress-tls-secret-prod`: Ingress TLS secret for production

#### Ingress Configuration

- `--ingress-enabled`: Enable ingress
- `--ingress-host`: Ingress host
- `--ingress-class`: Ingress class name (default: nginx)
- `--ingress-tls-secret`: Ingress TLS secret name

#### Image Configuration

- `--image-pull-secret`: Image pull secret name (can be repeated)

#### Service Account

- `--service-account`: Service account name
- `--create-service-account`: Create service account (default: true)

#### Resources

- `--resources-requests-cpu`: Resource requests CPU
- `--resources-requests-memory`: Resource requests memory
- `--resources-limits-cpu`: Resource limits CPU
- `--resources-limits-memory`: Resource limits memory

#### Advanced Features

- `--vpa-enabled`: Enable Vertical Pod Autoscaler
- `--resource-quota-enabled`: Enable resource quota

#### Output Modes

- `--render`: Render manifests to stdout
- `--output-dir`: Output directory for rendered manifests (creates files if not using `--render`)

## Generated Resources

The tool generates the following Kubernetes resources:

- **Namespace**: Creates namespace if specified
- **Deployment**: Main application deployment with configurable replicas
- **Service**: ClusterIP service for the deployment
- **ServiceAccount**: Service account for the pods (optional)
- **ConfigMap**: Environment-specific configuration
- **Secret**: Application secrets (production only)
- **Ingress**: HTTP/HTTPS ingress (optional)
- **ResourceQuota**: Resource quota limits (optional)
- **VPA**: Vertical Pod Autoscaler (optional)

## Output Structure

When using `--output-dir`, manifests are organized as follows:

### Single Environment

```
output-dir/
├── deployment-{app-name}-node.yaml
├── service-{app-name}.yaml
├── namespace-{namespace}.yaml
├── serviceaccount-{app-name}.yaml
├── configmap-{app-name}.yaml
├── ingress-{app-name}-ingress.yaml
└── ...
```

### Multiple Environments (`--all-environments`)

```
output-dir/
├── staging/
│   ├── deployment-{app-name}-node.yaml
│   ├── service-{app-name}.yaml
│   ├── namespace-{app-name}-staging.yaml
│   └── ...
└── production/
    ├── deployment-{app-name}-node.yaml
    ├── service-{app-name}.yaml
    ├── namespace-{app-name}-production.yaml
    └── ...
```

## CI/CD Integration

Example GitHub Actions workflow:

```yaml
name: Deploy to Kubernetes

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build k8s-config-generator
        run: |
          go build -o k8s-config-generator main.go
      
      - name: Generate Manifests
        run: |
          ./k8s-config-generator \
            --app-name myapp \
            --image-repo ${{ secrets.DOCKER_REGISTRY }}/myapp \
            --image-tag staging-${{ github.sha }} \
            --namespace myapp-staging \
            --env staging \
            --ingress-enabled \
            --ingress-host staging.example.com \
            --output-dir ./manifests
      
      - name: Apply Manifests
        run: |
          kubectl apply -f ./manifests
        env:
          KUBECONFIG: ${{ secrets.KUBECONFIG }}
```

### Generate and Apply Both Environments

```yaml
      - name: Generate Manifests for Both Environments
        run: |
          ./k8s-config-generator \
            --app-name myapp \
            --image-repo ${{ secrets.DOCKER_REGISTRY }}/myapp \
            --all-environments \
            --image-tag-stage staging-${{ github.sha }} \
            --image-tag-prod production-${{ github.sha }} \
            --namespace myapp-staging \
            --ingress-enabled \
            --ingress-host-stage staging.example.com \
            --ingress-host-prod prod.example.com \
            --ingress-tls-secret-prod k8s-tls-secret-replica \
            --output-dir ./manifests
      
      - name: Apply Staging Manifests
        run: |
          kubectl apply -f ./manifests/staging
        env:
          KUBECONFIG: ${{ secrets.KUBECONFIG_STAGING }}
      
      - name: Apply Production Manifests
        run: |
          kubectl apply -f ./manifests/production
        env:
          KUBECONFIG: ${{ secrets.KUBECONFIG_PROD }}
```

## Labels

The tool follows Kubernetes best practices by adding labels only where functionally required:

- **Deployment**: Labels on `spec.selector.matchLabels` and `spec.template.metadata.labels` (required for pod selection)
- **Service**: Labels on `spec.selector` (required for pod selection)
- **Other resources**: No labels (unless needed for functional purposes)

This keeps manifests clean and minimal while maintaining functionality.

## Development

### Building

```bash
go build -o k8s-config-generator main.go
```

### Testing

```bash
go test ./...
```

## License

[Add your license here]
