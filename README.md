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

### install goat-exporter

```bash
helm upgrade --install goat-exporter \
  ./helm/goat-exporter \
  --namespace goat-exporter \
  --create-namespace \
  --set secret.goatRpcNode='https://rpc.goat.network' \
  --set serviceMonitor.enabled='true'
```

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
