# Hockey Bot - Remote Management
# Использование: make <команда>

SERVER = yandex
REMOTE_DIR = /opt/hockey-bot
COMPOSE_FILE = docker-compose.yml

.DEFAULT_GOAL := help

# === ПОМОЩЬ ===
help: ## Показать все доступные команды
	@echo "🏒 Hockey Bot - Управление"
	@echo ""
	@echo "📋 Доступные команды:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "📖 Примеры использования:"
	@echo "  make status          # Проверить статус контейнеров"
	@echo "  make run-parser-bg   # Запустить парсер в фоне"
	@echo "  make logs-parser     # Посмотреть логи парсера"

# === ЛОКАЛЬНАЯ РАЗРАБОТКА ===

local-up: ## Запустить локальное окружение (PostgreSQL)
	docker compose -f deploy/compose/local/docker-compose.yml up -d postgres

local-migrate: ## Запустить миграции локально
	docker compose -f deploy/compose/local/docker-compose.yml --profile migrate run --rm migrate

local-bot: ## Запустить бота локально
	docker compose -f deploy/compose/local/docker-compose.yml --profile bot up -d bot

local-parser: ## Запустить парсер локально
	docker compose -f deploy/compose/local/docker-compose.yml --profile parser run --rm parser

local-down: ## Остановить локальное окружение
	docker compose -f deploy/compose/local/docker-compose.yml down

local-logs: ## Показать логи локального окружения
	docker compose -f deploy/compose/local/docker-compose.yml logs -f

# === ПРОДАКШН (УДАЛЕННЫЙ СЕРВЕР) ===

# === ПРОДАКШН (УДАЛЕННЫЙ СЕРВЕР) ===

logs: ## Показать логи бота в реальном времени
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) logs -f bot"

logs-tail: ## Показать последние 100 строк логов бота
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) logs --tail=100 bot"

status: ## Показать статус всех контейнеров
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) ps -a"

stats: ## Показать использование ресурсов Docker
	ssh $(SERVER) "docker stats --no-stream"

restart: ## Перезапустить бота
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) restart bot"

stop: ## Остановить все сервисы
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) down"

start: ## Запустить бота и базу данных
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) up -d postgres bot"

# === ПАРСЕРЫ ===

run-junior-parser: ## Запустить Junior парсер (игроки/команды)
	@echo "🚀 Запуск Junior парсера..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm junior-parser"

run-junior-stats: ## Запустить Junior Stats парсер (статистика)
	@echo "📊 Запуск Junior Stats парсера..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm junior-stats-parser"

run-fhspb-parser: ## Запустить FHSPB парсер (игроки/команды)
	@echo "🏒 Запуск FHSPB парсера..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm fhspb-parser"

run-fhspb-stats: ## Запустить FHSPB Stats парсер (статистика)
	@echo "📊 Запуск FHSPB Stats парсера..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm fhspb-stats-parser"

# Алиасы для обратной совместимости
run-parser: run-junior-parser ## Алиас для run-junior-parser
run-stats: run-junior-stats ## Алиас для run-junior-stats
run-fhspb: run-fhspb-parser ## Алиас для run-fhspb-parser

# === ФОНОВЫЙ ЗАПУСК ===

run-junior-parser-bg: ## Запустить Junior парсер в фоне
	@echo "🚀 Запуск Junior парсера в фоне..."
	nohup make run-junior-parser > junior-parser.log 2>&1 &
	@echo "✅ Junior парсер запущен в фоне. Логи: tail -f junior-parser.log"

run-junior-stats-bg: ## Запустить Junior Stats парсер в фоне
	@echo "📊 Запуск Junior Stats парсера в фоне..."
	nohup make run-junior-stats > junior-stats.log 2>&1 &
	@echo "✅ Junior Stats парсер запущен в фоне. Логи: tail -f junior-stats.log"

run-fhspb-parser-bg: ## Запустить FHSPB парсер в фоне
	@echo "🏒 Запуск FHSPB парсера в фоне..."
	nohup make run-fhspb-parser > fhspb-parser.log 2>&1 &
	@echo "✅ FHSPB парсер запущен в фоне. Логи: tail -f fhspb-parser.log"

run-fhspb-stats-bg: ## Запустить FHSPB Stats парсер в фоне
	@echo "📊 Запуск FHSPB Stats парсера в фоне..."
	nohup make run-fhspb-stats > fhspb-stats.log 2>&1 &
	@echo "✅ FHSPB Stats парсер запущен в фоне. Логи: tail -f fhspb-stats.log"

run-all-junior-bg: ## Запустить все Junior парсеры в фоне
	@echo "🚀 Запуск всех Junior парсеров в фоне..."
	nohup make run-junior-parser > junior-parser.log 2>&1 &
	nohup make run-junior-stats > junior-stats.log 2>&1 &
	@echo "✅ Все Junior парсеры запущены в фоне"

run-all-fhspb-bg: ## Запустить все FHSPB парсеры в фоне
	@echo "🏒 Запуск всех FHSPB парсеров в фоне..."
	nohup make run-fhspb-parser > fhspb-parser.log 2>&1 &
	nohup make run-fhspb-stats > fhspb-stats.log 2>&1 &
	@echo "✅ Все FHSPB парсеры запущены в фоне"

run-all-parsers-bg: ## Запустить ВСЕ парсеры в фоне
	@echo "🚀 Запуск ВСЕХ парсеров в фоне..."
	nohup make run-junior-parser > junior-parser.log 2>&1 &
	nohup make run-junior-stats > junior-stats.log 2>&1 &
	nohup make run-fhspb-parser > fhspb-parser.log 2>&1 &
	nohup make run-fhspb-stats > fhspb-stats.log 2>&1 &
	@echo "✅ Все парсеры запущены в фоне"

# Алиасы для обратной совместимости
run-parser-bg: run-junior-parser-bg ## Алиас для run-junior-parser-bg
run-stats-bg: run-junior-stats-bg ## Алиас для run-junior-stats-bg
run-all-bg: run-all-parsers-bg ## Алиас для run-all-parsers-bg

stop-parsers: ## Остановить все запущенные парсеры
	@echo "🛑 Остановка всех парсеров..."
	ssh $(SERVER) "docker ps | grep parser | awk '{print \$$1}' | xargs -r docker stop"

# === ЛОГИ ЛОКАЛЬНЫХ ФАЙЛОВ ===
logs-junior-parser-file: ## Показать логи Junior парсера из локального файла
	tail -f junior-parser.log

logs-junior-stats-file: ## Показать логи Junior Stats парсера из локального файла
	tail -f junior-stats.log

logs-fhspb-parser-file: ## Показать логи FHSPB парсера из локального файла
	tail -f fhspb-parser.log

logs-fhspb-stats-file: ## Показать логи FHSPB Stats парсера из локального файла
	tail -f fhspb-stats.log

# === ЛОГИ DOCKER КОНТЕЙНЕРОВ ===
logs-junior-parser: ## Показать логи Junior парсера из Docker
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) logs -f junior-parser"

logs-junior-stats: ## Показать логи Junior Stats парсера из Docker
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) logs -f junior-stats-parser"

logs-fhspb-parser: ## Показать логи FHSPB парсера из Docker
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) logs -f fhspb-parser"

logs-fhspb-stats: ## Показать логи FHSPB Stats парсера из Docker
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) logs -f fhspb-stats-parser"

# Алиасы для обратной совместимости
logs-parser: logs-junior-parser ## Алиас для logs-junior-parser
logs-stats: logs-junior-stats ## Алиас для logs-junior-stats
logs-parser-file: logs-junior-parser-file ## Алиас для logs-junior-parser-file
logs-stats-file: logs-junior-stats-file ## Алиас для logs-junior-stats-file
logs-fhspb-file: logs-fhspb-parser-file ## Алиас для logs-fhspb-parser-file
logs-fhspb: logs-fhspb-parser ## Алиас для logs-fhspb-parser

# === БАЗА ДАННЫХ ===

db-tunnel:
	@echo "🔗 Подключение к БД: localhost:5432"
	@echo "   User: hockey, DB: hockey_stats"
	@echo "   Ctrl+C для отключения"
	ssh -N -L 5432:localhost:5432 $(SERVER)

db-shell:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) exec postgres psql -U hockey -d hockey_stats"

# === ДЕПЛОЙ ===

deploy:
	scp deploy/compose/$(COMPOSE_FILE) $(SERVER):$(REMOTE_DIR)/
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) pull && docker compose -f $(COMPOSE_FILE) up -d postgres bot"

migrate:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) run --rm migrate"

show-env:
	ssh $(SERVER) "cat $(REMOTE_DIR)/.env"

edit-env:
	ssh $(SERVER) "nano $(REMOTE_DIR)/.env"

.PHONY: logs logs-tail status stats restart stop start \
        run-parser run-stats run-fhspb run-fhspb-stats \
        logs-parser logs-stats logs-fhspb \
        db-tunnel db-shell deploy migrate show-env edit-env
# === МОНИТОРИНГ И УПРАВЛЕНИЕ ===

deploy-monitoring: ## Установить Portainer для веб-управления
	@echo "🚀 Установка Portainer..."
	scp deploy/compose/monitoring/docker-compose.yml $(SERVER):/opt/hockey-bot/monitoring-compose.yml
	scp deploy/compose/monitoring/.env.example $(SERVER):/opt/hockey-bot/monitoring.env
	ssh $(SERVER) "cd /opt/hockey-bot && echo 'HockeyAdmin2024!' | docker run --rm -i portainer/portainer-ce --hash-password | head -1 > /tmp/portainer_password"
	ssh $(SERVER) "cd /opt/hockey-bot && docker compose -f monitoring-compose.yml up -d portainer"
	@echo "✅ Portainer установлен!"
	@echo "🌐 Доступ: http://158.160.211.62:9000"
	@echo "👤 Логин: admin"
	@echo "🔑 Пароль: HockeyAdmin2024!"

stop-monitoring: ## Остановить мониторинг
	ssh $(SERVER) "cd /opt/hockey-bot && docker compose -f monitoring-compose.yml down"

logs-monitoring: ## Показать логи мониторинга
	ssh $(SERVER) "cd /opt/hockey-bot && docker compose -f monitoring-compose.yml logs -f"
