# Hockey Bot - Remote Management
# Использование: make <команда>

SERVER = yandex
REMOTE_DIR = /opt/hockey-bot
COMPOSE_FILE = docker-compose.yml

.DEFAULT_GOAL := help

# Команды, которые не создают файлы
.PHONY: help local-up local-migrate local-bot local-parser local-down local-logs
.PHONY: logs logs-tail status stats restart stop start
.PHONY: run-junior-parser run-junior-stats run-fhspb-parser run-fhspb-stats
.PHONY: run-parser run-stats run-fhspb
.PHONY: run-junior-parser-bg run-junior-stats-bg run-fhspb-parser-bg run-fhspb-stats-bg
.PHONY: run-all-junior-bg run-all-fhspb-bg run-all-parsers-bg
.PHONY: run-parser-bg run-stats-bg run-all-bg stop-parsers
.PHONY: logs-junior-parser-file logs-junior-stats-file logs-fhspb-parser-file logs-fhspb-stats-file
.PHONY: logs-junior-parser logs-junior-stats logs-fhspb-parser logs-fhspb-stats
.PHONY: logs-parser logs-stats logs-parser-file logs-stats-file logs-fhspb-file logs-fhspb
.PHONY: db-tunnel db-shell deploy deploy-manual
.PHONY: deploy-monitoring logs-monitoring stop-monitoring

# === ПОМОЩЬ ===
help: ## Показать все доступные команды
	@echo "🏒 Hockey Bot - Управление"
	@echo ""
	@echo "📋 Доступные команды:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "📖 Примеры использования:"
	@echo "  make status          # Проверить статус контейнеров"
	@echo "  make run-junior-parser-bg   # Запустить Junior парсер в фоне"
	@echo "  make logs-junior-parser     # Посмотреть логи парсера"

# === ЛОКАЛЬНАЯ РАЗРАБОТКА ===

local-up: ## Запустить локальное окружение (PostgreSQL)
	docker compose -f deploy/compose/local/docker-compose.yml up -d

local-migrate: ## Запустить миграции локально
	@echo "🔄 Запуск миграций локально..."
	go run cmd/migrate/main.go

local-bot: ## Запустить бота локально
	@echo "🤖 Запуск бота локально..."
	go run cmd/bot/main.go

local-parser: ## Запустить парсер локально
	@echo "🏒 Запуск парсера локально..."
	go run cmd/parser/main.go

local-down: ## Остановить локальное окружение
	docker compose -f deploy/compose/local/docker-compose.yml down

local-logs: ## Показать логи локального окружения
	docker compose -f deploy/compose/local/docker-compose.yml logs -f

# === ПРОДАКШН (УДАЛЕННЫЙ СЕРВЕР) ===

logs: ## Показать логи бота в реальном времени
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) logs -f bot"

logs-tail: ## Показать последние 100 строк логов бота
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) logs --tail=100 bot"

status: ## Показать статус всех контейнеров
	@echo "📊 Статус контейнеров на сервере:"
	ssh $(SERVER) "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'"

stats: ## Показать использование ресурсов Docker
	ssh $(SERVER) "docker stats --no-stream"

restart: ## Перезапустить бота
	@echo "🔄 Перезапуск бота..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) restart bot"

stop: ## Остановить все сервисы
	@echo "🛑 Остановка всех сервисов..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) down"

start: ## Запустить бота и базу данных
	@echo "🚀 Запуск бота и базы данных..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) up -d postgres bot"

# === ПАРСЕРЫ (ПРОДАКШН) ===

run-junior-parser: ## Запустить Junior парсер (игроки/команды)
	@echo "🏒 Запуск Junior парсера..."
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

# === ФОНОВЫЕ ПАРСЕРЫ (ПРОДАКШН) ===

run-junior-parser-bg: ## Запустить Junior парсер в фоне на сервере
	@echo "🏒 Запуск Junior парсера в фоне на сервере..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && mkdir -p logs && nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm junior-parser > logs/junior-parser.log 2>&1 &"
	@echo "✅ Junior парсер запущен в фоне. Логи: make logs-junior-parser-file"

run-junior-stats-bg: ## Запустить Junior Stats парсер в фоне на сервере
	@echo "📊 Запуск Junior Stats парсера в фоне на сервере..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && mkdir -p logs && nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm junior-stats-parser > logs/junior-stats.log 2>&1 &"
	@echo "✅ Junior Stats парсер запущен в фоне. Логи: make logs-junior-stats-file"

run-fhspb-parser-bg: ## Запустить FHSPB парсер в фоне на сервере
	@echo "🏒 Запуск FHSPB парсера в фоне на сервере..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && mkdir -p logs && nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm fhspb-parser > logs/fhspb-parser.log 2>&1 &"
	@echo "✅ FHSPB парсер запущен в фоне. Логи: make logs-fhspb-parser-file"

run-fhspb-stats-bg: ## Запустить FHSPB Stats парсер в фоне на сервере
	@echo "📊 Запуск FHSPB Stats парсера в фоне на сервере..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && mkdir -p logs && nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm fhspb-stats-parser > logs/fhspb-stats.log 2>&1 &"
	@echo "✅ FHSPB Stats парсер запущен в фоне. Логи: make logs-fhspb-stats-file"

run-all-junior-bg: ## Запустить все Junior парсеры в фоне на сервере
	@echo "🚀 Запуск всех Junior парсеров в фоне на сервере..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && mkdir -p logs && nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm junior-parser > logs/junior-parser.log 2>&1 & nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm junior-stats-parser > logs/junior-stats.log 2>&1 &"
	@echo "✅ Все Junior парсеры запущены в фоне"

run-all-fhspb-bg: ## Запустить все FHSPB парсеры в фоне на сервере
	@echo "🏒 Запуск всех FHSPB парсеров в фоне на сервере..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && mkdir -p logs && nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm fhspb-parser > logs/fhspb-parser.log 2>&1 & nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm fhspb-stats-parser > logs/fhspb-stats.log 2>&1 &"
	@echo "✅ Все FHSPB парсеры запущены в фоне"

run-all-parsers-bg: ## Запустить ВСЕ парсеры в фоне на сервере
	@echo "🚀 Запуск ВСЕХ парсеров в фоне на сервере..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && mkdir -p logs && nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm junior-parser > logs/junior-parser.log 2>&1 & nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm junior-stats-parser > logs/junior-stats.log 2>&1 & nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm fhspb-parser > logs/fhspb-parser.log 2>&1 & nohup docker compose -f deploy/compose/production/$(COMPOSE_FILE) --profile parser run --rm fhspb-stats-parser > logs/fhspb-stats.log 2>&1 &"
	@echo "✅ Все парсеры запущены в фоне"

# Алиасы для обратной совместимости
run-parser-bg: run-junior-parser-bg ## Алиас для run-junior-parser-bg
run-stats-bg: run-junior-stats-bg ## Алиас для run-junior-stats-bg
run-all-bg: run-all-parsers-bg ## Алиас для run-all-parsers-bg

stop-parsers: ## Остановить все запущенные парсеры
	@echo "🛑 Остановка всех парсеров..."
	ssh $(SERVER) "docker ps | grep parser | awk '{print \$$1}' | xargs -r docker stop"

# === ЛОГИ ФАЙЛОВ С СЕРВЕРА ===

logs-junior-parser-file: ## Показать логи Junior парсера из файла на сервере
	ssh $(SERVER) "tail -f $(REMOTE_DIR)/logs/junior-parser.log"

logs-junior-stats-file: ## Показать логи Junior Stats парсера из файла на сервере
	ssh $(SERVER) "tail -f $(REMOTE_DIR)/logs/junior-stats.log"

logs-fhspb-parser-file: ## Показать логи FHSPB парсера из файла на сервере
	ssh $(SERVER) "tail -f $(REMOTE_DIR)/logs/fhspb-parser.log"

logs-fhspb-stats-file: ## Показать логи FHSPB Stats парсера из файла на сервере
	ssh $(SERVER) "tail -f $(REMOTE_DIR)/logs/fhspb-stats.log"

# === УПРАВЛЕНИЕ ЛОГАМИ ===

logs-list: ## Показать список всех файлов логов на сервере
	@echo "📋 Файлы логов на сервере:"
	ssh $(SERVER) "ls -la $(REMOTE_DIR)/logs/ 2>/dev/null || echo 'Папка logs не существует'"

logs-clean: ## Очистить все файлы логов на сервере
	@echo "🧹 Очистка логов на сервере..."
	ssh $(SERVER) "rm -f $(REMOTE_DIR)/logs/*.log && echo '✅ Логи очищены'"

logs-size: ## Показать размер файлов логов
	@echo "📊 Размер файлов логов:"
	ssh $(SERVER) "du -h $(REMOTE_DIR)/logs/*.log 2>/dev/null || echo 'Нет файлов логов'"

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

db-tunnel: ## Создать SSH туннель к базе данных
	@echo "🔗 Подключение к БД: localhost:5432"
	@echo "   User: hockey, DB: hockey_stats"
	@echo "   Ctrl+C для отключения"
	ssh -N -L 5432:localhost:5432 $(SERVER)

db-shell: ## Подключиться к PostgreSQL через psql
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) exec postgres psql -U hockey -d hockey_stats"

# === ДЕПЛОЙ ===

deploy: ## Деплой через GitHub Actions (автоматически)
	@echo "🚀 Деплой запускается автоматически через GitHub Actions при push в main"
	@echo "   Проверить статус: https://github.com/Daniil-Sakharov/hockey-project/actions"

deploy-manual: ## Ручной деплой (копирование файлов и перезапуск)
	@echo "📦 Копирование docker-compose на сервер..."
	scp deploy/compose/production/$(COMPOSE_FILE) $(SERVER):$(REMOTE_DIR)/deploy/compose/production/
	@echo "🔄 Перезапуск сервисов..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f deploy/compose/production/$(COMPOSE_FILE) pull && docker compose -f deploy/compose/production/$(COMPOSE_FILE) up -d postgres bot"
	@echo "✅ Деплой завершен"

# === МОНИТОРИНГ ===

deploy-monitoring: ## Установить Portainer для мониторинга
	@echo "📊 Установка Portainer..."
	scp deploy/compose/monitoring/docker-compose.yml $(SERVER):$(REMOTE_DIR)/monitoring-compose.yml
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f monitoring-compose.yml up -d"
	@echo "✅ Portainer установлен: http://your-server:9000"

logs-monitoring: ## Показать логи мониторинга
	ssh $(SERVER) "docker logs hockey-portainer -f"

stop-monitoring: ## Остановить мониторинг
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f monitoring-compose.yml down"
