.PHONY: up

dev/up:
	@docker compose \
		-f deployments/dev/services.yaml \
		-f deployments/dev/infra.yaml \
		-f deployments/dev/migration.yaml \
		up -d

dev/build:
	@docker compose \
		-f deployments/dev/services.yaml \
		-f deployments/dev/infra.yaml \
		-f deployments/dev/migration.yaml \
		up -d --build

dev/restart:
	@docker compose \
		-f deployments/dev/services.yaml \
		-f deployments/dev/infra.yaml \
		-f deployments/dev/migration.yaml \
		restart

dev/stop:
	@docker compose \
		-f deployments/dev/services.yaml \
		-f deployments/dev/infra.yaml \
		-f deployments/dev/migration.yaml \
		down

dev/purge:
	@docker compose \
		-f deployments/dev/services.yaml \
		-f deployments/dev/infra.yaml \
		-f deployments/dev/migration.yaml \
		down --rmi all

local/up:
	@docker compose \
		-f deployments/local/services.yaml \
		-f deployments/local/infra.yaml \
		-f deployments/local/migration.yaml \
		up -d

local/build:
	@docker compose \
		-f deployments/local/services.yaml \
		-f deployments/local/infra.yaml \
		-f deployments/local/migration.yaml \
		up -d --build

local/restart:
	@docker compose \
		-f deployments/local/services.yaml \
		-f deployments/local/infra.yaml \
		-f deployments/local/migration.yaml \
		restart

local/stop:
	@docker compose \
		-f deployments/local/services.yaml \
		-f deployments/local/infra.yaml \
		-f deployments/local/migration.yaml \
		down

local/purge:
	@docker compose \
		-f deployments/local/services.yaml \
		-f deployments/local/infra.yaml \
		-f deployments/local/migration.yaml \
		down --rmi all
