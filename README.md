# Varnish Cache Invalidator
[![CI](https://github.com/bilalcaliskan/varnish-cache-invalidator/workflows/CI/badge.svg?event=push)](https://github.com/bilalcaliskan/varnish-cache-invalidator/actions?query=workflow%3ACI)
[![Docker pulls](https://img.shields.io/docker/pulls/bilalcaliskan/varnish-cache-invalidator)](https://hub.docker.com/r/bilalcaliskan/varnish-cache-invalidator/)

This tool discovers kube-apiserver for running [Varnish](https://github.com/varnishcache/varnish-cache) pods inside 
Kubernetes and multiplexes `BAN` and `PURGE` requests on them at the same time to manage the cache properly. If you are 
using Varnish Enterprise, you already have that feature.

## Deployment
Varnish-cache-invalidator should be running inside the Kubernetes cluster to properly multiplex incoming requests. 
You can use [sample config file](config/sample.yaml) to deploy your Kubernetes cluster.

```shell
$ kubectl create -f config/sample.yaml
```

### Customization
Varnish-cache-invalidator can be customized with several environment variables. You can pass environment variables 
via [sample config file](config/sample.yaml). Here is the list of variables you can pass:

```
SERVER_PORT             Web server port to handle incoming http requests. Defaults to "3000".
METRICS_PORT            varnish-cache-invalidator exports prometheus metrics on specified port. Defaults 
                        to "3001".
WRITE_TIMEOUT_SECONDS   Maximum duration before timing out writes of the response. Defaults to "10".
READ_TIMEOUT_SECONDS    Maximum duration for reading the entire request, including the body.
VARNISH_NAMESPACE       Namespace of the Varnish pods. Defaults to "default".
VARNISH_LABEL           varnish-cache-invalidator fetches Add/Update/Delete events from kube-apiserver. 
                        Label will help us choosing correct Varnish pods. Defaults to "app=varnish".
MASTER_URL              Additional url to kube-apiserver. There is no need to specify anything here 
                        since we are using serviceAccount to access kube-apiserver. Defaults to "".
KUBE_CONFIG_PATH        Kubeconfig path to access to cluster while running the code out of cluster. 
                        Required for development purposes. Defaults to "~/.kube/config".
IN_CLUSTER              Specify if we are running the code out of cluster. Required for development purposes. 
                        Defaults to "false".
TARGET_HOSTS            Comma seperated lists of target hosts while our Varnish instances are 
                        not running in Kubernetes as pods. Designed for standalone installations but 
                        not actively used. Defaults to "" but not required when IN_CLUSTER=true.
```


### How varnish-cache-invalidator handles authentication with kube-apiserver?

varnish-cache-invalidator uses [client-go](https://github.com/kubernetes/client-go) to interact 
with `kube-apiserver`. [client-go](https://github.com/kubernetes/client-go) uses the [service account token](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/) 
mounted inside the Pod at the `/var/run/secrets/kubernetes.io/serviceaccount` path while initializing the client.

If you have RBAC enabled on your cluster, when you applied the sample deployment file [config/sample.yaml](config/sample.yaml), 
it will create required serviceaccount and clusterrolebinding and then use that serviceaccount to be used 
by our varnish-cache-invalidator pods.

If RBAC is not enabled on your cluster, please follow [that documentation](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) to enable it.
