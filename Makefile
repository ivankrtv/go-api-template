SERVICE_NAME=template

GOLANGCI_LINT_VERSION=v2.14-alpine

GOCACHE_HOST := $(CURDIR)/.cache/go-build
GOLANGCI_CACHE_HOST := $(CURDIR)/.cache/golangci-lint

# Generators
generate_oapi:
	go tool oapi-codegen -config=./internal/config/oapi-server-config.yaml docs/openapi.json
	go tool oapi-codegen -config=./internal/config/oapi-types-config.yaml docs/openapi.json

generate_mocks:
	docker run --rm -v $(CURDIR):/src -w /src vektra/mockery:3

# linters and tests
lint:
	docker run --rm -t -v $(CURDIR):/app -w /app \
		-v $(GOCACHE_HOST):/root/.cache/go-build \
		-v $(GOLANGCI_CACHE_HOST):/root/.cache/golangci-lint \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run

EXCLUDE_TEST_DIRS := test|cmd|gen|docs|deploy|migrations

unit_tests:
	go test -v -race $$(go list ./internal/service/... | grep -vE "$(EXCLUDE_TEST_DIRS)") -coverprofile=".coverage.out"

all_tests:
	set -a \
	. ./deploy/test/.env.test \
	set +a
	docker compose -f deploy/test/docker-compose.yaml -p $(SERVICE_NAME)_test down
	docker compose -f deploy/test/docker-compose.yaml -p $(SERVICE_NAME)_test up -d
	go clean --testcache
	go test -v -race $$(go list ./internal/service/... | grep -vE "$(EXCLUDE_TEST_DIRS)") -coverprofile=".coverage.out"
	docker compose -f deploy/test/docker-compose.yaml -p $(SERVICE_NAME)_test down

# Infrastructure
infra_up:
	docker-compose -f deploy/local/docker-compose.yaml -p $(SERVICE_NAME) up  -d

infra_down:
	docker-compose -f deploy/local/docker-compose.yaml -p $(SERVICE_NAME) down

# Test infrastructure
test_infra_up:
	docker-compose -f deploy/test/docker-compose.yaml -p $(SERVICE_NAME)_test --env-file deploy/test/.env.test up -d

test_infra_down:
	docker-compose -f deploy/test/docker-compose.yaml -p $(SERVICE_NAME)_test down