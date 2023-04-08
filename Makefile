BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git rev-parse --short HEAD)

default:
	@echo "Cmds: [build | run]"

test:
	@go test ./...

build:
	@go build -o bin/ksp-$(BRANCH)-$(HASH)

run:
	@air
