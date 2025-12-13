package module

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/closer"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Infrastructure содержит инфраструктурные зависимости
type Infrastructure struct {
	config       *config.Config
	db           *sqlx.DB
	juniorClient *junior.Client
	statsParser  *stats.Parser
}

func NewInfrastructure(cfg *config.Config) *Infrastructure {
	return &Infrastructure{config: cfg}
}

func (i *Infrastructure) PostgresDB(ctx context.Context) *sqlx.DB {
	if i.db == nil {
		db, err := sqlx.Connect("pgx", i.config.Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to PostgreSQL: %s", err.Error()))
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(1 * time.Minute)

		if err := db.Ping(); err != nil {
			panic(fmt.Sprintf("failed to ping PostgreSQL: %s", err.Error()))
		}

		closer.AddNamed("PostgreSQL", func(ctx context.Context) error {
			return db.Close()
		})

		logger.Info(ctx, "✅ PostgreSQL connected")
		i.db = db
	}
	return i.db
}

func (i *Infrastructure) JuniorClient() *junior.Client {
	if i.juniorClient == nil {
		i.juniorClient = junior.NewClient()
	}
	return i.juniorClient
}

func (i *Infrastructure) StatsParser() *stats.Parser {
	if i.statsParser == nil {
		httpClient := &http.Client{Timeout: 60 * time.Second}
		i.statsParser = stats.NewParser(httpClient)
	}
	return i.statsParser
}
