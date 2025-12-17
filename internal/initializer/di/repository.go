package di

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/fhspb"
	playerRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/junior/player"
	playerStatisticsRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/junior/player_statistics"
	playerTeamRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/junior/player_team"
	teamRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/junior/team"
	tournamentRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/junior/tournament"
)

// Repository содержит все репозитории
type Repository struct {
	infra                      *Infrastructure
	playerRepository           player.Repository
	teamRepository             team.Repository
	tournamentRepository       tournament.Repository
	playerTeamRepository       player_team.Repository
	playerStatisticsRepository player_statistics.Repository
	// FHSPB repositories
	fhspbTournamentRepo       *fhspb.TournamentRepository
	fhspbTeamRepo             *fhspb.TeamRepository
	fhspbPlayerRepo           *fhspb.PlayerRepository
	fhspbPlayerTeamRepo       *fhspb.PlayerTeamRepository
	fhspbPlayerStatisticsRepo *fhspb.PlayerStatisticsRepository
	fhspbGoalieStatisticsRepo *fhspb.GoalieStatisticsRepository
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

// FHSPB repositories

func (r *Repository) FHSPBTournament(ctx context.Context) *fhspb.TournamentRepository {
	if r.fhspbTournamentRepo == nil {
		r.fhspbTournamentRepo = fhspb.NewTournamentRepository(r.infra.PostgresDB(ctx))
	}
	return r.fhspbTournamentRepo
}

func (r *Repository) FHSPBTeam(ctx context.Context) *fhspb.TeamRepository {
	if r.fhspbTeamRepo == nil {
		r.fhspbTeamRepo = fhspb.NewTeamRepository(r.infra.PostgresDB(ctx))
	}
	return r.fhspbTeamRepo
}

func (r *Repository) FHSPBPlayer(ctx context.Context) *fhspb.PlayerRepository {
	if r.fhspbPlayerRepo == nil {
		r.fhspbPlayerRepo = fhspb.NewPlayerRepository(r.infra.PostgresDB(ctx))
	}
	return r.fhspbPlayerRepo
}

func (r *Repository) FHSPBPlayerTeam(ctx context.Context) *fhspb.PlayerTeamRepository {
	if r.fhspbPlayerTeamRepo == nil {
		r.fhspbPlayerTeamRepo = fhspb.NewPlayerTeamRepository(r.infra.PostgresDB(ctx))
	}
	return r.fhspbPlayerTeamRepo
}

func (r *Repository) FHSPBPlayerStatistics(ctx context.Context) *fhspb.PlayerStatisticsRepository {
	if r.fhspbPlayerStatisticsRepo == nil {
		r.fhspbPlayerStatisticsRepo = fhspb.NewPlayerStatisticsRepository(r.infra.PostgresDB(ctx))
	}
	return r.fhspbPlayerStatisticsRepo
}

func (r *Repository) FHSPBGoalieStatistics(ctx context.Context) *fhspb.GoalieStatisticsRepository {
	if r.fhspbGoalieStatisticsRepo == nil {
		r.fhspbGoalieStatisticsRepo = fhspb.NewGoalieStatisticsRepository(r.infra.PostgresDB(ctx))
	}
	return r.fhspbGoalieStatisticsRepo
}
