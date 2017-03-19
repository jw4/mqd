#!/bin/sh

set -e

tags=${BUILD_TAGS:--tags debug}
appver=${VERSION:-0.0.1}

export CGO_ENABLED=0

go clean    ./...
go generate ${tags} -ldflags "-X main.appVersion=${appver}" ./cmd/smtp-dispatcher/...
go build -i ${tags} -ldflags "-X main.appVersion=${appver}" ./...
