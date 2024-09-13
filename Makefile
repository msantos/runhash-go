.PHONY: all fmt format lint

all:
	CGO_ENABLED=0 go build -trimpath -ldflags "-w" ./cmd/runhash

fmt: format
format:
	go fmt -x ./...

lint:
	golangci-lint run
