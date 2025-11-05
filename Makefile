.PHONY: up

up:
	@docker compose -f deployments/dev/services.yaml -f deployments/dev/infra.yaml up -d

build:
	@docker compose -f deployments/dev/services.yaml -f deployments/dev/infra.yaml up -d --build

restart:
	@docker compose -f deployments/dev/services.yaml -f deployments/dev/infra.yaml restart

stop:
	@docker compose -f deployments/dev/services.yaml -f deployments/dev/infra.yaml down

purge:
	@docker compose -f deployments/dev/services.yaml -f deployments/dev/infra.yaml down --rmi all
