.PHONY: up

up:
	@docker compose \
		-f deployments/dev/services.yaml \
		-f deployments/dev/infra.yaml \
		-f deployments/dev/migration.yaml \
		up -d

build:
	@docker compose \
		-f deployments/dev/services.yaml \
		-f deployments/dev/infra.yaml \
		-f deployments/dev/migration.yaml \
		up -d --build

restart:
	@docker compose \
		-f deployments/dev/services.yaml \
		-f deployments/dev/infra.yaml \
		-f deployments/dev/migration.yaml \
		restart

stop:
	@docker compose \
		-f deployments/dev/services.yaml \
		-f deployments/dev/infra.yaml \
		-f deployments/dev/migration.yaml \
		down

purge:
	@docker compose \
		-f deployments/dev/services.yaml \
		-f deployments/dev/infra.yaml \
		-f deployments/dev/migration.yaml \
		down --rmi all
