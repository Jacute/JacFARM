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
	@echo "JACFARM_API_KEY=$(JACFARM_API_KEY)" > .env
	@echo "ADMIN_PASS=$(ADMIN_PASS)" >> .env
	@echo "$(PURPLE)Farm api key$(RESET) - $(JACFARM_API_KEY)"
	@echo "$(PURPLE)Farm admin creds$(RESET) - admin:$(ADMIN_PASS)"
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
