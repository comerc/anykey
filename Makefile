SHELL := /bin/bash

.DEFAULT_GOAL := all
GO ?= go

JSONF_EXAMPLE := '[{"id":2,"name":"Alice","age":30},{"age":25,"name":"Bob","id":1}]'

.PHONY: all jsonf limiter test

all: jsonf limiter

jsonf:
	@echo '== jsonf example =='
	@printf '%s\n' $(JSONF_EXAMPLE) | $(GO) run ./cmd/jsonf name age

limiter:
	@echo '== limiter example =='
	@$(GO) run ./cmd/limiter

test:
	@$(GO) test -count=1 ./...



