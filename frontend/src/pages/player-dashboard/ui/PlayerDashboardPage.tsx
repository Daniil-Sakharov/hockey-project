import { memo } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  Target,
  Users2,
  Zap,
  Award,
  TrendingUp,
  Calendar,
  MapPin,
  Trophy,
  ArrowRight,
  Eye,
} from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { KPICard } from '@/widgets/player-kpi'
import { cn } from '@/shared/lib/utils'
import { useAuthStore, usePlayerDashboardStore } from '@/shared/stores'
import { SubscriptionGate, FeatureLockedOverlay } from '@/features/subscription-gate'
import { getUpcomingMatches } from '@/shared/mocks'

export const PlayerDashboardPage = memo(function PlayerDashboardPage() {
  const { getSubscriptionTier, hasFeature } = useAuthStore()
  const linkedPlayer = usePlayerDashboardStore((state) => state.linkedPlayer)
  const teamMatches = usePlayerDashboardStore((state) => state.teamMatches)
  const achievements = usePlayerDashboardStore((state) => state.achievements)
  const tier = getSubscriptionTier()

  const upcomingMatches = getUpcomingMatches(teamMatches).slice(0, 3)
  const unlockedAchievements = achievements.filter((a) => !a.isLocked)
  const stats = linkedPlayer?.currentSeasonStats

  if (!linkedPlayer) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <div className="text-center">
          <h2 className="text-xl font-semibold text-white mb-2">Профиль не привязан</h2>
          <p className="text-gray-400 mb-4">
            Привяжите свой игровой профиль, чтобы увидеть статистику
          </p>
          <Link
            to="/link-player"
            className="inline-flex items-center gap-2 rounded-lg bg-[#00d4ff] px-6 py-3 font-medium text-white hover:bg-[#00d4ff]/90 transition-colors"
          >
            Привязать профиль
            <ArrowRight size={18} />
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
      >
        <div>
          <h1 className="text-2xl font-bold text-white">Мой профиль</h1>
          <p className="text-gray-400">Сезон 2024/25</p>
        </div>
        {tier !== 'free' && (
          <div
            className={cn(
              'inline-flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium',
              tier === 'ultra'
                ? 'bg-[#f59e0b]/20 text-[#f59e0b]'
                : 'bg-[#8b5cf6]/20 text-[#8b5cf6]'
            )}
          >
            <Award size={16} />
            {tier === 'ultra' ? 'ULTRA' : 'PRO'} подписка
          </div>
        )}
      </motion.div>

      {/* Player Profile Card */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        <GlassCard className="p-6" glowColor="cyan">
          <div className="flex flex-col gap-6 md:flex-row md:items-start">
            {/* Avatar */}
            <div className="flex-shrink-0">
              <div className="relative h-24 w-24 rounded-2xl bg-gradient-to-br from-[#00d4ff] to-[#8b5cf6] flex items-center justify-center text-white text-3xl font-bold">
                {linkedPlayer.name.charAt(0)}
                {linkedPlayer.isVerified && (
                  <div className="absolute -bottom-2 -right-2 h-8 w-8 rounded-full bg-[#f59e0b] flex items-center justify-center">
                    <Award size={16} className="text-white" />
                  </div>
                )}
              </div>
            </div>

            {/* Info */}
            <div className="flex-1 min-w-0">
              <div className="flex flex-wrap items-center gap-3">
                <h2 className="text-xl font-bold text-white">{linkedPlayer.name}</h2>
                {linkedPlayer.jerseyNumber && (
                  <span className="text-lg font-semibold text-[#00d4ff]">
                    #{linkedPlayer.jerseyNumber}
                  </span>
                )}
              </div>

              <div className="mt-2 flex flex-wrap gap-3 text-sm text-gray-400">
                <span className="flex items-center gap-1">
                  <Users2 size={14} />
                  {linkedPlayer.team}
                </span>
                <span className="flex items-center gap-1">
                  <MapPin size={14} />
                  {linkedPlayer.region}
                </span>
                <span>
                  {linkedPlayer.position === 'forward'
                    ? 'Нападающий'
                    : linkedPlayer.position === 'defender'
                      ? 'Защитник'
                      : 'Вратарь'}
                </span>
              </div>

              {/* Physical stats */}
              <div className="mt-3 flex flex-wrap gap-4 text-sm">
                {linkedPlayer.height && (
                  <span className="text-gray-300">
                    <span className="text-gray-500">Рост:</span> {linkedPlayer.height} см
                  </span>
                )}
                {linkedPlayer.weight && (
                  <span className="text-gray-300">
                    <span className="text-gray-500">Вес:</span> {linkedPlayer.weight} кг
                  </span>
                )}
                {linkedPlayer.handedness && (
                  <span className="text-gray-300">
                    <span className="text-gray-500">Хват:</span>{' '}
                    {linkedPlayer.handedness === 'left' ? 'Левый' : 'Правый'}
                  </span>
                )}
              </div>

              {/* Regional rank */}
              {linkedPlayer.regionalRank && (
                <div className="mt-4 inline-flex items-center gap-2 rounded-lg bg-[#00d4ff]/10 px-4 py-2">
                  <Trophy size={16} className="text-[#00d4ff]" />
                  <span className="text-sm text-gray-300">
                    <span className="font-semibold text-[#00d4ff]">
                      #{linkedPlayer.regionalRank}
                    </span>{' '}
                    место в регионе
                    <span className="text-gray-500">
                      {' '}
                      из {linkedPlayer.totalPlayersInRegion}
                    </span>
                  </span>
                </div>
              )}
            </div>

            {/* Scout views (PRO) */}
            <SubscriptionGate feature="scout_notifications" showUpgrade={false}>
              <div className="flex-shrink-0 rounded-xl bg-white/5 p-4 text-center">
                <div className="flex items-center justify-center gap-2 text-[#8b5cf6]">
                  <Eye size={20} />
                  <span className="text-2xl font-bold">{linkedPlayer.scoutViews || 0}</span>
                </div>
                <p className="mt-1 text-xs text-gray-400">Просмотров скаутами</p>
                {linkedPlayer.lastScoutView && (
                  <p className="mt-2 text-xs text-gray-500">
                    Последний: {new Date(linkedPlayer.lastScoutView).toLocaleDateString('ru-RU')}
                  </p>
                )}
              </div>
            </SubscriptionGate>
          </div>
        </GlassCard>
      </motion.div>

      {/* KPI Grid */}
      {stats && (
        <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
          <KPICard
            title="Голы"
            value={stats.goals}
            icon={<Target size={24} />}
            color="cyan"
            trend="up"
            trendValue="+5"
            delay={0.1}
          />
          <KPICard
            title="Передачи"
            value={stats.assists}
            icon={<Users2 size={24} />}
            color="purple"
            trend="up"
            trendValue="+8"
            delay={0.15}
          />
          <KPICard
            title="Очки"
            value={stats.points}
            icon={<Zap size={24} />}
            color="pink"
            trend="up"
            trendValue="+13"
            delay={0.2}
          />
          <KPICard
            title="+/-"
            value={stats.plusMinus}
            icon={<TrendingUp size={24} />}
            color="green"
            trend="up"
            trendValue="+3"
            delay={0.25}
          />
        </div>
      )}

      {/* Two Column Layout */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Upcoming Matches */}
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.3 }}
        >
          <GlassCard className="p-6" glowColor="blue">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-white flex items-center gap-2">
                <Calendar size={20} className="text-[#00d4ff]" />
                Ближайшие матчи
              </h3>
              <Link
                to="/player/calendar"
                className="text-sm text-[#00d4ff] hover:underline flex items-center gap-1"
              >
                Все матчи
                <ArrowRight size={14} />
              </Link>
            </div>

            <div className="space-y-3">
              {upcomingMatches.length > 0 ? (
                upcomingMatches.map((match) => (
                  <div
                    key={match.id}
                    className="flex items-center gap-4 rounded-lg bg-white/5 p-3"
                  >
                    <div className="flex-shrink-0 text-center">
                      <div className="text-sm font-semibold text-white">
                        {new Date(match.date).toLocaleDateString('ru-RU', {
                          day: 'numeric',
                          month: 'short',
                        })}
                      </div>
                      <div className="text-xs text-gray-500">{match.time}</div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <span
                          className={cn(
                            'text-xs px-2 py-0.5 rounded',
                            match.isHome
                              ? 'bg-green-500/20 text-green-400'
                              : 'bg-blue-500/20 text-blue-400'
                          )}
                        >
                          {match.isHome ? 'Дома' : 'В гостях'}
                        </span>
                        <span className="text-sm text-gray-300 truncate">
                          vs {match.opponent}
                        </span>
                      </div>
                      <div className="flex items-center gap-1 mt-1 text-xs text-gray-500">
                        <MapPin size={12} />
                        <span className="truncate">{match.location}</span>
                      </div>
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-center py-6 text-gray-500">
                  <Calendar size={32} className="mx-auto mb-2 opacity-50" />
                  <p>Нет предстоящих матчей</p>
                </div>
              )}
            </div>
          </GlassCard>
        </motion.div>

        {/* Achievements Preview */}
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.35 }}
        >
          <GlassCard className="p-6" glowColor="purple">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-white flex items-center gap-2">
                <Trophy size={20} className="text-[#8b5cf6]" />
                Достижения
              </h3>
              <Link
                to="/player/achievements"
                className="text-sm text-[#8b5cf6] hover:underline flex items-center gap-1"
              >
                Все достижения
                <ArrowRight size={14} />
              </Link>
            </div>

            <div className="space-y-3">
              <div className="flex items-center justify-between text-sm">
                <span className="text-gray-400">Разблокировано</span>
                <span className="font-semibold text-white">
                  {unlockedAchievements.length} из {achievements.length}
                </span>
              </div>

              {/* Progress bar */}
              <div className="h-2 rounded-full bg-white/10 overflow-hidden">
                <motion.div
                  initial={{ width: 0 }}
                  animate={{
                    width: `${(unlockedAchievements.length / achievements.length) * 100}%`,
                  }}
                  transition={{ delay: 0.5, duration: 0.8 }}
                  className="h-full bg-gradient-to-r from-[#8b5cf6] to-[#ec4899]"
                />
              </div>

              {/* Recent achievements */}
              <div className="mt-4 grid grid-cols-4 gap-2">
                {unlockedAchievements.slice(0, 4).map((ach) => (
                  <div
                    key={ach.id}
                    className="flex flex-col items-center gap-1 rounded-lg bg-white/5 p-2"
                    title={ach.title}
                  >
                    <div className="h-8 w-8 rounded-lg bg-[#8b5cf6]/20 flex items-center justify-center text-[#8b5cf6]">
                      <Award size={16} />
                    </div>
                    <span className="text-[10px] text-gray-400 text-center truncate w-full">
                      {ach.title}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </GlassCard>
        </motion.div>
      </div>

      {/* PRO Features Preview */}
      {!hasFeature('progress_charts') && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
        >
          <FeatureLockedOverlay
            feature="progress_charts"
            blurAmount="sm"
            previewContent={
              <GlassCard className="p-6">
                <h3 className="text-lg font-semibold text-white mb-4">
                  Графики прогресса
                </h3>
                <div className="h-48 flex items-center justify-center">
                  <div className="w-full h-full bg-gradient-to-r from-[#00d4ff]/10 via-[#8b5cf6]/10 to-[#ec4899]/10 rounded-lg" />
                </div>
              </GlassCard>
            }
          />
        </motion.div>
      )}
    </div>
  )
})
