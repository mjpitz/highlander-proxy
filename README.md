# highlander-proxy

highlander-proxy is a simple network proxy that encapsulates leadership election semantics.

## Try it out

Before trying this out, you must first pull down dependencies.

```
$ go mod vendor
```

Once you've resolved dependencies, you should be able to spin up several proxies using Kubernetes.

```
$ go run main.go \
    -bind-address localhost:1234 \
    -kubeconfig ~/.kube/minikube.yaml \
    -lock-namespace highlander-proxy \
    -lock-name demo \
    -protocol tcp \
    -remote-address host:port

$ go run main.go \
    -bind-address localhost:1235 \
    -kubeconfig ~/.kube/minikube.yaml \
    -lock-namespace highlander-proxy \
    -lock-name demo \
    -protocol tcp \
    -remote-address host:port
```
