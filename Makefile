.PHONY: fmt lint test build

TARGETS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64

default: all

fmt:
	go fmt ./...

lint:
	go vet ./...

test:
	go test ./...

build: $(TARGETS)

linux-amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/zeusctl-linux-amd64 ./main.go

linux-arm64:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o bin/zeusctl-linux-arm64 ./main.go

darwin-amd64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/zeusctl-darwin-amd64 ./main.go

darwin-arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o bin/zeusctl-darwin-arm64 ./main.go

all: fmt lint test build 
