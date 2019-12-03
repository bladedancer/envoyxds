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


the Redis Cache is used by the Authz service in order to perform some simple Authentication. When a frontend request is received it is delegated to the Authz service where the service will examine the header values to validate the appropriate API Key is present

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

### Prerequisites to Run the example
Most folks are likely to already have the dependencies installed, but...

- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Helm](https://helm.sh/docs/intro/install/)
- [k3d](hhttps://github.com/rancher/k3d)

### Prerequisites to Build
- gRPC

`curl -v --insecure -H "Host: test-gavin.bladedancer.dynu.net" -H "key: password" https://localhost:10000/api/music/instruments`

`curl --insecure -H "Host: test-0.bladedancer.dynu.net"  https://localhost:10000/route-0`
