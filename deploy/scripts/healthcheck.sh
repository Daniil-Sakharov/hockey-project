#!/bin/bash
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

check_service() {
    local name="$1"
    local container="$2"
    if docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
        local status=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null || echo "running")
        if [[ "$status" == "healthy" || "$status" == "running" ]]; then
            echo -e "${GREEN}✓${NC} $name"
            return 0
        fi
    fi
    echo -e "${RED}✗${NC} $name"
    return 1
}

echo "=== Hockey Bot Health Check ==="
echo ""

FAILED=0
check_service "PostgreSQL" "hockey-postgres" || FAILED=1
check_service "Bot" "hockey-bot" || FAILED=1

# Optional monitoring services
docker ps --format '{{.Names}}' | grep -q "hockey-prometheus" && check_service "Prometheus" "hockey-prometheus"
docker ps --format '{{.Names}}' | grep -q "hockey-grafana" && check_service "Grafana" "hockey-grafana"
docker ps --format '{{.Names}}' | grep -q "hockey-jaeger" && check_service "Jaeger" "hockey-jaeger"

echo ""
exit $FAILED
