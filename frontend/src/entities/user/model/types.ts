// Уровни подписки
export type SubscriptionTier = 'free' | 'pro' | 'ultra'

// Роли пользователей
export type UserRole = 'fan' | 'player' | 'parent' | 'scout' | 'coach'

// Подписка пользователя
export interface Subscription {
  tier: SubscriptionTier
  startDate: string
  endDate: string | null // null = бессрочная (для free)
  autoRenew: boolean
  price: number
}

// Профиль пользователя
export interface UserProfile {
  id: string
  email: string
  role: UserRole
  linkedPlayerId?: string // ID привязанного игрока из БД
  subscription: Subscription
  avatarUrl?: string
  createdAt: string
}

// Состояние авторизации
export interface AuthState {
  isAuthenticated: boolean
  user: UserProfile | null
  token: string | null
}

// Ключи фич для проверки доступа
export type FeatureKey =
  // FREE
  | 'basic_profile'
  | 'current_season_stats'
  | 'regional_ranking'
  | 'team_calendar'
  | 'basic_achievements'
  // PRO
  | 'all_seasons_history'
  | 'progress_charts'
  | 'player_comparison'
  | 'scout_visibility'
  | 'scout_notifications'
  | 'profile_photo'
  | 'all_achievements'
  // ULTRA
  | 'search_priority_max'
  | 'personal_url'
  | 'video_highlights'
  | 'ai_recommendations'
  | 'pdf_export'
  | 'scout_messages'
  | 'verified_badge'

// Маппинг фич к минимальному уровню подписки
export const FEATURE_TIERS: Record<FeatureKey, SubscriptionTier> = {
  // FREE
  basic_profile: 'free',
  current_season_stats: 'free',
  regional_ranking: 'free',
  team_calendar: 'free',
  basic_achievements: 'free',
  // PRO
  all_seasons_history: 'pro',
  progress_charts: 'pro',
  player_comparison: 'pro',
  scout_visibility: 'pro',
  scout_notifications: 'pro',
  profile_photo: 'pro',
  all_achievements: 'pro',
  // ULTRA
  search_priority_max: 'ultra',
  personal_url: 'ultra',
  video_highlights: 'ultra',
  ai_recommendations: 'ultra',
  pdf_export: 'ultra',
  scout_messages: 'ultra',
  verified_badge: 'ultra',
}

// Порядок уровней для сравнения
export const TIER_ORDER: Record<SubscriptionTier, number> = {
  free: 0,
  pro: 1,
  ultra: 2,
}

// Проверка доступа к фиче
export function hasFeatureAccess(
  userTier: SubscriptionTier,
  feature: FeatureKey
): boolean {
  const requiredTier = FEATURE_TIERS[feature]
  return TIER_ORDER[userTier] >= TIER_ORDER[requiredTier]
}

// Информация о планах подписок
export interface SubscriptionPlan {
  tier: SubscriptionTier
  name: string
  price: number
  priceMonthly: number
  features: string[]
  highlighted?: boolean
}

export const SUBSCRIPTION_PLANS: SubscriptionPlan[] = [
  {
    tier: 'free',
    name: 'Бесплатный',
    price: 0,
    priceMonthly: 0,
    features: [
      'Базовый профиль игрока',
      'Статистика текущего сезона',
      'Рейтинг по региону',
      'Календарь матчей',
      'Базовые достижения',
    ],
  },
  {
    tier: 'pro',
    name: 'PRO',
    price: 990,
    priceMonthly: 990,
    highlighted: true,
    features: [
      'Всё из бесплатного',
      'История за все сезоны',
      'Графики прогресса',
      'Сравнение с игроками',
      'Видимость для скаутов',
      'Уведомления о просмотрах',
      'Фото в профиле',
    ],
  },
  {
    tier: 'ultra',
    name: 'ULTRA',
    price: 2490,
    priceMonthly: 2490,
    features: [
      'Всё из PRO',
      'Приоритет в поиске скаутов',
      'Персональный URL',
      'Видео хайлайты',
      'AI рекомендации',
      'Экспорт в PDF',
      'Сообщения от скаутов',
      'Верифицированный badge',
    ],
  },
]
