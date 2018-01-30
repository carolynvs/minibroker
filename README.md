# InCluster Broker

This is an implementation of the [Open Service Broker API](https://openservicebrokerapi.org)
that is suited for local development. Instead of provisioning services
from a cloud provider, it creates the service in a container on the cluster using Helm.

# Install
TODO: Create a helm chart and repo.

# Use
TODO: Example manifests to create a mysql db.

# Local Development

## Requirements

* Docker
* [Minikube v0.24.1](https://github.com/kubernetes/minikube/releases/tag/v0.24.1)

On a Mac you will also need either VirtualBox installed,
or the [Minikube xhyve driver](https://github.com/kubernetes/minikube/blob/master/docs/drivers.md#xhyve-driver)
which uses the hypervisor that comes with Docker for Mac.

The default Minikube driver is virtualbox, to use xhyve specify it in
**~/.minikube/config/config.json**:

```json
{
    "vm-driver": "xhyve"
}
```

## Optional Tools
The Makefile tries to runs non-required tools in a Docker container. If you prefer to run
the commands locally, export `USE_DOCKER=false` and install the tools below:

* [Draft](https://draft.sh)
* [Helm](https://helm.sh)
* [Service Catalog CLI (svcat)](https://github.com/kubernetes-incubator/service-catalog/cmd/svcat)


## Initial Setup

1. Create a Minikube cluster for local development: `make init`.
2. Make sure everything is running: `kubectl get nodes --all-namespaces`.

## Deploy

Compile and deploy the broker to your local cluster: `make run`.

## Dependency Management

We use [dep](https://golang.github.io/dep) to manage our dependencies. Our vendor
directory is checked-in and kept up-to-date with Gopkg.lock, so unless you are
actively changing dependencies, you don't need to do anything extra.

### Add a new dependency

1. Add the dependency.
    * Import the dependency in the code OR
    * Run `dep ensure --add github.com/pkg/example@v1.0.0` to add an explicit constraint
       to Gopkg.toml.
       
       This is only necessary when we want to stick with a particular branch
       or version range, otherwise the lock will keep us on the same version and track what's used.
2. Run `dep ensure`.
3. Check in the changes to `Gopkg.lock` and `vendor/`.
