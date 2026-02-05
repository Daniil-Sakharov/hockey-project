import type { NavigateFunction } from 'react-router-dom'
import { Trophy } from 'lucide-react'
import { cleanTournamentName, formatGroupName } from '@/shared/lib/formatters'
import type { PlayerStatEntry } from '@/shared/api/exploreTypes'

export interface BirthYearGroup {
  birthYear: number
  entries: PlayerStatEntry[]
}

export interface TournamentGroup {
  id: string
  name: string
  birthYearGroups: BirthYearGroup[]
}

function sumEntries(entries: PlayerStatEntry[]) {
  return entries.reduce(
    (acc, e) => ({
      games: acc.games + e.games,
      goals: acc.goals + e.goals,
      assists: acc.assists + e.assists,
      points: acc.points + e.points,
      plusMinus: acc.plusMinus + e.plusMinus,
      penaltyMinutes: acc.penaltyMinutes + e.penaltyMinutes,
    }),
    { games: 0, goals: 0, assists: 0, points: 0, plusMinus: 0, penaltyMinutes: 0 },
  )
}

function BirthYearBadge({ year }: { year: number }) {
  if (!year) return null
  return (
    <span className="text-[10px] px-1.5 py-0.5 rounded bg-white/[0.06] text-gray-400 font-medium">
      {year} г.р.
    </span>
  )
}

function StatsTable({
  entries,
  tournamentId,
  navigate,
}: {
  entries: PlayerStatEntry[]
  tournamentId: string
  navigate: NavigateFunction
}) {
  return (
    <div className="space-y-1">
      <div className="grid grid-cols-[1fr_repeat(6,40px)] text-gray-500 text-xs px-3 pb-1">
        <span className="font-medium">Группа</span>
        <span className="text-center font-medium">И</span>
        <span className="text-center font-medium">Г</span>
        <span className="text-center font-medium">П</span>
        <span className="text-center font-medium">О</span>
        <span className="text-center font-medium">+/-</span>
        <span className="text-center font-medium">Штр</span>
      </div>
      {entries.map((e, i) => (
        <div
          key={i}
          onClick={() => navigate(`/explore/tournaments/detail/${tournamentId}?birthYear=${e.birthYear}&group=${encodeURIComponent(e.groupName)}`)}
          className="grid grid-cols-[1fr_repeat(6,40px)] text-xs px-3 py-1.5 rounded-lg cursor-pointer border border-transparent hover:border-[#8b5cf6]/50 transition-colors"
        >
          <span className="text-gray-400">{formatGroupName(e.groupName)}</span>
          <span className="text-center text-gray-300">{e.games}</span>
          <span className="text-center text-gray-300">{e.goals}</span>
          <span className="text-center text-gray-300">{e.assists}</span>
          <span className="text-center text-gray-300">{e.points}</span>
          <span className="text-center text-gray-300">
            {e.plusMinus > 0 ? `+${e.plusMinus}` : e.plusMinus}
          </span>
          <span className="text-center text-gray-300">{e.penaltyMinutes}</span>
        </div>
      ))}
      {entries.length > 1 && (() => {
        const total = sumEntries(entries)
        return (
          <div className="grid grid-cols-[1fr_repeat(6,40px)] text-xs px-3 py-1.5 rounded-lg bg-white/[0.03]">
            <span className="text-white font-semibold">Итого</span>
            <span className="text-center text-gray-300">{total.games}</span>
            <span className="text-center text-gray-300">{total.goals}</span>
            <span className="text-center text-gray-300">{total.assists}</span>
            <span className="text-center font-medium text-[#00d4ff]">{total.points}</span>
            <span className="text-center text-gray-300">
              {total.plusMinus > 0 ? `+${total.plusMinus}` : total.plusMinus}
            </span>
            <span className="text-center text-gray-300">{total.penaltyMinutes}</span>
          </div>
        )
      })()}
    </div>
  )
}

export function TournamentBlock({
  tg,
  navigate,
}: {
  tg: TournamentGroup
  navigate: NavigateFunction
}) {
  const singleYear = tg.birthYearGroups.length === 1

  return (
    <div>
      <div className="flex items-center gap-2 mb-2">
        <Trophy size={14} className="text-[#00d4ff]" />
        <span className="text-xs font-medium text-gray-300">{cleanTournamentName(tg.name)}</span>
        {singleYear && <BirthYearBadge year={tg.birthYearGroups[0].birthYear} />}
      </div>
      {singleYear ? (
        <StatsTable entries={tg.birthYearGroups[0].entries} tournamentId={tg.id} navigate={navigate} />
      ) : (
        <div className="space-y-3">
          {tg.birthYearGroups.map((byg) => (
            <div key={byg.birthYear}>
              <div className="flex items-center gap-1.5 mb-1 ml-1">
                <BirthYearBadge year={byg.birthYear} />
              </div>
              <StatsTable entries={byg.entries} tournamentId={tg.id} navigate={navigate} />
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
