import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Trophy, ArrowRight, Users } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import type { TournamentItem } from '@/shared/api/exploreTypes'
import { cleanTournamentName } from '@/shared/lib/formatters'

interface Props {
  tournament: TournamentItem
  groupName: string
  birthYear: number
  teamsCount: number
  matchesCount: number
  index: number
  region?: string
}

export function TournamentGroupCard({ tournament, groupName, birthYear, teamsCount, matchesCount, index, region }: Props) {
  const params = new URLSearchParams()
  params.set('birthYear', String(birthYear))
  if (groupName) params.set('group', groupName)
  if (region) params.set('from', region)

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.15 + index * 0.03 }}
    >
      <Link to={`/explore/tournaments/detail/${tournament.id}?${params.toString()}`}>
        <GlassCard className="p-5 h-full hover:border-[#00d4ff]/20 transition-colors cursor-pointer">
          <div className="flex items-start gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-[#00d4ff]/20 to-[#8b5cf6]/20 flex-shrink-0">
              {groupName ? <Users size={18} className="text-[#8b5cf6]" /> : <Trophy size={18} className="text-[#00d4ff]" />}
            </div>
            <div className="flex-1 min-w-0">
              <h3 className="text-sm font-semibold text-white leading-tight">{cleanTournamentName(tournament.name)}</h3>
              {groupName && <p className="text-xs text-[#8b5cf6] mt-1">{groupName}</p>}
            </div>
          </div>

          <div className="mt-4 flex items-center justify-between">
            <div className="flex gap-3 text-xs text-gray-400">
              <span>{teamsCount} команд</span>
              <span>{matchesCount} матчей</span>
            </div>
            <span
              className={cn(
                'text-xs px-2 py-0.5 rounded-full',
                !tournament.isEnded ? 'bg-green-500/20 text-green-400' : 'bg-gray-500/20 text-gray-400',
              )}
            >
              {!tournament.isEnded ? 'Активный' : 'Завершён'}
            </span>
          </div>

          <div className="mt-3 flex items-center justify-between">
            <span className="text-xs text-gray-600">Сезон {tournament.season}</span>
            <ArrowRight size={14} className="text-gray-600" />
          </div>
        </GlassCard>
      </Link>
    </motion.div>
  )
}
