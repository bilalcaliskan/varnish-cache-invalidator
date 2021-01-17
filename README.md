## Varnish Cache Invalidator
[![CI](https://github.com/bilalcaliskan/varnish-cache-invalidator/workflows/CI/badge.svg?event=push)](https://github.com/bilalcaliskan/varnish-cache-invalidator/actions?query=workflow%3ACI)

### Authenticating inside the cluster

client-go uses the [service account token](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/) 
mounted inside the Pod at the `/var/run/secrets/kubernetes.io/serviceaccount` path when the
`rest.InClusterConfig()` is used while initializing the client.

If you have RBAC enabled on your cluster, use the following
snippet to create the service account first, then create a role binding which 
will grant the previously created service account view
permissions. Finally use that serviceaccount in your deployment or deploymentconfig(Openshift).

```
kubectl create serviceaccount ${YOUR_SERVICE_ACCOUNT_NAME}
kubectl create clusterrolebinding ${YOUR_SERVICE_ACCOUNT_NAME}-view --clusterrole=view --serviceaccount=${YOUR_NAMESPACE}:${YOUR_SERVICE_ACCOUNT_NAME}
```