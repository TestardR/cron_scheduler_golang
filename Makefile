test: ## Test
	@go test -count=1 -race -coverprofile=${COVERAGE_FILE} -covermode=atomic ./...

lint: ## Lint
	@golangci-lint run

install: ## Download and install go mod
	@go mod download

build: ## Build App
	go build ./cmd/main.go
