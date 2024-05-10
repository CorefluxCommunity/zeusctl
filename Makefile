.PHONY: fmt lint test build

default: test

fmt:
	go fmt ./...

lint:
	go vet ./...

test:
	go test ./...

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cfctl ./main.go

