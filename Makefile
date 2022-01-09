# suppress output, run `make XXX V=` to be verbose
V := @

# Common
NAME = go.tgbot.crypto-currency-checker

.PHONY: lint
lint:
	$(V)golangci-lint run

run:
	CONFIG_PATH=configs/config.yaml go run cmd/go.tgbot.crypto-currency-checker/main.go