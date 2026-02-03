import { memo } from 'react'
import { HoloLineChart } from '@/widgets/holographic-charts'
import type { PlayerPerformancePoint } from '@/entities/player'
import { Skeleton } from '@/shared/ui'

interface PerformanceTimelineProps {
  data: PlayerPerformancePoint[] | null
  isLoading?: boolean
}

export const PerformanceTimeline = memo(function PerformanceTimeline({
  data,
  isLoading = false,
}: PerformanceTimelineProps) {
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
          Динамика очков
        </h3>
        <p className="text-center text-gray-500">Нет данных</p>
      </div>
    )
  }

  const chartData = data.map((point) => ({
    label: point.month,
    value: point.points,
  }))

  return <HoloLineChart data={chartData} title="Динамика очков" />
})
