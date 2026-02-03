import { memo, useState } from 'react'
import { motion } from 'framer-motion'
import {
  Calendar,
  MapPin,
  Clock,
  Target,
  Users2,
  TrendingUp,
} from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import { usePlayerDashboardStore } from '@/shared/stores'
import { getUpcomingMatches, getPastMatches } from '@/shared/mocks'

type TabType = 'upcoming' | 'past'

export const TeamCalendarPage = memo(function TeamCalendarPage() {
  const [activeTab, setActiveTab] = useState<TabType>('upcoming')
  const teamMatches = usePlayerDashboardStore((state) => state.teamMatches)

  const upcomingMatches = getUpcomingMatches(teamMatches)
  const pastMatches = getPastMatches(teamMatches)

  const displayMatches = activeTab === 'upcoming' ? upcomingMatches : pastMatches

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
            <Calendar className="text-[#00d4ff]" />
            Календарь матчей
          </h1>
          <p className="text-gray-400">Расписание игр вашей команды</p>
        </div>
      </motion.div>

      {/* Tabs */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="flex gap-2"
      >
        <button
          onClick={() => setActiveTab('upcoming')}
          className={cn(
            'px-4 py-2 rounded-lg text-sm font-medium transition-all',
            activeTab === 'upcoming'
              ? 'bg-[#00d4ff] text-white'
              : 'bg-white/5 text-gray-400 hover:bg-white/10 hover:text-white'
          )}
        >
          Предстоящие ({upcomingMatches.length})
        </button>
        <button
          onClick={() => setActiveTab('past')}
          className={cn(
            'px-4 py-2 rounded-lg text-sm font-medium transition-all',
            activeTab === 'past'
              ? 'bg-[#00d4ff] text-white'
              : 'bg-white/5 text-gray-400 hover:bg-white/10 hover:text-white'
          )}
        >
          Прошедшие ({pastMatches.length})
        </button>
      </motion.div>

      {/* Matches List */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.2 }}
        className="space-y-4"
      >
        {displayMatches.length > 0 ? (
          displayMatches.map((match, index) => (
            <motion.div
              key={match.id}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.1 + index * 0.05 }}
            >
              <GlassCard className="p-0 overflow-hidden" glowColor="blue">
                <div className="flex flex-col lg:flex-row">
                  {/* Date Section */}
                  <div
                    className={cn(
                      'flex-shrink-0 p-6 flex flex-col items-center justify-center',
                      'bg-gradient-to-br from-[#00d4ff]/20 to-[#8b5cf6]/20',
                      'lg:w-32'
                    )}
                  >
                    <div className="text-3xl font-bold text-white">
                      {new Date(match.date).getDate()}
                    </div>
                    <div className="text-sm text-gray-300 uppercase">
                      {new Date(match.date).toLocaleDateString('ru-RU', { month: 'short' })}
                    </div>
                    {match.time && (
                      <div className="mt-2 flex items-center gap-1 text-xs text-gray-400">
                        <Clock size={12} />
                        {match.time}
                      </div>
                    )}
                  </div>

                  {/* Match Info */}
                  <div className="flex-1 p-6">
                    <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
                      {/* Teams */}
                      <div className="flex-1">
                        <div className="flex items-center gap-3">
                          <span
                            className={cn(
                              'text-xs px-2 py-1 rounded font-medium',
                              match.isHome
                                ? 'bg-green-500/20 text-green-400'
                                : 'bg-blue-500/20 text-blue-400'
                            )}
                          >
                            {match.isHome ? 'ДОМА' : 'В ГОСТЯХ'}
                          </span>
                          {match.tournament && (
                            <span className="text-xs text-gray-500">{match.tournament}</span>
                          )}
                        </div>

                        <h3 className="mt-2 text-lg font-semibold text-white">
                          vs {match.opponent}
                        </h3>

                        <div className="mt-2 flex items-center gap-1 text-sm text-gray-400">
                          <MapPin size={14} />
                          {match.location}
                        </div>
                      </div>

                      {/* Result (for past matches) */}
                      {match.result && (
                        <div className="flex flex-col items-center gap-2">
                          <div
                            className={cn(
                              'text-3xl font-bold',
                              match.result.isWin ? 'text-green-400' : 'text-gray-400'
                            )}
                          >
                            {match.isHome
                              ? `${match.result.homeScore}:${match.result.awayScore}`
                              : `${match.result.awayScore}:${match.result.homeScore}`}
                          </div>
                          <span
                            className={cn(
                              'text-xs px-3 py-1 rounded-full font-medium',
                              match.result.isWin
                                ? 'bg-green-500/20 text-green-400'
                                : match.result.homeScore === match.result.awayScore
                                  ? 'bg-gray-500/20 text-gray-400'
                                  : 'bg-red-500/20 text-red-400'
                            )}
                          >
                            {match.result.isWin
                              ? 'ПОБЕДА'
                              : match.result.homeScore === match.result.awayScore
                                ? 'НИЧЬЯ'
                                : 'ПОРАЖЕНИЕ'}
                          </span>
                        </div>
                      )}
                    </div>

                    {/* Player Stats (for past matches) */}
                    {match.playerStats && (
                      <div className="mt-4 pt-4 border-t border-white/10">
                        <div className="text-xs text-gray-500 mb-2">Ваша статистика:</div>
                        <div className="flex gap-6">
                          <div className="flex items-center gap-2">
                            <div className="h-8 w-8 rounded-lg bg-[#00d4ff]/20 flex items-center justify-center">
                              <Target size={14} className="text-[#00d4ff]" />
                            </div>
                            <div>
                              <div className="text-lg font-bold text-white">
                                {match.playerStats.goals}
                              </div>
                              <div className="text-xs text-gray-500">Голы</div>
                            </div>
                          </div>
                          <div className="flex items-center gap-2">
                            <div className="h-8 w-8 rounded-lg bg-[#8b5cf6]/20 flex items-center justify-center">
                              <Users2 size={14} className="text-[#8b5cf6]" />
                            </div>
                            <div>
                              <div className="text-lg font-bold text-white">
                                {match.playerStats.assists}
                              </div>
                              <div className="text-xs text-gray-500">Передачи</div>
                            </div>
                          </div>
                          <div className="flex items-center gap-2">
                            <div className="h-8 w-8 rounded-lg bg-[#10b981]/20 flex items-center justify-center">
                              <TrendingUp size={14} className="text-[#10b981]" />
                            </div>
                            <div>
                              <div
                                className={cn(
                                  'text-lg font-bold',
                                  match.playerStats.plusMinus > 0
                                    ? 'text-green-400'
                                    : match.playerStats.plusMinus < 0
                                      ? 'text-red-400'
                                      : 'text-white'
                                )}
                              >
                                {match.playerStats.plusMinus > 0 ? '+' : ''}
                                {match.playerStats.plusMinus}
                              </div>
                              <div className="text-xs text-gray-500">+/-</div>
                            </div>
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </GlassCard>
            </motion.div>
          ))
        ) : (
          <GlassCard className="p-12 text-center">
            <Calendar size={48} className="mx-auto text-gray-600 mb-4" />
            <h3 className="text-lg font-semibold text-white mb-2">
              {activeTab === 'upcoming' ? 'Нет предстоящих матчей' : 'Нет прошедших матчей'}
            </h3>
            <p className="text-gray-400">
              {activeTab === 'upcoming'
                ? 'Расписание будет обновлено, когда появятся новые игры'
                : 'История матчей появится после первой игры'}
            </p>
          </GlassCard>
        )}
      </motion.div>
    </div>
  )
})
