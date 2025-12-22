# Hockey Bot Deployment

## Quick Start

```bash
# 1. Configure environment
cp deploy/compose/production/.env.example deploy/compose/production/.env
# Edit .env with your values

# 2. Start services
./deploy/scripts/deploy.sh up

# 3. Run migrations
./deploy/scripts/deploy.sh migrate
```

## Scripts

| Script | Description |
|--------|-------------|
| `deploy.sh up` | Start bot + postgres |
| `deploy.sh down` | Stop all services |
| `deploy.sh restart` | Restart services |
| `deploy.sh logs [service]` | View logs |
| `deploy.sh status` | Service status |
| `deploy.sh migrate` | Run DB migrations |
| `deploy.sh parser <type>` | Run parser (junior/fhspb/all) |
| `deploy.sh monitoring` | Start observability stack |
| `backup.sh` | Backup database |
| `healthcheck.sh` | Check service health |

## Monitoring Stack

Start monitoring:
```bash
./deploy/scripts/deploy.sh monitoring
```

Access:
- Grafana: http://localhost:3000 (admin/GRAFANA_PASSWORD)
- Prometheus: http://localhost:9090
- Jaeger: http://localhost:16686

### Dashboards

- **Hockey Overview** - requests, errors, response times
- **Parsers** - parsing stats, players/teams counts

### Alerts

| Alert | Condition |
|-------|-----------|
| HighErrorRate | >0.1 errors/sec for 5m |
| ServiceDown | Service unavailable for 1m |
| SlowResponses | Avg response >5s for 5m |
| ParsingFailed | >10 parsing errors in 1h |

## Environment Variables

Required in `.env`:
```
POSTGRES_USER=hockey
POSTGRES_PASSWORD=<secure-password>
POSTGRES_DB=hockey
TELEGRAM_BOT_TOKEN=<bot-token>
GITHUB_REPOSITORY=<owner/repo>
GRAFANA_PASSWORD=<grafana-password>
```

## Directory Structure

```
deploy/
├── compose/
│   ├── production/     # Production docker-compose
│   └── observability/  # Standalone monitoring stack
├── config/
│   ├── grafana/        # Dashboards & provisioning
│   ├── prometheus/     # Prometheus config & alerts
│   └── otel/           # OTEL Collector config
├── scripts/            # Deployment scripts
└── docker/             # Dockerfile
```

## Data Persistence

All data stored in `/opt/hockey-bot/data/`:
- `postgres/` - PostgreSQL data
- `prometheus/` - Metrics (15 days retention)
- `grafana/` - Grafana config
- `backups/` - Database backups (7 days rotation)
