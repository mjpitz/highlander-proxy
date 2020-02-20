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

	"github.com/sirupsen/logrus"
)

func exitIff(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func main() {
	protocol := "tcp"
	identity := ""
	bindAddress := ""
	remoteAddress := ""
	lockNamespace := ""
	lockName := ""
	leaseDuration := 5 * time.Second
	renewDuration := 2 * time.Second
	retryPeriod := 1 * time.Second
	kubeconfig := ""
	logLevel := "info"

	flag.StringVar(&protocol, "protocol", protocol, "The network protocol to operate on. Either 'tcp' or 'udp'.")
	flag.StringVar(&identity, "identity", identity, "The identity of this process (ip:port).")
	flag.StringVar(&bindAddress, "bind-address", bindAddress, "The address to bind the proxy to.")
	flag.StringVar(&remoteAddress, "remote-address", remoteAddress, "The destination address.")
	flag.StringVar(&lockNamespace, "lock-namespace", lockNamespace, "The namespace for the lock.")
	flag.StringVar(&lockName, "lock-name", lockName, "The name of the lock within the namespace.")
	flag.DurationVar(&leaseDuration, "lease-duration", leaseDuration, "The duration a lock is held.")
	flag.DurationVar(&renewDuration, "renew-duration", renewDuration, "The duration between lock renewal.")
	flag.DurationVar(&retryPeriod, "retry-period", retryPeriod, "The duration between retries.")
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "Location of the kubeconfig file.")
	flag.StringVar(&logLevel, "log-level", logLevel, "Verbosity of the log messages.")

	flag.Parse()

	if identity == "" {
		identity = bindAddress
	}

	level, err := logrus.ParseLevel(logLevel)
	exitIff(err)

	logrus.SetLevel(level)

	root, cancel := context.WithCancel(context.Background())
	defer cancel()

	electionConfig := &election.Config{
		Context:       root,
		Identity:      identity,
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
		Identity:      identity,
		RemoteAddress: remoteAddress,
		Leader:        leader,
	}

	listener, err := net.Listen(protocol, bindAddress)
	exitIff(err)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		logrus.Infof("received shutdown signal, terminating process")
		cancel()
		time.Sleep(time.Second)
		_ = listener.Close()
	}()

	logrus.Infof("listening on %s://%s", protocol, bindAddress)
	server.Serve(listener)
}
