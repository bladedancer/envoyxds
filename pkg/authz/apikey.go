package authz

import (
	"fmt"

	"context"

	"github.com/bladedancer/envoyxds/pkg/cache"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
)

const apiScheme = "apiKey"

type apikey struct {
	hdrs     map[string]string
	query    string
	authIn   string
	authName string
	auth     *cache.AuthEnvelope
}

//NewAPIKey returns new authorization interface
func NewAPIKey(req *auth.CheckRequest) Authorization {

	//    {"level":"info","msg":"Request Headers: map[content-length:0 x-shard:back-0 :method:GET x-forwarded-proto:https accept:*/* user-agent:curl/7.63.0 :path:/api/music/instruments :authority:test-gavin.bladedancer.dynu.net x-request-id:04bfdf64-9fba-43c7-9784-a3526b07673b]","package":"authz","time":"2019-11-15T20:08:58Z"}
	//    {"level":"info","msg":"Query String: ","package":"authz","time":"2019-11-15T20:08:58Z"}
	//    {"level":"info","msg":"Context Extensions: map[auth_type:apiKey auth_in:query tenant:gavin proxy:3 auth_name:key]","package":"authz","time":"2019-11-15T20:08:58Z"}

	hdrs := req.GetAttributes().GetRequest().GetHttp().GetHeaders()
	query := req.GetAttributes().GetRequest().GetHttp().GetQuery()
	ext := req.GetAttributes().GetContextExtensions()
	key := fmt.Sprintf("%s-%s-%s", ext["tenant"], ext["proxy"], apiScheme)
	authIn := ext["auth_in"]
	authName := ext["auth_name"]
	envelope := &cache.AuthEnvelope{}
	log.Infof("Looking in the Cache %s", key)
	c.Get(context.Background(), key, envelope, true)
	log.Infof("Result from Cache %v", envelope)

	return &apikey{hdrs: hdrs, query: query, authIn: authIn, authName: authName, auth: envelope}
}
func determineExtractKey(authIn string, authName string, hdrs map[string]string, query string) string {
	apiKey := ""
	switch authIn {
	case "header":
		apiKey = hdrs[authName]
		log.Infof("%s = %s ", authName, apiKey)
	case "query":
		//Not Implemented
	}
	return apiKey
}

//Authorize *auth.CheckResponse
func (a *apikey) Authorize() bool {
	apiKey := determineExtractKey(a.authIn, a.authName, a.hdrs, a.query)
	apiCtx := getAPIContext(a.auth)
	log.Infof("Comparing Key %s from Header to cached Key %s result %v", apiKey, apiCtx.Key, apiKey == apiCtx.Key)
	return apiKey == apiCtx.Key
}
