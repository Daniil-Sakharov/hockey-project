export interface Player {
  id: string
  name: string
  birthDate?: string
  position?: string
  height?: number
  weight?: number
  team?: string
  region?: string
}

export interface PlayerStats {
  goals: number
  assists: number
  games: number
  points?: number
  plusMinus?: number
}

export interface TopScorer {
  id: string
  name: string
  team: string
  goals: number
  assists: number
  games: number
}

// Dashboard types

export type PlayerPosition = 'forward' | 'defender' | 'goalie'

export interface PlayerProfile {
  id: string
  name: string
  birthDate?: string
  position?: PlayerPosition
  team?: string
  avatarUrl?: string
  jerseyNumber?: number
}

export interface PlayerDetailedStats {
  goals: number
  assists: number
  games: number
  points: number
  plusMinus: number
  penaltyMinutes: number
  powerplayGoals: number
  shorthandedGoals: number
  evenStrengthGoals: number
}

export interface PlayerSeasonStats {
  seasonId: string
  seasonName: string
  goals: number
  assists: number
  games: number
  points: number
}

export interface PlayerPerformancePoint {
  month: string
  goals: number
  points: number
}

export interface PlayerRankingEntry extends TopScorer {
  rank: number
  isCurrentPlayer?: boolean
}

// ===== Player Dashboard Types =====

// Расширенный профиль игрока для Player Dashboard
export interface PlayerProfileExtended extends PlayerProfile {
  // Базовые данные
  region?: string
  city?: string
  height?: number
  weight?: number
  handedness?: 'left' | 'right'

  // Текущий сезон (FREE)
  currentSeasonStats?: PlayerDetailedStats
  regionalRank?: number
  totalPlayersInRegion?: number

  // История (PRO)
  allSeasonsStats?: PlayerSeasonStats[]
  performanceHistory?: PlayerPerformancePoint[]

  // Скауты (PRO)
  scoutViews?: number
  lastScoutView?: string

  // ULTRA
  personalUrl?: string
  isVerified?: boolean
}

// Достижения
export type AchievementCategory = 'stats' | 'tournament' | 'milestone' | 'special'

export interface Achievement {
  id: string
  title: string
  description: string
  icon: string // имя иконки из Lucide
  category: AchievementCategory
  unlockedAt?: string // если разблокировано
  progress?: number // 0-100, если есть прогресс
  requirement?: string // что нужно для получения
  isLocked: boolean
}

// Матчи команды (календарь)
export interface TeamMatch {
  id: string
  date: string
  time?: string
  opponent: string
  opponentLogo?: string
  location: string
  isHome: boolean
  tournament?: string
  result?: {
    homeScore: number
    awayScore: number
    isWin: boolean
  }
  playerStats?: {
    goals: number
    assists: number
    plusMinus: number
  }
}

// Видео хайлайты (ULTRA)
export interface VideoHighlight {
  id: string
  title: string
  description?: string
  url: string
  thumbnailUrl: string
  duration: number // в секундах
  uploadedAt: string
  views?: number
}

// AI рекомендации (ULTRA)
export type RecommendationType = 'strength' | 'improvement' | 'training' | 'general'

export interface AIRecommendation {
  id: string
  type: RecommendationType
  title: string
  description: string
  priority: 'high' | 'medium' | 'low'
  basedOn?: string // на основе чего дана рекомендация
  createdAt: string
}

// Уведомления о просмотрах скаутами (PRO)
export interface ScoutNotification {
  id: string
  scoutName?: string // может быть анонимным
  scoutClub?: string
  viewedAt: string
  isRead: boolean
  message?: string // если скаут оставил комментарий
}
