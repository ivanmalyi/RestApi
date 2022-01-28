.PHONY: build
build:
		go build -v ./cmd/appserver
.DEFAULT_GOAL := build
