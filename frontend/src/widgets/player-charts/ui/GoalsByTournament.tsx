import { memo } from 'react'
import { HoloBarChart } from '@/widgets/holographic-charts'
import type { PlayerSeasonStats } from '@/entities/player'
import { Skeleton } from '@/shared/ui'

interface GoalsByTournamentProps {
  data: PlayerSeasonStats[] | null
  isLoading?: boolean
}

export const GoalsByTournament = memo(function GoalsByTournament({
  data,
  isLoading = false,
}: GoalsByTournamentProps) {
  if (isLoading) {
    return (
      <div className="glass-card rounded-xl p-6">
        <Skeleton className="mb-4 h-4 w-32" />
        <Skeleton className="h-32 w-full" />
      </div>
    )
  }

  if (!data || data.length === 0) {
    return (
      <div className="glass-card rounded-xl p-6">
        <h3 className="mb-4 text-sm font-medium uppercase tracking-wider text-gray-400">
          Голы по турнирам
        </h3>
        <p className="text-center text-gray-500">Нет данных</p>
      </div>
    )
  }

  const chartData = data.map((season) => ({
    label: season.seasonName.slice(0, 8),
    value: season.goals,
  }))

  return <HoloBarChart data={chartData} title="Голы по турнирам" />
})
