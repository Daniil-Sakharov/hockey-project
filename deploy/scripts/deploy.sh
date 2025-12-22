#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_DIR="$SCRIPT_DIR/../compose/production"
ENV_FILE="$COMPOSE_DIR/.env"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1" >&2; exit 1; }

usage() {
    cat << EOF
Usage: $0 <command> [options]

Commands:
  up              Start all services (bot + postgres)
  down            Stop all services
  restart         Restart all services
  logs [service]  Show logs (optionally for specific service)
  status          Show service status
  migrate         Run database migrations
  parser <type>   Run parser (junior|junior-stats|fhspb|fhspb-stats|all)
  monitoring      Start monitoring stack (prometheus, grafana, jaeger)
  pull            Pull latest images

Options:
  -h, --help      Show this help
EOF
}

check_env() {
    [[ -f "$ENV_FILE" ]] || error ".env file not found at $ENV_FILE"
}

compose() {
    docker compose -f "$COMPOSE_DIR/docker-compose.yml" --env-file "$ENV_FILE" "$@"
}

cmd_up() {
    check_env
    log "Starting services..."
    compose up -d bot postgres
    log "Services started"
}

cmd_down() {
    log "Stopping services..."
    compose down
    log "Services stopped"
}

cmd_restart() {
    cmd_down
    cmd_up
}

cmd_logs() {
    local service="${1:-}"
    if [[ -n "$service" ]]; then
        compose logs -f "$service"
    else
        compose logs -f
    fi
}

cmd_status() {
    compose ps
}

cmd_migrate() {
    check_env
    log "Running migrations..."
    compose --profile migrate up migrate
    log "Migrations completed"
}

cmd_parser() {
    check_env
    local type="${1:-}"
    case "$type" in
        junior)       compose --profile parser up junior-parser ;;
        junior-stats) compose --profile parser up junior-stats-parser ;;
        fhspb)        compose --profile parser up fhspb-parser ;;
        fhspb-stats)  compose --profile parser up fhspb-stats-parser ;;
        all)
            compose --profile parser up junior-parser
            compose --profile parser up junior-stats-parser
            compose --profile parser up fhspb-parser
            compose --profile parser up fhspb-stats-parser
            ;;
        *) error "Unknown parser type: $type. Use: junior|junior-stats|fhspb|fhspb-stats|all" ;;
    esac
}

cmd_monitoring() {
    check_env
    log "Starting monitoring stack..."
    compose --profile monitoring up -d
    log "Monitoring available at:"
    log "  Grafana:    http://localhost:3000"
    log "  Prometheus: http://localhost:9090"
    log "  Jaeger:     http://localhost:16686"
}

cmd_pull() {
    log "Pulling latest images..."
    compose pull
    log "Images updated"
}

# Main
case "${1:-}" in
    up)         cmd_up ;;
    down)       cmd_down ;;
    restart)    cmd_restart ;;
    logs)       cmd_logs "${2:-}" ;;
    status)     cmd_status ;;
    migrate)    cmd_migrate ;;
    parser)     cmd_parser "${2:-}" ;;
    monitoring) cmd_monitoring ;;
    pull)       cmd_pull ;;
    -h|--help)  usage ;;
    *)          usage; exit 1 ;;
esac
