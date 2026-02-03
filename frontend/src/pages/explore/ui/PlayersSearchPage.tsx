import { memo, useCallback } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Search, Filter, UserRound, ArrowRight, Loader2, Calendar } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import { usePlayersSearch, useSeasons } from '@/shared/api/useExploreQueries'

type PositionFilter = 'all' | 'forward' | 'defender' | 'goalie'

const POSITION_LABELS: Record<string, string> = {
  all: 'Все',
  forward: 'Нападающий',
  defender: 'Защитник',
  goalie: 'Вратарь',
}

export const PlayersSearchPage = memo(function PlayersSearchPage() {
  const [searchParams, setSearchParams] = useSearchParams()

  const search = searchParams.get('q') ?? ''
  const position = (searchParams.get('position') ?? 'all') as PositionFilter
  const yearFilter = searchParams.get('year') ?? 'all'
  const seasonFilter = searchParams.get('season') ?? null

  const { data: seasons } = useSeasons()
  const activeSeason = seasonFilter === 'all' ? '' : (seasonFilter ?? '')

  const birthYear = yearFilter !== 'all' ? Number(yearFilter) : 0
  const { data, isLoading } = usePlayersSearch(search, position, activeSeason, birthYear, 50, 0)
  const players = data?.players ?? []
  const total = data?.total ?? 0

  const updateParam = useCallback((key: string, value: string) => {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev)
      if (!value || value === 'all') {
        next.delete(key)
      } else {
        next.set(key, value)
      }
      return next
    }, { replace: true })
  }, [setSearchParams])

  // Build query string for links to player profiles
  const currentQuery = searchParams.toString()

  return (
    <div className="space-y-6">
      <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }}>
        <h1 className="text-2xl font-bold text-white">Поиск игроков</h1>
        <p className="text-gray-400">Найдите игроков по имени, команде или позиции</p>
      </motion.div>

      {/* Season filter */}
      {seasons && seasons.length > 0 && (
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.05 }}>
          <div className="flex items-center gap-2 flex-wrap">
            <Calendar size={16} className="text-gray-500" />
            <button
              onClick={() => updateParam('season', 'all')}
              className={cn(
                'px-3 py-1.5 rounded-lg text-xs font-medium transition-all',
                activeSeason === ''
                  ? 'bg-[#8b5cf6]/20 text-[#8b5cf6] border border-[#8b5cf6]/30'
                  : 'bg-white/5 text-gray-400 hover:bg-white/10'
              )}
            >
              Все
            </button>
            {seasons.map((s) => (
              <button
                key={s}
                onClick={() => updateParam('season', s)}
                className={cn(
                  'px-3 py-1.5 rounded-lg text-xs font-medium transition-all',
                  activeSeason === s
                    ? 'bg-[#8b5cf6]/20 text-[#8b5cf6] border border-[#8b5cf6]/30'
                    : 'bg-white/5 text-gray-400 hover:bg-white/10'
                )}
              >
                {s}
              </button>
            ))}
          </div>
        </motion.div>
      )}

      {/* Filters */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        <GlassCard className="p-4">
          <div className="flex flex-col gap-4 md:flex-row md:items-center">
            {/* Search input */}
            <div className="relative flex-1">
              <Search size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500" />
              <input
                type="text"
                placeholder="Имя игрока или команда..."
                value={search}
                onChange={(e) => updateParam('q', e.target.value)}
                className={cn(
                  'w-full pl-10 pr-4 py-2.5 rounded-lg',
                  'bg-white/5 border border-white/10',
                  'text-white placeholder-gray-500 text-sm',
                  'focus:outline-none focus:border-[#00d4ff]/50'
                )}
              />
            </div>

            {/* Position filter */}
            <div className="flex items-center gap-2">
              <Filter size={16} className="text-gray-500" />
              <div className="flex gap-1">
                {(['all', 'forward', 'defender', 'goalie'] as PositionFilter[]).map((pos) => (
                  <button
                    key={pos}
                    onClick={() => updateParam('position', pos)}
                    className={cn(
                      'px-3 py-1.5 rounded-lg text-xs font-medium transition-all',
                      position === pos
                        ? 'bg-[#00d4ff]/20 text-[#00d4ff] border border-[#00d4ff]/30'
                        : 'bg-white/5 text-gray-400 hover:bg-white/10'
                    )}
                  >
                    {POSITION_LABELS[pos]}
                  </button>
                ))}
              </div>
            </div>

            {/* Year filter */}
            <select
              value={yearFilter}
              onChange={(e) => updateParam('year', e.target.value)}
              className={cn(
                'px-3 py-2.5 rounded-lg text-sm',
                'bg-white/5 border border-white/10 text-gray-300',
                'focus:outline-none focus:border-[#00d4ff]/50'
              )}
            >
              <option value="all">Все года</option>
              {Array.from({ length: 15 }, (_, i) => 2016 - i).map((y) => (
                <option key={y} value={y}>{y} г.р.</option>
              ))}
            </select>
          </div>
        </GlassCard>
      </motion.div>

      {/* Results count */}
      <p className="text-sm text-gray-500">
        {isLoading ? 'Загрузка...' : `Найдено: ${total} игроков`}
      </p>

      {/* Loading */}
      {isLoading ? (
        <div className="flex justify-center py-16">
          <Loader2 size={32} className="animate-spin text-gray-500" />
        </div>
      ) : (
        <>
          {/* Players grid */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {players.map((player, i) => {
              const profileParams = new URLSearchParams()
              if (activeSeason) profileParams.set('season', activeSeason)
              if (currentQuery) profileParams.set('from', currentQuery)
              const profileQs = profileParams.toString()
              return (
                <motion.div
                  key={player.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.05 * Math.min(i, 10) }}
                >
                  <Link to={`/explore/players/${player.id}${profileQs ? `?${profileQs}` : ''}`}>
                    <GlassCard className="p-4 hover:bg-white/[0.04] transition-colors cursor-pointer">
                      <div className="flex items-start gap-3">
                        {player.photoUrl ? (
                          <img
                            src={player.photoUrl}
                            alt={player.name}
                            className="h-12 w-12 rounded-xl object-cover flex-shrink-0"
                          />
                        ) : (
                          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-[#00d4ff]/20 text-[#00d4ff] flex-shrink-0">
                            <UserRound size={24} />
                          </div>
                        )}
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center justify-between">
                            <h3 className="text-sm font-semibold text-white truncate">{player.name}</h3>
                            {player.jerseyNumber > 0 && (
                              <span className="text-xs text-gray-500">#{player.jerseyNumber}</span>
                            )}
                          </div>
                          <p className="text-xs text-gray-400 mt-0.5">{player.team}</p>
                          <div className="flex items-center gap-2 mt-1">
                            <span className={cn(
                              'text-[10px] px-1.5 py-0.5 rounded font-medium',
                              player.position === 'forward' ? 'bg-[#00d4ff]/20 text-[#00d4ff]' :
                              player.position === 'defender' ? 'bg-[#8b5cf6]/20 text-[#8b5cf6]' :
                              'bg-[#f59e0b]/20 text-[#f59e0b]'
                            )}>
                              {POSITION_LABELS[player.position] ?? player.position}
                            </span>
                            <span className="text-[10px] text-gray-600">{player.birthYear} г.р.</span>
                          </div>
                        </div>
                      </div>

                      {player.stats && (
                        <div className="mt-3 flex items-center justify-between rounded-lg bg-white/5 px-3 py-2">
                          <StatCell label="И" value={player.stats.games} />
                          <StatCell label="Г" value={player.stats.goals} />
                          <StatCell label="П" value={player.stats.assists} />
                          <StatCell label="О" value={player.stats.points} highlight />
                          <ArrowRight size={14} className="text-gray-600" />
                        </div>
                      )}
                    </GlassCard>
                  </Link>
                </motion.div>
              )
            })}
          </div>

          {players.length === 0 && (
            <div className="flex flex-col items-center justify-center py-16 text-center">
              <UserRound size={48} className="text-gray-700 mb-4" />
              <p className="text-gray-400">Игроки не найдены</p>
              <p className="text-sm text-gray-600 mt-1">Попробуйте изменить параметры поиска</p>
            </div>
          )}
        </>
      )}
    </div>
  )
})

function StatCell({ label, value, highlight }: { label: string; value: number; highlight?: boolean }) {
  return (
    <div className="text-center">
      <p className={cn('text-sm font-semibold', highlight ? 'text-[#00d4ff]' : 'text-white')}>
        {value}
      </p>
      <p className="text-[10px] text-gray-500">{label}</p>
    </div>
  )
}
