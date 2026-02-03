import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Loader2, Shield } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import type { Standing } from '@/shared/api/exploreTypes'

function TeamLogo({ url, name }: { url?: string; name: string }) {
  const [hasError, setHasError] = useState(false)

  if (!url || hasError) {
    return <Shield size={32} className="text-gray-600 shrink-0" />
  }

  return (
    <img
      src={url}
      alt={name}
      className="h-8 w-8 object-contain"
      onError={() => setHasError(true)}
    />
  )
}

interface Props {
  standings: Standing[]
  isLoading: boolean
}

export function StandingsTab({ standings, isLoading }: Props) {
  if (isLoading) {
    return (
      <div className="flex justify-center py-16">
        <Loader2 size={28} className="animate-spin text-gray-500" />
      </div>
    )
  }

  return (
    <GlassCard className="p-6" glowColor="cyan">
      <div className="overflow-x-auto">
        <table className="w-full text-base">
          <thead>
            <tr className="text-gray-500 border-b border-white/10">
              <th className="text-left py-4 pr-3 w-10">#</th>
              <th className="text-left py-4 pr-6">Команда</th>
              <th className="text-center py-4 px-3">И</th>
              <th className="text-center py-4 px-3">В</th>
              <th className="text-center py-4 px-3">Н</th>
              <th className="text-center py-4 px-3">П</th>
              <th className="text-center py-4 px-3">ШЗ</th>
              <th className="text-center py-4 px-3">ШП</th>
              <th className="text-center py-4 px-3 font-semibold text-[#00d4ff]">О</th>
            </tr>
          </thead>
          <tbody>
            {standings.map((row) => (
              <tr
                key={`${row.teamId}-${row.position}`}
                className={cn(
                  'border-b border-white/5 last:border-0',
                  row.position <= 3 && 'bg-white/[0.02]',
                )}
              >
                <td className="py-4 pr-3">
                  <span
                    className={cn(
                      'inline-flex h-8 w-8 items-center justify-center rounded-full text-sm font-bold',
                      row.position === 1 && 'bg-[#f59e0b]/20 text-[#f59e0b]',
                      row.position === 2 && 'bg-gray-400/20 text-gray-300',
                      row.position === 3 && 'bg-[#cd7f32]/20 text-[#cd7f32]',
                      row.position > 3 && 'text-gray-500',
                    )}
                  >
                    {row.position}
                  </span>
                </td>
                <td className="py-4 pr-6 font-medium text-white">
                  <Link to={`/explore/teams/${row.teamId}`} className="flex items-center gap-3 hover:text-[#00d4ff] transition-colors">
                    <TeamLogo url={row.logoUrl} name={row.team} />
                    {row.team}
                  </Link>
                </td>
                <td className="py-4 px-3 text-center text-gray-400">{row.games}</td>
                <td className="py-4 px-3 text-center text-green-400">{row.wins}</td>
                <td className="py-4 px-3 text-center text-gray-400">{row.draws}</td>
                <td className="py-4 px-3 text-center text-red-400">{row.losses}</td>
                <td className="py-4 px-3 text-center text-gray-300">{row.goalsFor}</td>
                <td className="py-4 px-3 text-center text-gray-300">{row.goalsAgainst}</td>
                <td className="py-4 px-3 text-center font-bold text-[#00d4ff]">{row.points}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </GlassCard>
  )
}
