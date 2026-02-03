package calendar

import (
	"context"
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	jrcal "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/calendar"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/helpers"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

// processCalendarWithFilters парсит календарь с итерацией по годам и группам через AJAX
func (o *Orchestrator) processCalendarWithFilters(ctx context.Context, tournament *entities.Tournament) error {
	calendarURL := buildCalendarURL(tournament.Domain + tournament.URL)

	logger.Info(ctx, "Loading calendar page for AJAX iteration",
		zap.String("url", calendarURL))

	resp, err := o.http.MakeRequest(calendarURL)
	if err != nil {
		return fmt.Errorf("failed to load calendar page: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return fmt.Errorf("calendar page returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse calendar HTML: %w", err)
	}

	// Извлекаем AJAX-ссылки на годы (используем competitions-calendar)
	yearLinks := helpers.ExtractYearLinksForCalendar(doc)
	logger.Info(ctx, "Found year links for calendar", zap.Int("count", len(yearLinks)))

	if len(yearLinks) == 0 {
		// Нет year dropdown — проверяем есть ли группы напрямую
		groupLinks := helpers.ExtractGroupLinksForCalendar(doc)
		if len(groupLinks) > 0 {
			// Есть группы без year dropdown — определяем год из birth_year_groups турнира
			birthYear := o.extractSingleBirthYear(tournament)
			logger.Info(ctx, "No year dropdown but found calendar groups, using birth year from tournament",
				zap.Int("groups", len(groupLinks)),
				zap.Int("birthYear", birthYear))
			return o.processCalendarGroupsWithoutYearDropdown(ctx, tournament, doc, groupLinks, birthYear)
		}
		// Нет ни годов, ни групп — парсим как есть (обычный календарь)
		return o.processCalendar(ctx, tournament.ID, tournament.Domain+tournament.URL)
	}

	// Итерируем по годам
	for _, yearLink := range yearLinks {
		if err := o.processYearCalendar(ctx, tournament, yearLink); err != nil {
			logger.Warn(ctx, "Failed to process year",
				zap.Int("year", yearLink.Year),
				zap.Error(err))
		}
	}

	return nil
}

// processYearCalendar обрабатывает календарь для конкретного года
func (o *Orchestrator) processYearCalendar(ctx context.Context, tournament *entities.Tournament, yearLink types.YearLink) error {
	fullURL := tournament.Domain + yearLink.AjaxURL

	logger.Debug(ctx, "Loading year calendar",
		zap.Int("year", yearLink.Year),
		zap.String("url", fullURL))

	resp, err := o.http.MakeRequest(fullURL)
	if err != nil {
		return fmt.Errorf("ajax request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ajax returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse ajax HTML: %w", err)
	}

	// Извлекаем группы для этого года (используем competitions-calendar)
	groupLinks := helpers.ExtractGroupLinksForCalendar(doc)

	if len(groupLinks) == 0 {
		// Год без групп — парсим матчи этого года напрямую
		return o.parseAndSaveMatchesFromAjax(ctx, tournament.Domain, yearLink.AjaxURL, tournament.ID, yearLink.Year, "")
	}

	logger.Info(ctx, "Found groups for year",
		zap.Int("year", yearLink.Year),
		zap.Int("groups", len(groupLinks)))

	// Парсим дефолтную (активную) группу из текущего AJAX-ответа
	activeGroupName := helpers.ExtractActiveGroupName(doc, "competitions-calendar")
	if activeGroupName != "" {
		if err := o.parseAndSaveMatchesFromAjax(ctx, tournament.Domain, yearLink.AjaxURL, tournament.ID, yearLink.Year, activeGroupName); err != nil {
			logger.Warn(ctx, "Failed to parse active group",
				zap.Int("year", yearLink.Year),
				zap.String("group", activeGroupName),
				zap.Error(err))
		}
	}

	// Итерируем по остальным группам через AJAX
	for _, group := range groupLinks {
		if err := o.parseAndSaveMatchesFromAjax(ctx, tournament.Domain, group.AjaxURL, tournament.ID, yearLink.Year, group.Name); err != nil {
			logger.Warn(ctx, "Failed to parse group",
				zap.Int("year", yearLink.Year),
				zap.String("group", group.Name),
				zap.Error(err))
		}
	}

	return nil
}

// parseAndSaveMatchesFromAjax парсит матчи через AJAX и сохраняет с правильным tournament_id
func (o *Orchestrator) parseAndSaveMatchesFromAjax(ctx context.Context, domain, ajaxURL, tournamentID string, birthYear int, groupName string) error {
	filter := jrcal.CalendarFilter{
		BirthYear: birthYear,
		GroupName: groupName,
	}

	matches, err := o.calendarParser.ParseWithFilter(domain, ajaxURL, filter)
	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	logger.Info(ctx, "Parsed matches from AJAX",
		zap.Int("count", len(matches)),
		zap.String("tournament_id", tournamentID),
		zap.Int("year", birthYear),
		zap.String("group", groupName))

	savedCount := 0
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
			continue
		}
		savedCount++
	}

	logger.Debug(ctx, "Saved matches",
		zap.Int("saved", savedCount),
		zap.Int("total", len(matches)))

	return nil
}

// processCalendarGroupsWithoutYearDropdown обрабатывает группы календаря когда нет year dropdown
func (o *Orchestrator) processCalendarGroupsWithoutYearDropdown(
	ctx context.Context,
	tournament *entities.Tournament,
	doc *goquery.Document,
	groupLinks []types.GroupLink,
	birthYear int,
) error {
	logger.Info(ctx, "Processing calendar groups without year dropdown",
		zap.String("tournament", tournament.Name),
		zap.Int("groups", len(groupLinks)),
		zap.Int("birthYear", birthYear))

	// Парсим активную группу (без data-ajax-link)
	activeGroupName := helpers.ExtractActiveGroupName(doc, "competitions-calendar")
	if activeGroupName != "" {
		// Для активной группы парсим текущую страницу
		if err := o.parseAndSaveMatchesFromDoc(ctx, doc, tournament.ID, tournament.Domain, birthYear, activeGroupName); err != nil {
			logger.Warn(ctx, "Failed to parse active calendar group",
				zap.String("group", activeGroupName),
				zap.Error(err))
		}
	}

	// Итерируем по остальным группам через AJAX
	for _, group := range groupLinks {
		if err := o.parseAndSaveMatchesFromAjax(ctx, tournament.Domain, group.AjaxURL, tournament.ID, birthYear, group.Name); err != nil {
			logger.Warn(ctx, "Failed to parse calendar group",
				zap.Int("year", birthYear),
				zap.String("group", group.Name),
				zap.Error(err))
		}
	}

	return nil
}

// parseAndSaveMatchesFromDoc парсит матчи из уже загруженного документа
func (o *Orchestrator) parseAndSaveMatchesFromDoc(ctx context.Context, doc *goquery.Document, tournamentID, domain string, birthYear int, groupName string) error {
	filter := jrcal.CalendarFilter{
		BirthYear: birthYear,
		GroupName: groupName,
	}

	matches, err := o.calendarParser.ParseFromDoc(doc, filter)
	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	logger.Info(ctx, "Parsed matches from doc",
		zap.Int("count", len(matches)),
		zap.String("tournament_id", tournamentID),
		zap.Int("year", birthYear),
		zap.String("group", groupName))

	savedCount := 0
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
			continue
		}
		savedCount++
	}

	return nil
}

// buildCalendarURL строит URL страницы календаря из URL турнира
func buildCalendarURL(tournamentURL string) string {
	return strings.TrimSuffix(tournamentURL, "/") + "/calendar/"
}
