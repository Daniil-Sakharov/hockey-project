package di

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	playerRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/player"
	playerStatisticsRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/player_statistics"
	playerTeamRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/player_team"
	teamRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/team"
	tournamentRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/tournament"
)

// Repository содержит все репозитории
type Repository struct {
	infra                      *Infrastructure
	playerRepository           player.Repository
	teamRepository             team.Repository
	tournamentRepository       tournament.Repository
	playerTeamRepository       player_team.Repository
	playerStatisticsRepository player_statistics.Repository
}

func NewRepository(infra *Infrastructure) *Repository {
	return &Repository{infra: infra}
}

func (r *Repository) Player(ctx context.Context) player.Repository {
	if r.playerRepository == nil {
		r.playerRepository = playerRepo.NewRepository(r.infra.PostgresDB(ctx))
	}
	return r.playerRepository
}

func (r *Repository) Team(ctx context.Context) team.Repository {
	if r.teamRepository == nil {
		r.teamRepository = teamRepo.NewRepository(r.infra.PostgresDB(ctx))
	}
	return r.teamRepository
}

func (r *Repository) Tournament(ctx context.Context) tournament.Repository {
	if r.tournamentRepository == nil {
		r.tournamentRepository = tournamentRepo.NewRepository(r.infra.PostgresDB(ctx))
	}
	return r.tournamentRepository
}

func (r *Repository) PlayerTeam(ctx context.Context) player_team.Repository {
	if r.playerTeamRepository == nil {
		r.playerTeamRepository = playerTeamRepo.NewRepository(r.infra.PostgresDB(ctx))
	}
	return r.playerTeamRepository
}

func (r *Repository) PlayerStatistics(ctx context.Context) player_statistics.Repository {
	if r.playerStatisticsRepository == nil {
		r.playerStatisticsRepository = playerStatisticsRepo.NewRepository(r.infra.PostgresDB(ctx))
	}
	return r.playerStatisticsRepository
}
