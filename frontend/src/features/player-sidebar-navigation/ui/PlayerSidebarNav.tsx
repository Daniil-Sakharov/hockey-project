import { memo } from 'react'
import type { ReactNode } from 'react'
import { NavLink } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  User,
  BarChart3,
  Calendar,
  Trophy,
  GitCompare,
  Eye,
  Video,
  Sparkles,
  CreditCard,
  Settings,
  Lock,
} from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useAuthStore, usePlayerDashboardStore } from '@/shared/stores'
import type { FeatureKey, SubscriptionTier } from '@/entities/user'

interface NavItem {
  id: string
  label: string
  icon: ReactNode
  path: string
  badge?: number
  requiredTier?: SubscriptionTier
  feature?: FeatureKey
  dividerBefore?: boolean
}

interface PlayerSidebarNavProps {
  isCollapsed: boolean
}

const TIER_COLORS: Record<SubscriptionTier, string> = {
  free: 'text-gray-500',
  pro: 'text-[#8b5cf6]',
  ultra: 'text-[#f59e0b]',
}

export const PlayerSidebarNav = memo(function PlayerSidebarNav({
  isCollapsed,
}: PlayerSidebarNavProps) {
  const hasFeature = useAuthStore((state) => state.hasFeature)
  const getSubscriptionTier = useAuthStore((state) => state.getSubscriptionTier)
  const unreadCount = usePlayerDashboardStore((state) => state.getUnreadNotificationsCount())
  const currentTier = getSubscriptionTier()

  const navItems: NavItem[] = [
    // FREE
    {
      id: 'profile',
      label: 'Мой профиль',
      icon: <User size={20} />,
      path: '/player',
    },
    {
      id: 'stats',
      label: 'Статистика',
      icon: <BarChart3 size={20} />,
      path: '/player/stats',
    },
    {
      id: 'calendar',
      label: 'Календарь',
      icon: <Calendar size={20} />,
      path: '/player/calendar',
    },
    {
      id: 'achievements',
      label: 'Достижения',
      icon: <Trophy size={20} />,
      path: '/player/achievements',
    },
    // PRO
    {
      id: 'compare',
      label: 'Сравнение',
      icon: <GitCompare size={20} />,
      path: '/player/compare',
      requiredTier: 'pro',
      feature: 'player_comparison',
      dividerBefore: true,
    },
    {
      id: 'notifications',
      label: 'Просмотры',
      icon: <Eye size={20} />,
      path: '/player/notifications',
      badge: unreadCount > 0 ? unreadCount : undefined,
      requiredTier: 'pro',
      feature: 'scout_notifications',
    },
    // ULTRA
    {
      id: 'highlights',
      label: 'Видео',
      icon: <Video size={20} />,
      path: '/player/highlights',
      requiredTier: 'ultra',
      feature: 'video_highlights',
      dividerBefore: true,
    },
    {
      id: 'recommendations',
      label: 'AI советы',
      icon: <Sparkles size={20} />,
      path: '/player/recommendations',
      requiredTier: 'ultra',
      feature: 'ai_recommendations',
    },
    // Settings
    {
      id: 'subscription',
      label: 'Подписка',
      icon: <CreditCard size={20} />,
      path: '/player/subscription',
      dividerBefore: true,
    },
    {
      id: 'settings',
      label: 'Настройки',
      icon: <Settings size={20} />,
      path: '/player/settings',
    },
  ]

  return (
    <nav className="mt-6 flex-1 overflow-y-auto px-2" aria-label="Навигация игрока">
      <ul className="space-y-1">
        {navItems.map((item, index) => {
          const hasAccess = !item.feature || hasFeature(item.feature)
          const isLocked = item.requiredTier && !hasAccess

          return (
            <motion.li
              key={item.id}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: index * 0.03 }}
            >
              {/* Divider */}
              {item.dividerBefore && (
                <div
                  className={cn('my-3 border-t border-white/5', isCollapsed ? 'mx-1' : 'mx-2')}
                />
              )}

              <NavLink
                to={item.path}
                className={({ isActive }) =>
                  cn(
                    'relative flex items-center gap-3 rounded-lg px-3 py-2.5',
                    'transition-all duration-200',
                    isLocked
                      ? 'cursor-default opacity-60'
                      : 'hover:bg-[#00d4ff]/10',
                    isActive && !isLocked
                      ? 'bg-[#00d4ff]/20 text-[#00d4ff] shadow-[0_0_15px_rgba(0,212,255,0.2)]'
                      : isLocked
                        ? 'text-gray-600'
                        : 'text-gray-400 hover:text-white',
                    isCollapsed && 'justify-center px-2'
                  )
                }
                onClick={(e) => {
                  if (isLocked) {
                    e.preventDefault()
                  }
                }}
              >
                <span className="relative flex-shrink-0">
                  {isLocked ? (
                    <Lock size={20} className={TIER_COLORS[item.requiredTier!]} />
                  ) : (
                    item.icon
                  )}
                  {/* Badge for collapsed state */}
                  {isCollapsed && item.badge && !isLocked && (
                    <span className="absolute -right-1.5 -top-1.5 flex h-4 w-4 items-center justify-center rounded-full bg-[#ec4899] text-[10px] font-bold text-white">
                      {item.badge > 9 ? '9+' : item.badge}
                    </span>
                  )}
                </span>
                {!isCollapsed && (
                  <>
                    <span className="flex-1 text-sm font-medium">{item.label}</span>
                    {/* Tier badge */}
                    {item.requiredTier && (
                      <span
                        className={cn(
                          'ml-auto text-[10px] font-bold uppercase',
                          TIER_COLORS[item.requiredTier]
                        )}
                      >
                        {item.requiredTier}
                      </span>
                    )}
                    {/* Count badge */}
                    {item.badge && !isLocked && (
                      <span className="ml-2 flex h-5 min-w-[20px] items-center justify-center rounded-full bg-[#ec4899] px-1.5 text-xs font-bold text-white">
                        {item.badge}
                      </span>
                    )}
                  </>
                )}
              </NavLink>
            </motion.li>
          )
        })}
      </ul>

      {/* Subscription badge at bottom */}
      {!isCollapsed && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="mt-4 px-2"
        >
          <div
            className={cn(
              'rounded-lg px-3 py-2 text-center text-xs font-medium',
              currentTier === 'ultra'
                ? 'bg-[#f59e0b]/20 text-[#f59e0b]'
                : currentTier === 'pro'
                  ? 'bg-[#8b5cf6]/20 text-[#8b5cf6]'
                  : 'bg-white/5 text-gray-500'
            )}
          >
            {currentTier === 'ultra'
              ? 'ULTRA подписка'
              : currentTier === 'pro'
                ? 'PRO подписка'
                : 'Бесплатный план'}
          </div>
        </motion.div>
      )}
    </nav>
  )
})
