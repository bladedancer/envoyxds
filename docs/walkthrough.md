# POC Reference
While the code itself is self documenting the following information should help for those wanting to play around with the example

## What are we running
Helm charts have been configured to startup the entire example.

### Front

The front envoy proxy runs a standard Envoy image from 

### Back

### Redis
```
version: 9.5.2
appVersion: 5.0.5
```


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
[{"id":"5cd2ffad4c48940220ade42e","type":"ukulele","price":150,"currency":"EUR"},{"id":"5cd2fa168a4dde021faae3f8","type":"clarinet","price":650,"currency":"USD"},{"id":"5cd2fa041782ec021aad73ee","type":"drums","price":1750,"currency":"USD"},{"id":"5cd2f9f68a4dde0218addde0","type":"piano","price":1100,"currency":"EUR"},{"id":"5cd2f9ca08ac9a0219adee3d","type":"guitar","price":400,"currency":"GBP"}]
```

**What happened?**


**Invoke an API for a given tenant**
`curl --insecure -H "Host: test-0.bladedancer.dynu.net"  https://localhost:10000/route-0`

The above command will call the frontend 

**Invoke an API for a given tenant**

`curl --insecure -H "Host: test-0.bladedancer.dynu.net"  https://localhost:10000/route-0`
