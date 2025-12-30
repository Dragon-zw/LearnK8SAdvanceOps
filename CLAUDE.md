# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a comprehensive Kubernetes learning project called LearnK8SAdvanceOps, designed to teach advanced Kubernetes operations and cloud-native technologies. The project covers various aspects of Kubernetes from basic concepts to advanced operators and application delivery platforms.

## Directory Structure

- `1.BaseK8S/` - Core Kubernetes fundamentals and building blocks
  - `1.BaseDocker/` - Docker basics
  - `2.Dockerfile/` - Dockerfile examples and best practices
  - `3.Kubernetes/` - Core Kubernetes resources (Pods, Controllers, Services, Volumes, RBAC)
  - `4.Helm/` - Helm chart examples
  - `5.Debug/` - Debugging techniques
  - `6.client-go/` - Kubernetes Go client examples
  - `7.CRD/` - Custom Resource Definition examples
  - `8.Kubebuilder/` - Operator framework with memcached-operator example
  - `9.Shell-Operator/` - Shell-based operator examples

- `1.Kustomize/` - Configuration management with Kustomize
- `2.KubeVela/` - Modern application delivery platform examples
- `3.Spinnaker/` - Continuous delivery platform configurations

## Key Technologies

- Kubernetes: Container orchestration
- Docker: Container technology
- Helm: Package management
- Kustomize: Configuration management
- Operator Framework: Using Kubebuilder to create custom controllers
- client-go: Official Go client library for Kubernetes
- CRD (Custom Resource Definitions): Extending Kubernetes API
- OAM (Open Application Model): Application delivery model used by KubeVela

## Development Commands

### For the memcached-operator (Go-based Operator):

```bash
# Build the operator
make build

# Build and push the controller image
make docker-build
make docker-push IMG=<your-image-tag>

# Generate manifests (CRDs, RBAC, etc.)
make manifests

# Generate code (DeepCopy methods)
make generate

# Run tests
make test

# Run linting
make lint
make lint-fix  # to fix issues

# Install CRDs to cluster
make install

# Deploy the operator to cluster
make deploy IMG=<your-image-tag>

# Undeploy the operator
make undeploy

# Run the operator locally (against a cluster)
make run

# Format code
make fmt

# Vet code
make vet

# Build installer (consolidated YAML)
make build-installer
```

### For Docker-based applications:

```bash
# Build Docker images
docker build -t <image-name> .

# Build with specific Dockerfile
docker build -f Dockerfile -t <image-name> .
```

### For Kustomize:

```bash
# Apply Kustomize configs
kubectl kustomize config/default | kubectl apply -f -

# Or with kubectl apply -k (if supported)
kubectl apply -k config/default
```

## Architecture & Design Patterns

### Operator Pattern (Kubebuilder/memcached-operator)
- Uses the controller pattern to watch for changes to custom resources
- Implements reconciliation loop to drive the cluster to the desired state
- Separates business logic (Controller) from the CRD definition (API)
- Uses client-go to interact with the Kubernetes API
- Includes webhooks for validation and mutation of custom resources

### Client-Go Examples
- Shows various patterns for interacting with Kubernetes API from outside and inside the cluster
- Demonstrates leader election for high availability
- Implements workqueues for reliable processing

### Multi-layer Architecture
- Infrastructure layer (Kubernetes, containers)
- Application layer (Operators managing applications)
- Configuration layer (Helm, Kustomize)
- Delivery layer (KubeVela, Spinnaker)

## Important Files & Locations

- `1.BaseK8S/8.Kubebuilder/memcached-operator/` - Main Go operator example with Makefile, controller logic, and CRDs
- `1.BaseK8S/6.client-go/client-go/examples/` - Various client-go usage patterns
- `1.Kustomize/k8s/base` and `1.Kustomize/k8s/overlays` - Configuration management examples
- `2.KubeVela/Docs/` - Documentation on modern application delivery
- `.claude/config.json` - Claude Code configuration for this repository

## Best Practices Demonstrated

- Use of structured logging with zap
- Proper error handling in reconcilers
- Health and readiness probes
- Leader election for HA deployments
- RBAC definitions for security
- Code generation for deepcopy methods
- Testing patterns for controllers
- Secure metrics endpoint configuration