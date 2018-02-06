FROM golang:1.9

ENV USE_DOCKER=false

WORKDIR /tmp

# Install Helm
RUN curl -sLO https://kubernetes-helm.storage.googleapis.com/helm-v2.8.0-linux-amd64.tar.gz
RUN tar -zxvf helm-v2.8.0-linux-amd64.tar.gz
RUN mv linux-amd64/helm /usr/local/bin/

# Install Draft
RUN curl -sLO https://azuredraft.blob.core.windows.net/draft/draft-v0.10.0-linux-amd64.tar.gz
RUN tar -zxvf draft-v0.10.0-linux-amd64.tar.gz
RUN mv linux-amd64/draft /usr/local/bin/

# Install Minikube
RUN curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
RUN chmod +x minikube
RUN mv minikube /usr/local/bin/

# Install svcat
RUN go get github.com/kubernetes-incubator/service-catalog/cmd/svcat

WORKDIR /go/src/github.com/deis/minibroker
