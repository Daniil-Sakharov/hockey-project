import { memo } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { useCountUp } from '@/shared/hooks'
import {
  Trophy,
  Users2,
  UserRound,
  Calendar,
  ArrowRight,
  Star,
  TrendingUp,
  Loader2,
} from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { useAuthStore } from '@/shared/stores'
import { cleanTournamentName } from '@/shared/lib/formatters'
import { FeatureLockedOverlay } from '@/features/subscription-gate'
import {
  useExploreOverview,
  useRecentResults,
  useTournaments,
  useRankings,
} from '@/shared/api/useExploreQueries'

export const ExploreDashboardPage = memo(function ExploreDashboardPage() {
  const user = useAuthStore((state) => state.user)
  const { data: overview, isLoading: overviewLoading } = useExploreOverview()
  const { data: recentMatches, isLoading: matchesLoading } = useRecentResults(undefined, 5)
  const { data: tournaments, isLoading: tournamentsLoading } = useTournaments()
  const { data: rankingsData, isLoading: scorersLoading } = useRankings('points', 5)
  const topScorers = rankingsData?.players
  const season = rankingsData?.season

  const kpiItems = [
    { label: 'Турниров', value: overview?.tournaments ?? 0, icon: <Trophy size={24} />, color: 'cyan' },
    { label: 'Команд', value: overview?.teams ?? 0, icon: <Users2 size={24} />, color: 'purple' },
    { label: 'Игроков', value: overview?.players ?? 0, icon: <UserRound size={24} />, color: 'pink' },
    { label: 'Матчей', value: overview?.matches ?? 0, icon: <Calendar size={24} />, color: 'green' },
  ]

  const colorMap: Record<string, string> = {
    cyan: '#00d4ff',
    purple: '#8b5cf6',
    pink: '#ec4899',
    green: '#10b981',
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <h1 className="text-2xl font-bold text-gradient-animated">Обзор</h1>
        <p className="text-gray-400">
          Добро пожаловать{user ? `, ${user.email}` : ''}
        </p>
      </motion.div>

      {/* KPI Grid */}
      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        {kpiItems.map((item, i) => (
          <motion.div
            key={item.label}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 + i * 0.05 }}
          >
            <GlassCard className="p-4">
              <div className="flex items-center gap-3">
                <div
                  className="flex h-10 w-10 items-center justify-center rounded-xl"
                  style={{ backgroundColor: `${colorMap[item.color]}20`, color: colorMap[item.color] }}
                >
                  {item.icon}
                </div>
                <div>
                  {overviewLoading ? (
                    <Loader2 size={20} className="animate-spin text-gray-500" />
                  ) : (
                    <KpiValue value={item.value} />
                  )}
                  <p className="text-xs text-gray-500">{item.label}</p>
                </div>
              </div>
            </GlassCard>
          </motion.div>
        ))}
      </div>

      {/* Two Column: Recent Matches + Top Tournaments */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Recent Matches */}
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.3 }}
        >
          <GlassCard className="p-6" glowColor="cyan">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold flex items-center gap-2">
                <Calendar size={20} className="text-[#00d4ff]" />
                <span className="text-gradient-animated">Последние матчи</span>
              </h3>
            </div>
            {matchesLoading ? (
              <div className="flex justify-center py-8">
                <Loader2 size={24} className="animate-spin text-gray-500" />
              </div>
            ) : (
              <div className="space-y-3">
                {(recentMatches ?? []).map((match) => (
                  <div
                    key={match.id}
                    className="flex items-center gap-3 rounded-lg bg-white/5 p-3"
                  >
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center justify-between">
                        <span className="text-sm font-medium text-white truncate">
                          {match.homeTeam}
                        </span>
                        <span className="text-lg font-bold text-[#00d4ff] mx-2">
                          {match.homeScore}:{match.awayScore}
                        </span>
                        <span className="text-sm font-medium text-white truncate text-right">
                          {match.awayTeam}
                        </span>
                      </div>
                      <div className="flex items-center justify-between mt-1">
                        <span className="text-xs text-gray-500">
                          {new Date(match.date).toLocaleDateString('ru-RU', {
                            day: 'numeric',
                            month: 'short',
                          })}
                        </span>
                        <span className="text-xs text-gray-600 truncate">{cleanTournamentName(match.tournament)}</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </GlassCard>
        </motion.div>

        {/* Top Tournaments */}
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.35 }}
        >
          <GlassCard className="p-6" glowColor="purple">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold flex items-center gap-2">
                <Trophy size={20} className="text-[#8b5cf6]" />
                <span className="text-gradient-animated">Активные турниры</span>
              </h3>
              <Link
                to="/explore/tournaments"
                className="text-sm text-[#8b5cf6] hover:underline flex items-center gap-1"
              >
                Все турниры
                <ArrowRight size={14} />
              </Link>
            </div>
            {tournamentsLoading ? (
              <div className="flex justify-center py-8">
                <Loader2 size={24} className="animate-spin text-gray-500" />
              </div>
            ) : (
              <div className="space-y-3">
                {(tournaments ?? [])
                  .filter((t) => !t.isEnded)
                  .slice(0, 4)
                  .map((tournament) => (
                    <Link
                      key={tournament.id}
                      to={`/explore/tournaments/detail/${tournament.id}`}
                      className="flex items-center gap-3 rounded-lg bg-white/5 p-3 hover:bg-white/[0.07] transition-colors"
                    >
                      <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-[#8b5cf6]/20 text-[#8b5cf6]">
                        <Trophy size={18} />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium text-white truncate">
                          {cleanTournamentName(tournament.name)}
                        </p>
                        <p className="text-xs text-gray-500">
                          {tournament.teamsCount} команд · {tournament.matchesCount} матчей
                        </p>
                      </div>
                      <ArrowRight size={16} className="text-gray-600 flex-shrink-0" />
                    </Link>
                  ))}
              </div>
            )}
          </GlassCard>
        </motion.div>
      </div>

      {/* Top Scorers */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
      >
        <GlassCard className="p-6" glowColor="blue">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-white flex items-center gap-2">
              <TrendingUp size={20} className="text-[#00d4ff]" />
              <span className="text-gradient-animated">Топ-бомбардиры</span>
            </h3>
            <div className="flex items-center gap-2">
              {season && (
                <span className="text-[11px] bg-[#00d4ff]/10 text-[#00d4ff] px-2 py-0.5 rounded border border-[#00d4ff]/20">
                  {season}
                </span>
              )}
              <span className="text-[11px] text-gray-500">По всей России</span>
            </div>
          </div>
          {scorersLoading ? (
            <div className="flex justify-center py-8">
              <Loader2 size={24} className="animate-spin text-gray-500" />
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="text-gray-500 border-b border-white/5">
                    <th className="text-left py-2 pr-4">#</th>
                    <th className="text-left py-2 pr-4">Игрок</th>
                    <th className="text-left py-2 pr-4">Команда</th>
                    <th className="text-center py-2 px-2">И</th>
                    <th className="text-center py-2 px-2">Г</th>
                    <th className="text-center py-2 px-2">П</th>
                    <th className="text-center py-2 px-2 font-semibold text-[#00d4ff]">О</th>
                  </tr>
                </thead>
                <tbody>
                  {(topScorers ?? []).map((scorer) => (
                    <tr
                      key={scorer.id}
                      className="border-b border-white/5 last:border-0"
                    >
                      <td className="py-2.5 pr-4 text-gray-500 font-medium">{scorer.rank}</td>
                      <td className="py-2.5 pr-4 font-medium text-white">{scorer.name}</td>
                      <td className="py-2.5 pr-4 text-gray-400">{scorer.team}</td>
                      <td className="py-2.5 px-2 text-center text-gray-400">{scorer.games}</td>
                      <td className="py-2.5 px-2 text-center text-gray-300">{scorer.goals}</td>
                      <td className="py-2.5 px-2 text-center text-gray-300">{scorer.assists}</td>
                      <td className="py-2.5 px-2 text-center font-semibold text-[#00d4ff]">{scorer.points}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </GlassCard>
      </motion.div>

      {/* Player Tracking (PRO locked) */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.45 }}
      >
        <FeatureLockedOverlay
          feature="player_comparison"
          blurAmount="sm"
          previewContent={
            <GlassCard className="p-6">
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                <Star size={20} className="text-[#f59e0b]" />
                Отслеживание игроков
              </h3>
              <div className="grid grid-cols-3 gap-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="rounded-lg bg-white/5 p-4 text-center">
                    <div className="h-12 w-12 rounded-full bg-white/10 mx-auto mb-2" />
                    <div className="h-3 w-20 bg-white/10 rounded mx-auto mb-1" />
                    <div className="h-2 w-14 bg-white/5 rounded mx-auto" />
                  </div>
                ))}
              </div>
            </GlassCard>
          }
        />
      </motion.div>
    </div>
  )
})

function KpiValue({ value }: { value: number }) {
  const animated = useCountUp(value)
  return (
    <p className="text-2xl font-bold text-white">
      {animated.toLocaleString('ru-RU')}
    </p>
  )
}
