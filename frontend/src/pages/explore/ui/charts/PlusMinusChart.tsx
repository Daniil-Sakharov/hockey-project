import { useMemo } from 'react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from 'recharts'
import { GlassCard } from '@/shared/ui'
import type { SeasonAggregated } from './PlayerChartsSection'
import { ChartTooltip } from './ChartTooltip'

interface Props {
  data: SeasonAggregated[]
}

export function PlusMinusChart({ data }: Props) {
  const chartData = useMemo(() =>
    data.map((s) => ({
      season: s.season,
      plusMinus: s.plusMinus,
    })),
  [data])

  return (
    <GlassCard className="p-4">
      <h4 className="text-sm font-medium text-gray-300 mb-3">Показатель +/-</h4>
      <div className="flex gap-4 mb-2 text-xs">
        <Legend color="#10b981" label="Положительный" />
        <Legend color="#ef4444" label="Отрицательный" />
      </div>
      <ResponsiveContainer width="100%" height={220}>
        <BarChart data={chartData} margin={{ top: 5, right: 5, left: -15, bottom: 0 }}>
          <CartesianGrid stroke="rgba(255,255,255,0.05)" />
          <XAxis dataKey="season" tick={{ fill: '#64748b', fontSize: 11 }} tickLine={false} axisLine={false} />
          <YAxis tick={{ fill: '#64748b', fontSize: 11 }} tickLine={false} axisLine={false} />
          <Tooltip content={<ChartTooltip />} />
          <Bar dataKey="plusMinus" radius={[4, 4, 0, 0]} barSize={28}>
            {chartData.map((entry, i) => (
              <Cell key={i} fill={entry.plusMinus >= 0 ? '#10b981' : '#ef4444'} />
            ))}
          </Bar>
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
