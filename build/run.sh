#!/usr/bin/env bash

set -xeuo pipefail

kubectl apply -f chart/broker.yaml

pushd cmd/broker
trap popd EXIT
draft up
