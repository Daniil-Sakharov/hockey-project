import { memo } from 'react'
import { Link } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { Star, ArrowLeft, Trash2 } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useScoutStore } from '@/shared/stores'
import { PlayerCard } from '@/shared/ui'
import { DashboardLayout } from '../../ui/DashboardLayout'

// Mock player data - in real app this would come from API based on watchlist IDs
const mockPlayerData: Record<
  string,
  {
    name: string
    team: string
    position: 'forward' | 'defender' | 'goalie'
    goals: number
    assists: number
    points: number
  }
> = {
  '1': { name: 'Иванов Александр', team: 'СКА-Юниор', position: 'forward', goals: 45, assists: 38, points: 83 },
  '2': { name: 'Петров Максим', team: 'Динамо-Юниор', position: 'forward', goals: 42, assists: 35, points: 77 },
  '3': { name: 'Сидоров Дмитрий', team: 'ЦСКА-Юниор', position: 'defender', goals: 38, assists: 36, points: 74 },
  '4': { name: 'Козлов Артём', team: 'Спартак-Юниор', position: 'forward', goals: 35, assists: 32, points: 67 },
  '5': { name: 'Смирнов Никита', team: 'Локомотив-Юниор', position: 'defender', goals: 33, assists: 30, points: 63 },
  '6': { name: 'Волков Егор', team: 'Авангард-Юниор', position: 'goalie', goals: 0, assists: 2, points: 2 },
  '7': { name: 'Новиков Павел', team: 'Металлург-Юниор', position: 'forward', goals: 29, assists: 27, points: 56 },
}

export function WatchlistPage() {
  const watchlist = useScoutStore((state) => state.watchlist)
  const clearWatchlist = useScoutStore((state) => state.clearWatchlist)

  return (
    <DashboardLayout>
      {/* Page header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="mb-8"
      >
        <div className="flex items-center gap-4">
          <Link
            to="/dashboard"
            className="flex h-10 w-10 items-center justify-center rounded-lg bg-white/5 text-gray-400 transition-colors hover:bg-white/10 hover:text-white"
          >
            <ArrowLeft size={20} />
          </Link>
          <div>
            <div className="flex items-center gap-3">
              <h1 className="text-gradient text-3xl font-bold">Избранные игроки</h1>
              {watchlist.length > 0 && (
                <span className="rounded-full bg-[#ec4899]/20 px-3 py-1 text-sm font-bold text-[#ec4899]">
                  {watchlist.length}
                </span>
              )}
            </div>
            <p className="mt-1 text-gray-400">
              Отслеживайте прогресс интересных вам игроков
            </p>
          </div>
        </div>
      </motion.div>

      {/* Actions */}
      {watchlist.length > 0 && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.1 }}
          className="mb-6 flex justify-end"
        >
          <button
            onClick={clearWatchlist}
            className={cn(
              'flex items-center gap-2 rounded-lg px-4 py-2',
              'bg-red-500/10 text-sm text-red-400',
              'transition-colors hover:bg-red-500/20'
            )}
          >
            <Trash2 size={16} />
            Очистить список
          </button>
        </motion.div>
      )}

      {/* Content */}
      {watchlist.length === 0 ? (
        <EmptyState />
      ) : (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2 }}
          className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3"
        >
          <AnimatePresence mode="popLayout">
            {watchlist.map((playerId, index) => {
              const player = mockPlayerData[playerId]
              if (!player) return null

              return (
                <motion.div
                  key={playerId}
                  layout
                  initial={{ opacity: 0, scale: 0.9 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.9 }}
                  transition={{ delay: index * 0.05 }}
                >
                  <PlayerCard
                    id={playerId}
                    name={player.name}
                    team={player.team}
                    position={player.position}
                    goals={player.goals}
                    assists={player.assists}
                    points={player.points}
                  />
                </motion.div>
              )
            })}
          </AnimatePresence>
        </motion.div>
      )}
    </DashboardLayout>
  )
}

const EmptyState = memo(function EmptyState() {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="glass-card flex flex-col items-center justify-center rounded-xl py-16 text-center"
    >
      <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-white/5">
        <Star size={32} className="text-gray-600" />
      </div>
      <h3 className="text-xl font-semibold text-white">Список пуст</h3>
      <p className="mt-2 max-w-md text-gray-400">
        Вы ещё не добавили ни одного игрока в избранное. Перейдите на страницу
        поиска или на главную, чтобы найти интересных игроков.
      </p>
      <div className="mt-6 flex gap-4">
        <Link
          to="/dashboard"
          className={cn(
            'flex items-center gap-2 rounded-lg px-6 py-3',
            'bg-[#00d4ff]/20 text-[#00d4ff]',
            'transition-colors hover:bg-[#00d4ff]/30'
          )}
        >
          На главную
        </Link>
        <Link
          to="/dashboard/search"
          className={cn(
            'flex items-center gap-2 rounded-lg px-6 py-3',
            'bg-white/5 text-white',
            'transition-colors hover:bg-white/10'
          )}
        >
          Поиск игроков
        </Link>
      </div>
    </motion.div>
  )
})
