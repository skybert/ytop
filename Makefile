
run: fmt
	go run src/main.go

fmt:
	gofmt -s -w ./src

build: fmt
	go build -o ytop src/*.go
