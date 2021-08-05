SHELL := /bin/bash
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org

cover: ## Run tests with coverage and creates cover.out profile
	@mkdir -p ./resources/cover
	@rm -f ./resources/cover/tmp-cover.log;
	@go get github.com/ory/go-acc
	@${GOPATH}/bin/go-acc ./... --output=resources/cover/cover.out --covermode=atomic

format: ## Format go code with goimports
	@go get golang.org/x/tools/cmd/goimports
	@${GOPATH}/bin/goimports -l -w .

format-check: ## Check if the code is formatted
	@go get golang.org/x/tools/cmd/goimports
	@for i in $$(${GOPATH}/bin/goimports -l .); do echo "[ERROR] Code is not formatted run 'make format'" && exit 1; done

check: format-check cyclo ## Linting and static analysis
	@if grep -r --include='*.go' -E "fmt.Print|spew.Dump" *; then \
		echo "code contains fmt.Print* or spew.Dump function"; \
		exit 1; \
	fi

	@go install honnef.co/go/tools/cmd/staticcheck@latest;
	@${GOPATH}/bin/staticcheck ./...;

cyclo: ## Cyclomatic complexities analysis
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@${GOPATH}/bin/gocyclo -over 100 .

init:
	cd ./example && go run github.com/99designs/gqlgen init

hitrix:
	./example/docker/services.sh hitrix
