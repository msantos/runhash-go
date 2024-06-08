.PHONY: all fmt format lint

all:
	cd cmd/runhash && \
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w"

fmt: format
format:
	go fmt -x ./...

lint:
	golangci-lint run
