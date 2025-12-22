#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_DIR="$SCRIPT_DIR/../compose/production"
ENV_FILE="$COMPOSE_DIR/.env"
BACKUP_DIR="${BACKUP_DIR:-/opt/hockey-bot/backups}"

# Load env
[[ -f "$ENV_FILE" ]] && source "$ENV_FILE"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/hockey_${TIMESTAMP}.sql.gz"

mkdir -p "$BACKUP_DIR"

echo "[INFO] Creating backup: $BACKUP_FILE"
docker exec hockey-postgres pg_dump -U "${POSTGRES_USER}" "${POSTGRES_DB}" | gzip > "$BACKUP_FILE"

# Keep only last 7 backups
find "$BACKUP_DIR" -name "hockey_*.sql.gz" -mtime +7 -delete 2>/dev/null || true

echo "[INFO] Backup completed: $BACKUP_FILE"
ls -lh "$BACKUP_FILE"
