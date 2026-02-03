import { memo, useState } from 'react'
import { motion } from 'framer-motion'
import {
  Trophy,
  Target,
  Medal,
  Star,
  Award,
  Lock,
  CheckCircle,
  Sparkles,
  Crown,
} from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import { useAuthStore, usePlayerDashboardStore } from '@/shared/stores'
import { SubscriptionGate } from '@/features/subscription-gate'
import type { AchievementCategory, Achievement } from '@/entities/player'

type FilterType = 'all' | AchievementCategory

const CATEGORY_CONFIG: Record<
  AchievementCategory,
  { label: string; icon: typeof Trophy; color: string }
> = {
  stats: { label: 'Статистика', icon: Target, color: 'text-[#00d4ff]' },
  tournament: { label: 'Турниры', icon: Medal, color: 'text-[#f59e0b]' },
  milestone: { label: 'Достижения', icon: Star, color: 'text-[#8b5cf6]' },
  special: { label: 'Особые', icon: Crown, color: 'text-[#ec4899]' },
}

const ICON_MAP: Record<string, typeof Trophy> = {
  Target,
  Crosshair: Target,
  Share2: Medal,
  Medal,
  Users: Trophy,
  Flame: Sparkles,
  Zap: Star,
  Trophy,
  Crown,
  Star,
  Eye: Target,
  Sparkles,
  BadgeCheck: Award,
  Video: Medal,
}

interface AchievementCardProps {
  achievement: Achievement
  delay: number
}

const AchievementCard = memo(function AchievementCard({
  achievement,
  delay,
}: AchievementCardProps) {
  const config = CATEGORY_CONFIG[achievement.category]
  const Icon = ICON_MAP[achievement.icon] || Trophy

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ delay }}
      whileHover={{ scale: 1.02 }}
    >
      <GlassCard
        className={cn(
          'p-4 h-full',
          achievement.isLocked && 'opacity-60'
        )}
        glowColor={achievement.isLocked ? undefined : 'purple'}
      >
        <div className="flex flex-col h-full">
          {/* Icon */}
          <div className="flex items-start justify-between mb-3">
            <div
              className={cn(
                'h-12 w-12 rounded-xl flex items-center justify-center',
                achievement.isLocked ? 'bg-white/5' : 'bg-[#8b5cf6]/20'
              )}
            >
              {achievement.isLocked ? (
                <Lock size={24} className="text-gray-500" />
              ) : (
                <Icon size={24} className="text-[#8b5cf6]" />
              )}
            </div>

            {/* Unlocked badge */}
            {!achievement.isLocked && (
              <div className="flex items-center gap-1 text-green-400">
                <CheckCircle size={16} />
              </div>
            )}
          </div>

          {/* Title & Description */}
          <h3
            className={cn(
              'font-semibold mb-1',
              achievement.isLocked ? 'text-gray-400' : 'text-white'
            )}
          >
            {achievement.title}
          </h3>
          <p className="text-xs text-gray-500 mb-3 flex-1">{achievement.description}</p>

          {/* Progress bar */}
          {achievement.progress !== undefined && achievement.isLocked && (
            <div className="mt-auto">
              <div className="flex items-center justify-between text-xs mb-1">
                <span className="text-gray-500">Прогресс</span>
                <span className="text-[#8b5cf6]">{achievement.progress}%</span>
              </div>
              <div className="h-1.5 rounded-full bg-white/10 overflow-hidden">
                <motion.div
                  initial={{ width: 0 }}
                  animate={{ width: `${achievement.progress}%` }}
                  transition={{ delay: delay + 0.2, duration: 0.5 }}
                  className="h-full bg-gradient-to-r from-[#8b5cf6] to-[#ec4899]"
                />
              </div>
            </div>
          )}

          {/* Unlocked date */}
          {achievement.unlockedAt && (
            <div className="mt-auto pt-2 border-t border-white/5">
              <span className="text-xs text-gray-500">
                Получено:{' '}
                {new Date(achievement.unlockedAt).toLocaleDateString('ru-RU', {
                  day: 'numeric',
                  month: 'long',
                  year: 'numeric',
                })}
              </span>
            </div>
          )}

          {/* Requirement */}
          {achievement.isLocked && !achievement.progress && achievement.requirement && (
            <div className="mt-auto pt-2 border-t border-white/5">
              <span className="text-xs text-gray-500">Требуется: {achievement.requirement}</span>
            </div>
          )}

          {/* Category badge */}
          <div className="mt-3">
            <span
              className={cn(
                'inline-flex items-center gap-1 text-[10px] px-2 py-1 rounded-full',
                'bg-white/5',
                config.color
              )}
            >
              <config.icon size={10} />
              {config.label}
            </span>
          </div>
        </div>
      </GlassCard>
    </motion.div>
  )
})

export const AchievementsPage = memo(function AchievementsPage() {
  const [filter, setFilter] = useState<FilterType>('all')
  const achievements = usePlayerDashboardStore((state) => state.achievements)
  const hasFeature = useAuthStore((state) => state.hasFeature)

  const canSeeAllAchievements = hasFeature('all_achievements')

  // Filter achievements
  const filteredAchievements = achievements.filter((a) => {
    if (filter === 'all') return true
    return a.category === filter
  })

  // For free users, only show basic achievements (no special category)
  const basicAchievements = filteredAchievements.filter((a) => a.category !== 'special')
  const specialAchievements = filteredAchievements.filter((a) => a.category === 'special')

  const displayAchievements = canSeeAllAchievements ? filteredAchievements : basicAchievements

  // Stats
  const totalUnlocked = achievements.filter((a) => !a.isLocked).length
  const totalAchievements = achievements.length

  return (
    <div className="space-y-6">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
      >
        <div>
          <h1 className="text-2xl font-bold text-white flex items-center gap-3">
            <Trophy className="text-[#8b5cf6]" />
            Достижения
          </h1>
          <p className="text-gray-400">
            Разблокировано {totalUnlocked} из {totalAchievements}
          </p>
        </div>

        {/* Overall progress */}
        <div className="flex items-center gap-4">
          <div className="w-32">
            <div className="h-2 rounded-full bg-white/10 overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                animate={{ width: `${(totalUnlocked / totalAchievements) * 100}%` }}
                transition={{ delay: 0.3, duration: 0.8 }}
                className="h-full bg-gradient-to-r from-[#8b5cf6] to-[#ec4899]"
              />
            </div>
          </div>
          <span className="text-lg font-bold text-[#8b5cf6]">
            {Math.round((totalUnlocked / totalAchievements) * 100)}%
          </span>
        </div>
      </motion.div>

      {/* Filters */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="flex flex-wrap gap-2"
      >
        <button
          onClick={() => setFilter('all')}
          className={cn(
            'px-4 py-2 rounded-lg text-sm font-medium transition-all',
            filter === 'all'
              ? 'bg-[#8b5cf6] text-white'
              : 'bg-white/5 text-gray-400 hover:bg-white/10 hover:text-white'
          )}
        >
          Все
        </button>
        {Object.entries(CATEGORY_CONFIG).map(([key, config]) => (
          <button
            key={key}
            onClick={() => setFilter(key as AchievementCategory)}
            className={cn(
              'px-4 py-2 rounded-lg text-sm font-medium transition-all flex items-center gap-2',
              filter === key
                ? 'bg-[#8b5cf6] text-white'
                : 'bg-white/5 text-gray-400 hover:bg-white/10 hover:text-white'
            )}
          >
            <config.icon size={14} />
            {config.label}
          </button>
        ))}
      </motion.div>

      {/* Achievements Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {displayAchievements.map((achievement, index) => (
          <AchievementCard
            key={achievement.id}
            achievement={achievement}
            delay={0.1 + index * 0.03}
          />
        ))}
      </div>

      {/* Special Achievements (PRO) */}
      {!canSeeAllAchievements && specialAchievements.length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
        >
          <SubscriptionGate feature="all_achievements">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {specialAchievements.map((achievement, index) => (
                <AchievementCard
                  key={achievement.id}
                  achievement={achievement}
                  delay={0.1 + index * 0.03}
                />
              ))}
            </div>
          </SubscriptionGate>
        </motion.div>
      )}

      {/* Empty state */}
      {displayAchievements.length === 0 && (
        <GlassCard className="p-12 text-center">
          <Trophy size={48} className="mx-auto text-gray-600 mb-4" />
          <h3 className="text-lg font-semibold text-white mb-2">
            Нет достижений в этой категории
          </h3>
          <p className="text-gray-400">Продолжайте играть, чтобы получить новые достижения</p>
        </GlassCard>
      )}
    </div>
  )
})
