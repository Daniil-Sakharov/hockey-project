import { useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Loader2 } from 'lucide-react'
import { GlassCard } from '@/shared/ui'

import { usePlayerStats } from '@/shared/api/useExploreQueries'
import type { PlayerStatEntry } from '@/shared/api/exploreTypes'
import { TournamentBlock } from './player/PlayerStatsTable'
import type { TournamentGroup } from './player/PlayerStatsTable'

interface SeasonGroup {
  season: string
  tournaments: TournamentGroup[]
}

function groupBySeason(stats: PlayerStatEntry[]): SeasonGroup[] {
  const map = new Map<string, Map<string, { id: string; name: string; byYear: Map<number, PlayerStatEntry[]> }>>()
  for (const s of stats) {
    if (s.groupName === 'Общая статистика') continue
    if (!map.has(s.season)) map.set(s.season, new Map())
    const tMap = map.get(s.season)!
    if (!tMap.has(s.tournamentId)) {
      tMap.set(s.tournamentId, { id: s.tournamentId, name: s.tournamentName, byYear: new Map() })
    }
    const t = tMap.get(s.tournamentId)!
    const year = s.birthYear || 0
    if (!t.byYear.has(year)) t.byYear.set(year, [])
    t.byYear.get(year)!.push(s)
  }
  return Array.from(map.entries()).map(([season, tMap]) => ({
    season,
    tournaments: Array.from(tMap.values()).map((t) => ({
      id: t.id,
      name: t.name,
      birthYearGroups: Array.from(t.byYear.entries())
        .sort(([a], [b]) => b - a)
        .map(([birthYear, entries]) => ({ birthYear, entries })),
    })),
  }))
}

export function PlayerStatsHistory({ playerId }: { playerId: string }) {
  const navigate = useNavigate()
  const { data: stats, isLoading } = usePlayerStats(playerId)

  const seasons = useMemo(() => {
    if (!stats || stats.length === 0) return []
    return groupBySeason(stats)
  }, [stats])

  if (isLoading) {
    return (
      <div className="flex justify-center py-8">
        <Loader2 size={24} className="animate-spin text-gray-500" />
      </div>
    )
  }

  if (seasons.length === 0) return null

  return (
    <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.5 }}>
      <h3 className="text-lg font-semibold text-white mb-4">Статистика по сезонам</h3>
      <div className="space-y-4">
        {seasons.map((sg) => (
          <GlassCard key={sg.season} className="p-4">
            <h4 className="text-sm font-semibold text-[#8b5cf6] mb-3">Сезон {sg.season}</h4>
            <div className="space-y-3">
              {sg.tournaments.map((tg) => (
                <div key={tg.id}>
                  <TournamentBlock tg={tg} navigate={navigate} />
                </div>
              ))}
            </div>
          </GlassCard>
        ))}
      </div>
    </motion.div>
  )
}
