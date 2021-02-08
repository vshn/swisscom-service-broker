MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

TESTDATA_DIR ?= ./testdata
TESTBIN_DIR ?= $(TESTDATA_DIR)/bin
KIND_BIN ?= $(TESTBIN_DIR)/kind
KIND_VERSION ?= 0.9.0
KIND_KUBECONFIG ?= $(TESTBIN_DIR)/kind-kubeconfig
KIND_NODE_VERSION ?= v1.18.8
KIND_CLUSTER ?= crossplane-service-broker
KIND_REGISTRY_NAME ?= kind-registry
KIND_REGISTRY_PORT ?= 5000

# Needs absolute path to setup env variables correctly.
ENVTEST_ASSETS_DIR = $(shell pwd)/testdata

DOCKER_CMD   ?= docker
DOCKER_ARGS  ?= --rm --user "$$(id -u)" --volume "$${PWD}:/src" --workdir /src

# Project parameters
BINARY_NAME ?= swisscom-service-broker

VERSION ?= $(shell git describe --tags --always --dirty --match=v* || (echo "command failed $$?"; exit 1))

IMAGE_NAME ?= docker.io/vshn/$(BINARY_NAME):$(VERSION)

ANTORA_PREVIEW_CMD ?= $(DOCKER_CMD) run --rm --publish 35729:35729 --publish 2020:2020 --volume "${PWD}":/preview/antora vshn/antora-preview:2.3.4 --style=syn --antora=docs

# Linting parameters
YAML_FILES      ?= $(shell git ls-files *.y*ml)
YAMLLINT_ARGS   ?= --no-warnings
YAMLLINT_CONFIG ?= .yamllint.yml
YAMLLINT_IMAGE  ?= docker.io/cytopia/yamllint:latest
YAMLLINT_DOCKER ?= $(DOCKER_CMD) run $(DOCKER_ARGS) $(YAMLLINT_IMAGE)

TESTDATA_CRD_DIR = $(TESTDATA_DIR)/crds
CROSSPLANE_VERSION = v1.0.0
CROSSPLANE_CRDS = $(addprefix $(TESTDATA_CRD_DIR)/, apiextensions.crossplane.io_compositeresourcedefinitions.yaml \
					apiextensions.crossplane.io_compositions.yaml \
					pkg.crossplane.io_configurationrevisions.yaml \
					pkg.crossplane.io_configurations.yaml \
					pkg.crossplane.io_controllerconfigs.yaml \
					pkg.crossplane.io_locks.yaml \
					pkg.crossplane.io_providerrevisions.yaml \
					pkg.crossplane.io_providers.yaml)

# Go parameters
GOCMD   ?= go
GOBUILD ?= $(GOCMD) build
GOCLEAN ?= $(GOCMD) clean
GOTEST  ?= $(GOCMD) test
GOGET   ?= $(GOCMD) get

BUILD_CMD ?= CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -v \
				-o $(BINARY_NAME) \
				-ldflags "-X main.Version=$(VERSION) -X 'main.BuildDate=$(shell date)'" \
				cmd/swisscom-service-broker/main.go

.PHONY: all
all: lint test build

.PHONY: build
build:
	$(BUILD_CMD)
	@echo built '$(VERSION)'

.PHONY: test
test:
	$(GOTEST) -v -cover ./...

.PHONY: run
run:
	go run cmd/swisscom-service-broker/main.go

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: docker
docker:
	DOCKER_BUILDKIT=1 docker build -t $(IMAGE_NAME) --build-arg VERSION="$(VERSION)" .
	@echo built image $(IMAGE_NAME)

.PHONY: lint
lint: fmt vet lint_yaml
	@echo 'Check for uncommitted changes ...'
	git diff --exit-code

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint_yaml
lint_yaml: $(YAML_FILES)
	$(YAMLLINT_DOCKER) -f parsable -c $(YAMLLINT_CONFIG) $(YAMLLINT_ARGS) -- $?

.PHONY: docs-serve
docs-serve:
	$(ANTORA_PREVIEW_CMD)

$(TESTBIN_DIR):
	mkdir $(TESTBIN_DIR)

# TODO(mw): something with this target is off, $@ should be used instead of $*.yaml but I can't seem to make it work.
$(TESTDATA_CRD_DIR)/%.yaml:
	curl -sSLo $@ https://raw.githubusercontent.com/crossplane/crossplane/$(CROSSPLANE_VERSION)/cluster/charts/crossplane/crds/$*.yaml

.PHONY: integration_test
integration_test: $(CROSSPLANE_CRDS)
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test -tags=integration -v ./... -coverprofile cover.out
