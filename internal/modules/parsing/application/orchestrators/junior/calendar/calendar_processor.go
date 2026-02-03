package calendar

import (
	"context"
	"net/url"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	jrcal "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/calendar"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (o *Orchestrator) processCalendar(ctx context.Context, tournamentID, tournamentURL string) error {
	matches, err := o.calendarParser.Parse(tournamentURL)
	if err != nil {
		return err
	}

	// Извлекаем домен из URL турнира
	domain := extractDomain(tournamentURL)

	logger.Debug(ctx, "Parsed calendar matches",
		zap.Int("count", len(matches)),
		zap.String("domain", domain))

	for _, m := range matches {
		if o.config.SkipExisting() {
			existing, _ := o.matchRepo.GetByExternalID(ctx, m.ExternalID, Source)
			if existing != nil {
				continue
			}
		}

		match := o.convertMatch(ctx, tournamentID, domain, m)
		if err := o.matchRepo.Upsert(ctx, match); err != nil {
			logger.Error(ctx, "Failed to save match",
				zap.String("external_id", m.ExternalID),
				zap.Error(err))
		}
	}

	return nil
}

// extractDomain извлекает базовый URL (scheme + host) из полного URL
func extractDomain(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(parsed.Scheme+"://"+parsed.Host, "/")
}

func (o *Orchestrator) convertMatch(ctx context.Context, tournamentID, domain string, m jrcal.MatchDTO) *entities.Match {
	match := &entities.Match{
		ID:           uuid.New().String(),
		ExternalID:   m.ExternalID,
		TournamentID: &tournamentID,
		Status:       m.Status,
		ScheduledAt:  m.ScheduledAt,
		HomeScore:    m.HomeScore,
		AwayScore:    m.AwayScore,
		ResultType:   strPtr(m.ResultType),
		MatchNumber:  m.MatchNumber,
		GroupName:    strPtr(m.GroupName),
		Venue:        strPtr(m.Venue),
		Source:       Source,
		Domain:       strPtr(domain),
	}

	// Копируем информацию о командах для возможного обогащения
	homeInfo := m.HomeTeam
	awayInfo := m.AwayTeam

	// Если у команд нет ID из календаря (нет логотипа), парсим страницу игры
	if homeInfo.ID == "" || awayInfo.ID == "" {
		o.enrichTeamInfoFromGame(ctx, domain, m.ExternalID, &homeInfo, &awayInfo)
	}

	// Находим команды с улучшенной логикой поиска
	if homeTeam := o.findTeam(ctx, homeInfo, tournamentID); homeTeam != nil {
		match.HomeTeamID = &homeTeam.ID
	}
	if awayTeam := o.findTeam(ctx, awayInfo, tournamentID); awayTeam != nil {
		match.AwayTeamID = &awayTeam.ID
	}

	if m.BirthYear > 0 {
		match.BirthYear = &m.BirthYear
	}

	return match
}

// enrichTeamInfoFromGame обогащает информацию о командах из страницы игры
func (o *Orchestrator) enrichTeamInfoFromGame(ctx context.Context, domain, externalID string, homeInfo, awayInfo *jrcal.TeamInfo) {
	gameURL := domain + "/games/" + externalID + "/"

	gameDetails, err := o.gameParser.Parse(gameURL)
	if err != nil {
		logger.Debug(ctx, "Failed to parse game for team URLs",
			zap.String("game_url", gameURL),
			zap.Error(err))
		return
	}

	// Обогащаем домашнюю команду
	if homeInfo.ID == "" && gameDetails.HomeTeamURL != "" {
		homeInfo.URL = gameDetails.HomeTeamURL
		homeInfo.ID = entities.ExtractTeamIDFromURLLegacy(gameDetails.HomeTeamURL)
		logger.Debug(ctx, "Enriched home team from game page",
			zap.String("url", gameDetails.HomeTeamURL),
			zap.String("id", homeInfo.ID))
	}

	// Обогащаем гостевую команду
	if awayInfo.ID == "" && gameDetails.AwayTeamURL != "" {
		awayInfo.URL = gameDetails.AwayTeamURL
		awayInfo.ID = entities.ExtractTeamIDFromURLLegacy(gameDetails.AwayTeamURL)
		logger.Debug(ctx, "Enriched away team from game page",
			zap.String("url", gameDetails.AwayTeamURL),
			zap.String("id", awayInfo.ID))
	}
}

// findTeam ищет команду с несколькими fallback'ами
func (o *Orchestrator) findTeam(ctx context.Context, teamInfo jrcal.TeamInfo, tournamentID string) *entities.Team {
	if teamInfo.ID == "" && teamInfo.URL == "" && teamInfo.Name == "" {
		return nil
	}

	// 1. Поиск по ID (извлечённому из лого или URL)
	if teamInfo.ID != "" {
		if team, _ := o.teamRepo.GetByID(ctx, teamInfo.ID); team != nil {
			return team
		}
	}

	// 2. Поиск по URL напрямую (fallback)
	if teamInfo.URL != "" {
		if team, _ := o.teamRepo.GetByURL(ctx, teamInfo.URL); team != nil {
			return team
		}
	}

	// 3. Поиск по извлечённому ID из URL (если ID не был установлен)
	if teamInfo.ID == "" && teamInfo.URL != "" {
		teamID := entities.ExtractTeamIDFromURLLegacy(teamInfo.URL)
		if teamID != "" {
			if team, _ := o.teamRepo.GetByID(ctx, teamID); team != nil {
				return team
			}
		}
	}

	// 4. Поиск по названию и городу
	if teamInfo.Name != "" {
		name, city := parseTeamNameAndCity(teamInfo.Name)
		if team, _ := o.teamRepo.GetByNameAndCity(ctx, name, city, Source); team != nil {
			return team
		}

		// 5. Поиск только по названию без города
		if team, _ := o.teamRepo.GetByName(ctx, name, Source); team != nil {
			return team
		}
	}

	// 6. Создаём новую команду если есть ID или URL, с привязкой к турниру
	if teamInfo.ID != "" || teamInfo.URL != "" {
		return o.createTeamFromInfo(ctx, teamInfo, tournamentID)
	}

	return nil
}

// createTeamFromInfo создаёт команду из информации календаря
func (o *Orchestrator) createTeamFromInfo(ctx context.Context, info jrcal.TeamInfo, tournamentID string) *entities.Team {
	name, city := parseTeamNameAndCity(info.Name)

	// Используем ID из info, если есть, иначе извлекаем из URL
	teamID := info.ID
	if teamID == "" && info.URL != "" {
		teamID = entities.ExtractTeamIDFromURLLegacy(info.URL)
	}

	if teamID == "" {
		logger.Warn(ctx, "Cannot create team without ID",
			zap.String("name", info.Name))
		return nil
	}

	// Если URL пустой, генерируем уникальный URL на основе ID
	// чтобы избежать duplicate key violation на teams_url_key
	teamURL := info.URL
	if teamURL == "" {
		teamURL = "/calendar-team/" + teamID + "/"
	}

	team := &entities.Team{
		ID:     teamID,
		URL:    teamURL,
		Name:   name,
		City:   city,
		Source: Source,
	}

	// Привязываем команду к турниру если указан
	if tournamentID != "" {
		team.TournamentID = &tournamentID
	}

	if err := o.teamRepo.Upsert(ctx, team); err != nil {
		// При ошибке duplicate key (race condition) пробуем найти уже созданную команду
		if existing, _ := o.teamRepo.GetByID(ctx, teamID); existing != nil {
			return existing
		}
		logger.Warn(ctx, "Failed to create team from calendar",
			zap.String("name", info.Name),
			zap.String("id", teamID),
			zap.Error(err))
		return nil
	}

	logger.Debug(ctx, "Created new team from calendar",
		zap.String("id", team.ID),
		zap.String("name", team.Name))

	return team
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// parseTeamNameAndCity разбирает строку "Название Город" на имя и город
func parseTeamNameAndCity(fullName string) (name, city string) {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return "", ""
	}

	// Последнее слово — это город
	parts := strings.Fields(fullName)
	if len(parts) < 2 {
		return fullName, ""
	}

	city = parts[len(parts)-1]
	name = strings.Join(parts[:len(parts)-1], " ")
	return name, city
}
