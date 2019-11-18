package authz

import (
	"fmt"

	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
)

const apiScheme = "apikey"

type apikey struct {
	key string
}

//NewAPIKey returns new authorization interface
func NewAPIKey(req *auth.CheckRequest) Authorization {

	//    {"level":"info","msg":"Request Headers: map[content-length:0 x-shard:back-0 :method:GET x-forwarded-proto:https accept:*/* user-agent:curl/7.63.0 :path:/api/music/instruments :authority:test-gavin.bladedancer.dynu.net x-request-id:04bfdf64-9fba-43c7-9784-a3526b07673b]","package":"authz","time":"2019-11-15T20:08:58Z"}
	//    {"level":"info","msg":"Query String: ","package":"authz","time":"2019-11-15T20:08:58Z"}
	//    {"level":"info","msg":"Context Extensions: map[auth_type:apiKey auth_in:query tenant:gavin proxy:3 auth_name:key]","package":"authz","time":"2019-11-15T20:08:58Z"}

	hdrs := req.GetAttributes().GetRequest().GetHttp().GetHeaders()
	query := req.GetAttributes().GetRequest().GetHttp().GetQuery()
	ext := req.GetAttributes().GetContextExtensions()
	//Arbitrary results, so far
	key := "unkonwn"
	// query is empty via curl at the moment
	if len(hdrs) > 0 {
		//todo
	}
	if len(query) > 0 {
		//todo
	}
	// This is where the crux lies for now
	if len(ext) > 0 {
		key = fmt.Sprintf("%s:%s:%s:", ext["tenant"], ext["proxy"], apiScheme)
	}

	return &apikey{key: key}
}

//Authorize *auth.CheckResponse
func (a *apikey) Authorize() bool {
	//TODO Check Redis Cache
	return true
}
