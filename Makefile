.DEFAULT_GOAL := help

DOCKER_HOST_PORT ?= 8080
DOCKER_pkg_PORT ?= 8080
DOCKER_HOST_ADDRESS ?= 0.0.0.0
DOCKER_NETWORK ?= esusu-net

PING_STATUS_OK ?= 200
GO_VERSION=1.22.1

# Pass args to the meme executable
MEME_ARGS ?=


##@ Build
build: meme ## build the binary

tidy:
	go mod tidy

build-light: 
	go build -ldflags="-X 'main.Version=dev'" -o meme ./cmd/api

meme: ./cmd/api format build-light

clean:  ## clean the build
	rm -f meme

#-------------End of Build section -------------------
##@ Run
run: meme ## Build as per tag and run the binary
	./meme $(MEME_ARGS)

# For faster repeated dev builds when you know the api hasn't changed
fast-run: ## Build the dev version and run the binary
	go build -ldflags="-X 'main.Version=dev'" -o meme ./cmd/api $(GO_BUILD_FLAGS)
	./meme $(MEME_ARGS)

build-docker:
	docker build . --build-arg GO_VERSION=$(GO_VERSION)

# The docker run arguments used here work with postgres and kafka running in all-in-1
run-docker: ## Run the latest docker image
	- docker rm -f meme 2> /dev/null
	docker run --name meme -dt -p $(DOCKER_HOST_ADDRESS):$(DOCKER_HOST_PORT):$(DOCKER_pkg_PORT) --network=$(DOCKER_NETWORK) -e MEME_DB_PG_HOST=postgres -e TOKEN_SIGNING_KEY="$$TOKEN_SIGNING_KEY" $(MEME_ARGS)

# ------------------ end of Run --------------------------

##@ Utility
linter: init-go ## Execute linters on codebase
	golangci-lint run --allow-parallel-runners --verbose --timeout=5m0s

format: init-go ## Format code using goimports and golines
	@echo "Formating .go files with golines"
	find . -name './*.go' -not -path "./vendor/*" -exec golines -m 100 --shorten-comments -t 4 -w {} \;

gosec: init-go ## Execute gosec safe-coding scan
	@echo "Running gosec..."
	gosec --quiet run --exclude-dir=vendor --timeout=3m0s ./...

check: tidy ## Check code quality before commit
	pre-commit run --all-files

cover: test ## Run the go test with coverage
	go tool cover -func coverage.out | grep 'total'

cover-html: test ## Run the go test with html coverage
	go tool cover -html coverage.out -o coverage.html

test-examine-postgres:
	psql -h localhost -U esusu -d postgres

# In the go test command below the ellipsis is a wildcard syntax that matches all sub-dirs. Also, -count=1 disables caching.
test-pkg: ## Check unit test with package level coverage
	go install github.com/mfridman/tparse@latest
	$(MAKE) tidy
	go test ./... -coverprofile coverage.out  -cover -json | tparse -all

test: test-pkg show_coverage ## Check unit test with package level and total coverage

# Run 'make init' the first time you are doing development
init: init-go ## Run this to install dependencies
	python3 -m pip install --upgrade pip
	python3 -m pip install pre-commit setuptools setuptools-rust
	# Install rust, required to install detect-secrets
	curl https://sh.rustup.rs -sSf | sh -s -- -y
	# shellcheck disable=SC1090,SC3046
	$(source ~/.cargo/env)
	pre-commit install
	pre-commit install-hooks

init-go: tidy ## Install golang dependencies and utilities
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.0
	go install github.com/securego/gosec/v2/cmd/gosec@v2.16.0
	go install golang.org/x/tools/cmd/goimports@v0.7.0
	go install github.com/segmentio/golines@v0.11.0
	go install github.com/mfridman/tparse@latest
	go install github.com/vektra/mockery/v2@v2.25.0
	go install go.uber.org/mock/mockgen@v0.4.0

infra_up: ## Start and stop postgres and kafka
	COMPOSE_PROFILES=complete docker-compose up -d
	docker ps


infra_down: 
	COMPOSE_PROFILES=complete docker-compose down -v
	docker container prune -f
	docker volume prune -f

show_coverage: ## Show total test coverage
	@echo "COVERAGE: "`go tool cover -func coverage.out | grep total | awk '{print $$3}'`


help:  ## Display this help.
ifeq ($(OS),Windows_NT)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-40s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
else
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-40s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
endif

mockgen: init-go ## mockgen generates `gomock` mocking functions in `mocks/`
	mockgen -source pkg/utils/subnetgroup/subnetgroup.go -destination mocks/mock_subnetgroup/mock_subnetgroup.go
	mockgen -source pkg/utils/imagemgr/image-mgr.go -destination mocks/mock_utils/mock_imagemgr.go

nancy: init-go ## Checks OSS Vulnerability
	go list -json -deps ./... | docker run --mount type=bind,source=$(shell pwd)/.nancy-ignore,target=/.nancy-ignore --rm -i sonatypecommunity/nancy:latest sleuth


.PHONY: build run fast-run docker check test run-docker init init-go
