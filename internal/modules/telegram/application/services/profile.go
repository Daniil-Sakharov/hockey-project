package services

import "context"

// PlayerProfile профиль игрока (соответствует шаблону player_profile.tmpl)
type PlayerProfile struct {
	BasicInfo           PlayerBasicInfo
	AllTimeStats        *PlayerStats
	CurrentSeasonStats  *SeasonStats
	RecentTournaments   []*ProfileTournamentStats
	TournamentsBySeason []*SeasonTournaments
}

// PlayerBasicInfo базовая информация об игроке
type PlayerBasicInfo struct {
	ID        string
	Name      string
	BirthYear int
	Position  string
	Height    *int
	Weight    *int
	Team      string
	Region    string
}

// PlayerStats статистика игрока за всё время
type PlayerStats struct {
	Tournaments      int
	Games            int
	Goals            int
	Assists          int
	Points           int
	PlusMinus        int
	Penalties        int
	GoalsPerGame     float64
	AssistsPerGame   float64
	PointsPerGame    float64
	PenaltiesPerGame float64
	HatTricks        int
	GameWinningGoals int
}

// SeasonStats статистика за сезон
type SeasonStats struct {
	Season           string
	Tournaments      int
	Games            int
	Goals            int
	Assists          int
	Points           int
	PlusMinus        int
	Penalties        int
	GoalsPerGame     float64
	AssistsPerGame   float64
	PointsPerGame    float64
	PenaltiesPerGame float64
	HatTricks        int
	GameWinningGoals int
}

// TournamentStats статистика по турниру
type ProfileTournamentStats struct {
	TournamentName string
	GroupName      string
	Season         string
	Games          int
	Goals          int
	Assists        int
	Points         int
	PlusMinus      int
	Penalties      int
	HatTricks      int
	WinningGoals   int
	IsChampionship bool
}

// SeasonTournaments турниры по сезону
type SeasonTournaments struct {
	Season      string
	Tournaments []*ProfileTournamentStats
}

// ProfileRepository интерфейс для получения профиля
type ProfileRepository interface {
	GetByID(ctx context.Context, playerID string) (*PlayerProfile, error)
	GetStats(ctx context.Context, playerID string) (*PlayerStats, error)
	GetSeasonStats(ctx context.Context, playerID, season string) (*SeasonStats, error)
	GetRecentTournaments(ctx context.Context, playerID string, limit int) ([]*ProfileTournamentStats, error)
	GetTournamentsBySeason(ctx context.Context, playerID string) ([]*SeasonTournaments, error)
	GetCurrentSeason() string
}

// ProfileService сервис работы с профилями
type ProfileService struct {
	repo ProfileRepository
}

// NewProfileService создает новый сервис профилей
func NewProfileService(repo ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

// GetProfile возвращает полный профиль игрока со статистикой
func (s *ProfileService) GetProfile(ctx context.Context, playerID string) (*PlayerProfile, error) {
	profile, err := s.repo.GetByID(ctx, playerID)
	if err != nil {
		return nil, err
	}

	// Получаем статистику за всё время
	stats, _ := s.repo.GetStats(ctx, playerID)
	profile.AllTimeStats = stats

	// Получаем статистику за текущий сезон
	season := s.repo.GetCurrentSeason()
	seasonStats, _ := s.repo.GetSeasonStats(ctx, playerID, season)
	profile.CurrentSeasonStats = seasonStats

	// Получаем последние 5 турниров
	recentTournaments, _ := s.repo.GetRecentTournaments(ctx, playerID, 5)
	profile.RecentTournaments = recentTournaments

	// Получаем турниры по сезонам
	tournamentsBySeason, _ := s.repo.GetTournamentsBySeason(ctx, playerID)
	profile.TournamentsBySeason = tournamentsBySeason

	return profile, nil
}
