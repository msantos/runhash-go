.PHONY: all fmt format lint

VERSION := 0.9.6

all:
	cd cmd/runhash && \
	CGO_ENABLED=0 go build -trimpath -ldflags "-X runhash/config.Version=$(VERSION) -s -w"

fmt: format
format:
	go fmt -x ./...

lint:
	golangci-lint run
