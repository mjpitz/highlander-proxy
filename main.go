package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/mjpitz/highlander-proxy/internal/election"
	"github.com/mjpitz/highlander-proxy/internal/election/k8s"
	"github.com/mjpitz/highlander-proxy/internal/proxy"
)

func exitIff(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func main() {
	protocol := "tcp"
	bindAddress := ""
	remoteAddress := ""
	lockNamespace := ""
	lockName := ""
	leaseDuration := 5 * time.Second
	renewDuration := 2 * time.Second
	retryPeriod := 1 * time.Second
	kubeconfig := ""

	flag.StringVar(&protocol, "protocol", protocol, "The network protocol to operate on. Either 'tcp' or 'udp'.")
	flag.StringVar(&bindAddress, "bind-address", bindAddress, "The address to bind the proxy to.")
	flag.StringVar(&remoteAddress, "remote-address", remoteAddress, "The destination address.")
	flag.StringVar(&lockNamespace, "lock-namespace", lockNamespace, "The namespace for the lock.")
	flag.StringVar(&lockName, "lock-name", lockName, "The name of the lock within the namespace.")
	flag.DurationVar(&leaseDuration, "lease-duration", leaseDuration, "The duration a lock is held.")
	flag.DurationVar(&renewDuration, "renew-duration", renewDuration, "The duration between lock renewal.")
	flag.DurationVar(&retryPeriod, "retry-period", retryPeriod, "The duration between retries.")
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "Location of the kubeconfig file.")

	flag.Parse()

	root, cancel := context.WithCancel(context.Background())
	defer cancel()

	electionConfig := &election.Config{
		Context:       root,
		Identity:      bindAddress,
		LockNamespace: lockNamespace,
		LockName:      lockName,
		LeaseDuration: leaseDuration,
		RenewDeadline: renewDuration,
		RetryPeriod:   retryPeriod,
	}

	leader, err := k8s.NewElector(electionConfig, kubeconfig)
	exitIff(err)

	server := &proxy.Server{
		Protocol:      protocol,
		BindAddress:   bindAddress,
		RemoteAddress: remoteAddress,
		Leader:        leader,
	}

	log.Println("listening on", bindAddress)
	err = server.Serve()
	exitIff(err)
}
