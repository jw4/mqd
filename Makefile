ARTIFACT=mqd

RELEASE=false
VERSION_FILE=VERSION
BUILD_TAGS=-tags debug
ifeq ($(strip $(RELEASE)),true)
	BUILD_TAGS=
endif

ifeq ($(strip $(NO_REV)),)
	DIRTY=$(shell git diff-index --quiet HEAD || echo "-dirty")
	REV=$(shell git rev-parse --short HEAD)$(DIRTY)
endif

ifeq ($(BUILD_VERSION),)
	VERSION=$(shell cat $(VERSION_FILE))
	ifeq ($(strip $(REV)),)
		BUILD_VERSION=$(VERSION)
	else
		BUILD_VERSION=$(VERSION)-$(REV)
	endif
endif

ifeq ($(TARGET_GOOS),)
	TARGET_GOOS=$(shell go env GOOS)
endif

ifeq ($(TARGET_GOARCH),)
	TARGET_GOARCH=$(shell go env GOARCH)
endif

PWD=$(shell pwd)
BUILD_DIR="/go/src/github.com/jw4/mqd"
DOCKER_WRAPPER=docker run \
	-v "$(PWD)":"$(BUILD_DIR)" \
	-w "$(BUILD_DIR)" \
	-e CGO_ENABLED=0 \
	-e GOOS="$(TARGET_GOOS)" \
	-e GOARCH="$(TARGET_GOARCH)" \
	-e ARTIFACT="$(ARTIFACT)" \
	-e VERSION="$(BUILD_VERSION)" \
	-e BUILD_TAGS="$(BUILD_TAGS)" \
	--rm golang:alpine



all: compile ## Alias for compile.

compile: test ## Compile with bin/build.sh && test.
	$(DOCKER_WRAPPER) bin/build.sh

test: ## Run tests with bin/verify.bash
	bin/verify.bash

clean: ## Run go clean.
	-rm $(ARTIFACT)
	go clean ./...

help: ## Display help text.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: all compile test clean help
.DEFAULT_GOAL := help
