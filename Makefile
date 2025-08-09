GREEN := \033[0;32m
PURPLE := \033[0;35m
RESET := \033[0m
RED := \033[0;31m

clean-db:
	@echo "$(GREEN)[$(PURPLE)cleaner$(GREEN)]$(RED) Clean db data$(RESET)"
	@sudo rm -rf volumes/postgres

clean-all:
	@echo "$(GREEN)[$(PURPLE)cleaner$(GREEN)]$(RED) Clean all data: db, rabbitmq, exploits$(RESET)"
	@sudo rm -rf volumes
