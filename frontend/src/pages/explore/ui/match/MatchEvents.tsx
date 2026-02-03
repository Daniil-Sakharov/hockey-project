import { memo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Target, Clock, User, Shield } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import type { MatchEvent, MatchDetail, MatchTeam } from '@/shared/api/exploreTypes'
import { cn } from '@/shared/lib/utils'

function PlayerPhoto({ url, name }: { url?: string; name: string }) {
  const [hasError, setHasError] = useState(false)
  if (!url || hasError) {
    return (
      <div className="flex items-center justify-center w-8 h-8 rounded-full bg-white/10">
        <User size={16} className="text-gray-600" />
      </div>
    )
  }
  return (
    <img
      src={url}
      alt={name}
      className="h-8 w-8 rounded-full object-cover"
      onError={() => setHasError(true)}
    />
  )
}

function TeamLogo({ url, name, size = 20 }: { url?: string; name: string; size?: number }) {
  const [hasError, setHasError] = useState(false)
  if (!url || hasError) {
    return <Shield size={size} className="text-gray-600" />
  }
  return (
    <img
      src={url}
      alt={name}
      className="object-contain"
      style={{ width: size, height: size }}
      onError={() => setHasError(true)}
    />
  )
}

interface GoalEventProps {
  event: MatchEvent
  team: MatchTeam
  isHomeTeam: boolean
}

function GoalEvent({ event, team, isHomeTeam }: GoalEventProps) {
  return (
    <div className="flex items-center gap-3 p-3 rounded-lg bg-white/5">
      {/* Time */}
      <div className="w-14 text-center">
        <span className="text-sm font-mono text-[#00d4ff]">{event.time}</span>
      </div>

      {/* Team indicator */}
      <div className="flex items-center gap-2 w-28 flex-shrink-0">
        <TeamLogo url={team.logoUrl} name={team.name} size={24} />
        <span className={cn(
          'text-xs font-medium truncate',
          isHomeTeam ? 'text-[#00d4ff]' : 'text-[#ec4899]'
        )}>
          {team.name}
        </span>
      </div>

      {/* Goal icon */}
      <div className="flex items-center justify-center w-8 h-8 rounded-full bg-[#10b981]/20 flex-shrink-0">
        <Target size={16} className="text-[#10b981]" />
      </div>

      {/* Player info */}
      <div className="flex items-center gap-2 flex-1 min-w-0">
        <PlayerPhoto url={event.playerPhoto} name={event.playerName || ''} />
        <div className="flex flex-col min-w-0">
          {event.playerId ? (
            <Link
              to={`/explore/players/${event.playerId}`}
              className="text-sm font-medium text-white hover:text-[#00d4ff] truncate"
            >
              {event.playerName}
            </Link>
          ) : (
            <span className="text-sm font-medium text-white truncate">{event.playerName || 'Гол'}</span>
          )}
          {(event.assist1Name || event.assist2Name) && (
            <span className="text-xs text-gray-500 truncate">
              Пас: {[event.assist1Name, event.assist2Name].filter(Boolean).join(', ')}
            </span>
          )}
        </div>
      </div>
    </div>
  )
}

interface PenaltyEventProps {
  event: MatchEvent
  team: MatchTeam
  isHomeTeam: boolean
}

function PenaltyEvent({ event, team, isHomeTeam }: PenaltyEventProps) {
  return (
    <div className="flex items-center gap-3 p-3 rounded-lg bg-white/5">
      {/* Time */}
      <div className="w-14 text-center">
        <span className="text-sm font-mono text-gray-400">{event.time}</span>
      </div>

      {/* Team indicator */}
      <div className="flex items-center gap-2 w-28 flex-shrink-0">
        <TeamLogo url={team.logoUrl} name={team.name} size={24} />
        <span className={cn(
          'text-xs font-medium truncate',
          isHomeTeam ? 'text-[#00d4ff]' : 'text-[#ec4899]'
        )}>
          {team.name}
        </span>
      </div>

      {/* Penalty icon */}
      <div className="flex items-center justify-center w-8 h-8 rounded-full bg-red-500/20 flex-shrink-0">
        <Clock size={16} className="text-red-400" />
      </div>

      {/* Player info */}
      <div className="flex items-center gap-2 flex-1 min-w-0">
        <PlayerPhoto url={event.playerPhoto} name={event.playerName || ''} />
        <div className="flex flex-col min-w-0">
          {event.playerName ? (
            event.playerId ? (
              <Link
                to={`/explore/players/${event.playerId}`}
                className="text-sm font-medium text-white hover:text-[#00d4ff] truncate"
              >
                {event.playerName}
              </Link>
            ) : (
              <span className="text-sm font-medium text-white truncate">{event.playerName}</span>
            )
          ) : (
            <span className="text-sm font-medium text-gray-400">—</span>
          )}
          {event.penaltyText && (
            <span className="text-xs text-gray-500 truncate">{event.penaltyText}</span>
          )}
        </div>
      </div>

      {/* Duration */}
      <div className="flex-shrink-0">
        <span className="text-sm font-bold text-red-400">{event.penaltyMins} мин</span>
      </div>
    </div>
  )
}

interface Props {
  match: MatchDetail
}

export const MatchEvents = memo(function MatchEvents({ match }: Props) {
  // Build player ID to team mapping from lineups
  const homePlayerIds = new Set(match.homeLineup.map(p => p.playerId))
  const awayPlayerIds = new Set(match.awayLineup.map(p => p.playerId))

  // Determine which team each event belongs to
  const getEventTeam = (event: MatchEvent): { team: MatchTeam; isHome: boolean } => {
    // For goals, goalType tells us home/away
    if (event.type === 'goal' && event.goalType) {
      const isHome = event.goalType === 'home'
      return { team: isHome ? match.homeTeam : match.awayTeam, isHome }
    }

    // For penalties, try to determine from player ID in lineup
    if (event.type === 'penalty' && event.playerId) {
      if (homePlayerIds.has(event.playerId)) {
        return { team: match.homeTeam, isHome: true }
      }
      if (awayPlayerIds.has(event.playerId)) {
        return { team: match.awayTeam, isHome: false }
      }
    }

    // Fallback to isHome field
    return { team: event.isHome ? match.homeTeam : match.awayTeam, isHome: event.isHome }
  }

  const goals = match.events.filter(e => e.type === 'goal')
  const penalties = match.events.filter(e => e.type === 'penalty')

  if (goals.length === 0 && penalties.length === 0) {
    return null
  }

  const periods = [1, 2, 3, 4]
  const periodLabels = ['1 период', '2 период', '3 период', 'Овертайм']

  return (
    <div className="space-y-4">
      {/* Goals */}
      {goals.length > 0 && (
        <GlassCard className="p-6" glowColor="cyan">
          <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Target size={20} className="text-[#10b981]" />
            Голы ({goals.length})
          </h3>
          <div className="space-y-4">
            {periods.map(period => {
              const periodGoals = goals.filter(g => g.period === period)
              if (periodGoals.length === 0) return null
              return (
                <div key={period}>
                  <div className="text-xs text-gray-500 uppercase tracking-wider mb-2 pb-1 border-b border-white/10">
                    {periodLabels[period - 1]}
                  </div>
                  <div className="space-y-2">
                    {periodGoals.map((goal, i) => {
                      const { team, isHome } = getEventTeam(goal)
                      return <GoalEvent key={i} event={goal} team={team} isHomeTeam={isHome} />
                    })}
                  </div>
                </div>
              )
            })}
          </div>
        </GlassCard>
      )}

      {/* Penalties */}
      {penalties.length > 0 && (
        <GlassCard className="p-6">
          <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Clock size={20} className="text-red-400" />
            Штрафы ({penalties.length})
          </h3>
          <div className="space-y-4">
            {periods.map(period => {
              const periodPenalties = penalties.filter(p => p.period === period)
              if (periodPenalties.length === 0) return null
              return (
                <div key={period}>
                  <div className="text-xs text-gray-500 uppercase tracking-wider mb-2 pb-1 border-b border-white/10">
                    {periodLabels[period - 1]}
                  </div>
                  <div className="space-y-2">
                    {periodPenalties.map((penalty, i) => {
                      const { team, isHome } = getEventTeam(penalty)
                      return <PenaltyEvent key={i} event={penalty} team={team} isHomeTeam={isHome} />
                    })}
                  </div>
                </div>
              )
            })}
          </div>
        </GlassCard>
      )}
    </div>
  )
})
