import { memo } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Lock, Sparkles, Crown, ArrowRight } from 'lucide-react'
import type { FeatureKey, SubscriptionTier } from '@/entities/user'
import { FEATURE_TIERS } from '@/entities/user'
import { cn } from '@/shared/lib/utils'

interface UpgradePromptProps {
  feature: FeatureKey
  className?: string
  compact?: boolean
}

// Названия фич для отображения
const FEATURE_NAMES: Record<FeatureKey, string> = {
  basic_profile: 'Базовый профиль',
  current_season_stats: 'Статистика сезона',
  regional_ranking: 'Региональный рейтинг',
  team_calendar: 'Календарь матчей',
  basic_achievements: 'Базовые достижения',
  all_seasons_history: 'История всех сезонов',
  progress_charts: 'Графики прогресса',
  player_comparison: 'Сравнение игроков',
  scout_visibility: 'Видимость для скаутов',
  scout_notifications: 'Уведомления о просмотрах',
  profile_photo: 'Фото в профиле',
  all_achievements: 'Все достижения',
  search_priority_max: 'Приоритет в поиске',
  personal_url: 'Персональный URL',
  video_highlights: 'Видео хайлайты',
  ai_recommendations: 'AI рекомендации',
  pdf_export: 'Экспорт в PDF',
  scout_messages: 'Сообщения от скаутов',
  verified_badge: 'Верификация профиля',
}

const TIER_CONFIG: Record<SubscriptionTier, { name: string; color: string; icon: typeof Lock }> = {
  free: { name: 'Бесплатный', color: 'gray', icon: Lock },
  pro: { name: 'PRO', color: 'purple', icon: Sparkles },
  ultra: { name: 'ULTRA', color: 'amber', icon: Crown },
}

export const UpgradePrompt = memo(function UpgradePrompt({
  feature,
  className,
  compact = false,
}: UpgradePromptProps) {
  const requiredTier = FEATURE_TIERS[feature]
  const tierConfig = TIER_CONFIG[requiredTier]
  const featureName = FEATURE_NAMES[feature]

  if (compact) {
    return (
      <Link
        to="/player/subscription"
        className={cn(
          'flex items-center gap-2 rounded-lg px-3 py-2',
          'bg-gradient-to-r from-[#8b5cf6]/20 to-[#ec4899]/20',
          'border border-[#8b5cf6]/30',
          'text-sm text-white',
          'transition-all hover:border-[#8b5cf6]/50',
          className
        )}
      >
        <tierConfig.icon size={16} className="text-[#8b5cf6]" />
        <span>Нужен {tierConfig.name}</span>
        <ArrowRight size={14} />
      </Link>
    )
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      className={cn(
        'relative overflow-hidden rounded-xl p-6',
        'bg-gradient-to-br from-[#8b5cf6]/10 to-[#ec4899]/10',
        'border border-[#8b5cf6]/20',
        className
      )}
    >
      {/* Background decoration */}
      <div className="absolute -right-10 -top-10 h-40 w-40 rounded-full bg-[#8b5cf6]/10 blur-3xl" />
      <div className="absolute -bottom-10 -left-10 h-40 w-40 rounded-full bg-[#ec4899]/10 blur-3xl" />

      <div className="relative flex flex-col items-center text-center">
        {/* Icon */}
        <div
          className={cn(
            'mb-4 flex h-16 w-16 items-center justify-center rounded-2xl',
            requiredTier === 'pro'
              ? 'bg-[#8b5cf6]/20 text-[#8b5cf6]'
              : 'bg-[#f59e0b]/20 text-[#f59e0b]'
          )}
        >
          <tierConfig.icon size={32} />
        </div>

        {/* Title */}
        <h3 className="mb-2 text-lg font-semibold text-white">{featureName}</h3>

        {/* Description */}
        <p className="mb-4 text-sm text-gray-400">
          Эта функция доступна на тарифе{' '}
          <span
            className={cn(
              'font-semibold',
              requiredTier === 'pro' ? 'text-[#8b5cf6]' : 'text-[#f59e0b]'
            )}
          >
            {tierConfig.name}
          </span>
        </p>

        {/* CTA Button */}
        <Link
          to="/player/subscription"
          className={cn(
            'flex items-center gap-2 rounded-lg px-6 py-3',
            'font-medium text-white',
            'transition-all hover:scale-105',
            requiredTier === 'pro'
              ? 'bg-gradient-to-r from-[#8b5cf6] to-[#7c3aed] hover:shadow-[0_0_30px_rgba(139,92,246,0.4)]'
              : 'bg-gradient-to-r from-[#f59e0b] to-[#d97706] hover:shadow-[0_0_30px_rgba(245,158,11,0.4)]'
          )}
        >
          <Sparkles size={18} />
          Улучшить подписку
          <ArrowRight size={18} />
        </Link>
      </div>
    </motion.div>
  )
})
