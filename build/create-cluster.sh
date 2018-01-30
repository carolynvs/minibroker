#!/usr/bin/env bash

set -xeuo pipefail

if ! type minikube &> /dev/null; then
    echo You must install minikube
    exit 1
fi

if ! type kubectl &> /dev/null; then
    echo You must install kubectl
    exit 1
fi


minikube start --profile icb --extra-config=apiserver.Authorization.Mode=RBAC
minikube addons enable registry --profile icb
minikube addons enable ingress --profile icb

kubectl create clusterrolebinding cluster-admin:kube-system \
    --clusterrole=cluster-admin --serviceaccount=kube-system:default
