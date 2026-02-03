package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

const (
	teamWorkers = 10 // Worker Pool для команд (10 параллельно)
)

// processTournaments обрабатывает турниры (команды → игроки) с Worker Pool
func (s *orchestratorService) processTournaments(
	ctx context.Context,
	tournaments []*entities.Tournament,
) error {
	logger.Info(ctx, "")
	logger.Info(ctx, "📊 STAGE 2: Processing tournaments...")

	totalTeams := 0
	totalErrors := 0

	for idx, t := range tournaments {
		logger.Info(ctx, fmt.Sprintf("  🏆 Tournament %d/%d: %s (ID: %s, URL: %s)",
			idx+1, len(tournaments), t.Name, t.ID, t.URL))

		// Парсим команды с контекстом года/группы
		logger.Info(ctx, "    🔍 Parsing teams with year/group context...")
		teamsWithContext, err := s.juniorService.ParseTeams(ctx, t.Domain, t.URL, t.FallbackBirthYears...)
		if err != nil {
			logger.Warn(ctx, fmt.Sprintf("    ⚠️  Failed to parse teams: %v", err))
			logger.Warn(ctx, "    ⏭️  SKIPPING tournament, continuing with next...")
			totalErrors++
			continue
		}

		logger.Info(ctx, fmt.Sprintf("    ✅ Found %d teams from page", len(teamsWithContext)))

		// Собираем birth_year_groups из контекста команд и сохраняем в турнир
		if byg := collectBirthYearGroups(teamsWithContext); byg != nil {
			if raw, err := json.Marshal(byg); err == nil {
				if err := s.tournamentRepo.UpdateBirthYearGroups(ctx, t.ID, string(raw)); err != nil {
					logger.Warn(ctx, fmt.Sprintf("    ⚠️  Failed to update birth_year_groups: %v", err))
				} else {
					logger.Info(ctx, fmt.Sprintf("    📅 Saved birth_year_groups: %d years", len(byg)))
				}
			}
		}

		// Дедупликация для сохранения в таблицу teams (одна запись на команду)
		uniqueTeams := uniqueTeamsForSave(teamsWithContext)
		logger.Info(ctx, fmt.Sprintf("    📊 Unique teams: %d (from %d contexts)", len(uniqueTeams), len(teamsWithContext)))

		// Сохраняем уникальные команды в БД
		logger.Info(ctx, fmt.Sprintf("    💾 Saving teams for %s...", t.Name))
		teams, err := s.SaveTeams(ctx, uniqueTeams, t.ID)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("    ❌ Failed to save teams: %v", err))
			totalErrors++
			continue
		}
		logger.Info(ctx, fmt.Sprintf("    ✅ Saved %d teams", len(teams)))

		if len(teams) == 0 {
			logger.Info(ctx, "    ℹ️  No teams to process, moving to next tournament")
			continue
		}

		// Строим задачи: одна на каждую комбинацию (team, year, group)
		teamTasks := buildTeamTasks(teamsWithContext, teams, t)
		logger.Info(ctx, fmt.Sprintf("    📋 Team tasks: %d (team × year × group combinations)", len(teamTasks)))

		// Worker Pool для параллельного парсинга команд
		logger.Info(ctx, fmt.Sprintf("    🚀 Starting team worker pool (%d workers) to parse players...", teamWorkers))
		pool := NewTeamWorkerPool(ctx, s, teamWorkers)
		pool.Start()

		teamProcessed := 0
		teamErrors := 0
		done := make(chan struct{})

		logger.Debug(ctx, "    📥 Starting result reader goroutine...")
		go func() {
			resultCount := 0
			for result := range pool.Results() {
				resultCount++
				if result.Error != nil {
					logger.Warn(ctx, fmt.Sprintf("      ⚠️  Team error: %s - %v", result.TeamName, result.Error))
					teamErrors++
				} else {
					teamProcessed++
				}
			}
			logger.Debug(ctx, fmt.Sprintf("    📥 Reader done, read %d results", resultCount))
			close(done)
		}()

		// Добавляем задачи в очередь (каждая комбинация team+year+group)
		for _, task := range teamTasks {
			pool.AddTask(task)
		}

		pool.Close()
		logger.Debug(ctx, "    ⏳ Waiting for workers to finish...")
		pool.Wait()
		logger.Debug(ctx, "    ✅ All workers finished!")
		<-done
		logger.Debug(ctx, "    ✅ Result reader finished!")

		totalTeams += teamProcessed
		totalErrors += teamErrors

		logger.Info(ctx, fmt.Sprintf("    📊 Tournament result: %d teams processed, %d errors", teamProcessed, teamErrors))
		logger.Info(ctx, "    ✅ Tournament COMPLETED")
	}

	logger.Info(ctx, "")
	logger.Info(ctx, "================================================================================")
	logger.Info(ctx, "📊 FINAL STATISTICS:")
	logger.Info(ctx, fmt.Sprintf("  Tournaments processed: %d", len(tournaments)))
	logger.Info(ctx, fmt.Sprintf("  Teams processed: %d", totalTeams))
	logger.Info(ctx, fmt.Sprintf("  Errors: %d", totalErrors))
	logger.Info(ctx, "================================================================================")

	return nil
}

// uniqueTeamsForSave дедуплицирует команды по ID (для сохранения в таблицу teams)
func uniqueTeamsForSave(teamsWithContext []types.TeamWithContext) []types.TeamDTO {
	seen := make(map[string]bool)
	var unique []types.TeamDTO

	for _, twc := range teamsWithContext {
		teamID := entities.ExtractTeamIDFromURLLegacy(twc.Team.URL)
		if !seen[teamID] {
			seen[teamID] = true
			unique = append(unique, twc.Team)
		}
	}

	return unique
}

// buildTeamTasks создаёт задачи для каждой комбинации (team, year, group)
func buildTeamTasks(
	teamsWithContext []types.TeamWithContext,
	savedTeams []*entities.Team,
	tournament *entities.Tournament,
) []TeamTask {
	// Маппинг teamID → сохранённая сущность
	teamByID := make(map[string]*entities.Team, len(savedTeams))
	for _, t := range savedTeams {
		teamByID[t.ID] = t
	}

	// Дедупликация по (teamID, year, group) — одна задача на комбинацию
	type taskKey struct {
		teamID    string
		birthYear int
		groupName string
	}
	seen := make(map[taskKey]bool)
	var tasks []TeamTask

	for _, twc := range teamsWithContext {
		teamID := entities.ExtractTeamIDFromURLLegacy(twc.Team.URL)
		team, ok := teamByID[teamID]
		if !ok {
			continue
		}

		year := 0
		if twc.BirthYear != nil {
			year = *twc.BirthYear
		}
		group := ""
		if twc.GroupName != nil {
			group = *twc.GroupName
		}

		key := taskKey{teamID: teamID, birthYear: year, groupName: group}
		if seen[key] {
			continue
		}
		seen[key] = true

		tasks = append(tasks, TeamTask{
			Team:       team,
			Tournament: tournament,
			BirthYear:  twc.BirthYear,
			GroupName:  twc.GroupName,
			Index:      len(tasks) + 1,
			Total:      0, // заполним после
		})
	}

	// Заполняем Total
	for i := range tasks {
		tasks[i].Total = len(tasks)
	}

	return tasks
}

// collectBirthYearGroups собирает map[year_string][]group_name из TeamWithContext
func collectBirthYearGroups(teamsWithContext []types.TeamWithContext) map[string][]string {
	// year → set of group names
	yearGroups := make(map[int]map[string]bool)

	for _, twc := range teamsWithContext {
		if twc.BirthYear == nil || *twc.BirthYear == 0 {
			continue
		}
		year := *twc.BirthYear
		if yearGroups[year] == nil {
			yearGroups[year] = make(map[string]bool)
		}
		if twc.GroupName != nil && *twc.GroupName != "" {
			yearGroups[year][*twc.GroupName] = true
		}
	}

	if len(yearGroups) == 0 {
		return nil
	}

	result := make(map[string][]string, len(yearGroups))
	for year, groups := range yearGroups {
		key := fmt.Sprintf("%d", year)
		names := make([]string, 0, len(groups))
		for g := range groups {
			names = append(names, g)
		}
		sort.Strings(names)
		result[key] = names
	}
	return result
}
