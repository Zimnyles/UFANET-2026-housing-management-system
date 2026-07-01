.PHONY: up down restart build logs ps \
        up-infra up-auth up-profile up-requests up-gateway \
        down-volumes clean \
        proto-auth proto-profile proto-requests proto-all \
        build-auth build-profile build-requests build-gateway \
        logs-auth logs-profile logs-requests logs-rabbitmq logs-gateway \
        config tidy-all

# ─── Docker Compose ───────────────────────────────────────────────────────────

## Поднять все сервисы (пересборка образов)
up:
	docker compose up --build -d

## Поднять только инфраструктуру (postgres + redis + rabbitmq)
up-infra:
	docker compose up -d postgres redis rabbitmq

## Поднять конкретные сервисы поверх инфраструктуры
up-auth:
	docker compose up --build -d auth-service

up-profile:
	docker compose up --build -d profile-service

up-requests:
	docker compose up --build -d requests-service

up-gateway:
	docker compose up --build -d api-gateway

## Остановить всё (тома не удаляются)
down:
	docker compose down

## Остановить и удалить тома (полная очистка БД)
down-volumes:
	docker compose down -v

## Пересобрать и перезапустить все сервисы
restart: down up

## Пересобрать образы без запуска
build:
	docker compose build

## Пересобрать отдельные образы
build-auth:
	docker compose build auth-service

build-profile:
	docker compose build profile-service

build-requests:
	docker compose build requests-service

build-gateway:
	docker compose build api-gateway

## Проверить docker-compose.yml
config:
	docker compose config --quiet

## Статус контейнеров
ps:
	docker compose ps

## Логи всех сервисов (follow)
logs:
	docker compose logs -f

## Логи конкретного сервиса: make logs-auth
logs-auth:
	docker compose logs -f auth-service

logs-profile:
	docker compose logs -f profile-service

logs-requests:
	docker compose logs -f requests-service

logs-rabbitmq:
	docker compose logs -f rabbitmq

logs-gateway:
	docker compose logs -f api-gateway

## Удалить все контейнеры, образы и тома проекта
clean:
	docker compose down -v --rmi local

# ─── Protobuf ─────────────────────────────────────────────────────────────────

## Сгенерировать код для auth контрактов
proto-auth:
	$(MAKE) -C contracts/auth generate-go

## Сгенерировать код для profile контрактов
proto-profile:
	$(MAKE) -C contracts/profile generate-go

## Сгенерировать код для requests контрактов
proto-requests:
	$(MAKE) -C contracts/requests generate-go

## Сгенерировать всё
proto-all: proto-auth proto-profile proto-requests

# ─── Go ───────────────────────────────────────────────────────────────────────

## go mod tidy во всех модулях
tidy-all:
	cd contracts           && go mod tidy
	cd auth-microservice   && go mod tidy
	cd profile-microservice && go mod tidy
	cd requests-microservice && go mod tidy
	cd api-gateway         && go mod tidy
