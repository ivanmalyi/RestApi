.PHONY: build
build:
		go build -v ./cmd/appserver
#название задачи
.PHONY: test
test:
		go test -v -race -timeout 30s ./...

#эта задача запускаеться по дефолту
.DEFAULT_GOAL := build
