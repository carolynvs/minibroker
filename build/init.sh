#!/usr/bin/env bash

set -xeuo pipefail

helm init
sleep 30
helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
helm install svc-cat/catalog --name catalog --namespace kube-system

draft init --auto-accept
