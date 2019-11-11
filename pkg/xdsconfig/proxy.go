package xdsconfig

import (
	"fmt"
	"strings"
	"time"

	"github.com/bladedancer/envoyxds/pkg/apimgmt"
	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/golang/protobuf/ptypes"
)

// Proxy is the representation of a "proxy" in envoy.
type Proxy struct {
	Routes  []*route.Route
	Cluster *api.Cluster
}

// MakeProxy create a proxy.
func MakeProxy(tenantName string, proxy *apimgmt.Proxy) *Proxy {
	return &Proxy{
		Routes:  makeRoute(tenantName, proxy),
		Cluster: makeCluster(tenantName, proxy),
	}
}

// makeRoute Creates two routes for the API - one that has an exact match on the basepath and one that has a prefix
// match on the basepath/
func makeRoute(tenantName string, proxy *apimgmt.Proxy) []*route.Route {
	clusterID := fmt.Sprintf("t_%s-p_%s", tenantName, proxy.Name)

	// Exact match on basepath
	exactRoute := &route.Route{
		Name: fmt.Sprintf("%s-%s", clusterID, "exact"),
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Path{
				Path: proxy.Frontend.BasePath,
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterID,
				},
				PrefixRewrite: proxy.Backend.Path,
				HostRewriteSpecifier: &route.RouteAction_HostRewrite{
					HostRewrite: proxy.Backend.Host,
				},
			},
		},
	}

	// Match basepath/
	target := proxy.Backend.Path
	if !strings.HasSuffix(target, "/") {
		target = target + "/"
	}

	prefixRoute := &route.Route{
		Name: fmt.Sprintf("%s-%s", clusterID, "prefix"),
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Prefix{
				Prefix: proxy.Frontend.BasePath + "/",
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterID,
				},
				PrefixRewrite: target,
				HostRewriteSpecifier: &route.RouteAction_HostRewrite{
					HostRewrite: proxy.Backend.Host,
				},
			},
		},
	}

	return []*route.Route{
		prefixRoute,
		exactRoute,
	}
}

// makeCluster Create a cluster for the proxy
func makeCluster(tenantName string, proxy *apimgmt.Proxy) *api.Cluster {
	//log.Infof("creating cluster for proxy: %s", proxy.Name)
	address := &core.Address{Address: &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Address: proxy.Backend.Host,
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: proxy.Backend.Port,
			},
		},
	}}

	var tlscontext *auth.UpstreamTlsContext
	if proxy.Backend.TLS {
		tlscontext = &auth.UpstreamTlsContext{
			CommonTlsContext: &auth.CommonTlsContext{
				TlsParams: &auth.TlsParameters{
					TlsMinimumProtocolVersion: auth.TlsParameters_TLSv1_2,
					TlsMaximumProtocolVersion: auth.TlsParameters_TLSv1_3,
				},
			},
			Sni: proxy.Backend.Host,
		}
	}

	clusterName := fmt.Sprintf("t_%s-p_%s", tenantName, proxy.Name)
	return &api.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       ptypes.DurationProto(5 * time.Second),
		ClusterDiscoveryType: &api.Cluster_Type{Type: api.Cluster_LOGICAL_DNS},
		DnsLookupFamily:      api.Cluster_V4_ONLY,
		RespectDnsTtl:        config.RespectDNSTTL,
		DnsRefreshRate:       ptypes.DurationProto(time.Duration(config.DNSRefreshRate) * time.Millisecond),
		LbPolicy:             api.Cluster_ROUND_ROBIN,
		LoadAssignment: &api.ClusterLoadAssignment{
			ClusterName: clusterName,
			Endpoints: []*endpoint.LocalityLbEndpoints{
				&endpoint.LocalityLbEndpoints{
					LbEndpoints: []*endpoint.LbEndpoint{
						&endpoint.LbEndpoint{
							HostIdentifier: &endpoint.LbEndpoint_Endpoint{
								Endpoint: &endpoint.Endpoint{
									Address: address,
								},
							},
						},
					},
				},
			},
		},
		TlsContext: tlscontext,
	}
}
