import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Loader2, Shield, User } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import type { Scorer } from '@/shared/api/exploreTypes'

function PlayerPhoto({ url, name }: { url?: string; name: string }) {
  const [hasError, setHasError] = useState(false)

  if (!url || hasError) {
    return (
      <div className="h-8 w-8 rounded-full bg-white/10 flex items-center justify-center shrink-0">
        <User size={18} className="text-gray-500" />
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

function TeamLogo({ url, name }: { url?: string; name: string }) {
  const [hasError, setHasError] = useState(false)

  if (!url || hasError) {
    return <Shield size={24} className="text-gray-600 shrink-0" />
  }

  return (
    <img
      src={url}
      alt={name}
      className="h-6 w-6 object-contain"
      onError={() => setHasError(true)}
    />
  )
}

interface Props {
  scorers: Scorer[]
  isLoading: boolean
}

export function ScorersTab({ scorers, isLoading }: Props) {
  if (isLoading) {
    return (
      <div className="flex justify-center py-16">
        <Loader2 size={28} className="animate-spin text-gray-500" />
      </div>
    )
  }

  return (
    <GlassCard className="p-6" glowColor="purple">
      <div className="overflow-x-auto">
        <table className="w-full text-base">
          <thead>
            <tr className="text-gray-500 border-b border-white/10">
              <th className="text-left py-4 pr-3 w-10">#</th>
              <th className="text-left py-4 pr-6">Игрок</th>
              <th className="text-left py-4 pr-6">Команда</th>
              <th className="text-center py-4 px-3">Г</th>
              <th className="text-center py-4 px-3">П</th>
              <th className="text-center py-4 px-3 font-semibold text-[#8b5cf6]">О</th>
            </tr>
          </thead>
          <tbody>
            {scorers.map((scorer) => (
              <tr key={scorer.playerId} className="border-b border-white/5 last:border-0">
                <td className="py-4 pr-3">
                  <span
                    className={cn(
                      'inline-flex h-8 w-8 items-center justify-center rounded-full text-sm font-bold',
                      scorer.position === 1 && 'bg-[#f59e0b]/20 text-[#f59e0b]',
                      scorer.position === 2 && 'bg-gray-400/20 text-gray-300',
                      scorer.position === 3 && 'bg-[#cd7f32]/20 text-[#cd7f32]',
                      scorer.position > 3 && 'text-gray-500',
                    )}
                  >
                    {scorer.position}
                  </span>
                </td>
                <td className="py-4 pr-6 font-medium text-white">
                  <Link to={`/explore/players/${scorer.playerId}`} className="flex items-center gap-3 hover:text-[#00d4ff] transition-colors">
                    <PlayerPhoto url={scorer.photoUrl} name={scorer.name} />
                    {scorer.name}
                  </Link>
                </td>
                <td className="py-4 pr-6 text-gray-400">
                  <Link to={`/explore/teams/${scorer.teamId}`} className="flex items-center gap-2 hover:text-[#00d4ff] transition-colors">
                    <TeamLogo url={scorer.logoUrl} name={scorer.team} />
                    {scorer.team}
                  </Link>
                </td>
                <td className="py-4 px-3 text-center text-gray-300">{scorer.goals}</td>
                <td className="py-4 px-3 text-center text-gray-300">{scorer.assists}</td>
                <td className="py-4 px-3 text-center font-bold text-[#8b5cf6]">{scorer.points}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </GlassCard>
  )
}
