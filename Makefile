
run: fmt
	go run src/main.go

fmt:
	gofmt -s -w ./src
