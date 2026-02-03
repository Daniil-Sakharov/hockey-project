import { cn } from '@/shared/lib/utils'

interface Props {
  seasons: string[]
  activeSeason: string
  onSeasonChange: (season: string) => void
  birthYears: number[]
  activeBirthYear: number | null
  onBirthYearChange: (year: number) => void
}

export function TournamentFilters({
  seasons,
  activeSeason,
  onSeasonChange,
  birthYears,
  activeBirthYear,
  onBirthYearChange,
}: Props) {
  return (
    <div className="space-y-3">
      {/* Season chips */}
      {seasons.length > 0 && (
        <div className="flex items-center gap-2 flex-wrap">
          <span className="text-xs text-gray-500 mr-1">Сезон</span>
          {seasons.map((s) => (
            <button
              key={s}
              onClick={() => onSeasonChange(s)}
              className={cn(
                'rounded-lg px-3 py-1.5 text-sm transition-all duration-200',
                activeSeason === s
                  ? 'bg-[#00d4ff]/20 text-[#00d4ff] border border-[#00d4ff]/30'
                  : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/[0.07]',
              )}
            >
              {s}
            </button>
          ))}
        </div>
      )}

      {/* Birth year chips */}
      {birthYears.length > 0 && (
        <div className="flex items-center gap-2 flex-wrap">
          <span className="text-xs text-gray-500 mr-1">Год рождения</span>
          {birthYears.map((y) => (
            <button
              key={y}
              onClick={() => onBirthYearChange(y)}
              className={cn(
                'rounded-lg px-3 py-1.5 text-sm transition-all duration-200',
                activeBirthYear === y
                  ? 'bg-[#8b5cf6]/20 text-[#8b5cf6] border border-[#8b5cf6]/30'
                  : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/[0.07]',
              )}
            >
              {y}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
