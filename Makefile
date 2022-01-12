test: ## Test
	@go test -count=1 -race -coverprofile=${COVERAGE_FILE} -covermode=atomic ./...

lint: ## Lint
	@golangci-lint run