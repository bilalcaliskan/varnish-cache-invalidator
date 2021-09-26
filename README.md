# Varnish Cache Invalidator
[![CI](https://github.com/bilalcaliskan/varnish-cache-invalidator/workflows/CI/badge.svg?event=push)](https://github.com/bilalcaliskan/varnish-cache-invalidator/actions?query=workflow%3ACI)
[![Docker pulls](https://img.shields.io/docker/pulls/bilalcaliskan/varnish-cache-invalidator)](https://hub.docker.com/r/bilalcaliskan/varnish-cache-invalidator/)
[![Go Report Card](https://goreportcard.com/badge/github.com/bilalcaliskan/varnish-cache-invalidator)](https://goreportcard.com/report/github.com/bilalcaliskan/varnish-cache-invalidator)
[![codecov](https://codecov.io/gh/bilalcaliskan/varnish-cache-invalidator/branch/master/graph/badge.svg)](https://codecov.io/gh/bilalcaliskan/varnish-cache-invalidator)

This tool basically multiplexes **BAN** and **PURGE** requests on [Varnish Cache](https://github.com/varnishcache/varnish-cache)
instances at the same time to manage the cache properly. If you are using Varnish Enterprise, you already have that feature.

varnish-cache-invalidator can be used for standalone Varnish instances or Varnish pods inside a Kubernetes cluster.

**Standalone mode**
In that mode, varnish-cache-invalidator multiplexes requests on static standalone Varnish instances which are provided
with **--targetHosts** flag. This flag gets comma seperated list of hosts.

Please check all the command line arguments on **Configuration** section. Required args for standalone mode:
```
--inCluster=false
--targetHosts
```

**Kubernetes mode**
In that mode, varnish-cache-invalidator discovers kube-apiserver for running [Varnish](https://github.com/varnishcache/varnish-cache) pods inside
Kubernetes and multiplexes **BAN** and **PURGE** requests on them at the same time to manage the cache properly.

Please check all the command line arguments on **Configuration** section. Required args for Kubernetes mode:
```
--inCluster=true
```

## Installation
### Kubernetes
Varnish-cache-invalidator can be run inside a Kubernetes cluster to multiplex requests for in-cluster Varnish containers.
You can use [sample deployment file](deployment/sample.yaml) to deploy your Kubernetes cluster.

```shell
$ kubectl create -f config/sample.yaml
```

### Binary
Binary can be downloaded from [Releases](https://github.com/bilalcaliskan/nginx-conf-generator/releases) page. You can
use binary method to manage standalone Varnish instances, not in Kubernetes.

After then, you can simply run binary by providing required command line arguments:
```shell
$ ./varnish-cache-invalidator --inCluster=false --targetHosts 10.0.0.100:6081,10.0.0.101:6081,10.0.0.102:6081
```

## Configuration
Varnish-cache-invalidator can be customized with several command line arguments. You can pass command line arguments via
[sample deployment file](deployment/sample.yaml). Here is the list of arguments you can pass:

```
--inCluster             InCluster is the boolean flag if varnish-cache-invalidator is running inside cluster or not,
                        defaults to true
--varnishNamespace      VarnishNamespace is the namespace of the target Varnish pods, defaults to default namespace
--varnishLabel          VarnishLabel is the label to select proper Varnish pods, defaults to app=varnish
--targetHosts           TargetHosts used when our Varnish instances(comma seperated) are not running in Kubernetes as
                        a pod, required for standalone Varnish instances, defaults to ''
--serverPort            ServerPort is the web server port of the varnish-cache-invalidator, defaults to 3000
--metricsPort           MetricsPort is the port of the metrics server, defaults to 3001
--writeTimeoutSeconds   WriteTimeoutSeconds is the write timeout of the both web server and metrics server, defaults to 10
--readTimeoutSeconds    ReadTimeoutSeconds is the read timeout of the both web server and metrics server, defaults to 10
```

## Examples
**TBD**

## Development
This project requires below tools while developing:
- [Golang 1.16](https://golang.org/doc/go1.16)
- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://golangci-lint.run/usage/install/) - required by [pre-commit](https://pre-commit.com/)

## License
Apache License 2.0

## How varnish-cache-invalidator handles authentication with kube-apiserver in Kubernetes mode?

varnish-cache-invalidator uses [client-go](https://github.com/kubernetes/client-go) to interact
with `kube-apiserver`. [client-go](https://github.com/kubernetes/client-go) uses the [service account token](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/)
mounted inside the Pod at the `/var/run/secrets/kubernetes.io/serviceaccount` path while initializing the client.

If you have RBAC enabled on your cluster, when you applied the sample deployment file [deployment/sample.yaml](deployment/sample.yaml),
it will create required serviceaccount and clusterrolebinding and then use that serviceaccount to be used
by our varnish-cache-invalidator pods.

If RBAC is not enabled on your cluster, please follow [that documentation](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) to enable it.
