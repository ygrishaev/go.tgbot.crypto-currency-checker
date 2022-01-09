# suppress output, run `make XXX V=` to be verbose
V := @

# Common
NAME = go.tgbot.crypto-currency-checker
VCS = gitlab.com
ORG = ygrishaev
VERSION := $(shell git describe --always --tags)
CURRENT_TIME := $(shell TZ="Europe/Moscow" date +"%d-%m-%y %T")

# Build
OUT_DIR = ./bin
MAIN_PKG = ./cmd/${NAME}
ACTION ?= build
GC_FLAGS = -gcflags 'all=-N -l'
LD_FLAGS = -ldflags "-s -v -w -X 'main.version=${VERSION}' -X 'main.buildTime=${CURRENT_TIME}'"
BUILD_CMD = CGO_ENABLED=1 go build -o ${OUT_DIR}/${NAME} ${LD_FLAGS} ${MAIN_PKG}

# Other
.DEFAULT_GOAL = build

.PHONY: build
build:
	@echo BUILDING PRODUCTION $(NAME)
	$(V)${BUILD_CMD}
	@echo DONE

.PHONY: lint
lint:
	$(V)golangci-lint run

run:
	CONFIG_PATH=configs/config.yaml go run cmd/go.tgbot.crypto-currency-checker/main.go

run-binary:
	CONFIG_PATH=configs/config.yaml ./bin/go.tgbot.crypto-currency-checker