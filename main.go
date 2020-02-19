package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	leaseDuration := 60 * time.Second
	renewDuration := 15 * time.Second
	retryPeriod := 5 * time.Second
	kubeconfig := ""

	flag.StringVar(&protocol, "protocol", protocol, "")
	flag.StringVar(&bindAddress, "bind-address", bindAddress, "")
	flag.StringVar(&remoteAddress, "remote-address", remoteAddress, "")
	flag.StringVar(&lockNamespace, "lock-namespace", lockNamespace, "")
	flag.StringVar(&lockName, "lock-name", lockName, "")
	flag.DurationVar(&leaseDuration, "lease-duration", leaseDuration, "")
	flag.DurationVar(&renewDuration, "renew-duration", renewDuration, "")
	flag.DurationVar(&retryPeriod, "retry-period", retryPeriod, "")
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "")

	flag.Parse()

	log.Println("listeninig on", bindAddress)
	listener, err := net.Listen(protocol, bindAddress)
	exitIff(err)

	root, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		log.Println("received termination, signaling shutdown")
		cancel()
	}()

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

	dialer := &proxy.Dialer{
		Leader:        leader,
		Protocol:      protocol,
		Identity:      bindAddress,
		RemoteAddress: remoteAddress,
	}

	connections := make(chan net.Conn, 10)

	go proxy.Connect(root, listener, connections)
	go proxy.Forward(root, dialer, connections)

	select {
	// let the workers be free
	}
}
