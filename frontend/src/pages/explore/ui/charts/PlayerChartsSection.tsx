import { useMemo, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { BarChart3, ChevronDown } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { usePlayerStats } from '@/shared/api/useExploreQueries'
import type { PlayerStatEntry } from '@/shared/api/exploreTypes'
import { SeasonProgressChart } from './SeasonProgressChart'
import { PointsBreakdownChart } from './PointsBreakdownChart'
import { AvgPerGameChart } from './AvgPerGameChart'
import { PenaltyChart } from './PenaltyChart'
import { PlusMinusChart } from './PlusMinusChart'
import { RadarCompareChart } from './RadarCompareChart'

export interface SeasonAggregated {
  season: string
  games: number
  goals: number
  assists: number
  points: number
  plusMinus: number
  penaltyMinutes: number
}

export function aggregateBySeason(stats: PlayerStatEntry[]): SeasonAggregated[] {
  const map = new Map<string, SeasonAggregated>()
  for (const s of stats) {
    if (s.groupName !== 'Общая статистика') continue
    const existing = map.get(s.season)
    if (existing) {
      existing.games += s.games
      existing.goals += s.goals
      existing.assists += s.assists
      existing.points += s.points
      existing.plusMinus += s.plusMinus
      existing.penaltyMinutes += s.penaltyMinutes
    } else {
      map.set(s.season, { season: s.season, ...s })
    }
  }
  return Array.from(map.values()).sort((a, b) => a.season.localeCompare(b.season))
}

interface Props {
  playerId: string
}

export function PlayerChartsSection({ playerId }: Props) {
  const { data: stats } = usePlayerStats(playerId)
  const [isOpen, setIsOpen] = useState(true)

  const seasonData = useMemo(() => {
    if (!stats || stats.length === 0) return []
    return aggregateBySeason(stats)
  }, [stats])

  if (seasonData.length < 2) return null

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.5 }}
      className="space-y-4"
    >
      <button
        onClick={() => setIsOpen((v) => !v)}
        className={cn(
          'flex items-center gap-3 w-full group',
          'rounded-xl px-4 py-3',
          'bg-white/[0.03] border border-white/5',
          'hover:bg-white/[0.06] hover:border-[#00d4ff]/20',
          'transition-all duration-200'
        )}
      >
        <BarChart3 size={20} className="text-[#00d4ff]" />
        <h3 className="text-lg font-semibold text-white">Аналитика</h3>
        <span className="bg-gradient-to-r from-[#f59e0b] to-[#ef4444] text-white text-xs px-2 py-0.5 rounded-full font-bold">
          PRO
        </span>
        <ChevronDown
          size={18}
          className={cn(
            'ml-auto text-gray-500 transition-transform duration-300',
            'group-hover:text-gray-300',
            isOpen && 'rotate-180'
          )}
        />
      </button>

      <AnimatePresence initial={false}>
        {isOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.3, ease: 'easeInOut' }}
            className="overflow-hidden"
          >
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
              <SeasonProgressChart data={seasonData} />
              <PointsBreakdownChart data={seasonData} />
              <AvgPerGameChart data={seasonData} />
              <PenaltyChart data={seasonData} />
              <PlusMinusChart data={seasonData} />
              <RadarCompareChart data={seasonData} />
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  )
}
