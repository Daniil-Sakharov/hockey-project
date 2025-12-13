package presenter

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

type ProfilePresenter struct {
	engine TemplateEngine
}

func NewProfilePresenter(engine TemplateEngine) *ProfilePresenter {
	return &ProfilePresenter{engine: engine}
}

func (p *ProfilePresenter) FormatPlayerProfile(profile *player.Profile) (string, error) {
	data := prepareProfileData(profile)

	result, err := p.engine.Render("player_profile.tmpl", data)
	if err != nil {
		return "", fmt.Errorf("failed to render profile: %w", err)
	}

	return result, nil
}

func prepareProfileData(profile *player.Profile) map[string]interface{} {
	data := map[string]interface{}{
		"BasicInfo":          profile.BasicInfo,
		"AllTimeStats":       profile.AllTimeStats,
		"CurrentSeasonStats": profile.CurrentSeasonStats,
		"RecentTournaments":  len(profile.RecentTournaments) > 0,
	}

	if len(profile.RecentTournaments) > 0 {
		data["TournamentsBySeason"] = groupTournamentsForTemplate(profile.RecentTournaments)
	}

	return data
}

type SeasonTournaments struct {
	Season      string
	Tournaments []player.TournamentStats
}

func groupTournamentsForTemplate(tournaments []player.TournamentStats) []SeasonTournaments {
	grouped := make(map[string][]player.TournamentStats)
	for _, t := range tournaments {
		grouped[t.Season] = append(grouped[t.Season], t)
	}

	seasons := make([]string, 0, len(grouped))
	for season := range grouped {
		seasons = append(seasons, season)
	}

	for i := 0; i < len(seasons); i++ {
		for j := i + 1; j < len(seasons); j++ {
			if seasons[i] < seasons[j] {
				seasons[i], seasons[j] = seasons[j], seasons[i]
			}
		}
	}

	result := make([]SeasonTournaments, 0, len(seasons))
	for _, season := range seasons {
		result = append(result, SeasonTournaments{
			Season:      season,
			Tournaments: grouped[season],
		})
	}

	return result
}
