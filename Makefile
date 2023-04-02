BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git rev-parse --short HEAD)

default:
	@echo "Cmds: [build | run]"

build:
	go build -o bin/asc-$(BRANCH)-$(HASH)
