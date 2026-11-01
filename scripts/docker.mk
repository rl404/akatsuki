# Compose command.
COMPOSE ?= docker compose

# Compose files helper.
compose = $(COMPOSE) -f deployment/compose.$(1).yml -p akatsuki-$(2)

.PHONY: docker-lint docker-test docker-build docker-migrate \
	docker-api docker-consumer docker-cron-update docker-cron-fill docker-stop

##@ Docker Tools

docker-lint: ## Run linter in docker.
	@$(call compose,lint,lint) --profile lint up
	@$(call compose,lint,lint) --profile lint rm -f

docker-test: ## Run test in docker.
	@$(call compose,test,test) --profile test up
	@$(call compose,test,test) --profile test rm -f

##@ Docker Build & Run

docker-build: ## Build docker image then prune the builder image.
	@$(call compose,build,build) --profile build build 
	@docker image prune -f --filter label=stage=akatsuki_builder

docker-migrate: ## Start migrate container.
	@$(call compose,migrate,migrate) --profile migrate up
	@$(call compose,migrate,migrate) --profile migrate rm -f

docker-api: ## Start API container.
	@$(call compose,api,api) up -d
	@$(call compose,api,api) logs --follow --tail 20

docker-consumer: ## Start consumer container.
	@$(call compose,consumer,consumer) up -d
	@$(call compose,consumer,consumer) logs --follow --tail 20

docker-cron-update: ## Start cron update container.
	@$(call compose,cron-update,cron-update) up
	@$(call compose,cron-update,cron-update) rm -f

docker-cron-fill: ## Start cron fill container.
	@$(call compose,cron-fill,cron-fill) up
	@$(call compose,cron-fill,cron-fill) rm -f

docker-stop: ## Stop running container.
	@$(call compose,api,api) stop
	@$(call compose,consumer,consumer) stop
	@$(call compose,api,api) rm -f
	@$(call compose,consumer,consumer) rm -f
