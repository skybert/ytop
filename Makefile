
run: build
	./ytop

fmt:
	gofmt -s -w ./cmd ./internal

build: fmt
	go build -o ytop cmd/*.go
