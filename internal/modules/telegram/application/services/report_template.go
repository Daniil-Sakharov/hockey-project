package services

const reportHTMLTemplate = `<!DOCTYPE html>
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
        }
        .header h1 { font-size: 24px; font-weight: 700; }
        .header p { color: var(--accent-light); font-size: 14px; }
        .player-card {
            background: var(--white);
            border-radius: 16px;
            padding: 24px;
            margin-bottom: 20px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
        }
        .player-name { font-size: 28px; font-weight: 700; color: var(--primary-dark); margin-bottom: 12px; }
        .player-details { display: flex; flex-wrap: wrap; gap: 16px; }
        .player-detail { display: flex; align-items: center; gap: 8px; color: var(--gray); font-size: 14px; }
        .stats-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 20px; }
        .stat-card {
            background: var(--white);
            border-radius: 12px;
            padding: 20px;
            text-align: center;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
        .stat-value { font-size: 32px; font-weight: 700; color: var(--accent); }
        .stat-label { font-size: 12px; color: var(--gray); text-transform: uppercase; }
        .stat-avg { font-size: 11px; color: var(--gray); margin-top: 4px; }
        .charts-section { display: grid; grid-template-columns: repeat(2, 1fr); gap: 20px; margin-bottom: 20px; }
        .chart-card {
            background: var(--white);
            border-radius: 12px;
            padding: 20px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        }
        .chart-title { font-size: 14px; font-weight: 600; color: var(--primary-dark); margin-bottom: 16px; text-align: center; }
        .chart-container { display: flex; justify-content: center; }
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
        .season-header {
            font-size: 16px; font-weight: 600; color: var(--accent);
            margin: 16px 0 12px; padding: 8px 12px;
            background: var(--ice); border-radius: 8px;
        }
        .tournament-card { border: 1px solid var(--ice); border-radius: 10px; padding: 16px; margin-bottom: 12px; }
        .tournament-header { display: flex; justify-content: space-between; margin-bottom: 8px; }
        .tournament-name { font-weight: 600; color: var(--primary-dark); font-size: 14px; }
        .tournament-team { font-size: 12px; color: var(--gray); }
        .tournament-stats { display: flex; gap: 8px; }
        .tournament-stat { background: var(--ice); padding: 4px 10px; border-radius: 6px; font-size: 12px; font-weight: 600; }
        .tournament-stat.goals { background: var(--accent); color: var(--white); }
        .empty-state { text-align: center; padding: 40px; color: var(--gray); }
        @media (max-width: 768px) {
            .stats-grid { grid-template-columns: repeat(2, 1fr); }
            .charts-section { grid-template-columns: 1fr; }
            .stats-table { grid-template-columns: repeat(2, 1fr); }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>HockeyStats</h1>
            <p>Полный отчет игрока</p>
        </div>

        <div class="player-card">
            <div class="player-name">{{.Report.Player.Name}}</div>
            <div class="player-details">
                <div class="player-detail">📅 {{.Report.Player.BirthYear}} г.р.</div>
                {{if .Report.Player.Position}}<div class="player-detail">🏒 {{.Report.Player.Position}}</div>{{end}}
                {{if .Report.Player.Height}}<div class="player-detail">📏 {{.Report.Player.Height}} см</div>{{end}}
                {{if .Report.Player.Weight}}<div class="player-detail">⚖️ {{.Report.Player.Weight}} кг</div>{{end}}
                {{if .Report.Player.Team}}<div class="player-detail">🏢 {{.Report.Player.Team}}</div>{{end}}
                {{if .Report.Player.Region}}<div class="player-detail">📍 {{.Report.Player.Region}}</div>{{end}}
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
            </div>
        </div>

        {{if .Report.Tournaments}}
        <div class="tournaments-section">
            <div class="section-title">История выступлений</div>
            {{$currentSeason := ""}}
            {{range .Report.Tournaments}}
                {{if ne .Season $currentSeason}}
                    <div class="season-header">Сезон {{.Season}}</div>
                    {{$currentSeason = .Season}}
                {{end}}
                <div class="tournament-card">
                    <div class="tournament-header">
                        <div>
                            <div class="tournament-name">{{.TournamentName}}</div>
                            <div class="tournament-team">{{.TeamName}}</div>
                        </div>
                        <div class="tournament-stats">
                            <span class="tournament-stat">{{.Games}} игр</span>
                            <span class="tournament-stat goals">{{.Goals}}+{{.Assists}}={{.Points}}</span>
                        </div>
                    </div>
                </div>
            {{end}}
        </div>
        {{end}}
        {{else}}
        <div class="detailed-section">
            <p class="empty-state">Нет статистики для отображения</p>
        </div>
        {{end}}
    </div>
</body>
</html>`
