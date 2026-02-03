import { memo, useState, useMemo } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Medal, ArrowUpDown, Loader2 } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import { useRankings, useRankingsFilters } from '@/shared/api/useExploreQueries'
import type { RankedPlayer } from '@/shared/api/exploreTypes'
import { RankingsFilters } from './player/RankingsFilters'
import { RankingsTable } from './player/RankingsTable'

type SortKey = 'points' | 'goals' | 'assists' | 'penaltyMinutes' | 'plusMinus'

const SORT_OPTIONS: { key: SortKey; label: string }[] = [
  { key: 'points', label: 'Очки' },
  { key: 'goals', label: 'Голы' },
  { key: 'assists', label: 'Передачи' },
  { key: 'plusMinus', label: '+/-' },
  { key: 'penaltyMinutes', label: 'Штраф' },
]

const MEDAL_COLORS = ['text-[#f59e0b]', 'text-gray-300', 'text-[#cd7f32]']

export const RankingsPage = memo(function RankingsPage() {
  const [sortBy, setSortBy] = useState<SortKey>('points')
  const [birthYear, setBirthYear] = useState<number | null>(null)
  const [domain, setDomain] = useState<string | null>(null)
  const [tournamentId, setTournamentId] = useState<string | null>(null)
  const [groupName, setGroupName] = useState<string | null>(null)

  const params = useMemo(() => ({
    sort: sortBy, limit: 50,
    birthYear: birthYear ?? undefined,
    domain: domain ?? undefined,
    tournamentId: tournamentId ?? undefined,
    groupName: groupName ?? undefined,
  }), [sortBy, birthYear, domain, tournamentId, groupName])

  const { data: rankingsData, isLoading } = useRankings(params)
  const { data: filters } = useRankingsFilters()
  const ranked = rankingsData?.players
  const season = rankingsData?.season

  const subtitle = useMemo(() => {
    const parts: string[] = []
    if (tournamentId && filters?.tournaments) {
      const t = filters.tournaments.find((x) => x.id === tournamentId)
      if (t) parts.push(t.name)
    } else if (domain && filters?.domains) {
      const d = filters.domains.find((x) => x.domain === domain)
      if (d) parts.push(d.label)
    }
    if (groupName) parts.push(groupName)
    return parts.length > 0 ? parts.join(' · ') : 'По всей России'
  }, [domain, tournamentId, groupName, filters])

  return (
    <div className="space-y-4">
      <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }}>
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-white">Рейтинг игроков</h1>
            <p className="text-gray-400">Лучшие игроки по статистике</p>
          </div>
          <div className="text-right">
            {season && (
              <span className="text-xs bg-[#00d4ff]/10 text-[#00d4ff] px-2.5 py-1 rounded-lg border border-[#00d4ff]/20">
                Сезон {season}
              </span>
            )}
            <p className="text-[11px] text-gray-500 mt-1.5">{subtitle}</p>
          </div>
        </div>
      </motion.div>

      {filters && (
        <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.05 }}>
          <RankingsFilters
            birthYears={filters.birthYears ?? []}
            activeBirthYear={birthYear}
            onBirthYearChange={(y) => { setBirthYear(y); setTournamentId(null); setGroupName(null) }}
            domains={filters.domains ?? []}
            activeDomain={domain}
            onDomainChange={(d) => { setDomain(d); setTournamentId(null); setGroupName(null) }}
            tournaments={filters.tournaments ?? []}
            activeTournament={tournamentId}
            onTournamentChange={(t) => { setTournamentId(t); setGroupName(null) }}
            groups={filters.groups ?? []}
            activeGroup={groupName}
            onGroupChange={setGroupName}
          />
        </motion.div>
      )}

      <motion.div initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.1 }}>
        <div className="flex items-center gap-2 flex-wrap">
          <ArrowUpDown size={16} className="text-gray-500" />
          {SORT_OPTIONS.map((opt) => (
            <button
              key={opt.key}
              onClick={() => setSortBy(opt.key)}
              className={cn(
                'px-3 py-1.5 rounded-lg text-xs font-medium transition-all',
                sortBy === opt.key
                  ? 'bg-[#00d4ff]/20 text-[#00d4ff] border border-[#00d4ff]/30'
                  : 'bg-white/5 text-gray-400 hover:bg-white/10',
              )}
            >
              {opt.label}
            </button>
          ))}
        </div>
      </motion.div>

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Loader2 size={32} className="animate-spin text-gray-500" />
        </div>
      ) : (
        <>
          <RankingsTable ranked={ranked ?? []} sortBy={sortBy} />
          {(ranked ?? []).length >= 3 && (
            <TopThreeCards ranked={ranked!} sortBy={sortBy} />
          )}
        </>
      )}
    </div>
  )
})

function TopThreeCards({ ranked, sortBy }: { ranked: RankedPlayer[]; sortBy: SortKey }) {
  return (
    <div className="grid grid-cols-3 gap-4">
      {ranked.slice(0, 3).map((player, i) => (
        <motion.div key={player.id} initial={{ opacity: 0, scale: 0.9 }} animate={{ opacity: 1, scale: 1 }} transition={{ delay: 0.3 + i * 0.1 }}>
          <Link to={`/explore/players/${player.id}`}>
            <GlassCard className={cn('p-4 text-center hover:bg-white/[0.04] transition-colors', i === 0 && 'ring-1 ring-[#f59e0b]/30')}>
              <Medal size={24} className={cn('mx-auto mb-2', MEDAL_COLORS[i])} />
              <p className="text-sm font-semibold text-white">{player.name}</p>
              <p className="text-xs text-gray-500 mt-0.5">{player.team}</p>
              <p className="text-lg font-bold text-[#00d4ff] mt-2">{player[sortBy]}</p>
              <p className="text-[10px] text-gray-600">{SORT_OPTIONS.find((o) => o.key === sortBy)?.label}</p>
            </GlassCard>
          </Link>
        </motion.div>
      ))}
    </div>
  )
}
