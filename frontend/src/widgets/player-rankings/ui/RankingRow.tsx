import { memo } from 'react'
import { motion } from 'framer-motion'
import type { PlayerRankingEntry } from '@/entities/player'
import { cn } from '@/shared/lib/utils'

interface RankingRowProps {
  player: PlayerRankingEntry
  index: number
}

const medals: Record<number, string> = {
  1: '🥇',
  2: '🥈',
  3: '🥉',
}

export const RankingRow = memo(function RankingRow({
  player,
  index,
}: RankingRowProps) {
  const medal = medals[player.rank]

  return (
    <motion.tr
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: index * 0.05 }}
      className={cn(
        'group transition-colors',
        player.isCurrentPlayer
          ? 'bg-[#00d4ff]/10'
          : 'hover:bg-white/5'
      )}
    >
      {/* Rank */}
      <td className="whitespace-nowrap py-3 pl-4 pr-2">
        <div
          className={cn(
            'flex h-8 w-8 items-center justify-center rounded-lg text-sm font-bold',
            player.isCurrentPlayer
              ? 'bg-[#00d4ff]/20 text-[#00d4ff]'
              : 'bg-white/5 text-gray-400'
          )}
        >
          {medal || player.rank}
        </div>
      </td>

      {/* Player */}
      <td className="py-3 pr-4">
        <div className="flex items-center gap-3">
          <div
            className={cn(
              'flex h-10 w-10 items-center justify-center rounded-full',
              'bg-gradient-to-br from-[#00d4ff]/30 to-[#8b5cf6]/30',
              'text-sm font-bold text-white',
              player.isCurrentPlayer && 'ring-2 ring-[#00d4ff] ring-offset-2 ring-offset-[#0a0e1a]'
            )}
          >
            {player.name
              .split(' ')
              .map((n: string) => n[0])
              .join('')
              .slice(0, 2)}
          </div>
          <div>
            <div
              className={cn(
                'font-medium',
                player.isCurrentPlayer ? 'text-[#00d4ff]' : 'text-white'
              )}
            >
              {player.name}
              {player.isCurrentPlayer && (
                <span className="ml-2 text-xs text-[#00d4ff]/70">(Вы)</span>
              )}
            </div>
            <div className="text-sm text-gray-500">{player.team}</div>
          </div>
        </div>
      </td>

      {/* Goals */}
      <td className="py-3 text-center">
        <span className="font-bold text-[#00d4ff]">{player.goals}</span>
      </td>

      {/* Assists */}
      <td className="py-3 text-center">
        <span className="font-medium text-[#8b5cf6]">{player.assists}</span>
      </td>

      {/* Games */}
      <td className="py-3 pr-4 text-center">
        <span className="text-gray-400">{player.games}</span>
      </td>
    </motion.tr>
  )
})
