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
	@for i in $$(${GOPATH}/bin/goimports -l .); do echo "[ERROR] Code is not formated run 'make format'" && exit 1; done

check: format-check ## Linting and static analysis
	@if grep -r --include='*.go'  -E "[^\/\/ ]+(fmt.Print|spew.Dump)"  *; then \
		echo "code contains fmt.Print* or spew.Dump function"; \
		exit 1; \
	fi

	@if test ! -e ./bin/golangci-lint; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh; \
	fi
	@./bin/golangci-lint run --timeout 180s -E gosec -E stylecheck -E golint -E goimports -E whitespace

init:
	cd ./example && go run github.com/99designs/gqlgen init

hitrix:
	./example/docker/services.sh hitrix