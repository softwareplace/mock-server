test:
	@go test ./...

build:
	@go build -o bin/$$(uname -m)/mock-server cmd/server/main.go

update:
	@go mod tidy
