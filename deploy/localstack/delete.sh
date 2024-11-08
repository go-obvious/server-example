#!/usr/bin/env bash
set -euo pipefail

helm uninstall localstack --namespace=localstack
