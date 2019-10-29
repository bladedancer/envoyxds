package xdsconfig

import (
	"fmt"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	access_config "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v2"
	access_filter "github.com/envoyproxy/go-control-plane/envoy/config/filter/accesslog/v2"
	http_conn "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/conversion"
)

// GetListenerResources Get the listener configuration data
func GetListenerResources() []cache.Resource {
	config := makeListenerConfiguration()
	resources := []cache.Resource{config}
	return resources
}

// GetListener Get a test listener
func makeListenerConfiguration() *api.Listener {
	var filterChains []*listener.FilterChain

	accessLogStruct, _ := conversion.MessageToStruct(&access_config.FileAccessLog{
		Path: "/dev/stdout",
	})

	filterConfig := &http_conn.HttpConnectionManager{
		RouteSpecifier: &http_conn.HttpConnectionManager_Rds{
			Rds: &http_conn.Rds{
				RouteConfigName: "local_route",
				ConfigSource: &core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
						ApiConfigSource: &core.ApiConfigSource{
							ApiType: core.ApiConfigSource_GRPC,
							GrpcServices: []*core.GrpcService{
								&core.GrpcService{
									TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
										EnvoyGrpc: &core.GrpcService_EnvoyGrpc{
											ClusterName: "service_xds",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		AccessLog: []*access_filter.AccessLog{
			&access_filter.AccessLog{
				Name: "envoy.file_access_log",
				ConfigType: &access_filter.AccessLog_Config{
					Config: accessLogStruct,
				},
			},
		},
		StatPrefix: "ingress_http",
		HttpFilters: []*http_conn.HttpFilter{
			&http_conn.HttpFilter{
				Name: "envoy.router",
			},
		},
	}
	filterConfigStruct, _ := conversion.MessageToStruct(filterConfig)

	filter := &listener.Filter{
		Name: "envoy.http_connection_manager",
		ConfigType: &listener.Filter_Config{
			Config: filterConfigStruct,
		},
	}

	filterChains = append(filterChains, &listener.FilterChain{
		Filters: []*listener.Filter{filter},
		TlsContext: &auth.DownstreamTlsContext{
			CommonTlsContext: &auth.CommonTlsContext{
				TlsCertificates: []*auth.TlsCertificate{
					&auth.TlsCertificate{
						CertificateChain: &core.DataSource{
							Specifier: &core.DataSource_Filename{
								Filename: fmt.Sprintf("%s/certificate", config.CertPath),
							},
						},
						PrivateKey: &core.DataSource{
							Specifier: &core.DataSource_Filename{
								Filename: fmt.Sprintf("%s/privateKey", config.CertPath),
							},
						},
					},
				},
			},
		},
	})

	return &api.Listener{
		Name: fmt.Sprintf("listener_%d", 42 /*todo*/),
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Address: "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: 443,
					},
				},
			},
		},
		FilterChains: filterChains,
	}
}
