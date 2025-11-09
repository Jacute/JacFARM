.PHONY: help up down clean-db clean-all reset coverage

# Цвета
GREEN  := \033[0;32m
PURPLE := \033[0;35m
RESET  := \033[0m
RED    := \033[0;31m
JACFARM_API_KEY := $(shell openssl rand -hex 32)
ADMIN_PASS     := $(shell openssl rand -hex 8)


# Автогенерируемый help: выводит список всех целей с описанием после ##
help: ## Show this help
	@echo "$(GREEN)Available commands:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  $(PURPLE)%-12s$(RESET) %s\n", $$1, $$2}'

up: ## Start services
	@echo "$(GREEN)[$(PURPLE)start$(GREEN)]$(RED) Starting services$(RESET)"
	@mkdir -p volumes/exploits
	@sudo chown -R 1000:1000 ./volumes/exploits
	@touch .env

	@if grep -q '^JACFARM_API_KEY=' .env; then \
		printf "$(PURPLE)Farm api key$(RESET) - %s\n" "$$(grep '^JACFARM_API_KEY=' .env | cut -d'=' -f2-)"; \
	else \
		echo "JACFARM_API_KEY=$(JACFARM_API_KEY)" >> .env; \
		printf "$(PURPLE)Farm api key$(RESET) - %s\n" "$(JACFARM_API_KEY)"; \
	fi

	@if grep -q '^ADMIN_PASS=' .env; then \
		printf "$(PURPLE)Farm creds$(RESET) - admin:%s\n" "$$(grep '^ADMIN_PASS=' .env | cut -d'=' -f2-)"; \
	else \
		echo "ADMIN_PASS=$(ADMIN_PASS)" >> .env; \
		printf "$(PURPLE)Farm creds$(RESET) - admin:%s\n" "$(ADMIN_PASS)"; \
	fi

	@sudo docker compose --env-file .env up --build -d
	@echo "$(GREEN)Services started$(RESET)"


down: ## Stop services
	@echo "$(GREEN)[$(PURPLE)stop$(GREEN)]$(RED) Stopping services$(RESET)"
	@sudo docker compose down

clean-db: ## Remove Postgres data
	@echo "$(GREEN)[$(PURPLE)cleaner$(GREEN)]$(RED) Clean db data$(RESET)"
	@sudo rm -rf volumes/postgres

clean-all: ## Remove all volumes (db, rabbitmq, exploits)
	@echo "$(GREEN)[$(PURPLE)cleaner$(GREEN)]$(RED) Clean all data$(RESET)"
	@sudo rm -rf volumes
	@rm .env

reset: down clean-all up ## Full reset: stop, clean and restart

coverage:
	@echo "$(GREEN)Running tests$(RESET)"
	
	@cd jacfarm-api && \
	PKGS=$$(go list ./... | grep -vE 'mocks|cmd/jacfarm-api' || true) && \
	if [ -z "$$PKGS" ]; then echo "No packages found in jacfarm-api"; exit 0; fi && \
	CSV=$$(echo $$PKGS | tr ' ' ',') && \
	echo "Running go test for: $$PKGS" && \
	go test -v -race -coverpkg=$$CSV -covermode=atomic -coverprofile=../jacfarm-api.out $$PKGS

	@cd workers/config_loader && \
	PKGS=$$(go list ./... | grep -vE 'mocks|cmd/config_loader' || true) && \
	if [ -z "$$PKGS" ]; then echo "No packages found in config_loader"; exit 0; fi && \
	CSV=$$(echo $$PKGS | tr ' ' ',') && \
	echo "Running go test for: $$PKGS" && \
	go test -v -race -coverpkg=$$CSV -covermode=atomic -coverprofile=../../config_loader.out $$PKGS

	@cd workers/exploit_runner && \
	PKGS=$$(go list ./... | grep -vE 'mocks|cmd/exploit_runner' || true) && \
	if [ -z "$$PKGS" ]; then echo "No packages found in exploit_runner"; exit 0; fi && \
	CSV=$$(echo $$PKGS | tr ' ' ',') && \
	echo "Running go test for: $$PKGS" && \
	go test -v -race -coverpkg=$$CSV -covermode=atomic -coverprofile=../../exploit_runner.out $$PKGS
	
	@cd workers/flag_sender && \
	PKGS=$$(go list ./... | grep -vE 'mocks|cmd/flag_sender' || true) && \
	if [ -z "$$PKGS" ]; then echo "No packages found in flag_sender"; exit 0; fi && \
	CSV=$$(echo $$PKGS | tr ' ' ',') && \
	echo "Running go test for: $$PKGS" && \
	go test -v -race -coverpkg=$$CSV -covermode=atomic -coverprofile=../../flag_sender.out $$PKGS
	
	gocovmerge jacfarm-api.out config_loader.out exploit_runner.out flag_sender.out > coverage.out
	rm exploit_runner.out flag_sender.out jacfarm-api.out config_loader.out