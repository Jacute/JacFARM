PKGS := $(shell go list ./... | grep -vE 'mocks|cmd/jacfarm-api|storage/|JacFARM/internal/config|JacFARM/internal/http/server|JacFARM/internal/rabbitmq|JacFARM/internal/http/validator')

coverage:
	@go test $(PKGS) -coverprofile=coverage.out -coverpkg=$(shell echo $(PKGS) | tr ' ' ',') -timeout 30s
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out

install-swag:
	@go install github.com/swaggo/swag/cmd/swag@latest

swagger:
	@swag init -g ./internal/http/server/router.go --parseDependency --parseInternal --generatedTime -o ../docs/jacfarm-api/swagger