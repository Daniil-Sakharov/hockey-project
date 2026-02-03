import { useState } from 'react'
import { Calendar, MapPin, Trophy, Users, ChevronDown, ChevronUp } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { cleanTournamentName } from '@/shared/lib/formatters'
import type { DomainOption, TournamentOption, GroupOption } from '@/shared/api/exploreTypes'

const COLLAPSED_TOURNAMENT_COUNT = 5

interface Props {
  birthYears: number[]
  activeBirthYear: number | null
  onBirthYearChange: (year: number | null) => void
  domains: DomainOption[]
  activeDomain: string | null
  onDomainChange: (domain: string | null) => void
  tournaments: TournamentOption[]
  activeTournament: string | null
  onTournamentChange: (id: string | null) => void
  groups: GroupOption[]
  activeGroup: string | null
  onGroupChange: (name: string | null) => void
}

function ChipButton({
  active,
  onClick,
  color,
  children,
}: {
  active: boolean
  onClick: () => void
  color: 'purple' | 'green' | 'orange' | 'blue'
  children: React.ReactNode
}) {
  const colors = {
    purple: 'bg-[#8b5cf6]/20 text-[#8b5cf6] border-[#8b5cf6]/30',
    green: 'bg-[#10b981]/20 text-[#10b981] border-[#10b981]/30',
    orange: 'bg-[#f59e0b]/20 text-[#f59e0b] border-[#f59e0b]/30',
    blue: 'bg-[#3b82f6]/20 text-[#3b82f6] border-[#3b82f6]/30',
  }
  return (
    <button
      onClick={onClick}
      className={cn(
        'px-2.5 py-1 rounded-lg text-xs font-medium transition-all border',
        active ? colors[color] : 'bg-white/5 text-gray-400 border-transparent hover:bg-white/10',
      )}
    >
      {children}
    </button>
  )
}

export function RankingsFilters({
  birthYears, activeBirthYear, onBirthYearChange,
  domains, activeDomain, onDomainChange,
  tournaments, activeTournament, onTournamentChange,
  groups, activeGroup, onGroupChange,
}: Props) {
  const [tournamentsExpanded, setTournamentsExpanded] = useState(false)

  // Filter tournaments by domain and birth year
  const filteredTournaments = tournaments.filter((t) => {
    if (activeDomain && t.domain !== activeDomain) return false
    if (activeBirthYear && t.birthYears && !t.birthYears.includes(activeBirthYear)) return false
    return true
  })

  // Filter groups by selected tournament
  const filteredGroups = activeTournament
    ? groups.filter((g) => g.tournamentId === activeTournament)
    : []

  // Determine if we should collapse tournaments (only when "All regions" and many tournaments)
  const shouldCollapse = !activeDomain && filteredTournaments.length > COLLAPSED_TOURNAMENT_COUNT
  const visibleTournaments = shouldCollapse && !tournamentsExpanded
    ? filteredTournaments.slice(0, COLLAPSED_TOURNAMENT_COUNT)
    : filteredTournaments
  const hiddenCount = filteredTournaments.length - COLLAPSED_TOURNAMENT_COUNT

  return (
    <div className="space-y-2">
      {birthYears.length > 0 && (
        <div className="flex items-center gap-2 flex-wrap">
          <Calendar size={14} className="text-gray-500" />
          <ChipButton active={!activeBirthYear} onClick={() => onBirthYearChange(null)} color="purple">
            Все
          </ChipButton>
          {birthYears.map((y) => (
            <ChipButton key={y} active={activeBirthYear === y} onClick={() => onBirthYearChange(y)} color="purple">
              {y}
            </ChipButton>
          ))}
        </div>
      )}

      {domains.length > 0 && (
        <div className="flex items-center gap-2 flex-wrap">
          <MapPin size={14} className="text-gray-500" />
          <ChipButton active={!activeDomain} onClick={() => { onDomainChange(null); onTournamentChange(null); setTournamentsExpanded(false) }} color="green">
            Все регионы
          </ChipButton>
          {domains.map((d) => (
            <ChipButton
              key={d.domain}
              active={activeDomain === d.domain}
              onClick={() => { onDomainChange(d.domain); onTournamentChange(null); setTournamentsExpanded(false) }}
              color="green"
            >
              {d.label}
            </ChipButton>
          ))}
        </div>
      )}

      {filteredTournaments.length > 0 && (
        <div className="flex items-center gap-2 flex-wrap">
          <Trophy size={14} className="text-gray-500" />
          <ChipButton active={!activeTournament} onClick={() => onTournamentChange(null)} color="orange">
            Все турниры
          </ChipButton>
          {visibleTournaments.map((t) => (
            <ChipButton
              key={t.id}
              active={activeTournament === t.id}
              onClick={() => onTournamentChange(t.id)}
              color="orange"
            >
              {cleanTournamentName(t.name)}
            </ChipButton>
          ))}
          {shouldCollapse && (
            <button
              onClick={() => setTournamentsExpanded(!tournamentsExpanded)}
              className="px-2.5 py-1 rounded-lg text-xs font-medium transition-all flex items-center gap-1 bg-white/5 text-gray-400 hover:bg-white/10 hover:text-gray-300"
            >
              {tournamentsExpanded ? (
                <>Свернуть <ChevronUp size={12} /></>
              ) : (
                <>Ещё {hiddenCount} <ChevronDown size={12} /></>
              )}
            </button>
          )}
        </div>
      )}

      {filteredGroups.length > 0 && (
        <div className="flex items-center gap-2 flex-wrap">
          <Users size={14} className="text-gray-500" />
          <ChipButton active={!activeGroup} onClick={() => onGroupChange(null)} color="blue">
            Все группы
          </ChipButton>
          {filteredGroups.map((g) => (
            <ChipButton
              key={g.name}
              active={activeGroup === g.name}
              onClick={() => onGroupChange(g.name)}
              color="blue"
            >
              {g.name}
            </ChipButton>
          ))}
        </div>
      )}
    </div>
  )
}
