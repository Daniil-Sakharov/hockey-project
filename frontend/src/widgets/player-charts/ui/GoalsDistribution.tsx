import { memo } from 'react'
import { HoloRingChart } from '@/widgets/holographic-charts'
import type { PlayerDetailedStats } from '@/entities/player'
import { Skeleton } from '@/shared/ui'

interface GoalsDistributionProps {
  stats: PlayerDetailedStats | null
  isLoading?: boolean
}

export const GoalsDistribution = memo(function GoalsDistribution({
  stats,
  isLoading = false,
}: GoalsDistributionProps) {
  if (isLoading) {
    return (
      <div className="glass-card rounded-xl p-6">
        <Skeleton className="mb-4 h-4 w-32" />
        <Skeleton className="mx-auto h-32 w-32 rounded-full" />
      </div>
    )
  }

  if (!stats) {
    return (
      <div className="glass-card rounded-xl p-6">
        <h3 className="mb-4 text-sm font-medium uppercase tracking-wider text-gray-400">
          Распределение голов
        </h3>
        <p className="text-center text-gray-500">Нет данных</p>
      </div>
    )
  }

  const chartData = [
    {
      label: 'Равные составы',
      value: stats.evenStrengthGoals,
      color: '#00d4ff',
    },
    {
      label: 'Большинство',
      value: stats.powerplayGoals,
      color: '#8b5cf6',
    },
    {
      label: 'Меньшинство',
      value: stats.shorthandedGoals,
      color: '#ec4899',
    },
  ]

  const totalGoals = stats.goals

  return (
    <HoloRingChart
      data={chartData}
      title="Распределение голов"
      centerValue={totalGoals.toString()}
      centerLabel="Всего"
    />
  )
})
