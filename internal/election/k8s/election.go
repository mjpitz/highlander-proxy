package k8s

import (
	"context"

	"github.com/mjpitz/highlander-proxy/internal/election"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

// NewElector creates a new Kubernetes based elector and returns a corresponding Leader.
func NewElector(config *election.Config, kubeconfig string) (*election.Leader, error) {
	var cfg *rest.Config
	var err error

	if kubeconfig != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		cfg, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, err
	}

	client := clientset.NewForConfigOrDie(cfg)

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Namespace: config.LockNamespace,
			Name:      config.LockName,
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: config.Identity,
		},
	}

	leader := election.NewLeader()

	go func() {
		leaderelection.RunOrDie(config.Context, leaderelection.LeaderElectionConfig{
			Lock:            lock,
			ReleaseOnCancel: false, // otherwise process won't shut down
			LeaseDuration:   config.LeaseDuration,
			RenewDeadline:   config.RenewDeadline,
			RetryPeriod:     config.RetryPeriod,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {},
				OnStoppedLeading: func() {},
				OnNewLeader: func(identity string) {
					leader.Update(identity)
				},
			},
		})
	}()

	return leader, nil
}
