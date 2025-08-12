PKGS := $(shell go list ./... | grep -vE 'mocks|cmd/flag_sender')

coverage:
	@go test $(PKGS) -coverprofile=coverage.out -coverpkg=$(shell echo $(PKGS) | tr ' ' ',') -timeout 30s
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out