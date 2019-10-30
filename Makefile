#!/usr/bin/make -f

export CGO_ENABLED=0

VERSION=$(shell git describe --tags --always)

# Builds the project.
build:
	gox -os='linux darwin' \
	    -arch='amd64' \
	    -output='bin/cluster-healthz_{{.OS}}_{{.Arch}}' \
	    -ldflags='-extldflags "-static"' \
	    github.com/skpr/cluster-healthz/cmd/cluster-healthz

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting.
test:
	go test -cover ./...

IMAGE=skpr/cluster-healthz

# Releases the project Docker Hub.
release:
	docker build -t ${IMAGE}:${VERSION} -t ${IMAGE}:latest .
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}:latest

.PHONY: *
