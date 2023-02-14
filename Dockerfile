############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

#FROM golang:1.17.2-stretch AS builder
# Install git.
# Git is required for fetching the dependencies.

WORKDIR $GOPATH/src/mypackage/
COPY . ./helm-api

WORKDIR $GOPATH/src/mypackage/helm-api/
# Fetch dependencies.
# Using go get.
RUN apk update \                                                                                                                                                                                                                        
  && apk add ca-certificates zip wget tar curl gcc musl-dev linux-headers\                                                                                                                                                                                                      
  && update-ca-certificates

RUN go install
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0  GOARCH=386 GOOS=linux go build -o /go/bin/helm-api
RUN cd /go/bin/
RUN mkdir /go/bin/charts
RUN mkdir /go/bin/config
RUN mkdir /go/bin/kubeconfig
WORKDIR /go/bin/

############################
# STEP 2 build a small image
############################
#FROM ubuntu 
FROM gcr.io/distroless/static
#FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/helm-api /go/bin/helm-api

EXPOSE 8000

VOLUME [ "/go/bin/charts" ]

LABEL version="1.0.0"
LABEL name="helm-api"
LABEL maintainer="Andrew Pye"
LABEL description="Helm-api is a CLI application written in Golang that gives the ability to perform Install, Uninstall and Upgrade of Helm Charts via Rest API endpoint."

ENV BASE_FOLDER=/charts
ENV WEB_IP=localhost
ENV WEB_PORT=8080
ENV WEB_CONFIG_PATH=/go/bin/kubeconfig/

ENTRYPOINT ["/go/bin/helm-api","web"]