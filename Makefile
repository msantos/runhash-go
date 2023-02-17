.PHONY: all fmt format lint

VERSION := 0.10.0

all:
	cd cmd/runhash && \
	CGO_ENABLED=0 go build -trimpath -ldflags "-X codeberg.org/msantos/runhash-go/config.Version=$(VERSION) -s -w"

fmt: format
format:
	go fmt -x ./...

lint:
	golangci-lint run
