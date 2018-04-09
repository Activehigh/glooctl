package consul

import (
	"fmt"
	"time"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/solo-io/gloo/pkg/endpointdiscovery"
)

func NewEndpointDiscovery(masterUrl, kubeconfigPath string, resyncDuration time.Duration) (endpointdiscovery.Interface, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build rest config: %v", err)
	}

	ctl, err := newEndpointController(cfg, resyncDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize endpoint controller: %v", err)
	}

	return ctl, nil
}
