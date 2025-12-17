# Hockey Bot - Remote Management
# Использование: make <команда>

SERVER = yandex
REMOTE_DIR = /opt/hockey-bot
COMPOSE_FILE = docker-compose.prod.yml

# === МОНИТОРИНГ ===

logs:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) logs -f bot"

logs-tail:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) logs --tail=100 bot"

status:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) ps -a"

stats:
	ssh $(SERVER) "docker stats --no-stream"

# === УПРАВЛЕНИЕ БОТОМ ===

restart:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) restart bot"

stop:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) down"

start:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) up -d postgres bot"

# === ПАРСЕРЫ ===

run-parser:
	@echo "🚀 Запуск Junior парсера..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) --profile parser run --rm parser"

run-stats:
	@echo "📊 Запуск Stats парсера..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) --profile parser run --rm stats-parser"

run-fhspb:
	@echo "🏒 Запуск FHSPB парсера..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) --profile parser run --rm fhspb-parser"

run-fhspb-stats:
	@echo "📊 Запуск FHSPB Stats парсера..."
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) --profile parser run --rm fhspb-stats-parser"

# Логи парсеров
logs-parser:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) logs -f parser"

logs-stats:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) logs -f stats-parser"

logs-fhspb:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker compose -f $(COMPOSE_FILE) logs -f fhspb-parser"

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
