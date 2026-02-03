import type { UserProfile, SubscriptionTier } from '@/entities/user'
import type {
  PlayerProfileExtended,
  Achievement,
  TeamMatch,
  ScoutNotification,
  VideoHighlight,
  AIRecommendation,
} from '@/entities/player'

// ===== Mock Users =====

export const createMockUser = (
  tier: SubscriptionTier = 'free',
  linkedPlayerId?: string
): UserProfile => ({
  id: 'user-1',
  email: 'player@example.com',
  role: 'player',
  linkedPlayerId,
  subscription: {
    tier,
    startDate: '2024-01-15',
    endDate: tier === 'free' ? null : '2025-01-15',
    autoRenew: tier !== 'free',
    price: tier === 'free' ? 0 : tier === 'pro' ? 990 : 2490,
  },
  createdAt: '2024-01-15',
})

export const MOCK_USER_FREE = createMockUser('free', 'player-1')
export const MOCK_USER_PRO = createMockUser('pro', 'player-1')
export const MOCK_USER_ULTRA = createMockUser('ultra', 'player-1')

// ===== Mock Player Profile =====

export const MOCK_PLAYER_PROFILE: PlayerProfileExtended = {
  id: 'player-1',
  name: 'Иванов Александр Сергеевич',
  birthDate: '2008-03-15',
  position: 'forward',
  team: 'СКА-Юниор',
  avatarUrl: undefined,
  jerseyNumber: 17,

  // Базовые данные
  region: 'Санкт-Петербург',
  city: 'Санкт-Петербург',
  height: 175,
  weight: 68,
  handedness: 'left',

  // Текущий сезон (FREE)
  currentSeasonStats: {
    goals: 23,
    assists: 31,
    games: 42,
    points: 54,
    plusMinus: 18,
    penaltyMinutes: 24,
    powerplayGoals: 7,
    shorthandedGoals: 1,
    evenStrengthGoals: 15,
  },
  regionalRank: 12,
  totalPlayersInRegion: 847,

  // История (PRO)
  allSeasonsStats: [
    { seasonId: '2024-25', seasonName: '2024/25', goals: 23, assists: 31, games: 42, points: 54 },
    { seasonId: '2023-24', seasonName: '2023/24', goals: 18, assists: 24, games: 38, points: 42 },
    { seasonId: '2022-23', seasonName: '2022/23', goals: 12, assists: 18, games: 35, points: 30 },
    { seasonId: '2021-22', seasonName: '2021/22', goals: 8, assists: 11, games: 32, points: 19 },
  ],
  performanceHistory: [
    { month: 'Сен', goals: 3, points: 7 },
    { month: 'Окт', goals: 5, points: 11 },
    { month: 'Ноя', goals: 4, points: 9 },
    { month: 'Дек', goals: 6, points: 14 },
    { month: 'Янв', goals: 3, points: 8 },
    { month: 'Фев', goals: 2, points: 5 },
  ],

  // Скауты (PRO)
  scoutViews: 34,
  lastScoutView: '2024-12-28',

  // ULTRA
  personalUrl: 'ivanov-alexander',
  isVerified: false,
}

// ===== Mock Team Matches =====

export const MOCK_TEAM_MATCHES: TeamMatch[] = [
  // Прошедшие матчи
  {
    id: 'match-1',
    date: '2024-12-20',
    time: '18:00',
    opponent: 'Динамо-Юниор',
    location: 'Ледовый дворец СКА',
    isHome: true,
    tournament: 'ЮХЛ',
    result: { homeScore: 4, awayScore: 2, isWin: true },
    playerStats: { goals: 1, assists: 2, plusMinus: 2 },
  },
  {
    id: 'match-2',
    date: '2024-12-23',
    time: '15:00',
    opponent: 'ЦСКА-Юниор',
    location: 'ЛДС ЦСКА',
    isHome: false,
    tournament: 'ЮХЛ',
    result: { homeScore: 3, awayScore: 3, isWin: false },
    playerStats: { goals: 0, assists: 1, plusMinus: 0 },
  },
  {
    id: 'match-3',
    date: '2024-12-27',
    time: '19:00',
    opponent: 'Спартак-Юниор',
    location: 'Ледовый дворец СКА',
    isHome: true,
    tournament: 'ЮХЛ',
    result: { homeScore: 5, awayScore: 1, isWin: true },
    playerStats: { goals: 2, assists: 1, plusMinus: 3 },
  },
  // Предстоящие матчи
  {
    id: 'match-4',
    date: '2025-01-03',
    time: '17:00',
    opponent: 'Локомотив-Юниор',
    location: 'Арена Локомотив',
    isHome: false,
    tournament: 'ЮХЛ',
  },
  {
    id: 'match-5',
    date: '2025-01-07',
    time: '18:30',
    opponent: 'Авангард-Юниор',
    location: 'Ледовый дворец СКА',
    isHome: true,
    tournament: 'ЮХЛ',
  },
  {
    id: 'match-6',
    date: '2025-01-12',
    time: '16:00',
    opponent: 'Металлург-Юниор',
    location: 'Арена Металлург',
    isHome: false,
    tournament: 'ЮХЛ',
  },
  {
    id: 'match-7',
    date: '2025-01-15',
    time: '19:00',
    opponent: 'Ак Барс-Юниор',
    location: 'Ледовый дворец СКА',
    isHome: true,
    tournament: 'ЮХЛ',
  },
  {
    id: 'match-8',
    date: '2025-01-20',
    time: '18:00',
    opponent: 'Салават Юлаев-Юниор',
    location: 'Уфа-Арена',
    isHome: false,
    tournament: 'ЮХЛ',
  },
]

// ===== Mock Achievements =====

export const MOCK_ACHIEVEMENTS: Achievement[] = [
  // Разблокированные (FREE)
  {
    id: 'ach-1',
    title: 'Первый гол',
    description: 'Забей свой первый гол в официальном матче',
    icon: 'Target',
    category: 'milestone',
    unlockedAt: '2021-09-15',
    isLocked: false,
  },
  {
    id: 'ach-2',
    title: 'Снайпер',
    description: 'Забей 10 голов за сезон',
    icon: 'Crosshair',
    category: 'stats',
    unlockedAt: '2022-12-10',
    isLocked: false,
  },
  {
    id: 'ach-3',
    title: 'Плеймейкер',
    description: 'Сделай 20 результативных передач за сезон',
    icon: 'Share2',
    category: 'stats',
    unlockedAt: '2023-02-20',
    isLocked: false,
  },
  {
    id: 'ach-4',
    title: '100 матчей',
    description: 'Сыграй 100 официальных матчей',
    icon: 'Medal',
    category: 'milestone',
    unlockedAt: '2024-01-15',
    isLocked: false,
  },
  {
    id: 'ach-5',
    title: 'Командный игрок',
    description: 'Заверши сезон с положительным показателем +/-',
    icon: 'Users',
    category: 'stats',
    unlockedAt: '2024-03-30',
    isLocked: false,
  },
  // С прогрессом
  {
    id: 'ach-6',
    title: 'Бомбардир',
    description: 'Набери 50 очков за сезон',
    icon: 'Flame',
    category: 'stats',
    progress: 92, // 46 из 50
    requirement: '50 очков',
    isLocked: true,
  },
  {
    id: 'ach-7',
    title: 'Суперснайпер',
    description: 'Забей 25 голов за сезон',
    icon: 'Zap',
    category: 'stats',
    progress: 92, // 23 из 25
    requirement: '25 голов',
    isLocked: true,
  },
  // Заблокированные (без прогресса)
  {
    id: 'ach-8',
    title: 'MVP турнира',
    description: 'Получи звание MVP на турнире',
    icon: 'Trophy',
    category: 'tournament',
    requirement: 'Стать MVP турнира',
    isLocked: true,
  },
  {
    id: 'ach-9',
    title: 'Чемпион региона',
    description: 'Выиграй региональный чемпионат',
    icon: 'Crown',
    category: 'tournament',
    requirement: 'Победа в региональном чемпионате',
    isLocked: true,
  },
  {
    id: 'ach-10',
    title: 'Хет-трик',
    description: 'Забей 3 гола в одном матче',
    icon: 'Star',
    category: 'stats',
    requirement: '3 гола в матче',
    isLocked: true,
  },
  // PRO достижения
  {
    id: 'ach-11',
    title: 'На радаре',
    description: 'Твой профиль просмотрели 10 скаутов',
    icon: 'Eye',
    category: 'special',
    progress: 100,
    unlockedAt: '2024-11-20',
    isLocked: false,
  },
  {
    id: 'ach-12',
    title: 'PRO подписчик',
    description: 'Активируй PRO подписку',
    icon: 'Sparkles',
    category: 'special',
    requirement: 'PRO подписка',
    isLocked: true,
  },
  // ULTRA достижения
  {
    id: 'ach-13',
    title: 'Верифицированный',
    description: 'Получи верификацию профиля',
    icon: 'BadgeCheck',
    category: 'special',
    requirement: 'ULTRA подписка',
    isLocked: true,
  },
  {
    id: 'ach-14',
    title: 'Видеозвезда',
    description: 'Загрузи 5 видео хайлайтов',
    icon: 'Video',
    category: 'special',
    requirement: '5 видео',
    isLocked: true,
  },
]

// ===== Mock Scout Notifications (PRO) =====

export const MOCK_SCOUT_NOTIFICATIONS: ScoutNotification[] = [
  {
    id: 'notif-1',
    scoutName: 'Петров А.В.',
    scoutClub: 'СКА',
    viewedAt: '2024-12-28T14:30:00',
    isRead: false,
  },
  {
    id: 'notif-2',
    scoutName: undefined,
    scoutClub: 'ЦСКА',
    viewedAt: '2024-12-27T10:15:00',
    isRead: false,
  },
  {
    id: 'notif-3',
    scoutName: 'Сидоров К.М.',
    scoutClub: 'Динамо',
    viewedAt: '2024-12-25T16:45:00',
    isRead: true,
    message: 'Интересный игрок, буду следить за развитием',
  },
  {
    id: 'notif-4',
    scoutName: undefined,
    scoutClub: undefined,
    viewedAt: '2024-12-23T09:20:00',
    isRead: true,
  },
  {
    id: 'notif-5',
    scoutName: 'Козлов Д.Н.',
    scoutClub: 'Локомотив',
    viewedAt: '2024-12-20T11:00:00',
    isRead: true,
  },
  {
    id: 'notif-6',
    scoutName: undefined,
    scoutClub: 'Ак Барс',
    viewedAt: '2024-12-18T15:30:00',
    isRead: true,
  },
  {
    id: 'notif-7',
    scoutName: 'Морозов И.П.',
    scoutClub: 'Металлург',
    viewedAt: '2024-12-15T13:45:00',
    isRead: true,
    message: 'Хорошая техника катания',
  },
]

// ===== Mock Video Highlights (ULTRA) =====

export const MOCK_VIDEO_HIGHLIGHTS: VideoHighlight[] = [
  {
    id: 'video-1',
    title: 'Гол в верхний угол против Динамо',
    description: 'Красивый бросок с кистей в верхний угол',
    url: 'https://example.com/video1.mp4',
    thumbnailUrl: '/thumbnails/goal1.jpg',
    duration: 45,
    uploadedAt: '2024-12-21',
    views: 234,
  },
  {
    id: 'video-2',
    title: 'Хет-трик против Спартака',
    description: 'Все три гола в одном видео',
    url: 'https://example.com/video2.mp4',
    thumbnailUrl: '/thumbnails/goal2.jpg',
    duration: 120,
    uploadedAt: '2024-12-15',
    views: 567,
  },
  {
    id: 'video-3',
    title: 'Голевая передача на пустые ворота',
    description: 'Точный пас через всю зону',
    url: 'https://example.com/video3.mp4',
    thumbnailUrl: '/thumbnails/assist1.jpg',
    duration: 30,
    uploadedAt: '2024-12-10',
    views: 189,
  },
]

// ===== Mock AI Recommendations (ULTRA) =====

export const MOCK_AI_RECOMMENDATIONS: AIRecommendation[] = [
  {
    id: 'rec-1',
    type: 'strength',
    title: 'Отличная результативность',
    description:
      'Ваш показатель голов за игру (0.55) выше среднего по позиции. Продолжайте в том же духе!',
    priority: 'low',
    basedOn: 'Статистика текущего сезона',
    createdAt: '2024-12-28',
  },
  {
    id: 'rec-2',
    type: 'improvement',
    title: 'Работа над вбрасываниями',
    description:
      'Процент выигранных вбрасываний (48%) можно улучшить. Рекомендуем отработать технику на тренировках.',
    priority: 'medium',
    basedOn: 'Анализ игровых действий',
    createdAt: '2024-12-27',
  },
  {
    id: 'rec-3',
    type: 'training',
    title: 'Развитие скорости',
    description:
      'Для вашей позиции важна скорость перехода в атаку. Добавьте интервальные тренировки 2-3 раза в неделю.',
    priority: 'high',
    basedOn: 'Анализ позиционных требований',
    createdAt: '2024-12-25',
  },
  {
    id: 'rec-4',
    type: 'general',
    title: 'Подготовка к плей-офф',
    description:
      'До плей-офф осталось 2 месяца. Самое время сфокусироваться на командных взаимодействиях.',
    priority: 'medium',
    basedOn: 'Календарь соревнований',
    createdAt: '2024-12-20',
  },
]

// ===== Helper Functions =====

// Получить предстоящие матчи
export const getUpcomingMatches = (matches: TeamMatch[]): TeamMatch[] =>
  matches
    .filter((m) => !m.result)
    .sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime())

// Получить прошедшие матчи
export const getPastMatches = (matches: TeamMatch[]): TeamMatch[] =>
  matches
    .filter((m) => m.result)
    .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())

// Получить разблокированные достижения
export const getUnlockedAchievements = (achievements: Achievement[]): Achievement[] =>
  achievements.filter((a) => !a.isLocked)

// Получить достижения с прогрессом
export const getInProgressAchievements = (achievements: Achievement[]): Achievement[] =>
  achievements.filter((a) => a.isLocked && a.progress !== undefined && a.progress > 0)

// Получить непрочитанные уведомления
export const getUnreadNotifications = (notifications: ScoutNotification[]): ScoutNotification[] =>
  notifications.filter((n) => !n.isRead)
