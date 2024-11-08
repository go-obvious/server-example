#!/usr/bin/env bash
set -euo pipefail

helm repo add localstack https://localstack.github.io/helm-charts
helm repo update
helm upgrade --install localstack localstack/localstack --create-namespace --namespace=localstack -f localstack-values.yaml
