import { memo } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { LogOut, Crown, Sparkles } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useAuthStore, usePlayerDashboardStore } from '@/shared/stores'

interface PlayerSidebarProfileProps {
  isCollapsed: boolean
}

export const PlayerSidebarProfile = memo(function PlayerSidebarProfile({
  isCollapsed,
}: PlayerSidebarProfileProps) {
  const { user, logout, getSubscriptionTier } = useAuthStore()
  const linkedPlayer = usePlayerDashboardStore((state) => state.linkedPlayer)
  const tier = getSubscriptionTier()

  const displayName = linkedPlayer?.name || user?.email?.split('@')[0] || 'Игрок'
  const teamName = linkedPlayer?.team || 'Команда не привязана'
  const position = linkedPlayer?.position
    ? linkedPlayer.position === 'forward'
      ? 'Нападающий'
      : linkedPlayer.position === 'defender'
        ? 'Защитник'
        : 'Вратарь'
    : null

  const TierIcon = tier === 'ultra' ? Crown : tier === 'pro' ? Sparkles : null
  const tierColor =
    tier === 'ultra' ? 'text-[#f59e0b]' : tier === 'pro' ? 'text-[#8b5cf6]' : null

  if (isCollapsed) {
    return (
      <div className="flex flex-col items-center py-4">
        <Link
          to="/player"
          className="relative flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-[#00d4ff] to-[#8b5cf6] text-white font-bold text-sm"
        >
          {displayName.charAt(0).toUpperCase()}
          {TierIcon && (
            <span
              className={cn(
                'absolute -bottom-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-[#0d1224]',
                tierColor
              )}
            >
              <TierIcon size={10} />
            </span>
          )}
        </Link>
      </div>
    )
  }

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="p-4"
    >
      <div className="flex items-start gap-3">
        {/* Avatar */}
        <Link
          to="/player"
          className="relative flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-[#00d4ff] to-[#8b5cf6] text-white font-bold text-lg hover:scale-105 transition-transform"
        >
          {displayName.charAt(0).toUpperCase()}
          {TierIcon && (
            <span
              className={cn(
                'absolute -bottom-1 -right-1 flex h-5 w-5 items-center justify-center rounded-full bg-[#0d1224] border-2 border-[#0d1224]',
                tierColor
              )}
            >
              <TierIcon size={12} />
            </span>
          )}
        </Link>

        {/* Info */}
        <div className="flex-1 min-w-0">
          <Link to="/player" className="block">
            <h3 className="truncate text-sm font-semibold text-white hover:text-[#00d4ff] transition-colors">
              {displayName.split(' ').slice(0, 2).join(' ')}
            </h3>
          </Link>
          <p className="truncate text-xs text-gray-500">{teamName}</p>
          {position && linkedPlayer?.jerseyNumber && (
            <p className="text-xs text-gray-600">
              {position} • #{linkedPlayer.jerseyNumber}
            </p>
          )}
        </div>

        {/* Logout */}
        <button
          onClick={logout}
          className="flex-shrink-0 rounded-lg p-2 text-gray-500 hover:bg-white/5 hover:text-white transition-colors"
          title="Выйти"
        >
          <LogOut size={16} />
        </button>
      </div>

      {/* Regional rank (if available) */}
      {linkedPlayer?.regionalRank && linkedPlayer?.totalPlayersInRegion && (
        <div className="mt-3 flex items-center gap-2 rounded-lg bg-white/5 px-3 py-2">
          <span className="text-xs text-gray-400">Рейтинг региона:</span>
          <span className="text-xs font-semibold text-[#00d4ff]">
            #{linkedPlayer.regionalRank}
          </span>
          <span className="text-xs text-gray-600">
            из {linkedPlayer.totalPlayersInRegion}
          </span>
        </div>
      )}
    </motion.div>
  )
})
