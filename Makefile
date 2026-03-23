.PHONY: generate
generate:
	go generate ./...
	go fmt ./...

.PHONY: api
api:
	@echo "openapi: getting latest Huma REST API"
	@curl https://openapi.humahr.com/v3/api-docs | jq . | yq -Poy > ./api/openapi.yaml

.PHONY: check
check:
	go vet ./...
	go tool -modfile=go.tool.mod golangci-lint run ./...

.PHONY: bump
bump:
	go tool -modfile=go.tool.mod git-bump 
