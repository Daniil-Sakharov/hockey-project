import { memo } from 'react'
import { motion } from 'framer-motion'
import type { PlayerDetailedStats, PlayerSeasonStats, PlayerPerformancePoint } from '@/entities/player'
import { GoalsByTournament } from './GoalsByTournament'
import { PerformanceTimeline } from './PerformanceTimeline'
import { GoalsDistribution } from './GoalsDistribution'

interface PlayerChartsSectionProps {
  stats: PlayerDetailedStats | null
  seasonStats: PlayerSeasonStats[] | null
  performanceData: PlayerPerformancePoint[] | null
  isLoading?: boolean
}

export const PlayerChartsSection = memo(function PlayerChartsSection({
  stats,
  seasonStats,
  performanceData,
  isLoading = false,
}: PlayerChartsSectionProps) {
  return (
    <motion.section
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.3 }}
      className="mb-8"
    >
      <h2 className="mb-4 text-lg font-semibold text-white">
        Аналитика
      </h2>

      <div className="grid gap-6 lg:grid-cols-2 xl:grid-cols-3">
        <GoalsByTournament data={seasonStats} isLoading={isLoading} />
        <PerformanceTimeline data={performanceData} isLoading={isLoading} />
        <GoalsDistribution stats={stats} isLoading={isLoading} />
      </div>
    </motion.section>
  )
})
