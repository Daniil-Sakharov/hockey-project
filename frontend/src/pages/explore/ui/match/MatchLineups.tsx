import { memo, useState } from 'react'
import { Link } from 'react-router-dom'
import { User, Shield } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import type { LineupPlayer, MatchTeam } from '@/shared/api/exploreTypes'
import { cn } from '@/shared/lib/utils'

function PlayerPhoto({ url, name }: { url?: string; name: string }) {
  const [hasError, setHasError] = useState(false)
  if (!url || hasError) {
    return (
      <div className="w-10 h-10 rounded-full bg-white/10 flex items-center justify-center">
        <User size={20} className="text-gray-600" />
      </div>
    )
  }
  return (
    <img
      src={url}
      alt={name}
      className="w-10 h-10 rounded-full object-cover"
      onError={() => setHasError(true)}
    />
  )
}

function TeamLogo({ url, name }: { url?: string; name: string }) {
  const [hasError, setHasError] = useState(false)
  if (!url || hasError) {
    return <Shield size={20} className="text-gray-600" />
  }
  return <img src={url} alt={name} className="h-6 w-6 object-contain" onError={() => setHasError(true)} />
}

function PositionBadge({ position }: { position?: string }) {
  const labels: Record<string, { text: string; color: string }> = {
    forward: { text: 'Нап', color: 'bg-[#00d4ff]/20 text-[#00d4ff]' },
    defender: { text: 'Защ', color: 'bg-[#8b5cf6]/20 text-[#8b5cf6]' },
    goalie: { text: 'Вр', color: 'bg-[#f59e0b]/20 text-[#f59e0b]' },
  }
  const badge = labels[position || '']
  if (!badge) return null
  return (
    <span className={cn('text-[10px] px-1.5 py-0.5 rounded font-medium', badge.color)}>
      {badge.text}
    </span>
  )
}

function LineupTable({ players, isGoalie }: { players: LineupPlayer[]; isGoalie?: boolean }) {
  if (players.length === 0) return null

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="text-xs text-gray-500 border-b border-white/10">
            <th className="text-left py-2 pl-2">#</th>
            <th className="text-left py-2">Игрок</th>
            <th className="text-center py-2">Г</th>
            <th className="text-center py-2">П</th>
            <th className="text-center py-2">О</th>
            <th className="text-center py-2">+/-</th>
            <th className="text-center py-2 pr-2">Штр</th>
            {isGoalie && (
              <>
                <th className="text-center py-2">Сэйвы</th>
                <th className="text-center py-2 pr-2">Пропущ</th>
              </>
            )}
          </tr>
        </thead>
        <tbody>
          {players.map((player) => (
            <tr key={player.playerId} className="border-b border-white/5 hover:bg-white/5">
              <td className="py-2 pl-2 text-gray-500">
                {player.jerseyNumber || '-'}
              </td>
              <td className="py-2">
                <Link
                  to={`/explore/players/${player.playerId}`}
                  className="flex items-center gap-2 hover:text-[#00d4ff]"
                >
                  <PlayerPhoto url={player.playerPhoto} name={player.playerName} />
                  <div className="flex flex-col">
                    <span className="text-white font-medium">{player.playerName}</span>
                    <PositionBadge position={player.position} />
                  </div>
                </Link>
              </td>
              <td className="text-center py-2 text-[#00d4ff] font-bold">{player.goals}</td>
              <td className="text-center py-2 text-gray-400">{player.assists}</td>
              <td className="text-center py-2 text-white font-semibold">{player.points}</td>
              <td className={cn(
                'text-center py-2 font-medium',
                player.plusMinus > 0 ? 'text-green-400' : player.plusMinus < 0 ? 'text-red-400' : 'text-gray-500'
              )}>
                {player.plusMinus > 0 ? '+' : ''}{player.plusMinus}
              </td>
              <td className="text-center py-2 pr-2 text-gray-500">{player.penaltyMinutes}</td>
              {isGoalie && (
                <>
                  <td className="text-center py-2 text-gray-400">{player.saves ?? '-'}</td>
                  <td className="text-center py-2 pr-2 text-gray-400">{player.goalsAgainst ?? '-'}</td>
                </>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

interface Props {
  homeLineup: LineupPlayer[]
  awayLineup: LineupPlayer[]
  homeTeam: MatchTeam
  awayTeam: MatchTeam
}

export const MatchLineups = memo(function MatchLineups({ homeLineup, awayLineup, homeTeam, awayTeam }: Props) {
  if (homeLineup.length === 0 && awayLineup.length === 0) {
    return null
  }

  const groupByPosition = (players: LineupPlayer[]) => {
    const goalies = players.filter(p => p.position === 'goalie')
    const skaters = players.filter(p => p.position !== 'goalie')
    return { goalies, skaters }
  }

  const homeGroups = groupByPosition(homeLineup)
  const awayGroups = groupByPosition(awayLineup)

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
      {/* Home team */}
      <GlassCard className="p-6" glowColor="cyan">
        <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <TeamLogo url={homeTeam.logoUrl} name={homeTeam.name} />
          {homeTeam.name}
          <span className="text-sm text-gray-500 font-normal">({homeLineup.length})</span>
        </h3>
        {homeGroups.goalies.length > 0 && (
          <div className="mb-4">
            <div className="text-xs text-gray-500 uppercase mb-2">Вратари</div>
            <LineupTable players={homeGroups.goalies} isGoalie />
          </div>
        )}
        {homeGroups.skaters.length > 0 && (
          <div>
            <div className="text-xs text-gray-500 uppercase mb-2">Полевые игроки</div>
            <LineupTable players={homeGroups.skaters} />
          </div>
        )}
      </GlassCard>

      {/* Away team */}
      <GlassCard className="p-6" glowColor="purple">
        <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <TeamLogo url={awayTeam.logoUrl} name={awayTeam.name} />
          {awayTeam.name}
          <span className="text-sm text-gray-500 font-normal">({awayLineup.length})</span>
        </h3>
        {awayGroups.goalies.length > 0 && (
          <div className="mb-4">
            <div className="text-xs text-gray-500 uppercase mb-2">Вратари</div>
            <LineupTable players={awayGroups.goalies} isGoalie />
          </div>
        )}
        {awayGroups.skaters.length > 0 && (
          <div>
            <div className="text-xs text-gray-500 uppercase mb-2">Полевые игроки</div>
            <LineupTable players={awayGroups.skaters} />
          </div>
        )}
      </GlassCard>
    </div>
  )
})
