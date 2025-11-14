#!/bin/bash
# Add the Jetstack Helm repository
helm repo add jetstack https://charts.jetstack.io --force-update

# Install the cert-manager helm chart
helm upgrade --install \
    cert-manager jetstack/cert-manager \
    --namespace cert-manager \
    --create-namespace \
    --version v1.19.1 \
    --set crds.enabled=true

kubectl apply -f certs.yaml

helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm upgrade --install \
    nats nats/nats \
    --namespace nats \
    --create-namespace \
    -f nats-values.yaml
