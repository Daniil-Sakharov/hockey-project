import { memo } from 'react'
import { motion } from 'framer-motion'
import { Star } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useScoutStore } from '@/shared/stores'

interface PlayerCardProps {
  id: string
  name: string
  team?: string
  position?: 'forward' | 'defender' | 'goalie'
  goals?: number
  assists?: number
  points?: number
  avatarUrl?: string
  className?: string
  onClick?: () => void
}

const positionLabels: Record<string, string> = {
  forward: 'Нападающий',
  defender: 'Защитник',
  goalie: 'Вратарь',
}

export const PlayerCard = memo(function PlayerCard({
  id,
  name,
  team,
  position,
  goals = 0,
  assists = 0,
  points = 0,
  avatarUrl,
  className,
  onClick,
}: PlayerCardProps) {
  const isInWatchlist = useScoutStore((state) => state.isInWatchlist(id))
  const addToWatchlist = useScoutStore((state) => state.addToWatchlist)
  const removeFromWatchlist = useScoutStore((state) => state.removeFromWatchlist)

  const handleWatchlistToggle = (e: React.MouseEvent) => {
    e.stopPropagation()
    if (isInWatchlist) {
      removeFromWatchlist(id)
    } else {
      addToWatchlist(id)
    }
  }

  const initials = name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()

  return (
    <motion.div
      className={cn(
        'glass-card relative cursor-pointer rounded-xl p-4',
        'transition-all duration-300',
        'hover:border-[#00d4ff]/30 hover:shadow-[0_0_30px_rgba(0,212,255,0.15)]',
        className
      )}
      whileHover={{ scale: 1.02, y: -2 }}
      whileTap={{ scale: 0.98 }}
      onClick={onClick}
    >
      {/* Watchlist button */}
      <motion.button
        className={cn(
          'absolute right-3 top-3 rounded-full p-1.5',
          'transition-colors duration-200',
          isInWatchlist
            ? 'bg-[#ec4899]/20 text-[#ec4899]'
            : 'bg-white/5 text-gray-500 hover:bg-white/10 hover:text-[#ec4899]'
        )}
        onClick={handleWatchlistToggle}
        whileHover={{ scale: 1.1 }}
        whileTap={{ scale: 0.9 }}
        title={isInWatchlist ? 'Удалить из избранного' : 'Добавить в избранное'}
      >
        <Star
          size={16}
          fill={isInWatchlist ? 'currentColor' : 'none'}
        />
      </motion.button>

      {/* Player info */}
      <div className="flex items-center gap-3">
        {/* Avatar */}
        <div
          className={cn(
            'flex h-12 w-12 items-center justify-center rounded-full',
            'bg-gradient-to-br from-[#00d4ff]/30 to-[#8b5cf6]/30',
            'border border-[#00d4ff]/20'
          )}
        >
          {avatarUrl ? (
            <img
              src={avatarUrl}
              alt={name}
              className="h-full w-full rounded-full object-cover"
            />
          ) : (
            <span className="text-sm font-bold text-[#00d4ff]">{initials}</span>
          )}
        </div>

        {/* Name and details */}
        <div className="min-w-0 flex-1">
          <h4 className="truncate font-semibold text-white">{name}</h4>
          {team && (
            <p className="truncate text-xs text-gray-400">{team}</p>
          )}
          {position && (
            <span className="mt-1 inline-block rounded bg-[#00d4ff]/10 px-1.5 py-0.5 text-xs text-[#00d4ff]">
              {positionLabels[position] || position}
            </span>
          )}
        </div>
      </div>

      {/* Stats */}
      <div className="mt-3 flex items-center gap-4 border-t border-white/5 pt-3">
        <StatItem label="Г" value={goals} />
        <StatItem label="П" value={assists} />
        <StatItem label="О" value={points} highlight />
      </div>
    </motion.div>
  )
})

const StatItem = memo(function StatItem({
  label,
  value,
  highlight = false,
}: {
  label: string
  value: number
  highlight?: boolean
}) {
  return (
    <div className="flex items-center gap-1.5">
      <span className="text-xs text-gray-500">{label}</span>
      <span
        className={cn(
          'text-sm font-bold',
          highlight ? 'text-[#00d4ff]' : 'text-white'
        )}
      >
        {value}
      </span>
    </div>
  )
})
