FROM alpine:3.7

RUN apk add --no-cache --update -t deps curl ca-certificates bash

ENV HELM_VERSION="v2.8.1"
ENV FILENAME="helm-${HELM_VERSION}-linux-amd64.tar.gz"

RUN curl -L http://storage.googleapis.com/kubernetes-helm/${FILENAME} -o /tmp/${FILENAME} && \
    tar -zxvf /tmp/${FILENAME} -C /tmp && \
    mv /tmp/linux-amd64/helm /usr/local/bin/

COPY servicebroker /usr/local/bin/

RUN adduser -D minibroker
USER minibroker

RUN helm init --client-only

CMD ["servicebroker", "--help"]
