SHELL := /bin/bash

.DEFAULT_GOAL := all

JSONF_EXAMPLE := '[{"id":2,"name":"Alice","age":30},{"age":25,"name":"Bob","id":1}]'

.PHONY: all jsonf limiter test

all: jsonf limiter

jsonf:
	@echo '== jsonf example =='
	@printf '%s\n' $(JSONF_EXAMPLE) | go run ./cmd/jsonf name age
	@echo  

limiter:
	@echo '== limiter example =='
	@go run ./cmd/limiter
	@echo 

test:
	@GOEXPERIMENT=synctest go test -count=1 ./...


