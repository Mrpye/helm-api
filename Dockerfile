############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

#FROM golang:1.17.2-stretch AS builder
# Install git.
# Git is required for fetching the dependencies.

WORKDIR $GOPATH/src/mypackage/
COPY . ./cimpex

WORKDIR $GOPATH/src/mypackage/cimpex/
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

ENV BASE_FOLDER=/go/bin/charts
ENV WEB_IP=localhost
ENV WEB_PORT=8080

ENTRYPOINT ["/go/bin/helm-api","web"]