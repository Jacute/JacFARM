GREEN := \033[0;32m
PURPLE := \033[0;35m
RESET := \033[0m
RED := \033[0;31m

up:
	@echo "$(GREEN)[$(PURPLE)start$(GREEN)]$(RED) Starting services$(RESET)"
	@mkdir -p volumes/exploits
	@sudo chown -R 1000:1000 ./volumes/exploits
	@sudo docker compose up --build -d

down:
	@sudo docker compose down

clean-db:
	@echo "$(GREEN)[$(PURPLE)cleaner$(GREEN)]$(RED) Clean db data$(RESET)"
	@sudo rm -rf volumes/postgres

clean-all:
	@echo "$(GREEN)[$(PURPLE)cleaner$(GREEN)]$(RED) Clean all data: db, rabbitmq, exploits$(RESET)"
	@sudo rm -rf volumes
