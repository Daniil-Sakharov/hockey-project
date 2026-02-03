import { memo } from 'react'
import type { ReactNode } from 'react'
import { NavLink } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  LayoutDashboard,
  Search,
  Star,
  GitCompare,
  Bell,
  Settings,
} from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useScoutStore } from '@/shared/stores'

interface NavItem {
  id: string
  label: string
  icon: ReactNode
  path: string
  badge?: number
}

interface SidebarNavProps {
  isCollapsed: boolean
}

export const SidebarNav = memo(function SidebarNav({
  isCollapsed,
}: SidebarNavProps) {
  const watchlistCount = useScoutStore((state) => state.watchlist.length)

  const navItems: NavItem[] = [
    {
      id: 'dashboard',
      label: 'Dashboard',
      icon: <LayoutDashboard size={20} />,
      path: '/dashboard',
    },
    {
      id: 'search',
      label: 'Поиск игроков',
      icon: <Search size={20} />,
      path: '/dashboard/search',
    },
    {
      id: 'watchlist',
      label: 'Избранные',
      icon: <Star size={20} />,
      path: '/dashboard/watchlist',
      badge: watchlistCount > 0 ? watchlistCount : undefined,
    },
    {
      id: 'compare',
      label: 'Сравнение',
      icon: <GitCompare size={20} />,
      path: '/dashboard/compare',
    },
    {
      id: 'notifications',
      label: 'Уведомления',
      icon: <Bell size={20} />,
      path: '/dashboard/notifications',
    },
    {
      id: 'settings',
      label: 'Настройки',
      icon: <Settings size={20} />,
      path: '/dashboard/settings',
    },
  ]

  return (
    <nav className="mt-6 flex-1 px-2" aria-label="Основная навигация">
      <ul className="space-y-1">
        {navItems.map((item, index) => (
          <motion.li
            key={item.id}
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: index * 0.05 }}
          >
            <NavLink
              to={item.path}
              className={({ isActive }) =>
                cn(
                  'relative flex items-center gap-3 rounded-lg px-3 py-2.5',
                  'transition-all duration-200',
                  'hover:bg-[#00d4ff]/10',
                  isActive
                    ? 'bg-[#00d4ff]/20 text-[#00d4ff] shadow-[0_0_15px_rgba(0,212,255,0.2)]'
                    : 'text-gray-400 hover:text-white',
                  isCollapsed && 'justify-center px-2'
                )
              }
            >
              <span className="relative flex-shrink-0">
                {item.icon}
                {/* Badge for collapsed state */}
                {isCollapsed && item.badge && (
                  <span className="absolute -right-1.5 -top-1.5 flex h-4 w-4 items-center justify-center rounded-full bg-[#ec4899] text-[10px] font-bold text-white">
                    {item.badge > 9 ? '9+' : item.badge}
                  </span>
                )}
              </span>
              {!isCollapsed && (
                <>
                  <span className="text-sm font-medium">{item.label}</span>
                  {/* Badge for expanded state */}
                  {item.badge && (
                    <span className="ml-auto flex h-5 min-w-[20px] items-center justify-center rounded-full bg-[#ec4899] px-1.5 text-xs font-bold text-white">
                      {item.badge}
                    </span>
                  )}
                </>
              )}
            </NavLink>
          </motion.li>
        ))}
      </ul>
    </nav>
  )
})
