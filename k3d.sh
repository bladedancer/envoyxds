#!/bin/bash
REPO=${REGISTRY:-bladedancer}
echo "using $REPO for image pull "
echo ======================
echo === Create Cluster ===
echo ======================
k3d delete --name xds
k3d create --server-arg '--no-deploy=servicelb' --server-arg '--no-deploy=traefik' --name xds --port 7443
sleep 15
export KUBECONFIG="$(k3d get-kubeconfig --name='xds')"
kubectl cluster-info

echo ======================
echo ===   Setup Helm   ===
echo ======================

cat > rbac.config << EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tiller
    namespace: kube-system
EOF

kubectl apply -f rbac.config
helm init --service-account tiller --history-max 200
helm repo up

echo ======================
echo === Install Images ===
echo ======================

docker pull $REPO/envoyxds:latest
docker pull $REPO/authz:latest
k3d i --name=xds  $REPO/envoyxds:latest
k3d i --name=xds  $REPO/authz:latest

echo ======================
echo ===    Wait     ===
echo ======================

while [[ $(kubectl -n kube-system get pods -l app=helm -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}') != "True" ]]; do echo "waiting for helm" && sleep 3; done

echo ======================
echo ===    Install     ===
echo ======================

helm dep update ./helm/saas
helm install --name=saas ./helm/saas
