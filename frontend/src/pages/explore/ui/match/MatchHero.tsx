import { memo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Shield, MapPin, Calendar, Clock, Trophy } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import type { MatchDetail } from '@/shared/api/exploreTypes'

function TeamLogo({ url, name, size = 64 }: { url?: string; name: string; size?: number }) {
  const [hasError, setHasError] = useState(false)
  if (!url || hasError) {
    return <Shield size={size} className="text-gray-600" />
  }
  return (
    <img
      src={url}
      alt={name}
      className="object-contain"
      style={{ height: size, width: size }}
      onError={() => setHasError(true)}
    />
  )
}

function ResultTypeBadge({ type }: { type?: string }) {
  if (!type || type === 'regular') return null
  const label = type === 'OT' ? 'ОТ' : type === 'SO' ? 'Б' : null
  if (!label) return null
  return (
    <span className="text-sm font-semibold text-orange-400 ml-2">({label})</span>
  )
}

// Calculate shootout winner from final score vs regulation score
function getShootoutScore(match: MatchDetail, isHome: boolean): number | null {
  if (!match.scoreByPeriod || match.homeScore === null || match.awayScore === null) {
    return null
  }

  // Calculate regulation score (sum of all periods including OT)
  const regHome = (match.scoreByPeriod.homeP1 ?? 0) +
                  (match.scoreByPeriod.homeP2 ?? 0) +
                  (match.scoreByPeriod.homeP3 ?? 0) +
                  (match.scoreByPeriod.homeOt ?? 0)
  const regAway = (match.scoreByPeriod.awayP1 ?? 0) +
                  (match.scoreByPeriod.awayP2 ?? 0) +
                  (match.scoreByPeriod.awayP3 ?? 0) +
                  (match.scoreByPeriod.awayOt ?? 0)

  // If final score differs from regulation, we can determine SO winner
  const homeWon = match.homeScore > regHome
  const awayWon = match.awayScore > regAway

  if (homeWon) return isHome ? 1 : 0
  if (awayWon) return isHome ? 0 : 1

  // If scores are equal (data issue), return null
  return null
}

interface Props {
  match: MatchDetail
}

export const MatchHero = memo(function MatchHero({ match }: Props) {
  const isFinished = match.status === 'finished'

  return (
    <GlassCard className="p-6" glowColor="cyan">
      {/* Tournament info */}
      <div className="flex items-center justify-center gap-2 mb-6">
        <Trophy size={16} className="text-[#8b5cf6]" />
        <Link
          to={`/explore/tournaments/detail/${match.tournament.id}`}
          className="text-sm text-[#8b5cf6] hover:underline"
        >
          {match.tournament.name}
        </Link>
        {match.groupName && (
          <span className="text-sm text-gray-500">• {match.groupName}</span>
        )}
        {match.birthYear && (
          <span className="text-sm text-gray-500">• {match.birthYear} г.р.</span>
        )}
      </div>

      {/* Teams and score */}
      <div className="flex items-center justify-center gap-6 md:gap-12">
        {/* Home team */}
        <Link
          to={`/explore/teams/${match.homeTeam.id}`}
          className="flex flex-col items-center gap-2 hover:opacity-80 transition-opacity flex-1"
        >
          <TeamLogo url={match.homeTeam.logoUrl} name={match.homeTeam.name} size={80} />
          <span className="text-lg font-bold text-white text-center">{match.homeTeam.name}</span>
          {match.homeTeam.city && (
            <span className="text-xs text-gray-500">{match.homeTeam.city}</span>
          )}
        </Link>

        {/* Score */}
        <div className="flex flex-col items-center">
          {isFinished ? (
            <>
              <div className="flex items-center gap-3">
                <span className="text-4xl md:text-5xl font-bold text-white">{match.homeScore}</span>
                <span className="text-2xl text-gray-500">:</span>
                <span className="text-4xl md:text-5xl font-bold text-white">{match.awayScore}</span>
              </div>
              <ResultTypeBadge type={match.resultType} />
            </>
          ) : (
            <span className="text-2xl font-semibold text-gray-500">vs</span>
          )}
          <span className="text-xs text-gray-500 mt-2 uppercase">
            {isFinished ? 'Завершён' : 'Запланирован'}
          </span>
        </div>

        {/* Away team */}
        <Link
          to={`/explore/teams/${match.awayTeam.id}`}
          className="flex flex-col items-center gap-2 hover:opacity-80 transition-opacity flex-1"
        >
          <TeamLogo url={match.awayTeam.logoUrl} name={match.awayTeam.name} size={80} />
          <span className="text-lg font-bold text-white text-center">{match.awayTeam.name}</span>
          {match.awayTeam.city && (
            <span className="text-xs text-gray-500">{match.awayTeam.city}</span>
          )}
        </Link>
      </div>

      {/* Score by period */}
      {match.scoreByPeriod && isFinished && (
        <div className="flex justify-center gap-4 mt-6 pt-4 border-t border-white/10">
          {[
            { label: '1 период', home: match.scoreByPeriod.homeP1, away: match.scoreByPeriod.awayP1 },
            { label: '2 период', home: match.scoreByPeriod.homeP2, away: match.scoreByPeriod.awayP2 },
            { label: '3 период', home: match.scoreByPeriod.homeP3, away: match.scoreByPeriod.awayP3 },
            ...(match.resultType === 'OT' || (match.scoreByPeriod.homeOt !== undefined && match.scoreByPeriod.homeOt !== null)
              ? [{ label: 'ОТ', home: match.scoreByPeriod.homeOt ?? 0, away: match.scoreByPeriod.awayOt ?? 0 }]
              : []),
            ...(match.resultType === 'SO'
              ? [{ label: 'Б', home: getShootoutScore(match, true), away: getShootoutScore(match, false), isSO: true }]
              : []),
          ].map((period) => (
            <div key={period.label} className="text-center">
              <div className="text-xs text-gray-500 mb-1">{period.label}</div>
              <div className="text-sm text-white font-medium">
                {'isSO' in period ? (period.home !== null ? `${period.home}:${period.away}` : '—') : `${period.home ?? '-'}:${period.away ?? '-'}`}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Match info */}
      <div className="flex flex-wrap justify-center gap-4 mt-6 pt-4 border-t border-white/10">
        <div className="flex items-center gap-1.5 text-sm text-gray-400">
          <Calendar size={14} />
          <span>{new Date(match.date).toLocaleDateString('ru-RU', {
            day: 'numeric', month: 'long', year: 'numeric'
          })}</span>
        </div>
        {match.time && (
          <div className="flex items-center gap-1.5 text-sm text-gray-400">
            <Clock size={14} />
            <span>{match.time} МСК</span>
          </div>
        )}
        {match.venue && (
          <div className="flex items-center gap-1.5 text-sm text-gray-400">
            <MapPin size={14} />
            <span>{match.venue}</span>
          </div>
        )}
        {match.matchNumber && (
          <div className="text-sm text-gray-400">
            Матч #{match.matchNumber}
          </div>
        )}
      </div>
    </GlassCard>
  )
})
