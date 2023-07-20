GOCMD=go

run: ## Start application
	set -a && \
	. ./config.env && \
	set +a && \
	$(GOCMD) run main.go
