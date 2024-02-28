build:
	@go build -o .

run: build
	@./gobank

test:
	@go test -v ./...
