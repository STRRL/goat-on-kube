# Goat Exporter - Prometheus exporter for GOAT blockchain metrics

```
helm upgrade --install goat-exporter \
  ./helm/goat-exporter \
  --namespace goat-exporter \
  --create-namespace \
  --set secret.goatRpcNode='https://rpc.goat.network'
```
