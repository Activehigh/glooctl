package install

import (
	_ "github.com/solo-io/gloo-plugins/aws"
	_ "github.com/solo-io/gloo-plugins/google"
	_ "github.com/solo-io/gloo-plugins/grpc"
	_ "github.com/solo-io/gloo-plugins/kubernetes"
	_ "github.com/solo-io/gloo-plugins/nats-streaming"
	_ "github.com/solo-io/gloo-plugins/rest"
)
