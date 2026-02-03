import { useMemo } from 'react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { GlassCard } from '@/shared/ui'
import type { SeasonAggregated } from './PlayerChartsSection'
import { ChartTooltip } from './ChartTooltip'

interface Props {
  data: SeasonAggregated[]
}

export function AvgPerGameChart({ data }: Props) {
  const avgData = useMemo(() =>
    data.map((s) => ({
      season: s.season,
      goalsPerGame: s.games > 0 ? +(s.goals / s.games).toFixed(2) : 0,
      assistsPerGame: s.games > 0 ? +(s.assists / s.games).toFixed(2) : 0,
      pointsPerGame: s.games > 0 ? +(s.points / s.games).toFixed(2) : 0,
    })),
  [data])

  return (
    <GlassCard className="p-4">
      <h4 className="text-sm font-medium text-gray-300 mb-3">Средние за игру</h4>
      <div className="flex gap-4 mb-2 text-xs">
        <Legend color="#ec4899" label="Очков" />
        <Legend color="#00d4ff" label="Голов" />
        <Legend color="#8b5cf6" label="Передач" />
      </div>
      <ResponsiveContainer width="100%" height={220}>
        <BarChart data={avgData} margin={{ top: 5, right: 5, left: -15, bottom: 0 }}>
          <CartesianGrid stroke="rgba(255,255,255,0.05)" />
          <XAxis dataKey="season" tick={{ fill: '#64748b', fontSize: 11 }} tickLine={false} axisLine={false} />
          <YAxis tick={{ fill: '#64748b', fontSize: 11 }} tickLine={false} axisLine={false} />
          <Tooltip content={<ChartTooltip />} />
          <Bar dataKey="pointsPerGame" fill="#ec4899" radius={[4, 4, 0, 0]} barSize={14} />
          <Bar dataKey="goalsPerGame" fill="#00d4ff" radius={[4, 4, 0, 0]} barSize={14} />
          <Bar dataKey="assistsPerGame" fill="#8b5cf6" radius={[4, 4, 0, 0]} barSize={14} />
        </BarChart>
      </ResponsiveContainer>
    </GlassCard>
  )
}

function Legend({ color, label }: { color: string; label: string }) {
  return (
    <span className="flex items-center gap-1.5 text-gray-400">
      <span className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: color }} />
      {label}
    </span>
  )
}
