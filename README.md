# Envoy XDS

This is a POC of an Envoy XDS control plane.


## Building
```
make all
```

## Installing
If you're ushing k3d you can use the k3d.sh.

Install k3d:
```
wget -q -O - https://raw.githubusercontent.com/rancher/k3d/master/install.sh | bash
```

Build and run:
```
make clean docker-build && . k3d.sh
```

When you edit the helm chart you don't need to rebuild again, you can just use helm directly.

```
helm delete --purge xds
helm install --name=xds ./helm
```

## Testing
The service is deployed as a loadbalancer but as k3d isn't run with external dns there is no public IP assigned so you have to port forward manually.

Envoy is configured with domain match rules - you can override this domain when deploying:

```
helm install --name=xds ./helm --set xds.domain=mydomain.com
```

The only requirement is that *.mydomain.com will resolve to the server run k3d.

```
kubectl port-forward --address 0.0.0.0 svc/envoy 443:443
kubectl port-forward --address 0.0.0.0 deployment/xds-envoy 9901
```

Then configure your firewall to route incoming requests to that machine, such that urls like: `https://test-9.bladedancer.dynu.net/route-1` resolve.
