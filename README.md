# certdeploy

[![Build Status](https://github.com/crallen/certdeploy/workflows/Build/badge.svg?branch=master)](https://github.com/crallen/certdeploy/actions) ![Version v0.1.2](https://img.shields.io/badge/Version-v0.1.2-blue)

A small utility intended to assist with pushing TLS certificates (most notably 
wildcard certificates produced by Let's Encrypt) to multiple Kubernetes
clusters. It was built with the intention of having one place where
certificates are managed and renewed separate from the clusters.

## Usage

```
Deploy certificates to multiple Kubernetes clusters

Usage:
  certdeploy [flags]

Flags:
  -c, --config string       path to config file
  -h, --help                help for certdeploy
  -k, --kubeconfig string   path to kubeconfig file (default "$HOME/.kube/config")
```

## Configuration

A simple example of a configuration file containing one cluster might look
something like this:

```yaml
secrets:
  # Key to reference the secret by
  tls-dev:
    # The name of the secret to create
    name: tls-dev
    # Key value pairs representing the data the secret will be created with
    files:
      tls.crt: /etc/letsencrypt/live/myawesomedevdomain.com/fullchain.pem
      tls.key: /etc/letsencrypt/live/myawesomedevdomain.com/privkey.pem
    # A list of namespaces where this secret will be created
    namespaces:
      - kube-system
      - my-ns

clusters:
    # A label to use for the cluster in output
  - name: dev
    # The context to use from the kubeconfig file
    context: dev-cluster
    # A list of secrets (defined above) to create in this cluster
    secrets:
      - tls-dev
```

## Kubernetes Integration

It is highly recommended you use a Kubernetes service account when connecting
to the clusters. There are some scripts included to assist with this in the
`scripts` directory. These scripts may not be suitable for all use cases, so
please modify them as necessary for your particular situation.

The `service-account.yml` file is an example of such a service account. It
creates a cluster-wide role that allows the service account to get secrets
by name, create secrets, and update secrets. This is the basic
functionality that certdeploy uses. You can secure this even further if you
already know what the names of your secrets will be:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: certdeploy
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    resourceNames: ["my-tls-secret"]
    verbs: ["get", "create", "update"]
```

This will limit the service account to _only_ being able to get, create, or
update a secret with the name `my-tls-secret`.

To use this service account when connecting to clusters, you will need a
`kubeconfig` file with the necessary credentials. The script
`create-kubeconfig.sh` can assist with this once the service account has been
created in your cluster(s). This script assumes that the machine you're
running it from already has one or more Kubernetes contexts available to
connect to clusters via tools like `kubectl`. Using that, it can then
retrieve the service account's credentials and generate a new `kubeconfig`
based on that information. It can also be run multiple times against the same
`kubeconfig` file to generate multiple contexts, which is useful when using
certdeploy for multiple clusters.

The following command will create a `kubeconfig` file in the current directory, 
where `certdeploy` is the name of the service account:

```
scripts/create-kubeconfig.sh certdeploy
```

Optionally, you can also control where the output will be placed:

```
scripts/create-kubeconfig.sh -o ./cfg/kubeconfig certdeploy
``` 

More information is available by running `scripts/create-kubeconfig.sh --help`.
