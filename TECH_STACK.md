# Hockey Stats Parser - Tech Stack

## Технологии

| Категория | Технология | Статус |
|-----------|-----------|--------|
| **Язык** | Go 1.24+ | ✅ |
| **Архитектура** | Clean Architecture + DDD | ✅ |
| **Тип приложения** | Монолит (2 процесса) | ✅ |
| **База данных** | PostgreSQL | ✅ |
| **БД библиотека** | sqlx + pgx (драйвер) | ✅ |
| **Парсер** | chromedp | ✅ |
| **UI** | Telegram Bot (Inline кнопки) | ✅ |
| **Telegram библиотека** | go-telegram-bot-api/telegram-bot-api/v5 | ✅ |
| **Миграции** | pkg/migrator/pg (goose) | ✅ |
| **Логирование** | pkg/logger (zap) | ✅ |
| **Graceful shutdown** | pkg/closer | ✅ |
| **HTTP router** | Не нужен (MVP) | ⏭️ |
| **Cron scheduler** | Не нужен (MVP) | ⏭️ |
| **Запуск парсера** | Вручную: `go run cmd/parser` | ✅ |
| **Запуск бота** | Вручную: `go run cmd/bot` (long polling) | ✅ |

---

## Структура проекта

```
HockeyProject/
├── cmd/                          # Точки входа
│   ├── parser/main.go           # Parser процесс (парсинг → БД)
│   └── bot/main.go              # Bot процесс (горутина + polling)
│
├── internal/                     # Приватный код
│   ├── domain/                  # Domain layer (чистый, без зависимостей)
│   │   ├── entity/              # Player, School, ScraperRun
│   │   └── vo/                  # Rank, ScraperStatus
│   │
│   ├── service/                 # Application layer
│   │   ├── player/              # PlayerService
│   │   └── scraper/             # ScraperService
│   │
│   ├── repository/              # Infrastructure: DB (sqlx + pgx)
│   │   └── postgres/            # PostgreSQL реализации
│   │
│   ├── client/                  # Infrastructure: External
│   │   ├── scraper/fhr/         # FHR парсер (chromedp)
│   │   │   └── dto/             # DTO для Vaadin протокола
│   │   └── telegram/            # Telegram Bot API клиент
│   │
│   ├── api/                     # Presentation layer
│   │   └── telegram/            # Handlers + Callbacks + Inline клавиатуры
│   │       └── dto/             # DTO для Telegram запросов/ответов
│   │
│   ├── config/                  # Конфигурация из env
│   └── app/                     # DI container
│
├── pkg/                         # Личная библиотека
│   ├── logger/                  # Zap логирование
│   ├── closer/                  # Graceful shutdown
│   ├── migrator/pg/             # Goose миграции
│   ├── cache/redis/             # Redis кеш (опционально)
│   ├── metrics/                 # Prometheus (опционально)
│   └── tracing/                 # OpenTelemetry (опционально)
│
├── migrations/                  # SQL миграции
├── docs/                        # Документация
└── go.mod
```

## Процессы

**Parser** - Парсинг в БД:
```bash
go run cmd/parser/main.go
# → Chromedp → registry.fhr.ru → PostgreSQL
```

**Bot** - Чтение из БД + Telegram:
```bash
go run cmd/bot/main.go
# → Горутина → Long polling → Inline кнопки → PostgreSQL → Ответ
```
