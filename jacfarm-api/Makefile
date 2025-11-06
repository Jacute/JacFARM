.PHONY: coverage install-swag swagger
PKGS := $(shell go list ./... | grep -vE 'mocks|cmd/jacfarm-api')

coverage:
	@go test $(PKGS) -coverprofile=coverage.out -coverpkg=$(shell echo $(PKGS) | tr ' ' ',') -timeout 30s
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out

install-swag:
	@go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	@swag init -g ./internal/http/server/router.go --parseDependency --parseInternal --generatedTime -o ../docs/jacfarm-api/swagger

wrk-send-flags:
	@echo "Sending many flags with wrk..."
	wrk -t2 -c30 -d20s http://localhost:1337 -s flags.lua
