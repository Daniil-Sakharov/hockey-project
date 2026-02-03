const LABELS: Record<string, string> = {
  points: 'Очки',
  goals: 'Голы',
  assists: 'Передачи',
  plusMinus: '+/-',
  penaltyMinutes: 'Штраф',
  games: 'Игры',
  goalsPerGame: 'Голов/игра',
  assistsPerGame: 'Передач/игра',
  pointsPerGame: 'Очков/игра',
}

interface TooltipProps {
  active?: boolean
  payload?: Array<{ name: string; value: number; color: string }>
  label?: string
}

export function ChartTooltip({ active, payload, label }: TooltipProps) {
  if (!active || !payload?.length) return null

  return (
    <div className="rounded-lg border border-white/10 bg-[#0d1224] px-3 py-2 shadow-xl">
      <p className="text-xs text-gray-400 mb-1">{label}</p>
      {payload.map((entry) => (
        <p key={entry.name} className="text-xs font-medium" style={{ color: entry.color }}>
          {LABELS[entry.name] ?? entry.name}: {typeof entry.value === 'number' ? entry.value.toLocaleString('ru-RU', { maximumFractionDigits: 2 }) : entry.value}
        </p>
      ))}
    </div>
  )
}
