import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Loader2, Shield, MapPin, ChevronRight } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import type { MatchItem } from '@/shared/api/exploreTypes'

function TeamLogo({ url, name }: { url?: string; name: string }) {
  const [hasError, setHasError] = useState(false)

  if (!url || hasError) {
    return <Shield size={20} className="text-gray-600 shrink-0" />
  }

  return (
    <img
      src={url}
      alt={name}
      className="h-5 w-5 object-contain"
      onError={() => setHasError(true)}
    />
  )
}

function ResultTypeBadge({ type }: { type?: string }) {
  if (!type || type === 'regular') return null

  const label = type === 'OT' ? 'ОТ' : type === 'SO' ? 'Б' : null
  if (!label) return null

  return (
    <span className="text-xs font-semibold text-orange-400">
      {label}
    </span>
  )
}

interface Props {
  matches: MatchItem[]
  isLoading: boolean
}

export function MatchesTab({ matches, isLoading }: Props) {
  if (isLoading) {
    return (
      <div className="flex justify-center py-16">
        <Loader2 size={28} className="animate-spin text-gray-500" />
      </div>
    )
  }

  const finished = matches.filter((m) => m.status === 'finished')
  const scheduled = matches.filter((m) => m.status === 'scheduled')

  return (
    <div className="space-y-6">
      {finished.length > 0 && (
        <GlassCard className="p-6" glowColor="cyan">
          <h4 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-4">Завершённые</h4>
          <div className="space-y-2">
            {finished.map((match) => (
              <Link
                key={match.id}
                to={`/explore/matches/${match.id}`}
                className="block rounded-lg bg-white/5 p-3 hover:bg-white/[0.08] transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="text-xs text-gray-500 w-16 flex-shrink-0">
                    <div>{new Date(match.date).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })}</div>
                    <div className="text-gray-600">{match.time}</div>
                  </div>
                  <div className="flex-1 flex items-center justify-center gap-2">
                    <span className="text-sm text-white text-right flex-1 truncate">{match.homeTeam}</span>
                    <TeamLogo url={match.homeLogoUrl} name={match.homeTeam} />
                    <div className="flex flex-col items-center w-16">
                      <span className="text-lg font-bold text-[#00d4ff]">
                        {match.homeScore}:{match.awayScore}
                      </span>
                      <ResultTypeBadge type={match.resultType} />
                    </div>
                    <TeamLogo url={match.awayLogoUrl} name={match.awayTeam} />
                    <span className="text-sm text-white text-left flex-1 truncate">{match.awayTeam}</span>
                  </div>
                  <ChevronRight size={16} className="text-gray-600 flex-shrink-0" />
                </div>
                {match.venue && (
                  <div className="flex items-center gap-1 mt-2 ml-16 text-xs text-gray-500">
                    <MapPin size={12} />
                    <span>{match.venue}</span>
                  </div>
                )}
              </Link>
            ))}
          </div>
        </GlassCard>
      )}

      {scheduled.length > 0 && (
        <GlassCard className="p-6" glowColor="blue">
          <h4 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-4">Предстоящие</h4>
          <div className="space-y-2">
            {scheduled.map((match) => (
              <Link
                key={match.id}
                to={`/explore/matches/${match.id}`}
                className="block rounded-lg bg-white/5 p-3 hover:bg-white/[0.08] transition-colors"
              >
                <div className="flex items-center gap-3">
                  <div className="text-xs text-gray-500 w-16 flex-shrink-0">
                    <div>{new Date(match.date).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })}</div>
                    <div className="text-gray-600">{match.time}</div>
                  </div>
                  <div className="flex-1 flex items-center justify-center gap-2">
                    <span className="text-sm text-white text-right flex-1 truncate">{match.homeTeam}</span>
                    <TeamLogo url={match.homeLogoUrl} name={match.homeTeam} />
                    <span className="text-sm text-gray-600 w-16 text-center">vs</span>
                    <TeamLogo url={match.awayLogoUrl} name={match.awayTeam} />
                    <span className="text-sm text-white text-left flex-1 truncate">{match.awayTeam}</span>
                  </div>
                  <ChevronRight size={16} className="text-gray-600 flex-shrink-0" />
                </div>
                {match.venue && (
                  <div className="flex items-center gap-1 mt-2 ml-16 text-xs text-gray-500">
                    <MapPin size={12} />
                    <span>{match.venue}</span>
                  </div>
                )}
              </Link>
            ))}
          </div>
        </GlassCard>
      )}
    </div>
  )
}
