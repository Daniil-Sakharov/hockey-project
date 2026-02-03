// Полные профили игроков для поиска и страниц профилей
export interface MockPlayerProfile {
  id: string
  name: string
  birthDate: string
  birthYear: number
  team: string
  teamId: string
  position: 'forward' | 'defender' | 'goalie'
  jerseyNumber: number
  height?: number
  weight?: number
  handedness?: 'left' | 'right'
  city?: string
  stats: {
    games: number
    goals: number
    assists: number
    points: number
    plusMinus: number
    penaltyMinutes: number
  }
}

export const MOCK_PLAYERS: MockPlayerProfile[] = [
  { id: 'p1', name: 'Петров Дмитрий', birthDate: '2008-03-15', birthYear: 2008, team: 'АК Барс', teamId: 't1', position: 'forward', jerseyNumber: 17, height: 172, weight: 65, handedness: 'left', city: 'Казань', stats: { games: 22, goals: 28, assists: 19, points: 47, plusMinus: 18, penaltyMinutes: 12 } },
  { id: 'p2', name: 'Сидоров Артём', birthDate: '2008-07-22', birthYear: 2008, team: 'ЦСКА', teamId: 't2', position: 'forward', jerseyNumber: 9, height: 175, weight: 68, handedness: 'right', city: 'Москва', stats: { games: 24, goals: 22, assists: 24, points: 46, plusMinus: 14, penaltyMinutes: 8 } },
  { id: 'p3', name: 'Козлов Максим', birthDate: '2008-01-10', birthYear: 2008, team: 'СКА-Юность', teamId: 't3', position: 'forward', jerseyNumber: 11, height: 178, weight: 70, handedness: 'left', city: 'Санкт-Петербург', stats: { games: 23, goals: 25, assists: 18, points: 43, plusMinus: 16, penaltyMinutes: 14 } },
  { id: 'p4', name: 'Морозов Иван', birthDate: '2008-11-05', birthYear: 2008, team: 'Торос', teamId: 't4', position: 'forward', jerseyNumber: 22, height: 170, weight: 63, handedness: 'right', city: 'Нефтекамск', stats: { games: 21, goals: 20, assists: 22, points: 42, plusMinus: 12, penaltyMinutes: 6 } },
  { id: 'p5', name: 'Новиков Егор', birthDate: '2009-05-18', birthYear: 2009, team: 'Динамо СПб', teamId: 't5', position: 'forward', jerseyNumber: 7, height: 168, weight: 60, handedness: 'left', city: 'Санкт-Петербург', stats: { games: 22, goals: 24, assists: 15, points: 39, plusMinus: 10, penaltyMinutes: 10 } },
  { id: 'p6', name: 'Волков Кирилл', birthDate: '2008-09-30', birthYear: 2008, team: 'Спартак', teamId: 't6', position: 'forward', jerseyNumber: 14, height: 174, weight: 66, handedness: 'right', city: 'Москва', stats: { games: 24, goals: 18, assists: 20, points: 38, plusMinus: 8, penaltyMinutes: 18 } },
  { id: 'p7', name: 'Фёдоров Даниил', birthDate: '2008-06-12', birthYear: 2008, team: 'Саров', teamId: 't7', position: 'defender', jerseyNumber: 4, height: 180, weight: 75, handedness: 'left', city: 'Саров', stats: { games: 20, goals: 21, assists: 14, points: 35, plusMinus: 15, penaltyMinutes: 22 } },
  { id: 'p8', name: 'Орлов Тимофей', birthDate: '2009-02-28', birthYear: 2009, team: 'Кристалл-2', teamId: 't8', position: 'forward', jerseyNumber: 19, height: 166, weight: 58, handedness: 'left', city: 'Электросталь', stats: { games: 23, goals: 16, assists: 18, points: 34, plusMinus: -4, penaltyMinutes: 8 } },
  { id: 'p9', name: 'Беляев Роман', birthDate: '2008-12-01', birthYear: 2008, team: 'Мордовия', teamId: 't9', position: 'forward', jerseyNumber: 10, height: 171, weight: 64, handedness: 'right', city: 'Саранск', stats: { games: 21, goals: 19, assists: 13, points: 32, plusMinus: 6, penaltyMinutes: 14 } },
  { id: 'p10', name: 'Зайцев Матвей', birthDate: '2008-08-14', birthYear: 2008, team: 'АК Барс', teamId: 't1', position: 'forward', jerseyNumber: 23, height: 169, weight: 62, handedness: 'left', city: 'Казань', stats: { games: 22, goals: 15, assists: 16, points: 31, plusMinus: 14, penaltyMinutes: 4 } },
  { id: 'p11', name: 'Григорьев Алексей', birthDate: '2008-04-20', birthYear: 2008, team: 'АК Барс', teamId: 't1', position: 'defender', jerseyNumber: 3, height: 182, weight: 78, handedness: 'right', city: 'Казань', stats: { games: 22, goals: 5, assists: 18, points: 23, plusMinus: 20, penaltyMinutes: 28 } },
  { id: 'p12', name: 'Кузнецов Артём', birthDate: '2008-10-08', birthYear: 2008, team: 'Торос', teamId: 't4', position: 'goalie', jerseyNumber: 1, height: 176, weight: 72, handedness: 'left', city: 'Нефтекамск', stats: { games: 18, goals: 0, assists: 1, points: 1, plusMinus: 0, penaltyMinutes: 2 } },
  { id: 'p13', name: 'Смирнов Никита', birthDate: '2009-01-25', birthYear: 2009, team: 'ЦСКА', teamId: 't2', position: 'defender', jerseyNumber: 5, height: 179, weight: 73, handedness: 'left', city: 'Москва', stats: { games: 24, goals: 3, assists: 15, points: 18, plusMinus: 12, penaltyMinutes: 16 } },
  { id: 'p14', name: 'Попов Илья', birthDate: '2009-07-11', birthYear: 2009, team: 'СКА-Юность', teamId: 't3', position: 'goalie', jerseyNumber: 30, height: 180, weight: 74, handedness: 'left', city: 'Санкт-Петербург', stats: { games: 20, goals: 0, assists: 0, points: 0, plusMinus: 0, penaltyMinutes: 4 } },
  { id: 'p15', name: 'Лебедев Денис', birthDate: '2008-02-17', birthYear: 2008, team: 'Спартак', teamId: 't6', position: 'defender', jerseyNumber: 44, height: 183, weight: 80, handedness: 'right', city: 'Москва', stats: { games: 23, goals: 7, assists: 12, points: 19, plusMinus: 4, penaltyMinutes: 32 } },
]

// Команды
export interface MockTeam {
  id: string
  name: string
  city: string
  tournaments: string[]
  playersCount: number
  stats: { wins: number; losses: number; draws: number; goalsFor: number; goalsAgainst: number }
}

export const MOCK_TEAMS: MockTeam[] = [
  { id: 't1', name: 'АК Барс', city: 'Казань', tournaments: ['Первенство ПФО'], playersCount: 22, stats: { wins: 18, losses: 2, draws: 2, goalsFor: 98, goalsAgainst: 34 } },
  { id: 't2', name: 'ЦСКА', city: 'Москва', tournaments: ['Чемпионат Москвы'], playersCount: 24, stats: { wins: 16, losses: 5, draws: 3, goalsFor: 85, goalsAgainst: 42 } },
  { id: 't3', name: 'СКА-Юность', city: 'Санкт-Петербург', tournaments: ['Первенство Санкт-Петербурга'], playersCount: 23, stats: { wins: 15, losses: 4, draws: 3, goalsFor: 78, goalsAgainst: 38 } },
  { id: 't4', name: 'Торос', city: 'Нефтекамск', tournaments: ['Первенство ПФО'], playersCount: 21, stats: { wins: 17, losses: 4, draws: 1, goalsFor: 89, goalsAgainst: 41 } },
  { id: 't5', name: 'Динамо СПб', city: 'Санкт-Петербург', tournaments: ['Первенство Санкт-Петербурга'], playersCount: 22, stats: { wins: 12, losses: 7, draws: 3, goalsFor: 65, goalsAgainst: 52 } },
  { id: 't6', name: 'Спартак', city: 'Москва', tournaments: ['Чемпионат Москвы'], playersCount: 23, stats: { wins: 10, losses: 9, draws: 5, goalsFor: 58, goalsAgainst: 55 } },
  { id: 't7', name: 'Саров', city: 'Саров', tournaments: ['Первенство ПФО'], playersCount: 20, stats: { wins: 11, losses: 8, draws: 3, goalsFor: 62, goalsAgainst: 55 } },
  { id: 't8', name: 'Кристалл-2', city: 'Электросталь', tournaments: ['Первенство ПФО'], playersCount: 20, stats: { wins: 5, losses: 16, draws: 1, goalsFor: 38, goalsAgainst: 82 } },
  { id: 't9', name: 'Мордовия', city: 'Саранск', tournaments: ['Первенство ПФО'], playersCount: 21, stats: { wins: 14, losses: 5, draws: 3, goalsFor: 75, goalsAgainst: 48 } },
]

// Календарь матчей (расширенный)
export interface MockCalendarMatch {
  id: string
  homeTeam: string
  awayTeam: string
  homeScore: number | null
  awayScore: number | null
  date: string
  time: string
  tournament: string
  status: 'finished' | 'scheduled'
  venue?: string
}

export const MOCK_CALENDAR_MATCHES: MockCalendarMatch[] = [
  { id: 'c1', homeTeam: 'АК Барс', awayTeam: 'Сокол', homeScore: 7, awayScore: 3, date: '2025-01-28', time: '15:00', tournament: 'Первенство ПФО', status: 'finished', venue: 'Татнефть Арена' },
  { id: 'c2', homeTeam: 'Торос', awayTeam: 'Кристалл-2', homeScore: 9, awayScore: 3, date: '2025-01-28', time: '17:00', tournament: 'Первенство ПФО', status: 'finished', venue: 'ЛД Нефтекамск' },
  { id: 'c3', homeTeam: 'СКА-Юность', awayTeam: 'Динамо СПб', homeScore: 4, awayScore: 2, date: '2025-01-27', time: '14:00', tournament: 'Первенство Санкт-Петербурга', status: 'finished', venue: 'Хоккейный Город' },
  { id: 'c4', homeTeam: 'ЦСКА', awayTeam: 'Спартак', homeScore: 5, awayScore: 1, date: '2025-01-26', time: '12:00', tournament: 'Чемпионат Москвы', status: 'finished', venue: 'ЛД ЦСКА' },
  { id: 'c5', homeTeam: 'Мордовия', awayTeam: 'Саров', homeScore: 6, awayScore: 3, date: '2025-01-27', time: '16:00', tournament: 'Первенство ПФО', status: 'finished', venue: 'ЛД Саранск' },
  { id: 'c6', homeTeam: 'Сокол', awayTeam: 'Торос', homeScore: null, awayScore: null, date: '2025-02-03', time: '15:00', tournament: 'Первенство ПФО', status: 'scheduled', venue: 'ЛД Сокол' },
  { id: 'c7', homeTeam: 'Кристалл-2', awayTeam: 'АК Барс', homeScore: null, awayScore: null, date: '2025-02-03', time: '17:00', tournament: 'Первенство ПФО', status: 'scheduled', venue: 'ЛД Электросталь' },
  { id: 'c8', homeTeam: 'Динамо СПб', awayTeam: 'СКА-Юность', homeScore: null, awayScore: null, date: '2025-02-04', time: '14:00', tournament: 'Первенство Санкт-Петербурга', status: 'scheduled', venue: 'ЛД Динамо' },
  { id: 'c9', homeTeam: 'Спартак', awayTeam: 'ЦСКА', homeScore: null, awayScore: null, date: '2025-02-05', time: '16:00', tournament: 'Чемпионат Москвы', status: 'scheduled', venue: 'ЛД Сокольники' },
  { id: 'c10', homeTeam: 'АК Барс', awayTeam: 'Мордовия', homeScore: null, awayScore: null, date: '2025-02-07', time: '15:00', tournament: 'Первенство ПФО', status: 'scheduled', venue: 'Татнефть Арена' },
  { id: 'c11', homeTeam: 'Торос', awayTeam: 'Саров', homeScore: null, awayScore: null, date: '2025-02-08', time: '12:00', tournament: 'Первенство ПФО', status: 'scheduled', venue: 'ЛД Нефтекамск' },
  { id: 'c12', homeTeam: 'СКА-Юность', awayTeam: 'ЦСКА', homeScore: null, awayScore: null, date: '2025-02-10', time: '14:00', tournament: 'Кубок Регионов', status: 'scheduled', venue: 'Хоккейный Город' },
]

// Прогнозы
export interface MockPredictionMatch {
  id: string
  homeTeam: string
  awayTeam: string
  date: string
  time: string
  tournament: string
  venue?: string
  userPrediction?: { homeScore: number; awayScore: number }
  actualResult?: { homeScore: number; awayScore: number }
  pointsEarned?: number
}

export const MOCK_PREDICTION_MATCHES: MockPredictionMatch[] = [
  { id: 'pr1', homeTeam: 'Сокол', awayTeam: 'Торос', date: '2025-02-03', time: '15:00', tournament: 'Первенство ПФО' },
  { id: 'pr2', homeTeam: 'Кристалл-2', awayTeam: 'АК Барс', date: '2025-02-03', time: '17:00', tournament: 'Первенство ПФО' },
  { id: 'pr3', homeTeam: 'Динамо СПб', awayTeam: 'СКА-Юность', date: '2025-02-04', time: '14:00', tournament: 'Первенство Санкт-Петербурга' },
  { id: 'pr4', homeTeam: 'Спартак', awayTeam: 'ЦСКА', date: '2025-02-05', time: '16:00', tournament: 'Чемпионат Москвы' },
]

// Рейтинг прогнозистов
export interface MockPredictor {
  position: number
  name: string
  correctScores: number
  correctOutcomes: number
  totalPredictions: number
  points: number
}

export const MOCK_PREDICTORS_RANKING: MockPredictor[] = [
  { position: 1, name: 'HockeyFan2008', correctScores: 12, correctOutcomes: 45, totalPredictions: 60, points: 159 },
  { position: 2, name: 'IceKing', correctScores: 10, correctOutcomes: 42, totalPredictions: 58, points: 144 },
  { position: 3, name: 'PuckMaster', correctScores: 8, correctOutcomes: 40, totalPredictions: 55, points: 128 },
  { position: 4, name: 'GoalPredictor', correctScores: 9, correctOutcomes: 38, totalPredictions: 56, points: 123 },
  { position: 5, name: 'SnipeCity', correctScores: 7, correctOutcomes: 36, totalPredictions: 52, points: 115 },
]
