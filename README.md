# highlander-proxy

highlander-proxy is a simple network proxy that encapsulates leadership election semantics.

## Try it out

```
$ go mod vendor
$ go run main.go \
    -bind-address localhost:1234 \
    -kubeconfig ~/.kube/minikube.yaml \
    -lock-namespace highlander-proxy \
    -lock-name demo \
    -protocol tcp \
    -remote-address host:port
```
