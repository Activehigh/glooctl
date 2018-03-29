package natsstreaming

import (
	"errors"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"

	"github.com/solo-io/gloo/pkg/protoutil"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/coreplugins/common"
	"github.com/solo-io/gloo/pkg/plugin"
)

//go:generate protoc -I=. -I=${GOPATH}/src/github.com/gogo/protobuf/ --gogo_out=. nats_streaming_filter.proto

func init() {
	plugin.Register(&Plugin{}, nil)
}

type Plugin struct {
	filters []plugin.StagedFilter
}

const (
	ServiceTypeNatsStreaming = "nats-streaming"

	// generic plugin info
	filterName  = "io.solo.nats_streaming"
	pluginStage = plugin.OutAuth

	clusterId                 = "cluster_id"
	clusterIdAnnotations      = "gloo.solo.io/cluster_id"
	defaultClusterId          = "test-cluster"
	discoverPrefix            = "discover_prefix"
	discoverPrefixAnnotations = "gloo.solo.io/discover_prefix"
	defaultDiscoverPrefix     = "_STAN.discover"
)

func (p *Plugin) GetDependencies(cfg *v1.Config) *plugin.Dependencies {
	return nil
}

func (p *Plugin) HttpFilters(params *plugin.FilterPluginParams) []plugin.StagedFilter {
	filters := p.filters
	p.filters = nil
	return filters
}

func (p *Plugin) ProcessUpstream(params *plugin.UpstreamPluginParams, in *v1.Upstream, out *envoyapi.Cluster) error {
	if in.ServiceInfo == nil || in.ServiceInfo.Type != ServiceTypeNatsStreaming {
		return nil
	}
	//    string nats_streaming_cluster_id = 3;
	//    string discover_prefix = 4;
	// in.Metadata

	cid := in.Metadata.Annotations[clusterIdAnnotations]
	if cid == "" {
		cid = defaultClusterId
	}
	dp := in.Metadata.Annotations[clusterIdAnnotations]
	if dp == "" {
		dp = defaultDiscoverPrefix
	}
	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{}
	}
	common.InitFilterMetadataField(filterName, clusterId, out.Metadata).Kind = &types.Value_StringValue{StringValue: defaultClusterId}
	common.InitFilterMetadataField(filterName, discoverPrefix, out.Metadata).Kind = &types.Value_StringValue{StringValue: dp}

	p.filters = append(p.filters, plugin.StagedFilter{HttpFilter: &envoyhttp.HttpFilter{Name: filterName, Config: natsConfig(out.Name)}, Stage: pluginStage})

	return nil
}

func natsConfig(cluster string) *types.Struct {
	natsStreaming := NatsStreaming{
		MaxConnections: 1,
		Cluster:        cluster,
	}

	filterConfig, err := protoutil.MarshalStruct(&natsStreaming)
	if err != nil {
		runtime.HandleError(err)
		return nil
	}
	return filterConfig
}

func (p *Plugin) ParseFunctionSpec(params *plugin.FunctionPluginParams, in v1.FunctionSpec) (*types.Struct, error) {
	if params.ServiceType != ServiceTypeNatsStreaming {
		return nil, nil
	}
	return nil, errors.New("functions are not required for service type " + ServiceTypeNatsStreaming)
}
