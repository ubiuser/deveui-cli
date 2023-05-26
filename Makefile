.PHONY: run
run: ## Run the app from source code
	go run main.go

.PHONY: docker
docker: ## Run the app in a container built from docker file
	docker compose up app_docker

.PHONY: local
local: ## Run the app in a container using source code
	docker compose up app_local
