#!/usr/bin/env bash
set -xeuo pipefail

go build -i -o bin/incluster-broker ./cmd/broker
