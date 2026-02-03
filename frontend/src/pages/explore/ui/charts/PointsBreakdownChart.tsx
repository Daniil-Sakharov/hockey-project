import { useMemo } from 'react'
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts'
import { GlassCard } from '@/shared/ui'
import type { SeasonAggregated } from './PlayerChartsSection'

interface Props {
  data: SeasonAggregated[]
}

const COLORS = ['#00d4ff', '#8b5cf6']

export function PointsBreakdownChart({ data }: Props) {
  const breakdown = useMemo(() => {
    const totalGoals = data.reduce((sum, s) => sum + s.goals, 0)
    const totalAssists = data.reduce((sum, s) => sum + s.assists, 0)
    return [
      { name: 'Голы', value: totalGoals },
      { name: 'Передачи', value: totalAssists },
    ]
  }, [data])

  const total = breakdown[0].value + breakdown[1].value

  return (
    <GlassCard className="p-4">
      <h4 className="text-sm font-medium text-gray-300 mb-3">Разбивка очков</h4>
      <div className="flex items-center">
        <ResponsiveContainer width="50%" height={180}>
          <PieChart>
            <Pie
              data={breakdown}
              cx="50%"
              cy="50%"
              innerRadius={50}
              outerRadius={75}
              paddingAngle={3}
              dataKey="value"
              stroke="none"
            >
              {breakdown.map((_, i) => (
                <Cell key={i} fill={COLORS[i]} />
              ))}
            </Pie>
            <Tooltip
              content={({ active, payload }) => {
                if (!active || !payload?.length) return null
                const item = payload[0]
                return (
                  <div className="rounded-lg border border-white/10 bg-[#0d1224] px-3 py-2 shadow-xl">
                    <p className="text-xs font-medium" style={{ color: item.payload.fill }}>
                      {item.name}: {item.value}
                    </p>
                  </div>
                )
              }}
            />
          </PieChart>
        </ResponsiveContainer>
        <div className="flex-1 space-y-3">
          {breakdown.map((item, i) => {
            const pct = total > 0 ? ((item.value / total) * 100).toFixed(1) : '0'
            return (
              <div key={item.name}>
                <div className="flex items-center justify-between mb-1">
                  <span className="text-xs text-gray-400 flex items-center gap-1.5">
                    <span className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: COLORS[i] }} />
                    {item.name}
                  </span>
                  <span className="text-xs font-medium text-white">{item.value}</span>
                </div>
                <div className="h-1.5 rounded-full bg-white/5 overflow-hidden">
                  <div
                    className="h-full rounded-full transition-all duration-700"
                    style={{ width: `${pct}%`, backgroundColor: COLORS[i] }}
                  />
                </div>
                <p className="text-[10px] text-gray-500 mt-0.5">{pct}%</p>
              </div>
            )
          })}
        </div>
      </div>
    </GlassCard>
  )
}
