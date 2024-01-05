REGISTRY ?= harbor.infra.cluster.ionos.com
REGISTRY_PROJECT ?= $(USER)/simpleapp
VERSION ?= dev
IMAGE_FRONTEND := $(REGISTRY)/$(REGISTRY_PROJECT)/frontend:$(VERSION)
IMAGE_BACKEND := $(REGISTRY)/$(REGISTRY_PROJECT)/backend:$(VERSION)
BUILD_DATE ?= $(shell date -Iseconds)
GITHUB_SHA ?= $(shell git rev-parse HEAD 2>/dev/null)
GOOS ?= $(shell go env GOOS 2>/dev/null)
LDFLAGS := '-extldflags "-static" \
	-X github.com/dsteininger86/go-appinfo.version=$(VERSION) \
	-X github.com/dsteininger86/go-appinfo.buildDate=$(BUILD_DATE) \
	-X github.com/dsteininger86/go-appinfo.gitCommit=$(GITHUB_SHA) \
	-X github.com/dsteininger86/go-appinfo.name=kw-token-cleanup \
	'

.PHONY: compile
compile: compile_backend compile_frontend

.PHONY: test
test: lint

.PHONY: compile_backend
compile_backend:
	CGO_ENABLED=0 GOOS="$(GOOS)" \
	go build -trimpath -ldflags $(LDFLAGS) -o bin/backend cmd/backend/main.go

.PHONY: compile_frontend
compile_frontend:
	CGO_ENABLED=0 GOOS="$(GOOS)" \
	go build -trimpath -ldflags $(LDFLAGS) -o bin/frontend cmd/frontend/main.go

.PHONY: docker_compose_build
docker_compose_build: compile
	IMAGE_FRONTEND="$(IMAGE_FRONTEND)" \
	IMAGE_BACKEND="$(IMAGE_BACKEND)" \
	docker compose build

.PHONY: docker_compose_push
docker_compose_push: docker_compose_build
	IMAGE_FRONTEND="$(IMAGE_FRONTEND)" \
	IMAGE_BACKEND="$(IMAGE_BACKEND)" \
	docker compose push

.PHONY: docker_compose_up
docker_compose_up: docker_compose_build
	IMAGE_FRONTEND="$(IMAGE_FRONTEND)" \
	IMAGE_BACKEND="$(IMAGE_BACKEND)" \
	docker compose up -d

.PHONY: docker_compose_down
docker_compose_down:
	IMAGE_FRONTEND="$(IMAGE_FRONTEND)" \
	IMAGE_BACKEND="$(IMAGE_BACKEND)" \
	docker compose down

.PHONY: docker_compose_rm
docker_compose_rm:
	IMAGE_FRONTEND="$(IMAGE_FRONTEND)" \
	IMAGE_BACKEND="$(IMAGE_BACKEND)" \
	docker compose rm


.PHONY: lint
lint:
	golangci-lint run

.PHONY: protoc
protoc:
	protoc \
	--go_out=envlookup --go_opt=paths=source_relative \
	--go-grpc_out=envlookup --go-grpc_opt=paths=source_relative \
	envlookup.proto

.PHONY: clean
clean: docker_compose_down docker_compose_rm
	rm -rf bin/
