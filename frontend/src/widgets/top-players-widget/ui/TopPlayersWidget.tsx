import { memo, useState } from 'react'
import { Link } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { Star, ChevronRight, Trophy } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useScoutStore } from '@/shared/stores'

interface Player {
  id: string
  rank: number
  name: string
  team: string
  goals: number
  assists: number
  points: number
}

interface TopPlayersWidgetProps {
  className?: string
}

// Mock data
const mockPlayers: Player[] = [
  { id: '1', rank: 1, name: 'Иванов Александр', team: 'СКА-Юниор', goals: 45, assists: 38, points: 83 },
  { id: '2', rank: 2, name: 'Петров Максим', team: 'Динамо-Юниор', goals: 42, assists: 35, points: 77 },
  { id: '3', rank: 3, name: 'Сидоров Дмитрий', team: 'ЦСКА-Юниор', goals: 38, assists: 36, points: 74 },
  { id: '4', rank: 4, name: 'Козлов Артём', team: 'Спартак-Юниор', goals: 35, assists: 32, points: 67 },
  { id: '5', rank: 5, name: 'Смирнов Никита', team: 'Локомотив-Юниор', goals: 33, assists: 30, points: 63 },
  { id: '6', rank: 6, name: 'Волков Егор', team: 'Авангард-Юниор', goals: 31, assists: 28, points: 59 },
  { id: '7', rank: 7, name: 'Новиков Павел', team: 'Металлург-Юниор', goals: 29, assists: 27, points: 56 },
]

export const TopPlayersWidget = memo(function TopPlayersWidget({
  className,
}: TopPlayersWidgetProps) {
  const [hoveredRow, setHoveredRow] = useState<string | null>(null)

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.4 }}
      className={cn('glass-card rounded-xl p-6', className)}
    >
      {/* Header */}
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Trophy size={20} className="text-[#fbbf24]" />
          <h3 className="text-sm font-medium uppercase tracking-wider text-gray-400">
            Топ бомбардиров
          </h3>
        </div>
        <Link
          to="/dashboard/search?sort=points"
          className="flex items-center gap-1 text-xs text-[#00d4ff] transition-colors hover:text-[#00d4ff]/80"
        >
          Все игроки
          <ChevronRight size={14} />
        </Link>
      </div>

      {/* Table */}
      <div className="overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-white/5 text-xs text-gray-500">
              <th className="pb-3 text-left font-medium">#</th>
              <th className="pb-3 text-left font-medium">Игрок</th>
              <th className="pb-3 text-center font-medium">Г</th>
              <th className="pb-3 text-center font-medium">П</th>
              <th className="pb-3 text-center font-medium">О</th>
              <th className="pb-3 text-right font-medium"></th>
            </tr>
          </thead>
          <tbody>
            <AnimatePresence>
              {mockPlayers.map((player, index) => (
                <PlayerRow
                  key={player.id}
                  player={player}
                  index={index}
                  isHovered={hoveredRow === player.id}
                  onHover={() => setHoveredRow(player.id)}
                  onLeave={() => setHoveredRow(null)}
                />
              ))}
            </AnimatePresence>
          </tbody>
        </table>
      </div>

      {/* Bottom glow line */}
      <div className="mt-4 h-px bg-gradient-to-r from-transparent via-[#fbbf24]/50 to-transparent" />
    </motion.div>
  )
})

const PlayerRow = memo(function PlayerRow({
  player,
  index,
  isHovered,
  onHover,
  onLeave,
}: {
  player: Player
  index: number
  isHovered: boolean
  onHover: () => void
  onLeave: () => void
}) {
  const isInWatchlist = useScoutStore((state) => state.isInWatchlist(player.id))
  const addToWatchlist = useScoutStore((state) => state.addToWatchlist)
  const removeFromWatchlist = useScoutStore((state) => state.removeFromWatchlist)

  const handleWatchlistToggle = (e: React.MouseEvent) => {
    e.stopPropagation()
    if (isInWatchlist) {
      removeFromWatchlist(player.id)
    } else {
      addToWatchlist(player.id)
    }
  }

  const getRankBadgeColor = (rank: number) => {
    switch (rank) {
      case 1:
        return 'bg-[#fbbf24]/20 text-[#fbbf24] border-[#fbbf24]/30'
      case 2:
        return 'bg-[#94a3b8]/20 text-[#94a3b8] border-[#94a3b8]/30'
      case 3:
        return 'bg-[#cd7c32]/20 text-[#cd7c32] border-[#cd7c32]/30'
      default:
        return 'bg-white/5 text-gray-400 border-white/10'
    }
  }

  return (
    <motion.tr
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: index * 0.05 }}
      className={cn(
        'cursor-pointer border-b border-white/5 transition-colors',
        isHovered && 'bg-white/5'
      )}
      onMouseEnter={onHover}
      onMouseLeave={onLeave}
    >
      {/* Rank */}
      <td className="py-3">
        <span
          className={cn(
            'inline-flex h-6 w-6 items-center justify-center rounded border text-xs font-bold',
            getRankBadgeColor(player.rank)
          )}
        >
          {player.rank}
        </span>
      </td>

      {/* Player info */}
      <td className="py-3">
        <div>
          <p className="font-medium text-white">{player.name}</p>
          <p className="text-xs text-gray-500">{player.team}</p>
        </div>
      </td>

      {/* Goals */}
      <td className="py-3 text-center">
        <span className="text-sm font-medium text-white">{player.goals}</span>
      </td>

      {/* Assists */}
      <td className="py-3 text-center">
        <span className="text-sm font-medium text-white">{player.assists}</span>
      </td>

      {/* Points */}
      <td className="py-3 text-center">
        <span className="text-sm font-bold text-[#00d4ff]">{player.points}</span>
      </td>

      {/* Watchlist button */}
      <td className="py-3 text-right">
        <motion.button
          onClick={handleWatchlistToggle}
          className={cn(
            'rounded-full p-1.5 transition-colors',
            isInWatchlist
              ? 'bg-[#ec4899]/20 text-[#ec4899]'
              : 'bg-white/5 text-gray-500 hover:bg-white/10 hover:text-[#ec4899]'
          )}
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.9 }}
          title={isInWatchlist ? 'Удалить из избранного' : 'Добавить в избранное'}
        >
          <Star size={14} fill={isInWatchlist ? 'currentColor' : 'none'} />
        </motion.button>
      </td>
    </motion.tr>
  )
})
