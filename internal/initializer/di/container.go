package di

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/filter"
	profileHandler "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/profile"
	reportHandler "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/report"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/search"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/command"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/keyboard"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/message"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router"
	routerMessage "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router/message"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/template"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot"
	reportService "github.com/Daniil-Sakharov/HockeyProject/internal/service/bot/report"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser"
	"github.com/jmoiron/sqlx"
)

// Container DI контейнер - фасад над модулями
type Container struct {
	infra    *Infrastructure
	repo     *Repository
	service  *Service
	telegram *Telegram
}

// NewContainer создает новый DI контейнер
func NewContainer(cfg *config.Config) *Container {
	infra := NewInfrastructure(cfg)
	repo := NewRepository(infra)
	service := NewService(cfg, infra, repo)
	tg := NewTelegram(cfg, service)

	return &Container{
		infra:    infra,
		repo:     repo,
		service:  service,
		telegram: tg,
	}
}

// Infrastructure
func (d *Container) PostgresDB(ctx context.Context) *sqlx.DB    { return d.infra.PostgresDB(ctx) }
func (d *Container) JuniorClient() *junior.Client               { return d.infra.JuniorClient() }
func (d *Container) StatsParser() *stats.Parser                 { return d.infra.StatsParser() }

// Repositories
func (d *Container) PlayerRepository(ctx context.Context) player.Repository                      { return d.repo.Player(ctx) }
func (d *Container) TeamRepository(ctx context.Context) team.Repository                          { return d.repo.Team(ctx) }
func (d *Container) TournamentRepository(ctx context.Context) tournament.Repository              { return d.repo.Tournament(ctx) }
func (d *Container) PlayerTeamRepository(ctx context.Context) player_team.Repository             { return d.repo.PlayerTeam(ctx) }
func (d *Container) PlayerStatisticsRepository(ctx context.Context) player_statistics.Repository { return d.repo.PlayerStatistics(ctx) }

// Services
func (d *Container) JuniorParserService(ctx context.Context) parser.JuniorParserService           { return d.service.JuniorParser(ctx) }
func (d *Container) OrchestratorService(ctx context.Context) parser.OrchestratorService           { return d.service.Orchestrator(ctx) }
func (d *Container) StatsParserService(ctx context.Context) parser.StatsParserService             { return d.service.StatsParser(ctx) }
func (d *Container) StatsOrchestratorService(ctx context.Context) parser.StatsOrchestratorService { return d.service.StatsOrchestrator(ctx) }
func (d *Container) StateManager() bot.StateManager                                               { return d.service.StateManager() }
func (d *Container) SearchPlayerService(ctx context.Context) bot.SearchPlayerService              { return d.service.SearchPlayer(ctx) }
func (d *Container) ProfileService(ctx context.Context) bot.ProfileService                        { return d.service.Profile(ctx) }
func (d *Container) ReportService(ctx context.Context) *reportService.Service                     { return d.service.Report(ctx) }

// Telegram
func (d *Container) TemplateEngine() *template.Engine                                     { return d.telegram.TemplateEngine() }
func (d *Container) MessagePresenter() *message.MessagePresenter                          { return d.telegram.MessagePresenter() }
func (d *Container) KeyboardPresenter() *keyboard.KeyboardPresenter                       { return d.telegram.KeyboardPresenter() }
func (d *Container) StartHandler(ctx context.Context) *command.StartHandler               { return d.telegram.StartHandler(ctx) }
func (d *Container) FilterHandler(ctx context.Context) *filter.FilterHandler              { return d.telegram.FilterHandler(ctx) }
func (d *Container) SearchHandler(ctx context.Context) *search.Handler                    { return d.telegram.SearchHandler(ctx) }
func (d *Container) ProfileHandler(ctx context.Context) *profileHandler.Handler           { return d.telegram.ProfileHandler(ctx) }
func (d *Container) ReportHandler(ctx context.Context) *reportHandler.Handler             { return d.telegram.ReportHandler(ctx) }
func (d *Container) FioInputHandler(ctx context.Context) *routerMessage.FioInputHandler   { return d.telegram.FioInputHandler(ctx) }
func (d *Container) TelegramRouter(ctx context.Context) *router.Router                    { return d.telegram.Router(ctx) }
func (d *Container) TelegramBot(ctx context.Context) *telegram.Bot                        { return d.telegram.Bot(ctx) }
