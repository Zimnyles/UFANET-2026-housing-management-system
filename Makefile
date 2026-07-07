.PHONY: up down restart build logs ps \
        up-infra up-auth up-profile up-news up-requests up-notifications up-gateway up-frontend \
        down-volumes clean \
        proto-auth proto-profile proto-requests proto-all \
        build-auth build-profile build-news build-requests build-notifications build-gateway build-frontend \
        logs-auth logs-profile logs-news logs-requests logs-rabbitmq logs-gateway logs-frontend \
        config tidy-all

up:
	docker compose config --quiet
	docker compose up --build -d --wait
	@echo "Приложение запущено: http://zhkh.localhost:$${FRONTEND_PORT:-3000}"

up-infra:
	docker compose up -d postgres redis rabbitmq

up-auth:
	docker compose up --build -d auth-service

up-profile:
	docker compose up --build -d profile-service

up-news:
	docker compose up --build -d news-service

up-requests:
	docker compose up --build -d requests-service

up-notifications:
	docker compose up --build -d notification-service

up-gateway:
	docker compose up --build -d api-gateway

up-frontend:
	docker compose up --build -d frontend

down:
	docker compose down

down-volumes:
	docker compose down -v

restart: down up

build:
	docker compose build

build-auth:
	docker compose build auth-service

build-profile:
	docker compose build profile-service

build-news:
	docker compose build news-service

build-requests:
	docker compose build requests-service

build-notifications:
	docker compose build notification-service

build-gateway:
	docker compose build api-gateway

build-frontend:
	docker compose build frontend

config:
	docker compose config --quiet

ps:
	docker compose ps

logs:
	docker compose logs -f

logs-auth:
	docker compose logs -f auth-service

logs-profile:
	docker compose logs -f profile-service

logs-news:
	docker compose logs -f news-service

logs-requests:
	docker compose logs -f requests-service

logs-rabbitmq:
	docker compose logs -f rabbitmq

logs-gateway:
	docker compose logs -f api-gateway

logs-frontend:
	docker compose logs -f frontend

clean:
	docker compose down -v --rmi local

proto-auth:
	$(MAKE) -C contracts/auth generate-go

proto-profile:
	$(MAKE) -C contracts/profile generate-go

proto-requests:
	$(MAKE) -C contracts/requests generate-go

proto-all: proto-auth proto-profile proto-requests

tidy-all:
	cd contracts           && go mod tidy
	cd auth-microservice   && go mod tidy
	cd profile-microservice && go mod tidy
	cd requests-microservice && go mod tidy
	cd news-microservice && go mod tidy
	cd api-gateway         && go mod tidy
	cd notification-microservice && go mod tidy
