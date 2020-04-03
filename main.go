package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mjpitz/highlander-proxy/internal/config"
	"github.com/mjpitz/highlander-proxy/internal/election"
	"github.com/mjpitz/highlander-proxy/internal/election/k8s"
	"github.com/mjpitz/highlander-proxy/internal/proxy"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func exitIff(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

var longDescription = `
highlander-proxy is a simple to use, leader elected network proxy. It works by
running alongside applications as a sidecar, ensuring that only one instance 
handles traffic at a time. As soon as a new leader is elected, all connections 
to the previous leader are terminated.

To configure routes, use the the --routes flag. This flag can be provided
multiple times for multiple port bindings. The format of the value is as such:

  <bind-address>|<forward-address>[,<bind-address>|<forward-address>,...]

Some examples of route configurations include:

  tcp://0.0.0.0:8080|tcp://localhost:8080
  tcp://0.0.0.0:8080|unix:///path/to/file
`

func main() {
	identity := ""
	routes := config.RouteSlice{}
	logLevel := "info"

	// Kubernetes Configuration
	kubeconfig := ""
	lockNamespace := ""
	lockName := ""
	leaseDuration := 5 * time.Second
	renewDuration := 2 * time.Second
	retryPeriod := 1 * time.Second

	cmd := &cobra.Command{
		Use:   "highlander-proxy",
		Short: "A leader elected network proxy.",
		Long:  longDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			stopCh := make(chan struct{})

			for _, route := range routes.Routes() {
				server := &proxy.Server{
					Route:    route,
					Identity: identity,
					Leader:   leader,
				}

				err := server.Start(stopCh)
				exitIff(err)
			}

			ch := make(chan os.Signal, 2)
			signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-ch
				logrus.Infof("received shutdown signal, terminating process")
				close(stopCh)

				<-ch
				logrus.Infof("forcing termination of process")
				os.Exit(1)
			}()

			select {
			// let the workers be free
			}
		},
	}

	flags := cmd.Flags()

	// Route Configuration
	flags.StringVar(&identity, "identity", identity, "The identity of this process.")
	flags.Var(&routes, "routes", "Configures the underlying route table.")

	// Debug Configuration
	flags.StringVar(&logLevel, "log-level", logLevel, "Verbosity of the log messages.")

	// Kubernetes Configuration
	flags.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "Location of the kubeconfig file.")
	flags.StringVar(&lockNamespace, "lock-namespace", lockNamespace, "The namespace for the lock.")
	flags.StringVar(&lockName, "lock-name", lockName, "The name of the lock within the namespace.")
	flags.DurationVar(&leaseDuration, "lease-duration", leaseDuration, "The duration a lock is held.")
	flags.DurationVar(&renewDuration, "renew-duration", renewDuration, "The duration between lock renewal.")
	flags.DurationVar(&retryPeriod, "retry-period", retryPeriod, "The duration between retries.")

	err := cmd.Execute()
	exitIff(err)
}
