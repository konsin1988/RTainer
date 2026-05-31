ifneq (,$(wildcard .env))
include .env
export
endif

dev:
	cd ~/keycloak && docker compose up -d; 
	progress_bar 10;
	COMPOSE_PROFILES=dev docker compose up -d;
	cd app && npm run dev

prod:
	@set -e; docker image rm konsin1988/$(PROJECT_NAME):$(PROJECT_VERSION) 2>/dev/null || true; COMPOSE_PROFILES=prod docker compose up -d --build

down:
	COMPOSE_PROFILES=dev,prod docker compose down;
	cd ~/keycloak && docker compose down;

