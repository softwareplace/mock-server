impl_test:
	@go test ./...

server-test:
	@make codegen
	@cd test && go test

codegen:
	 @rm -rf test/gen/api.gen.go
	 @oapi-codegen --config ./test/resource/config.yaml ./test/resource/pet-store.yaml

pet-store:
	 @oapi-codegen --config ./test/resource/local-config.yaml ./test/resource/pet-store.yaml  2>&1 | pbcopy

