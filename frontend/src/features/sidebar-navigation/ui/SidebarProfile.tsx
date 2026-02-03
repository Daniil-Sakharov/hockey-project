import { memo } from 'react'
import { motion } from 'framer-motion'
import { Binoculars, Megaphone, BarChart3 } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useScoutStore, type ScoutProfile } from '@/shared/stores'

interface SidebarProfileProps {
  isCollapsed: boolean
}

const roleConfig: Record<ScoutProfile['role'], { label: string; icon: React.ReactNode }> = {
  scout: { label: 'Скаут', icon: <Binoculars size={14} /> },
  coach: { label: 'Тренер', icon: <Megaphone size={14} /> },
  analyst: { label: 'Аналитик', icon: <BarChart3 size={14} /> },
}

export const SidebarProfile = memo(function SidebarProfile({
  isCollapsed,
}: SidebarProfileProps) {
  const profile = useScoutStore((state) => state.profile)
  const watchlistCount = useScoutStore((state) => state.watchlist.length)

  const initials = profile.name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()

  const roleInfo = roleConfig[profile.role]

  return (
    <div className={cn('p-4', isCollapsed ? 'px-2' : 'px-4')}>
      {/* Avatar */}
      <div className="flex flex-col items-center">
        <div
          className={cn(
            'relative rounded-full bg-gradient-to-br from-[#00d4ff] to-[#8b5cf6]',
            'flex items-center justify-center',
            'shadow-[0_0_20px_rgba(0,212,255,0.5)]',
            'transition-all duration-300',
            isCollapsed ? 'h-10 w-10' : 'h-20 w-20'
          )}
        >
          {profile.avatarUrl ? (
            <img
              src={profile.avatarUrl}
              alt={profile.name}
              className="h-full w-full rounded-full object-cover"
            />
          ) : (
            <span
              className={cn(
                'font-bold text-white',
                isCollapsed ? 'text-xs' : 'text-2xl'
              )}
            >
              {initials}
            </span>
          )}
          {/* Online indicator */}
          <span className="absolute bottom-0 right-0 h-3 w-3 rounded-full border-2 border-[#0a0e1a] bg-green-500" />
        </div>

        {/* Scout info */}
        {!isCollapsed && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="mt-3 text-center"
          >
            <h3 className="text-lg font-semibold text-white">{profile.name}</h3>

            {/* Role badge */}
            <div className="mt-1 flex items-center justify-center gap-1.5">
              <span className="text-[#00d4ff]">{roleInfo.icon}</span>
              <span className="text-sm text-gray-400">{roleInfo.label}</span>
            </div>

            {/* Club */}
            {profile.club && (
              <p className="mt-1 text-xs text-gray-500">{profile.club}</p>
            )}
          </motion.div>
        )}
      </div>

      {/* Watchlist count */}
      {!isCollapsed && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="mt-4"
        >
          <div className="rounded-lg bg-white/5 p-3 text-center">
            <div className="text-2xl font-bold text-[#00d4ff]">{watchlistCount}</div>
            <div className="text-xs text-gray-500">Игроков в избранном</div>
          </div>
        </motion.div>
      )}
    </div>
  )
})
