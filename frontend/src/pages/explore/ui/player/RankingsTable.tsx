import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Medal, UserRound } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import type { RankedPlayer } from '@/shared/api/exploreTypes'

type SortKey = 'points' | 'goals' | 'assists' | 'penaltyMinutes' | 'plusMinus'

const MEDAL_COLORS = ['text-[#f59e0b]', 'text-gray-300', 'text-[#cd7f32]']

export function RankingsTable({ ranked, sortBy }: { ranked: RankedPlayer[]; sortBy: SortKey }) {
  return (
    <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.15 }}>
      <GlassCard className="p-0 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-gray-500 border-b border-white/10 bg-white/[0.02]">
                <th className="text-left py-3 pl-4 pr-2 w-12">#</th>
                <th className="w-10 py-3 px-1" />
                <th className="text-left py-3 px-3">Игрок</th>
                <th className="text-left py-3 px-3 hidden md:table-cell">Команда</th>
                <th className="text-left py-3 px-3 hidden lg:table-cell">Поз.</th>
                <th className="text-center py-3 px-2">И</th>
                <th className={cn('text-center py-3 px-2', sortBy === 'goals' && 'text-[#00d4ff]')}>Г</th>
                <th className={cn('text-center py-3 px-2', sortBy === 'assists' && 'text-[#00d4ff]')}>П</th>
                <th className={cn('text-center py-3 px-2 font-semibold', sortBy === 'points' && 'text-[#00d4ff]')}>О</th>
                <th className={cn('text-center py-3 px-2 hidden md:table-cell', sortBy === 'plusMinus' && 'text-[#00d4ff]')}>+/-</th>
                <th className={cn('text-center py-3 pr-4 px-2 hidden md:table-cell', sortBy === 'penaltyMinutes' && 'text-[#00d4ff]')}>ШМ</th>
              </tr>
            </thead>
            <tbody>
              {ranked.map((player, i) => (
                <motion.tr
                  key={player.id}
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.02 * i }}
                  className="border-b border-white/5 hover:bg-white/[0.03] transition-colors"
                >
                  <td className="py-3 pl-4 pr-2">
                    {i < 3 ? (
                      <Medal size={18} className={MEDAL_COLORS[i]} />
                    ) : (
                      <span className="text-gray-500 font-medium">{player.rank}</span>
                    )}
                  </td>
                  <td className="py-2 px-1">
                    {player.photoUrl ? (
                      <img src={player.photoUrl} alt={player.name} className="h-8 w-8 rounded-full object-cover" />
                    ) : (
                      <div className="flex h-8 w-8 items-center justify-center rounded-full bg-white/10">
                        <UserRound size={14} className="text-gray-500" />
                      </div>
                    )}
                  </td>
                  <td className="py-3 px-3">
                    <Link to={`/explore/players/${player.id}`} className="font-medium text-white hover:text-[#00d4ff] transition-colors">
                      {player.name}
                    </Link>
                  </td>
                  <td className="py-3 px-3 text-gray-400 hidden md:table-cell">
                    <Link to={`/explore/teams/${player.teamId}`} className="hover:text-[#8b5cf6] transition-colors">
                      {player.team}
                    </Link>
                  </td>
                  <td className="py-3 px-3 hidden lg:table-cell">
                    <PositionBadge position={player.position} />
                  </td>
                  <td className="py-3 px-2 text-center text-gray-400">{player.games}</td>
                  <td className={cn('py-3 px-2 text-center', sortBy === 'goals' ? 'text-[#00d4ff] font-semibold' : 'text-gray-300')}>
                    {player.goals}
                  </td>
                  <td className={cn('py-3 px-2 text-center', sortBy === 'assists' ? 'text-[#00d4ff] font-semibold' : 'text-gray-300')}>
                    {player.assists}
                  </td>
                  <td className={cn('py-3 px-2 text-center font-semibold', sortBy === 'points' ? 'text-[#00d4ff]' : 'text-white')}>
                    {player.points}
                  </td>
                  <td className={cn('py-3 px-2 text-center hidden md:table-cell', sortBy === 'plusMinus' ? 'text-[#00d4ff] font-semibold' : player.plusMinus >= 0 ? 'text-[#10b981]' : 'text-[#ef4444]')}>
                    {player.plusMinus > 0 ? `+${player.plusMinus}` : player.plusMinus}
                  </td>
                  <td className={cn('py-3 pr-4 px-2 text-center hidden md:table-cell', sortBy === 'penaltyMinutes' ? 'text-[#00d4ff] font-semibold' : 'text-gray-400')}>
                    {player.penaltyMinutes}
                  </td>
                </motion.tr>
              ))}
            </tbody>
          </table>
        </div>
      </GlassCard>
    </motion.div>
  )
}

function PositionBadge({ position }: { position: string }) {
  const labels: Record<string, string> = { forward: 'Нап', defender: 'Защ', goalie: 'Вр' }
  return (
    <span className={cn(
      'text-[10px] px-1.5 py-0.5 rounded font-medium',
      position === 'forward' ? 'bg-[#00d4ff]/20 text-[#00d4ff]' :
      position === 'defender' ? 'bg-[#8b5cf6]/20 text-[#8b5cf6]' :
      'bg-[#f59e0b]/20 text-[#f59e0b]'
    )}>
      {labels[position] || position}
    </span>
  )
}
