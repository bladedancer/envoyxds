package authz

import (
auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
)
//Authorization interface for all security schemes
type Authorization interface {

Authorize() *auth.CheckResponse
}
