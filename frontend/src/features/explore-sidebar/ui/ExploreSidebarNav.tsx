import { memo } from 'react'
import type { ReactNode } from 'react'
import { NavLink } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  Home,
  Trophy,
  Star,
  Settings,
  Lock,
  Search,
  Calendar,
  ClipboardList,
  Medal,
  Zap,
} from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useAuthStore } from '@/shared/stores'
import type { FeatureKey, SubscriptionTier } from '@/entities/user'

interface NavItem {
  id: string
  label: string
  icon: ReactNode
  path: string
  requiredTier?: SubscriptionTier
  feature?: FeatureKey
  dividerBefore?: boolean
}

interface ExploreSidebarNavProps {
  isCollapsed: boolean
}

const TIER_COLORS: Record<SubscriptionTier, string> = {
  free: 'text-gray-500',
  pro: 'text-[#8b5cf6]',
  ultra: 'text-[#f59e0b]',
}

export const ExploreSidebarNav = memo(function ExploreSidebarNav({
  isCollapsed,
}: ExploreSidebarNavProps) {
  const hasFeature = useAuthStore((state) => state.hasFeature)
  const getSubscriptionTier = useAuthStore((state) => state.getSubscriptionTier)
  const currentTier = getSubscriptionTier()

  const navItems: NavItem[] = [
    {
      id: 'overview',
      label: 'Обзор',
      icon: <Home size={20} />,
      path: '/explore',
    },
    {
      id: 'tournaments',
      label: 'Турниры',
      icon: <Trophy size={20} />,
      path: '/explore/tournaments',
    },
    {
      id: 'players',
      label: 'Игроки',
      icon: <Search size={20} />,
      path: '/explore/players',
    },
    {
      id: 'rankings',
      label: 'Рейтинг',
      icon: <Medal size={20} />,
      path: '/explore/rankings',
    },
    {
      id: 'results',
      label: 'Результаты',
      icon: <ClipboardList size={20} />,
      path: '/explore/results',
      dividerBefore: true,
    },
    {
      id: 'calendar',
      label: 'Календарь',
      icon: <Calendar size={20} />,
      path: '/explore/calendar',
    },
    {
      id: 'predictions',
      label: 'Прогнозы',
      icon: <Zap size={20} />,
      path: '/explore/predictions',
    },
    {
      id: 'favorites',
      label: 'Избранное',
      icon: <Star size={20} />,
      path: '/explore/favorites',
      requiredTier: 'pro',
      feature: 'player_comparison',
      dividerBefore: true,
    },
    {
      id: 'settings',
      label: 'Настройки',
      icon: <Settings size={20} />,
      path: '/explore/settings',
      dividerBefore: true,
    },
  ]

  return (
    <nav className="mt-6 flex-1 overflow-y-auto px-2" aria-label="Навигация">
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
              whileHover={isLocked ? undefined : { x: 4, transition: { type: 'spring' as const, stiffness: 400, damping: 15 } }}
            >
              {item.dividerBefore && (
                <div
                  className={cn('my-3 border-t border-white/5', isCollapsed ? 'mx-1' : 'mx-2')}
                />
              )}

              <NavLink
                to={item.path}
                end={item.path === '/explore'}
                className={({ isActive }) =>
                  cn(
                    'group relative flex items-center gap-3 rounded-lg px-3 py-2.5',
                    'transition-all duration-200',
                    isLocked
                      ? 'cursor-default opacity-60'
                      : 'hover:bg-white/[0.06]',
                    isActive && !isLocked
                      ? 'text-[#00d4ff] nav-active'
                      : isLocked
                        ? 'text-gray-600'
                        : 'text-gray-400 hover:text-white',
                    isCollapsed && 'justify-center px-2'
                  )
                }
                onClick={(e) => {
                  if (isLocked) e.preventDefault()
                }}
              >
                <span className="icon-glow relative flex-shrink-0 transition-transform duration-200 group-hover:scale-110">
                  {isLocked ? (
                    <Lock size={20} className={TIER_COLORS[item.requiredTier!]} />
                  ) : (
                    item.icon
                  )}
                </span>
                {!isCollapsed && (
                  <>
                    <span className="flex-1 text-sm font-medium">{item.label}</span>
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
                  </>
                )}
              </NavLink>
            </motion.li>
          )
        })}
      </ul>

      {/* Subscription badge */}
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
