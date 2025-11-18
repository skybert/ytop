
run: build
	./ytop

fmt:
	gofmt -s -w ./cmd ./internal

.phony: version
version ?= $(shell git describe --tags --always --dirty)

.phony: build
build: fmt
	go build \
	  -o ytop \
	  -ldflags "-X main.Version=$(version)" \
	  cmd/*.go

install: build
	mkdir -p ~/.local/bin
	cp ytop ~/.local/bin/.
