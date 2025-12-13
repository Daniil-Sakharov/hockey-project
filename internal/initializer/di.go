package initializer

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
	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer/module"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot"
	reportService "github.com/Daniil-Sakharov/HockeyProject/internal/service/bot/report"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser"
	"github.com/jmoiron/sqlx"
)

// diContainer DI контейнер - фасад над модулями
type diContainer struct {
	infra    *module.Infrastructure
	repo     *module.Repository
	service  *module.Service
	telegram *module.Telegram
}

// NewDiContainer создает новый DI контейнер
func NewDiContainer(cfg *config.Config) *diContainer {
	infra := module.NewInfrastructure(cfg)
	repo := module.NewRepository(infra)
	service := module.NewService(cfg, infra, repo)
	tg := module.NewTelegram(cfg, service)

	return &diContainer{
		infra:    infra,
		repo:     repo,
		service:  service,
		telegram: tg,
	}
}

// Infrastructure
func (d *diContainer) PostgresDB(ctx context.Context) *sqlx.DB       { return d.infra.PostgresDB(ctx) }
func (d *diContainer) JuniorClient() *junior.Client                  { return d.infra.JuniorClient() }
func (d *diContainer) StatsParser() *stats.Parser                    { return d.infra.StatsParser() }

// Repositories
func (d *diContainer) PlayerRepository(ctx context.Context) player.Repository                       { return d.repo.Player(ctx) }
func (d *diContainer) TeamRepository(ctx context.Context) team.Repository                           { return d.repo.Team(ctx) }
func (d *diContainer) TournamentRepository(ctx context.Context) tournament.Repository               { return d.repo.Tournament(ctx) }
func (d *diContainer) PlayerTeamRepository(ctx context.Context) player_team.Repository              { return d.repo.PlayerTeam(ctx) }
func (d *diContainer) PlayerStatisticsRepository(ctx context.Context) player_statistics.Repository  { return d.repo.PlayerStatistics(ctx) }

// Services
func (d *diContainer) JuniorParserService(ctx context.Context) parser.JuniorParserService           { return d.service.JuniorParser(ctx) }
func (d *diContainer) OrchestratorService(ctx context.Context) parser.OrchestratorService           { return d.service.Orchestrator(ctx) }
func (d *diContainer) StatsParserService(ctx context.Context) parser.StatsParserService             { return d.service.StatsParser(ctx) }
func (d *diContainer) StatsOrchestratorService(ctx context.Context) parser.StatsOrchestratorService { return d.service.StatsOrchestrator(ctx) }
func (d *diContainer) StateManager() bot.StateManager                                               { return d.service.StateManager() }
func (d *diContainer) SearchPlayerService(ctx context.Context) bot.SearchPlayerService              { return d.service.SearchPlayer(ctx) }
func (d *diContainer) ProfileService(ctx context.Context) bot.ProfileService                        { return d.service.Profile(ctx) }
func (d *diContainer) ReportService(ctx context.Context) *reportService.Service                     { return d.service.Report(ctx) }

// Telegram
func (d *diContainer) TemplateEngine() *template.Engine                                { return d.telegram.TemplateEngine() }
func (d *diContainer) MessagePresenter() *message.MessagePresenter                     { return d.telegram.MessagePresenter() }
func (d *diContainer) KeyboardPresenter() *keyboard.KeyboardPresenter                  { return d.telegram.KeyboardPresenter() }
func (d *diContainer) StartHandler(ctx context.Context) *command.StartHandler          { return d.telegram.StartHandler(ctx) }
func (d *diContainer) FilterHandler(ctx context.Context) *filter.FilterHandler         { return d.telegram.FilterHandler(ctx) }
func (d *diContainer) SearchHandler(ctx context.Context) *search.Handler               { return d.telegram.SearchHandler(ctx) }
func (d *diContainer) ProfileHandler(ctx context.Context) *profileHandler.Handler      { return d.telegram.ProfileHandler(ctx) }
func (d *diContainer) ReportHandler(ctx context.Context) *reportHandler.Handler        { return d.telegram.ReportHandler(ctx) }
func (d *diContainer) FioInputHandler(ctx context.Context) *routerMessage.FioInputHandler { return d.telegram.FioInputHandler(ctx) }
func (d *diContainer) TelegramRouter(ctx context.Context) *router.Router               { return d.telegram.Router(ctx) }
func (d *diContainer) TelegramBot(ctx context.Context) *telegram.Bot                   { return d.telegram.Bot(ctx) }
