package presenter

import (
	"fmt"

	domainPlayer "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// ProfilePresenter presenter для профиля игрока
type ProfilePresenter struct {
	engine TemplateEngine
}

// NewProfilePresenter создает новый presenter
func NewProfilePresenter(engine TemplateEngine) *ProfilePresenter {
	return &ProfilePresenter{
		engine: engine,
	}
}

// FormatPlayerProfile форматирует профиль игрока для отображения в Telegram
func (p *ProfilePresenter) FormatPlayerProfile(profile *domainPlayer.Profile) (string, error) {
	// Подготавливаем данные для шаблона
	data := prepareProfileData(profile)

	// Рендерим через template
	result, err := p.engine.Render("player_profile.tmpl", data)
	if err != nil {
		return "", fmt.Errorf("failed to render profile: %w", err)
	}

	return result, nil
}

// prepareProfileData подготавливает данные для шаблона
func prepareProfileData(profile *domainPlayer.Profile) map[string]interface{} {
	data := map[string]interface{}{
		"BasicInfo":         profile.BasicInfo,
		"AllTimeStats":      profile.AllTimeStats,
		"CurrentSeasonStats": profile.CurrentSeasonStats,
		"RecentTournaments": len(profile.RecentTournaments) > 0,
	}

	// Группируем турниры по сезонам для шаблона
	if len(profile.RecentTournaments) > 0 {
		data["TournamentsBySeason"] = groupTournamentsForTemplate(profile.RecentTournaments)
	}

	return data
}

// SeasonTournaments группа турниров одного сезона
type SeasonTournaments struct {
	Season      string
	Tournaments []domainPlayer.TournamentStats
}

// groupTournamentsForTemplate группирует турниры по сезонам для шаблона
func groupTournamentsForTemplate(tournaments []domainPlayer.TournamentStats) []SeasonTournaments {
	// Группируем по сезонам
	grouped := groupTournamentsBySeason(tournaments)
	
	// Сортируем сезоны
	seasons := getSortedSeasons(grouped)
	
	// Формируем результат
	result := make([]SeasonTournaments, 0, len(seasons))
	for _, season := range seasons {
		result = append(result, SeasonTournaments{
			Season:      season,
			Tournaments: grouped[season],
		})
	}
	
	return result
}

// groupTournamentsBySeason группирует турниры по сезонам
func groupTournamentsBySeason(tournaments []domainPlayer.TournamentStats) map[string][]domainPlayer.TournamentStats {
	result := make(map[string][]domainPlayer.TournamentStats)
	for _, t := range tournaments {
		result[t.Season] = append(result[t.Season], t)
	}
	return result
}

// getSortedSeasons возвращает отсортированные сезоны (от нового к старому)
func getSortedSeasons(tournamentsBySeason map[string][]domainPlayer.TournamentStats) []string {
	seasons := make([]string, 0, len(tournamentsBySeason))
	for season := range tournamentsBySeason {
		seasons = append(seasons, season)
	}
	
	// Простая сортировка по убыванию (2024-2025 > 2023-2024)
	for i := 0; i < len(seasons); i++ {
		for j := i + 1; j < len(seasons); j++ {
			if seasons[i] < seasons[j] {
				seasons[i], seasons[j] = seasons[j], seasons[i]
			}
		}
	}
	
	return seasons
}
