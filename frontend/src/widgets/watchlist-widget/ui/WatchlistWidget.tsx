import { memo } from 'react'
import { Link } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { Star, ChevronRight, Plus, X } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useScoutStore } from '@/shared/stores'

interface WatchlistWidgetProps {
  className?: string
}

// Mock player data - in real app this would come from API based on watchlist IDs
const mockPlayerData: Record<string, { name: string; team: string; goals: number; assists: number }> = {
  '1': { name: 'Иванов Александр', team: 'СКА-Юниор', goals: 45, assists: 38 },
  '2': { name: 'Петров Максим', team: 'Динамо-Юниор', goals: 42, assists: 35 },
  '3': { name: 'Сидоров Дмитрий', team: 'ЦСКА-Юниор', goals: 38, assists: 36 },
  '4': { name: 'Козлов Артём', team: 'Спартак-Юниор', goals: 35, assists: 32 },
  '5': { name: 'Смирнов Никита', team: 'Локомотив-Юниор', goals: 33, assists: 30 },
}

export const WatchlistWidget = memo(function WatchlistWidget({
  className,
}: WatchlistWidgetProps) {
  const watchlist = useScoutStore((state) => state.watchlist)
  const removeFromWatchlist = useScoutStore((state) => state.removeFromWatchlist)

  const displayedPlayers = watchlist.slice(0, 5)
  const hasMore = watchlist.length > 5

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.5 }}
      className={cn('glass-card rounded-xl p-6', className)}
    >
      {/* Header */}
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Star size={20} className="text-[#ec4899]" />
          <h3 className="text-sm font-medium uppercase tracking-wider text-gray-400">
            Избранные
          </h3>
          {watchlist.length > 0 && (
            <span className="rounded-full bg-[#ec4899]/20 px-2 py-0.5 text-xs font-bold text-[#ec4899]">
              {watchlist.length}
            </span>
          )}
        </div>
        {watchlist.length > 0 && (
          <Link
            to="/dashboard/watchlist"
            className="flex items-center gap-1 text-xs text-[#00d4ff] transition-colors hover:text-[#00d4ff]/80"
          >
            Все
            <ChevronRight size={14} />
          </Link>
        )}
      </div>

      {/* Content */}
      {watchlist.length === 0 ? (
        <EmptyState />
      ) : (
        <div className="space-y-2">
          <AnimatePresence mode="popLayout">
            {displayedPlayers.map((playerId, index) => {
              const player = mockPlayerData[playerId]
              if (!player) return null

              return (
                <motion.div
                  key={playerId}
                  layout
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: 20 }}
                  transition={{ delay: index * 0.05 }}
                  className="group flex items-center gap-3 rounded-lg bg-white/5 p-3 transition-colors hover:bg-white/10"
                >
                  {/* Player info */}
                  <div className="min-w-0 flex-1">
                    <p className="truncate font-medium text-white">{player.name}</p>
                    <p className="truncate text-xs text-gray-500">{player.team}</p>
                  </div>

                  {/* Stats */}
                  <div className="flex items-center gap-3 text-xs">
                    <span className="text-gray-400">
                      <span className="font-bold text-white">{player.goals}</span> Г
                    </span>
                    <span className="text-gray-400">
                      <span className="font-bold text-white">{player.assists}</span> П
                    </span>
                  </div>

                  {/* Remove button */}
                  <motion.button
                    onClick={() => removeFromWatchlist(playerId)}
                    className={cn(
                      'rounded-full p-1 opacity-0 transition-opacity',
                      'bg-white/5 text-gray-500 hover:bg-red-500/20 hover:text-red-400',
                      'group-hover:opacity-100'
                    )}
                    whileHover={{ scale: 1.1 }}
                    whileTap={{ scale: 0.9 }}
                    title="Удалить из избранного"
                  >
                    <X size={14} />
                  </motion.button>
                </motion.div>
              )
            })}
          </AnimatePresence>

          {hasMore && (
            <Link
              to="/dashboard/watchlist"
              className="block text-center text-xs text-gray-500 transition-colors hover:text-[#00d4ff]"
            >
              + ещё {watchlist.length - 5} игроков
            </Link>
          )}
        </div>
      )}

      {/* Bottom glow line */}
      <div className="mt-4 h-px bg-gradient-to-r from-transparent via-[#ec4899]/50 to-transparent" />
    </motion.div>
  )
})

const EmptyState = memo(function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center py-8 text-center">
      <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-white/5">
        <Star size={24} className="text-gray-600" />
      </div>
      <p className="text-sm text-gray-400">Пока нет избранных игроков</p>
      <p className="mt-1 text-xs text-gray-500">
        Нажмите на звёздочку рядом с игроком,
        <br />
        чтобы добавить его в избранное
      </p>
      <Link
        to="/dashboard/search"
        className={cn(
          'mt-4 flex items-center gap-2 rounded-lg px-4 py-2',
          'bg-[#00d4ff]/20 text-sm text-[#00d4ff]',
          'transition-colors hover:bg-[#00d4ff]/30'
        )}
      >
        <Plus size={16} />
        Найти игроков
      </Link>
    </div>
  )
})
