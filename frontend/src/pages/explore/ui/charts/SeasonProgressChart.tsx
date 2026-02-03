import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { GlassCard } from '@/shared/ui'
import type { SeasonAggregated } from './PlayerChartsSection'
import { ChartTooltip } from './ChartTooltip'

interface Props {
  data: SeasonAggregated[]
}

export function SeasonProgressChart({ data }: Props) {
  return (
    <GlassCard className="p-4">
      <h4 className="text-sm font-medium text-gray-300 mb-3">Прогресс по сезонам</h4>
      <div className="flex gap-4 mb-2 text-xs">
        <Legend color="#ec4899" label="Очки" />
        <Legend color="#00d4ff" label="Голы" />
        <Legend color="#8b5cf6" label="Передачи" />
      </div>
      <ResponsiveContainer width="100%" height={220}>
        <AreaChart data={data} margin={{ top: 5, right: 5, left: -15, bottom: 0 }}>
          <defs>
            <linearGradient id="gradPoints" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="#ec4899" stopOpacity={0.3} />
              <stop offset="100%" stopColor="#ec4899" stopOpacity={0} />
            </linearGradient>
            <linearGradient id="gradGoals" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="#00d4ff" stopOpacity={0.3} />
              <stop offset="100%" stopColor="#00d4ff" stopOpacity={0} />
            </linearGradient>
            <linearGradient id="gradAssists" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="#8b5cf6" stopOpacity={0.3} />
              <stop offset="100%" stopColor="#8b5cf6" stopOpacity={0} />
            </linearGradient>
          </defs>
          <CartesianGrid stroke="rgba(255,255,255,0.05)" />
          <XAxis dataKey="season" tick={{ fill: '#64748b', fontSize: 11 }} tickLine={false} axisLine={false} />
          <YAxis tick={{ fill: '#64748b', fontSize: 11 }} tickLine={false} axisLine={false} />
          <Tooltip content={<ChartTooltip />} />
          <Area type="monotone" dataKey="points" stroke="#ec4899" fill="url(#gradPoints)" strokeWidth={2} />
          <Area type="monotone" dataKey="goals" stroke="#00d4ff" fill="url(#gradGoals)" strokeWidth={2} />
          <Area type="monotone" dataKey="assists" stroke="#8b5cf6" fill="url(#gradAssists)" strokeWidth={2} />
        </AreaChart>
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
