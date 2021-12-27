# suppress output, run `make XXX V=` to be verbose
V := @

# Common
NAME = go.tgbot.crypto-currency-checker

.PHONY: lint
lint:
	$(V)golangci-lint run