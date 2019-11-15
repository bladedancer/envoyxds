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
	ext_authz "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/ext_authz/v2"
	"github.com/envoyproxy/go-control-plane/pkg/conversion"
	"github.com/golang/protobuf/ptypes"
	_struct "github.com/golang/protobuf/ptypes/struct"
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
	extAuthConfig := makeExtAuthzConfig(tenantName, proxy)

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
		PerFilterConfig: map[string]*_struct.Struct{
			"envoy.ext_authz": extAuthConfig,
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
		PerFilterConfig: map[string]*_struct.Struct{
			"envoy.ext_authz": extAuthConfig,
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
					EcdhCurves: []string{
						"P-256",
						"P-384",
						"P-521",
					},
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

// makeExtAuthzConfig Create the per route config overrides for external authorization.
func makeExtAuthzConfig(tenantName string, proxy *apimgmt.Proxy) *_struct.Struct {
	var config *ext_authz.ExtAuthzPerRoute
	callAuthz := false // Only need to call it if we have authorization

	// Common ExtAuthz params
	checkSettings := map[string]string{
		"tenant": tenantName,
		"proxy":  proxy.Name,
	}

	// Frontend authorization
	if proxy.Frontend.Authorization != nil && len(proxy.Frontend.Authorization) > 0 {
		// TODO: Handle scenario where we have multiple frontend auth profiles....
		// will it still be more efficient to store the data statically on the route
		// and pass it to the auth server?
		authorization := proxy.Frontend.Authorization[0]
		switch authorization.Type() {
		case apimgmt.AuthorizationTypePassthrough:
			// Noop
		case apimgmt.AuthorizationTypeAPIKey:
			typedAuth := authorization.(*apimgmt.APIKeyAuthorization)
			checkSettings["auth_type"] = (string)(typedAuth.Type())
			checkSettings["auth_in"] = typedAuth.Location
			checkSettings["auth_name"] = typedAuth.Name
			callAuthz = true
		case apimgmt.AuthorizationTypeHTTP:
			typedAuth := authorization.(*apimgmt.HTTPAuthorization)
			checkSettings["auth_type"] = (string)(typedAuth.Type())
			checkSettings["auth_scheme"] = typedAuth.Scheme
			callAuthz = true
		}
	}

	// Backend authorization
	if proxy.Backend.Authorization != nil {
		// This is very crude, in reality the user probably has a credential store and
		// as part of the proxy configuration they'll associate them with the proxy. Also
		// need to support multiple backend credentials.
		switch proxy.Backend.Authorization.Type() {
		case apimgmt.AuthorizationTypePassthrough:
			// Noop
		case apimgmt.AuthorizationTypeAPIKey:
			typedAuth := proxy.Backend.Authorization.(*apimgmt.APIKeyAuthorization)
			checkSettings["be_type"] = (string)(typedAuth.Type())
			checkSettings["be_in"] = typedAuth.Location
			checkSettings["be_name"] = typedAuth.Name
			callAuthz = true
		case apimgmt.AuthorizationTypeHTTP:
			typedAuth := proxy.Backend.Authorization.(*apimgmt.HTTPAuthorization)
			checkSettings["be_type"] = (string)(typedAuth.Type())
			checkSettings["be_scheme"] = typedAuth.Scheme
			callAuthz = true
		}
	}

	// Override the authorizaiton
	if callAuthz {
		config = &ext_authz.ExtAuthzPerRoute{
			Override: &ext_authz.ExtAuthzPerRoute_CheckSettings{
				CheckSettings: &ext_authz.CheckSettings{
					ContextExtensions: checkSettings,
				},
			},
		}
	} else {
		config = &ext_authz.ExtAuthzPerRoute{
			Override: &ext_authz.ExtAuthzPerRoute_Disabled{
				Disabled: true,
			},
		}
	}

	perRoute, _ := conversion.MessageToStruct(config)
	return perRoute
}
