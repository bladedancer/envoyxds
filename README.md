# Envoy XDS

This is a POC of an Envoy XDS control plane.

## Building
```
make all
```

## Installing
If you're ushing k3d you can use the k3d.sh.

```
. ./k3d.sh
```

## Testing
kubectl port-forward svc/envoy 10000:10000
