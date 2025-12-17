package report

// HTMLTemplate содержит HTML шаблон отчёта
const HTMLTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Report.Player.Name}} - HockeyStats</title>
    <style>
        :root {
            --primary-dark: #0a1628;
            --primary: #1a3a5c;
            --accent: #4a90d9;
            --accent-light: #7bb8e8;
            --white: #ffffff;
            --ice: #e8f4fc;
            --gray: #6b7280;
            --gray-light: #f3f4f6;
        }
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', Roboto, -apple-system, sans-serif;
            background: var(--gray-light);
            color: var(--primary-dark);
            line-height: 1.5;
        }
        .container { max-width: 1000px; margin: 0 auto; padding: 20px; }

        .header {
            background: linear-gradient(135deg, var(--primary-dark) 0%, var(--primary) 100%);
            color: var(--white);
            padding: 24px;
            border-radius: 16px;
            margin-bottom: 20px;
            display: flex;
            align-items: center;
            gap: 16px;
        }
        .header .logo { width: 60px; height: 60px; }
        .header h1 { font-size: 24px; font-weight: 700; }
        .header p { color: var(--accent-light); font-size: 14px; }

        .player-card {
            background: var(--white);
            border-radius: 16px;
            padding: 24px;
            margin-bottom: 20px;
            display: flex;
            gap: 24px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
        }
        .player-photo {
            width: 120px; height: 150px;
            background: var(--ice);
            border-radius: 12px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--accent);
            font-size: 48px;
            flex-shrink: 0;
        }
        .player-info { flex: 1; }
        .player-name { font-size: 28px; font-weight: 700; color: var(--primary-dark); margin-bottom: 12px; }
        .player-details { display: flex; flex-wrap: wrap; gap: 16px; }
        .player-detail { display: flex; align-items: center; gap: 8px; color: var(--gray); font-size: 14px; }
        .player-detail .icon { width: 20px; text-align: center; }

        .stats-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 20px; }
        .stat-card {
            background: var(--white);
            border-radius: 12px;
            padding: 20px;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
        .stat-value { font-size: 32px; font-weight: 700; color: var(--accent); }
        .stat-label { font-size: 12px; color: var(--gray); text-transform: uppercase; letter-spacing: 0.5px; }
        .stat-avg { font-size: 11px; color: var(--gray); margin-top: 4px; }

        .charts-section { display: grid; grid-template-columns: repeat(2, 1fr); gap: 20px; margin-bottom: 20px; }
        .chart-card {
            background: var(--white);
            border-radius: 12px;
            padding: 20px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
        .chart-title { font-size: 14px; font-weight: 600; color: var(--primary-dark); margin-bottom: 16px; text-align: center; }
        .chart-container { display: flex; justify-content: center; align-items: center; }
        .chart-container svg { max-width: 100%; height: auto; }

        .detailed-section {
            background: var(--white);
            border-radius: 12px;
            padding: 24px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
        .section-title {
            font-size: 18px; font-weight: 600; color: var(--primary-dark);
            margin-bottom: 16px; padding-bottom: 8px; border-bottom: 2px solid var(--ice);
        }
        .stats-table { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
        .stat-row {
            display: flex; justify-content: space-between;
            padding: 8px 12px; background: var(--ice); border-radius: 8px;
        }
        .stat-row .label { color: var(--gray); font-size: 13px; }
        .stat-row .value { font-weight: 600; color: var(--primary-dark); }

        .tournaments-section {
            background: var(--white);
            border-radius: 12px;
            padding: 24px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
        .tournament-group { margin-bottom: 24px; }
        .tournament-group:last-child { margin-bottom: 0; }
        .season-header {
            font-size: 16px; font-weight: 600; color: var(--accent);
            margin-bottom: 12px; padding: 8px 12px;
            background: var(--ice); border-radius: 8px;
        }
        .tournament-card { border: 1px solid var(--ice); border-radius: 10px; padding: 16px; margin-bottom: 12px; }
        .tournament-card:last-child { margin-bottom: 0; }
        .tournament-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 12px; }
        .tournament-name { font-weight: 600; color: var(--primary-dark); font-size: 14px; }
        .tournament-team { font-size: 12px; color: var(--gray); margin-top: 4px; }
        .tournament-main-stats { display: flex; gap: 8px; }
        .tournament-stat { background: var(--ice); padding: 4px 10px; border-radius: 6px; font-size: 12px; font-weight: 600; }
        .tournament-stat.goals { background: var(--accent); color: var(--white); }
        .tournament-detailed {
            display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px;
            padding-top: 12px; border-top: 1px solid var(--ice);
        }
        .tournament-detail-item { font-size: 11px; color: var(--gray); }
        .tournament-detail-item span { font-weight: 600; color: var(--primary-dark); }

        .empty-state { text-align: center; padding: 40px; color: var(--gray); }

        @media (max-width: 768px) {
            .stats-grid { grid-template-columns: repeat(2, 1fr); }
            .charts-section { grid-template-columns: 1fr; }
            .stats-table { grid-template-columns: repeat(2, 1fr); }
            .tournament-detailed { grid-template-columns: repeat(2, 1fr); }
            .player-card { flex-direction: column; align-items: center; text-align: center; }
            .player-details { justify-content: center; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <svg class="logo" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
                <circle cx="50" cy="50" r="48" fill="#1a3a5c" stroke="#4a90d9" stroke-width="2"/>
                <path d="M30 70 L50 25 L70 70 M35 55 L65 55" stroke="#7bb8e8" stroke-width="4" fill="none" stroke-linecap="round"/>
                <circle cx="50" cy="75" r="4" fill="#4a90d9"/>
            </svg>
            <div>
                <h1>HockeyStats</h1>
                <p>Полный отчет игрока</p>
            </div>
        </div>

        <div class="player-card">
            <div class="player-photo">
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                    <circle cx="12" cy="7" r="4"/>
                </svg>
            </div>
            <div class="player-info">
                <div class="player-name">{{.Report.Player.Name}}</div>
                <div class="player-details">
                    <div class="player-detail"><span class="icon">📅</span><span>{{.Report.Player.BirthYear}} г.р.</span></div>
                    {{if .Report.Player.Position}}<div class="player-detail"><span class="icon">🏒</span><span>{{.Report.Player.Position}}</span></div>{{end}}
                    {{if .Report.Player.Height}}<div class="player-detail"><span class="icon">📏</span><span>{{.Report.Player.Height}} см</span></div>{{end}}
                    {{if .Report.Player.Weight}}<div class="player-detail"><span class="icon">⚖️</span><span>{{.Report.Player.Weight}} кг</span></div>{{end}}
                    <div class="player-detail"><span class="icon">🏢</span><span>{{.Report.Player.Team}}</span></div>
                    <div class="player-detail"><span class="icon">📍</span><span>{{.Report.Player.Region}}</span></div>
                </div>
            </div>
        </div>

        {{if .Report.HasStats}}
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-value">{{.Report.TotalStats.TotalGames}}</div>
                <div class="stat-label">Игр</div>
                <div class="stat-avg">{{.Report.TotalStats.TotalTournaments}} турниров</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{.Report.TotalStats.TotalGoals}}</div>
                <div class="stat-label">Голов</div>
                <div class="stat-avg">{{formatFloat .Report.TotalStats.GoalsPerGame}} за игру</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{.Report.TotalStats.TotalAssists}}</div>
                <div class="stat-label">Передач</div>
                <div class="stat-avg">{{formatFloat .Report.TotalStats.AssistsPerGame}} за игру</div>
            </div>
            <div class="stat-card">
                <div class="stat-value">{{.Report.TotalStats.TotalPoints}}</div>
                <div class="stat-label">Очков</div>
                <div class="stat-avg">{{formatFloat .Report.TotalStats.PointsPerGame}} за игру</div>
            </div>
        </div>

        <div class="charts-section">
            {{if .Report.HasDetailedStats}}
            <div class="chart-card">
                <div class="chart-title">Распределение голов по типу</div>
                <div class="chart-container">{{.Charts.GoalsTypePie}}</div>
            </div>
            {{end}}

            <div class="chart-card">
                <div class="chart-title">Голы по периодам</div>
                <div class="chart-container">{{.Charts.PeriodBar}}</div>
            </div>

            {{if .Report.HasMultipleSeasons}}
            <div class="chart-card">
                <div class="chart-title">Прогресс по сезонам</div>
                <div class="chart-container">{{.Charts.ProgressLine}}</div>
            </div>
            {{end}}

            <div class="chart-card">
                <div class="chart-title">Профиль игрока</div>
                <div class="chart-container">{{.Charts.ProfileRadar}}</div>
            </div>
        </div>

        <div class="detailed-section">
            <div class="section-title">Детальная статистика</div>
            <div class="stats-table">
                <div class="stat-row"><span class="label">+/-</span><span class="value">{{plusMinusFormat .Report.TotalStats.TotalPlusMinus}}</span></div>
                <div class="stat-row"><span class="label">Штраф. минут</span><span class="value">{{.Report.TotalStats.TotalPenalties}}</span></div>
                <div class="stat-row"><span class="label">Хет-трики</span><span class="value">{{.Report.TotalStats.TotalHatTricks}}</span></div>
                <div class="stat-row"><span class="label">Победные голы</span><span class="value">{{.Report.TotalStats.TotalWinningGoals}}</span></div>
                {{if .Report.HasDetailedStats}}
                <div class="stat-row"><span class="label">Голы в равных</span><span class="value">{{.Report.TotalStats.GoalsEvenStrength}}</span></div>
                <div class="stat-row"><span class="label">Голы в большинстве</span><span class="value">{{.Report.TotalStats.GoalsPowerPlay}}</span></div>
                <div class="stat-row"><span class="label">Голы в меньшинстве</span><span class="value">{{.Report.TotalStats.GoalsShortHanded}}</span></div>
                <div class="stat-row"><span class="label">Голы в 1 периоде</span><span class="value">{{.Report.TotalStats.GoalsPeriod1}}</span></div>
                <div class="stat-row"><span class="label">Голы во 2 периоде</span><span class="value">{{.Report.TotalStats.GoalsPeriod2}}</span></div>
                <div class="stat-row"><span class="label">Голы в 3 периоде</span><span class="value">{{.Report.TotalStats.GoalsPeriod3}}</span></div>
                <div class="stat-row"><span class="label">Голы в овертайме</span><span class="value">{{.Report.TotalStats.GoalsOvertime}}</span></div>
                {{end}}
            </div>
        </div>

        <div class="tournaments-section">
            <div class="section-title">История выступлений</div>
            {{if .Report.Tournaments}}
                {{$currentSeason := ""}}
                {{range .Report.Tournaments}}
                    {{if ne .Season $currentSeason}}
                        {{if ne $currentSeason ""}}</div>{{end}}
                        <div class="tournament-group">
                        <div class="season-header">Сезон {{.Season}}</div>
                        {{$currentSeason = .Season}}
                    {{end}}
                    <div class="tournament-card">
                        <div class="tournament-header">
                            <div>
                                <div class="tournament-name">{{.TournamentName}}</div>
                                <div class="tournament-team">{{.TeamName}}{{if .GroupName}} • {{.GroupName}}{{end}}</div>
                            </div>
                            <div class="tournament-main-stats">
                                <span class="tournament-stat">{{.Games}} игр</span>
                                <span class="tournament-stat goals">{{.Goals}}+{{.Assists}}={{.Points}}</span>
                            </div>
                        </div>
                        <div class="tournament-detailed">
                            <div class="tournament-detail-item">+/-: <span>{{plusMinusFormat .PlusMinus}}</span></div>
                            <div class="tournament-detail-item">Штраф: <span>{{.PenaltyMinutes}} мин</span></div>
                            <div class="tournament-detail-item">Хет-трики: <span>{{.HatTricks}}</span></div>
                            <div class="tournament-detail-item">Поб. голы: <span>{{.GameWinningGoals}}</span></div>
                        </div>
                    </div>
                {{end}}
                </div>
            {{else}}
                <p class="empty-state">Нет данных о турнирах</p>
            {{end}}
        </div>
        {{else}}
        <div class="detailed-section">
            <p class="empty-state">Нет статистики для отображения</p>
        </div>
        {{end}}
    </div>
</body>
</html>`
