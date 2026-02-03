import { memo, useState, useMemo } from 'react'
import { motion } from 'framer-motion'
import { Calendar, Filter, Loader2 } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import { cleanTournamentName } from '@/shared/lib/formatters'
import { useRecentResults } from '@/shared/api/useExploreQueries'

export const MatchResultsPage = memo(function MatchResultsPage() {
  const [tournamentFilter, setTournamentFilter] = useState('all')
  const { data: results, isLoading } = useRecentResults(tournamentFilter, 50)

  const tournaments = useMemo(() => {
    if (!results) return []
    return Array.from(new Set(results.map((m) => m.tournament)))
  }, [results])

  // Group by date
  const grouped = useMemo(() => {
    if (!results) return []
    const map = new Map<string, typeof results>()
    for (const match of results) {
      const existing = map.get(match.date) || []
      existing.push(match)
      map.set(match.date, existing)
    }
    return Array.from(map.entries()).sort(([a], [b]) => b.localeCompare(a))
  }, [results])

  return (
    <div className="space-y-6">
      <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }}>
        <h1 className="text-2xl font-bold text-white">Результаты матчей</h1>
        <p className="text-gray-400">Лента завершённых матчей</p>
      </motion.div>

      {/* Tournament filter */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        <div className="flex items-center gap-2 flex-wrap">
          <Filter size={16} className="text-gray-500" />
          <button
            onClick={() => setTournamentFilter('all')}
            className={cn(
              'px-3 py-1.5 rounded-lg text-xs font-medium transition-all',
              tournamentFilter === 'all'
                ? 'bg-[#00d4ff]/20 text-[#00d4ff] border border-[#00d4ff]/30'
                : 'bg-white/5 text-gray-400 hover:bg-white/10'
            )}
          >
            Все
          </button>
          {tournaments.map((t) => (
            <button
              key={t}
              onClick={() => setTournamentFilter(t)}
              className={cn(
                'px-3 py-1.5 rounded-lg text-xs font-medium transition-all',
                tournamentFilter === t
                  ? 'bg-[#00d4ff]/20 text-[#00d4ff] border border-[#00d4ff]/30'
                  : 'bg-white/5 text-gray-400 hover:bg-white/10'
              )}
            >
              {cleanTournamentName(t)}
            </button>
          ))}
        </div>
      </motion.div>

      {/* Loading */}
      {isLoading ? (
        <div className="flex justify-center py-16">
          <Loader2 size={32} className="animate-spin text-gray-500" />
        </div>
      ) : (
        <>
          {/* Results grouped by date */}
          {grouped.map(([date, matches], groupIdx) => (
            <motion.div
              key={date}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.15 + groupIdx * 0.05 }}
            >
              <div className="flex items-center gap-2 mb-3">
                <Calendar size={16} className="text-[#00d4ff]" />
                <h3 className="text-sm font-medium text-gray-300">
                  {new Date(date).toLocaleDateString('ru-RU', {
                    weekday: 'long',
                    day: 'numeric',
                    month: 'long',
                  })}
                </h3>
              </div>

              <div className="space-y-2">
                {matches.map((match) => (
                  <GlassCard key={match.id} className="p-4">
                    <div className="flex items-center gap-3">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between">
                          <span className="text-sm font-medium text-white truncate flex-1">
                            {match.homeTeam}
                          </span>
                          <div className="flex items-center gap-1 mx-3">
                            <span className={cn(
                              'text-lg font-bold',
                              (match.homeScore ?? 0) > (match.awayScore ?? 0) ? 'text-[#10b981]' : 'text-white'
                            )}>
                              {match.homeScore}
                            </span>
                            <span className="text-gray-600">:</span>
                            <span className={cn(
                              'text-lg font-bold',
                              (match.awayScore ?? 0) > (match.homeScore ?? 0) ? 'text-[#10b981]' : 'text-white'
                            )}>
                              {match.awayScore}
                            </span>
                          </div>
                          <span className="text-sm font-medium text-white truncate flex-1 text-right">
                            {match.awayTeam}
                          </span>
                        </div>
                        <div className="flex items-center justify-between mt-1.5">
                          <span className="text-xs text-gray-600">{match.time}</span>
                          <span className="text-xs text-gray-500">{cleanTournamentName(match.tournament)}</span>
                          {match.venue && (
                            <span className="text-xs text-gray-600">{match.venue}</span>
                          )}
                        </div>
                      </div>
                    </div>
                  </GlassCard>
                ))}
              </div>
            </motion.div>
          ))}

          {(results ?? []).length === 0 && (
            <div className="flex flex-col items-center justify-center py-16 text-center">
              <Calendar size={48} className="text-gray-700 mb-4" />
              <p className="text-gray-400">Нет результатов</p>
            </div>
          )}
        </>
      )}
    </div>
  )
})
