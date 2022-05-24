generate:
	@if [ -d tmp ]; then rm -r tmp; fi;
	@go run github.com/99designs/gqlgen generate

.PHONY: generate