GOCMD=go
docker-command=docker-compose -f docker-compose.yml

run: ## Start application
	set -a && \
	. ./config.env && \
	set +a && \
	$(GOCMD) run main.go

dep-build-up:
	${docker-command} up -d --build

dep-up:
	${docker-command} up -d

dep-down:
	${docker-command} down

dep-stop:
	${docker-command} stop

swag: ## Generate swagger docs
	swag init -o ./docs
