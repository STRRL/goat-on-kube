# Goat on Kube

This is a monorepo which contains

- goat-exporter, a prometheus exporter for goat network rpc node
- helm chart to deploy goat-exporter
- helm chart to deploy goat rpc node

## Get Started

### setup kubernetes cluster

```bash
minikube start
minikube addons enable ingress
minikbue addons enable storage-provisioner
minikube addons enable default-storageclass
```

### setup prometheus / grafana stack

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install kube-prometheus-stack \
  prometheus-community/kube-prometheus-stack \
  --create-namespace \
  --namespace monitoring \
  --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false
```

### install goat rpc node by helm

```bash
helm upgrade --install goat \
  ./helm/goat \
  --namespace goat \
  --create-namespace \
  --set ingress.enabled=true
```

**Access via Ingress** (requires `minikube tunnel`):
- Geth HTTP-RPC API: `http://geth.127.0.0.1.nip.io`
- GOAT REST API: `http://goat.127.0.0.1.nip.io`


### install goat-exporter by helm

```bash
helm upgrade --install goat-exporter \
  ./helm/goat-exporter \
  --namespace goat-exporter \
  --create-namespace \
  --set secret.goatRpcNode='http://goat.goat.svc.cluster.local:8545' \
  --set serviceMonitor.enabled='true' \
  --set ingress.enabled='true'
```

Access goat-exporter at `http://goat-exporter.127.0.0.1.nip.io/metrics` (requires `minikube tunnel`)

### access the grafana

```bash
# create an adhoc ingress
kubectl create ingress grafana -n monitoring --class=nginx --rule="grafana.127.0.0.1.nip.io/*=kube-prometheus-stack-grafana:80"
# need sudo on macOS for binding on 80
sudo minikube tunnel
```

Access Grafana at `http://grafana.127.0.0.1.nip.io`

Get the default password:

```bash
kubectl get secret -n monitoring kube-prometheus-stack-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```
