# POC Reference
While the code itself is self documenting the following information should help for those wanting to dig deeper and potentially make modifications to the example

## High Level Architecture

<img src="https://github.com/bladedancer/envoyxds/raw/master/docs/envoy-shard-2.png" width="800">

## What are we running

[Helm charts](../helm/saas/requirements.yaml) have been configured to startup the entire example. 

### Redis

### Postgres

### XDS

The [XDS Service](../pkg/xds/envoyxds.go) is responsible for providing configuration to envoys dynamically. Upon startup an envoy's 
configuration will specify dynamic resource sections and indicate that these dependencies are 
to be fulfilled by an XDS service.

```json
{
     "lds_config": {
      "api_config_source": {
       "api_type": "GRPC",
       "grpc_services": [
        {
         "envoy_grpc": {
          "cluster_name": "service_xds"
         }
        }
       ]
      }
     },
     "cds_config": {
      "api_config_source": {
       "api_type": "GRPC",
       "grpc_services": [
        {
         "envoy_grpc": {
          "cluster_name": "service_xds"
         }
        }
       ]
      }
     }
}
```

> The XDS service is defined as a [CLI](../cmd/xds/cmd/cmd.go) using Cobra and runs just the same as any other defined service. (Using Cobra allows for easier local debug sessions and clear understanding of parameters )
> Dependencies on [Redis](#Redis) and [Postgres](#Postgres) 

### Authz

The [Authz Service](../pkg/authz/authz.go) is responsible for providing an external authentication point for an envoy proxy. 
 configuration to envoys dynamically. Upon startup an envoy's 
configuration will specify dynamic resource sections and indicate that these dependencies are 
to be fulfilled by an XDS service.


### Frontend Envoy

The frontend envoy acts as sort of loadbalancer to determine the appropriate backend cluster to allow requests to travel upstream. In addition, Authorization schemes, such as API-KEY, can be applied here to short circuit invalid requests

> Dependencies on [XDS](#XDS) and [AUTHZ](#Authz) services


### Backenend Envoy

The backenend envoy is seperated into different partitions, called shards, which allows for multiple configurations to be applied by tenant, while still achieving the goal of horizontal scale. In addition, here it is possible to apply specific credential requiements to authenticate with backend resources, using OAuth or other similar mechanisms.

> Dependencies on [XDS](#XDS) and [AUTHZ](#Authz) services



The Redis Cache is used by the Authz service in order to perform some simple Authentication. When a frontend request is received it is delegated to the Authz service where the service will examine the header values to validate the appropriate API Key is present. The Cache could be used for other items such as rate limiting


###  Launching the example

**Prerequisites**

Most folks are likely to already have the dependencies installed, but...

- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Helm](https://helm.sh/docs/intro/install/)
- [k3d](hhttps://github.com/rancher/k3d)

**Launch the example**

`. ./k3d.sh`

The startup sequence should take several minutes, but...

`kubectl get pods`

Should eventually yield the following results

```
NAME                     READY   STATUS    RESTARTS   AGE
saas-redis-master-0      1/1     Running   0          13m
saas-postgresql-0        1/1     Running   0          13m
authz-57fc7dd7bc-sqjxc   1/1     Running   0          13m
saas-redis-slave-0       1/1     Running   0          13m
xds-5c99dcf48d-z5wbm     1/1     Running   0          13m
back-2                   1/1     Running   0          13m
saas-redis-slave-1       1/1     Running   0          13m
back-0                   1/1     Running   0          13m
back-1                   1/1     Running   0          13m
front-585987cfd6-bwq69   1/1     Running   0          13m
front-585987cfd6-6q54z   1/1     Running   0          13m

```

### Build the Example: \(_Optional_\)

If anyone should want to tweak the example on their own...

**Prerequisites**

- [gRPC](https://grpc.io/docs/quickstart/go/)
- **make** \-There are many different flavors of make and most OS should already have this installed

**Building**

By default, the images are pulled from the bladedancer registry. In order to build your own versions, you will need to make sure to change the Docker Registry to match your registry. For Helm, you will need to make sure to edit the [values file](../helm/saas/values.yaml)

```
export REGISTRY=<dockerhub registry name eg: bladedancer>
make all
. ./k3d.sh
```


### Running the Example

**Enable port forwarding**

`kubectl port-forward --address 0.0.0.0  svc/front 10000:443`

This will forward all requests on your host machine to the frontend service on port 443

\(_**Optionally**_\) setup DNS entry to resolve _bladedancer.dynu.net_ to host

If you are not ambitious enough to setup the DNS resolution, you can add the Host header to your requests and achieve the same results. The following examples are setup using this approach

**Invoke an API secured by Authz**

```
curl -v --insecure -H "Host: test-gavin.bladedancer.dynu.net" -H "key: password" https://localhost:10000/api/music/instruments
```

> The above command should yield the following...

```json
[
  {
    "id": "5cd2ffad4c48940220ade42e",
    "type": "ukulele",
    "price": 150,
    "currency": "EUR"
  },
  {
    "id": "5cd2fa168a4dde021faae3f8",
    "type": "clarinet",
    "price": 650,
    "currency": "USD"
  },
  {
    "id": "5cd2fa041782ec021aad73ee",
    "type": "drums",
    "price": 1750,
    "currency": "USD"
  },
  {
    "id": "5cd2f9f68a4dde0218addde0",
    "type": "piano",
    "price": 1100,
    "currency": "EUR"
  },
  {
    "id": "5cd2f9ca08ac9a0219adee3d",
    "type": "guitar",
    "price": 400,
    "currency": "GBP"
  }
]
```

#### What happened?

<img src="https://github.com/bladedancer/envoyxds/raw/master/docs/music.png" width="800">

##### Frontend Envoy
**Listen to all traffic on 443**

```json
{
    "name": "listener_42",
    "address": {
      "socket_address": {
        "address": "0.0.0.0",
        "port_value": 443
      }
    }
}
```

**Apply the Filters**

````json
[
  {
    "name": "envoy.lua",
    "config": {
      "inline_code": "\n    function envoy_on_request(request_handle)\n       for key, value in pairs(request_handle:headers()) do\n          request_handle:logInfo(key .. \": \" .. value)\n       end\n       local headers, body = request_handle:httpCall(\n         \"service_xds_shard\",\n         {\n          [\":method\"] = \"GET\",\n          [\":path\"] = \"\/shard\",\n          [\":authority\"] = \"service_xds_shard\"\n        },\n        request_handle:headers():get(\":authority\") .. \":\" .. request_handle:headers():get(\":path\") ,\n        5000)\n      request_handle:logInfo(\"Adding Shard via Lua \" .. body)\n      request_handle:headers():add(\"x-shard\", body)\n    end"
    }
  },
  {
    "name": "envoy.router"
  }
]
````
The above filter chain will invoke the Lua Filter, followed by the built-in envoy.router.  The Lua filter is defined with inline_code and will call a simple service defined [within this module](../pkg/xds/tennantroute.go) 


**Perform the Route**

`request_handle:headers():add(\"x-shard\", body)`

The Lua Filter has now created a header entry called x-shard, which will now be used to perform the route from the frontend to the appropriate backend cluster
```json
[
  {
    "match": {
      "prefix": "/"
    },
    "route": {
      "cluster_header": "x-shard"
    },
    "name": "front"
  }
]
```

##### Backend Envoy

**Listen to all traffic on 80**

```json
{
    "name": "listener_42",
    "address": {
      "socket_address": {
        "address": "0.0.0.0",
        "port_value": 80
      }
    }
}
```

**Apply the Filters**

````json
[
                  {
                    "name": "envoy.ext_authz",
                    "config": {
                      "grpc_service": {
                        "envoy_grpc": {
                          "cluster_name": "service_authz"
                        },
                        "timeout": "5s"
                      }
                    }
                  },
                  {
                    "name": "envoy.router"
                  }
]
````
The above filter chain will invoke the [authz](../pkg/authz/authz.go) grpc service, followed by the built-in evnoy.router filter. This configuration is not intuitive, because it actually pulls config from the route filter using _"per_filter_config"_ and passes this to the authz service.


**Perform the Route**

The Authz filter has now responded with the authoriztion decision. Most commonly a **200** for _authorized_ or **403** for _forbidden_. If authorized the request will continue upstream to the Music API, but if forbidden it will reject back to the requester
```json
{
  "match": {
    "prefix": "/api/music/"
  },
  "route": {
    "cluster": "t_gavin-p_3",
    "prefix_rewrite": "/envoy/music/v2/",
    "host_rewrite": "prod-e4ec6c3369cdafa50169ce18e33d00bb.apicentral.axwayamplify.com"
  },
  "per_filter_config": {
    "envoy.ext_authz": {
      "check_settings": {
        "context_extensions": {
          "auth_type": "apiKey",
          "auth_in": "header",
          "tenant": "gavin",
          "proxy": "3",
          "auth_name": "key"
        }
      }
    }
  },
  "name": "t_gavin-p_3-prefix"
}
```
