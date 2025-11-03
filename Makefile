.PHONY: up

up:
	@docker compose -f deployments/dev/compose.yaml up -d

build:
	@docker compose -f deployments/dev/compose.yaml up -d --build

restart:
	@docker compose -f deployments/dev/compose.yaml restart

stop:
	@docker compose -f deployments/dev/compose.yaml down

purge:
	@docker compose -f deployments/dev/compose.yaml down --rmi all
