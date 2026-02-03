import { memo } from 'react'
import { motion } from 'framer-motion'
import type { PlayerRankingEntry } from '@/entities/player'
import { Skeleton } from '@/shared/ui'
import { RankingRow } from './RankingRow'

interface PlayerRankingsTableProps {
  rankings: PlayerRankingEntry[] | null
  isLoading?: boolean
}

export const PlayerRankingsTable = memo(function PlayerRankingsTable({
  rankings,
  isLoading = false,
}: PlayerRankingsTableProps) {
  if (isLoading) {
    return (
      <div className="glass-card overflow-hidden rounded-xl">
        <div className="border-b border-white/5 p-4">
          <Skeleton className="h-5 w-32" />
        </div>
        <div className="p-4">
          {[...Array(5)].map((_, i) => (
            <Skeleton key={i} className="mb-3 h-14 w-full" />
          ))}
        </div>
      </div>
    )
  }

  if (!rankings || rankings.length === 0) {
    return (
      <div className="glass-card rounded-xl p-6">
        <h3 className="mb-4 text-lg font-semibold text-white">
          Рейтинг бомбардиров
        </h3>
        <p className="text-center text-gray-500">Нет данных</p>
      </div>
    )
  }

  return (
    <motion.section
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.4 }}
    >
      <h2 className="mb-4 text-lg font-semibold text-white">
        Рейтинг бомбардиров
      </h2>

      <div className="glass-card overflow-hidden rounded-xl">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-white/5 text-left text-xs uppercase tracking-wider text-gray-500">
                <th className="py-3 pl-4 pr-2">#</th>
                <th className="py-3 pr-4">Игрок</th>
                <th className="py-3 text-center">Голы</th>
                <th className="py-3 text-center">Передачи</th>
                <th className="py-3 pr-4 text-center">Игры</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {rankings.map((player, index) => (
                <RankingRow key={player.id} player={player} index={index} />
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </motion.section>
  )
})
