#!/usr/bin/env bash

set -xeuo pipefail

helm upgrade --install minibroker charts/minibroker \
    --namespace kube-system \
    --set imagePullPolicy=Never

#draft up
