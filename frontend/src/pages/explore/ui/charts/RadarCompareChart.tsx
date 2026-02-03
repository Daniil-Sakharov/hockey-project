import { useMemo } from 'react'
import { RadarChart, Radar, PolarGrid, PolarAngleAxis, PolarRadiusAxis, ResponsiveContainer, Tooltip } from 'recharts'
import { GlassCard } from '@/shared/ui'
import type { SeasonAggregated } from './PlayerChartsSection'

interface Props {
  data: SeasonAggregated[]
}

const METRICS = [
  { key: 'goals', label: 'Голы' },
  { key: 'assists', label: 'Передачи' },
  { key: 'plusMinus', label: '+/-' },
  { key: 'penaltyMinutes', label: 'Штраф' },
  { key: 'games', label: 'Игры' },
] as const

export function RadarCompareChart({ data }: Props) {
  const radarData = useMemo(() => {
    if (data.length < 2) return []

    const current = data[data.length - 1]
    const previous = data[data.length - 2]

    // Find max for each metric across all seasons for normalization
    const maxVals: Record<string, number> = {}
    for (const m of METRICS) {
      maxVals[m.key] = Math.max(...data.map((s) => Math.abs(s[m.key])), 1)
    }

    return METRICS.map((m) => ({
      metric: m.label,
      current: Math.round((Math.abs(current[m.key]) / maxVals[m.key]) * 100),
      previous: Math.round((Math.abs(previous[m.key]) / maxVals[m.key]) * 100),
    }))
  }, [data])

  if (radarData.length === 0) return null

  const currentSeason = data[data.length - 1]?.season
  const prevSeason = data[data.length - 2]?.season

  return (
    <GlassCard className="p-4">
      <h4 className="text-sm font-medium text-gray-300 mb-3">Сравнение сезонов</h4>
      <div className="flex gap-4 mb-2 text-xs">
        <Legend color="#00d4ff" label={currentSeason} />
        <Legend color="#8b5cf6" label={prevSeason} />
      </div>
      <ResponsiveContainer width="100%" height={220}>
        <RadarChart data={radarData} cx="50%" cy="50%" outerRadius="70%">
          <PolarGrid stroke="rgba(255,255,255,0.1)" />
          <PolarAngleAxis dataKey="metric" tick={{ fill: '#64748b', fontSize: 11 }} />
          <PolarRadiusAxis tick={false} axisLine={false} domain={[0, 100]} />
          <Radar
            name={currentSeason}
            dataKey="current"
            stroke="#00d4ff"
            fill="#00d4ff"
            fillOpacity={0.2}
            strokeWidth={2}
          />
          <Radar
            name={prevSeason}
            dataKey="previous"
            stroke="#8b5cf6"
            fill="#8b5cf6"
            fillOpacity={0.1}
            strokeWidth={2}
          />
          <Tooltip
            content={({ active, payload }) => {
              if (!active || !payload?.length) return null
              return (
                <div className="rounded-lg border border-white/10 bg-[#0d1224] px-3 py-2 shadow-xl">
                  {payload.map((entry) => (
                    <p key={entry.name} className="text-xs font-medium" style={{ color: entry.color as string }}>
                      {entry.name}: {entry.value}%
                    </p>
                  ))}
                </div>
              )
            }}
          />
        </RadarChart>
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
