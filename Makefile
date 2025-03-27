test:
	@make update
	@go test ./...

build:
	@make test
	@go build -o bin/$$(uname -m)/mock-server cmd/server/main.go

update:
	@go mod tidy

run:
	@go run cmd/server/main.go -config ./dev/config.yml
