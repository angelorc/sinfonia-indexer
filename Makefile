gql-init:
	@go run github.com/99designs/gqlgen init

gql-generate:
	@go run github.com/99designs/gqlgen generate

.PHONY: gql-init gql-generate