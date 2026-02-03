// Статистика платформы
export const MOCK_PLATFORM_STATS = {
  tournaments: 47,
  teams: 312,
  players: 4850,
  matches: 1236,
}

// Турниры
export interface MockTournament {
  id: string
  name: string
  source: string
  domain: string
  teamsCount: number
  matchesCount: number
  status: 'active' | 'finished'
  season: string
}

export const MOCK_TOURNAMENTS: MockTournament[] = [
  {
    id: 'tour-1',
    name: 'Первенство Санкт-Петербурга',
    source: 'fhspb',
    domain: 'fhspb.ru',
    teamsCount: 24,
    matchesCount: 156,
    status: 'active',
    season: '2024/25',
  },
  {
    id: 'tour-2',
    name: 'Первенство ПФО',
    source: 'junior',
    domain: 'junior.fhr.ru',
    teamsCount: 32,
    matchesCount: 248,
    status: 'active',
    season: '2024/25',
  },
  {
    id: 'tour-3',
    name: 'Кубок Регионов',
    source: 'junior',
    domain: 'junior.fhr.ru',
    teamsCount: 16,
    matchesCount: 64,
    status: 'finished',
    season: '2024/25',
  },
  {
    id: 'tour-4',
    name: 'Чемпионат Москвы',
    source: 'fhmoscow',
    domain: 'fhmoscow.ru',
    teamsCount: 20,
    matchesCount: 190,
    status: 'active',
    season: '2024/25',
  },
  {
    id: 'tour-5',
    name: 'Открытое Первенство Республики Татарстан',
    source: 'junior',
    domain: 'junior.fhr.ru',
    teamsCount: 18,
    matchesCount: 102,
    status: 'active',
    season: '2024/25',
  },
  {
    id: 'tour-6',
    name: 'Турнир памяти В.В. Тихонова',
    source: 'junior',
    domain: 'junior.fhr.ru',
    teamsCount: 8,
    matchesCount: 28,
    status: 'finished',
    season: '2024/25',
  },
]

// Последние матчи
export interface MockRecentMatch {
  id: string
  homeTeam: string
  awayTeam: string
  homeScore: number
  awayScore: number
  date: string
  tournament: string
  status: 'finished' | 'live' | 'scheduled'
}

export const MOCK_RECENT_MATCHES: MockRecentMatch[] = [
  {
    id: 'match-1',
    homeTeam: 'АК Барс',
    awayTeam: 'Сокол',
    homeScore: 7,
    awayScore: 3,
    date: '2025-01-28',
    tournament: 'Первенство ПФО',
    status: 'finished',
  },
  {
    id: 'match-2',
    homeTeam: 'Торос',
    awayTeam: 'Кристалл-2',
    homeScore: 9,
    awayScore: 3,
    date: '2025-01-28',
    tournament: 'Первенство ПФО',
    status: 'finished',
  },
  {
    id: 'match-3',
    homeTeam: 'Саров',
    awayTeam: 'Мордовия',
    homeScore: 3,
    awayScore: 6,
    date: '2025-01-27',
    tournament: 'Первенство ПФО',
    status: 'finished',
  },
  {
    id: 'match-4',
    homeTeam: 'СКА-Юность',
    awayTeam: 'Динамо СПб',
    homeScore: 4,
    awayScore: 2,
    date: '2025-01-27',
    tournament: 'Первенство Санкт-Петербурга',
    status: 'finished',
  },
  {
    id: 'match-5',
    homeTeam: 'ЦСКА',
    awayTeam: 'Спартак',
    homeScore: 5,
    awayScore: 1,
    date: '2025-01-26',
    tournament: 'Чемпионат Москвы',
    status: 'finished',
  },
]

// Топ бомбардиры
export interface MockTopScorer {
  id: string
  name: string
  team: string
  goals: number
  assists: number
  points: number
  games: number
}

export const MOCK_TOP_SCORERS: MockTopScorer[] = [
  { id: 'p1', name: 'Петров Дмитрий', team: 'АК Барс', goals: 28, assists: 19, points: 47, games: 22 },
  { id: 'p2', name: 'Сидоров Артём', team: 'ЦСКА', goals: 22, assists: 24, points: 46, games: 24 },
  { id: 'p3', name: 'Козлов Максим', team: 'СКА-Юность', goals: 25, assists: 18, points: 43, games: 23 },
  { id: 'p4', name: 'Морозов Иван', team: 'Торос', goals: 20, assists: 22, points: 42, games: 21 },
  { id: 'p5', name: 'Новиков Егор', team: 'Динамо СПб', goals: 24, assists: 15, points: 39, games: 22 },
  { id: 'p6', name: 'Волков Кирилл', team: 'Спартак', goals: 18, assists: 20, points: 38, games: 24 },
  { id: 'p7', name: 'Фёдоров Даниил', team: 'Саров', goals: 21, assists: 14, points: 35, games: 20 },
  { id: 'p8', name: 'Орлов Тимофей', team: 'Кристалл-2', goals: 16, assists: 18, points: 34, games: 23 },
  { id: 'p9', name: 'Беляев Роман', team: 'Мордовия', goals: 19, assists: 13, points: 32, games: 21 },
  { id: 'p10', name: 'Зайцев Матвей', team: 'АК Барс', goals: 15, assists: 16, points: 31, games: 22 },
]

// Турнирная таблица
export interface MockStanding {
  position: number
  team: string
  games: number
  wins: number
  draws: number
  losses: number
  goalsFor: number
  goalsAgainst: number
  points: number
}

export const MOCK_TOURNAMENT_STANDINGS: MockStanding[] = [
  { position: 1, team: 'АК Барс', games: 22, wins: 18, draws: 2, losses: 2, goalsFor: 98, goalsAgainst: 34, points: 56 },
  { position: 2, team: 'Торос', games: 22, wins: 17, draws: 1, losses: 4, goalsFor: 89, goalsAgainst: 41, points: 52 },
  { position: 3, team: 'Мордовия', games: 22, wins: 14, draws: 3, losses: 5, goalsFor: 75, goalsAgainst: 48, points: 45 },
  { position: 4, team: 'Нефтяник', games: 22, wins: 13, draws: 2, losses: 7, goalsFor: 68, goalsAgainst: 52, points: 41 },
  { position: 5, team: 'Саров', games: 22, wins: 11, draws: 3, losses: 8, goalsFor: 62, goalsAgainst: 55, points: 36 },
  { position: 6, team: 'Сокол', games: 22, wins: 8, draws: 2, losses: 12, goalsFor: 51, goalsAgainst: 67, points: 26 },
  { position: 7, team: 'Кристалл-2', games: 22, wins: 5, draws: 1, losses: 16, goalsFor: 38, goalsAgainst: 82, points: 16 },
  { position: 8, team: 'Спутник', games: 22, wins: 2, draws: 0, losses: 20, goalsFor: 25, goalsAgainst: 95, points: 6 },
]

// Матчи турнира
export interface MockTournamentMatch {
  id: string
  homeTeam: string
  awayTeam: string
  homeScore: number | null
  awayScore: number | null
  date: string
  time: string
  status: 'finished' | 'scheduled'
  groupName?: string
}

export const MOCK_TOURNAMENT_MATCHES: MockTournamentMatch[] = [
  { id: 'tm-1', homeTeam: 'АК Барс', awayTeam: 'Сокол', homeScore: 7, awayScore: 3, date: '2025-01-28', time: '15:00', status: 'finished' },
  { id: 'tm-2', homeTeam: 'Торос', awayTeam: 'Кристалл-2', homeScore: 9, awayScore: 3, date: '2025-01-28', time: '17:00', status: 'finished' },
  { id: 'tm-3', homeTeam: 'Мордовия', awayTeam: 'Нефтяник', homeScore: 4, awayScore: 2, date: '2025-01-27', time: '14:00', status: 'finished' },
  { id: 'tm-4', homeTeam: 'Саров', awayTeam: 'Спутник', homeScore: 6, awayScore: 1, date: '2025-01-27', time: '16:00', status: 'finished' },
  { id: 'tm-5', homeTeam: 'Сокол', awayTeam: 'Торос', homeScore: null, awayScore: null, date: '2025-02-03', time: '15:00', status: 'scheduled' },
  { id: 'tm-6', homeTeam: 'Кристалл-2', awayTeam: 'АК Барс', homeScore: null, awayScore: null, date: '2025-02-03', time: '17:00', status: 'scheduled' },
  { id: 'tm-7', homeTeam: 'Нефтяник', awayTeam: 'Саров', homeScore: null, awayScore: null, date: '2025-02-04', time: '14:00', status: 'scheduled' },
  { id: 'tm-8', homeTeam: 'Спутник', awayTeam: 'Мордовия', homeScore: null, awayScore: null, date: '2025-02-04', time: '16:00', status: 'scheduled' },
]

// Бомбардиры турнира
export interface MockTournamentScorer {
  position: number
  name: string
  team: string
  goals: number
  assists: number
  points: number
}

export const MOCK_TOURNAMENT_SCORERS: MockTournamentScorer[] = [
  { position: 1, name: 'Петров Дмитрий', team: 'АК Барс', goals: 28, assists: 19, points: 47 },
  { position: 2, name: 'Морозов Иван', team: 'Торос', goals: 20, assists: 22, points: 42 },
  { position: 3, name: 'Беляев Роман', team: 'Мордовия', goals: 19, assists: 13, points: 32 },
  { position: 4, name: 'Фёдоров Даниил', team: 'Саров', goals: 21, assists: 14, points: 35 },
  { position: 5, name: 'Зайцев Матвей', team: 'АК Барс', goals: 15, assists: 16, points: 31 },
]
