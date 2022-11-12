export GO111MODULE=on

format-check: ## Format the code and run linters
	@if test ! -e ./bin/golangci-lint; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.50.1; \
	fi
	@./bin/golangci-lint run --fix

test-cover: ## Run tests with coverage
	@go install github.com/ory/go-acc@latest
	@go-acc ./... --output=coverage.out --covermode=atomic -- -race -p 1

init:
	cd ./example && go run github.com/99designs/gqlgen init

hitrix:
	./example/docker/services.sh hitrix
