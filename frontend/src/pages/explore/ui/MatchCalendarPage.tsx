import { memo, useState, useMemo } from 'react'
import { motion } from 'framer-motion'
import { Calendar, Clock, MapPin, Filter, Loader2 } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import { cleanTournamentName } from '@/shared/lib/formatters'
import { useUpcomingMatches } from '@/shared/api/useExploreQueries'

export const MatchCalendarPage = memo(function MatchCalendarPage() {
  const [tournamentFilter, setTournamentFilter] = useState('all')
  const { data: upcoming, isLoading } = useUpcomingMatches(tournamentFilter, 50)

  const tournaments = useMemo(() => {
    if (!upcoming) return []
    return Array.from(new Set(upcoming.map((m) => m.tournament)))
  }, [upcoming])

  // Group by date
  const grouped = useMemo(() => {
    if (!upcoming) return []
    const map = new Map<string, typeof upcoming>()
    for (const match of upcoming) {
      const existing = map.get(match.date) || []
      existing.push(match)
      map.set(match.date, existing)
    }
    return Array.from(map.entries()).sort(([a], [b]) => a.localeCompare(b))
  }, [upcoming])

  return (
    <div className="space-y-6">
      <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }}>
        <h1 className="text-2xl font-bold text-white">Календарь матчей</h1>
        <p className="text-gray-400">Предстоящие игры</p>
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
                ? 'bg-[#8b5cf6]/20 text-[#8b5cf6] border border-[#8b5cf6]/30'
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
                  ? 'bg-[#8b5cf6]/20 text-[#8b5cf6] border border-[#8b5cf6]/30'
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
          {/* Upcoming grouped by date */}
          {grouped.map(([date, matches], groupIdx) => (
            <motion.div
              key={date}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.15 + groupIdx * 0.05 }}
            >
              <div className="flex items-center gap-2 mb-3">
                <Calendar size={16} className="text-[#8b5cf6]" />
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
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium text-white truncate flex-1">
                        {match.homeTeam}
                      </span>
                      <div className="flex items-center gap-2 mx-3">
                        <span className="text-sm font-semibold text-[#8b5cf6]">vs</span>
                      </div>
                      <span className="text-sm font-medium text-white truncate flex-1 text-right">
                        {match.awayTeam}
                      </span>
                    </div>
                    <div className="flex items-center justify-between mt-2">
                      <span className="inline-flex items-center gap-1 text-xs text-gray-500">
                        <Clock size={12} />
                        {match.time}
                      </span>
                      <span className="text-xs text-gray-600">{cleanTournamentName(match.tournament)}</span>
                      {match.venue && (
                        <span className="inline-flex items-center gap-1 text-xs text-gray-500">
                          <MapPin size={12} />
                          {match.venue}
                        </span>
                      )}
                    </div>
                  </GlassCard>
                ))}
              </div>
            </motion.div>
          ))}

          {(upcoming ?? []).length === 0 && (
            <div className="flex flex-col items-center justify-center py-16 text-center">
              <Calendar size={48} className="text-gray-700 mb-4" />
              <p className="text-gray-400">Нет запланированных матчей</p>
            </div>
          )}
        </>
      )}
    </div>
  )
})
