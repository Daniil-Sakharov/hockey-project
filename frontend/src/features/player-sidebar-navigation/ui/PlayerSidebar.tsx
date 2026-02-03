import { memo } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { cn } from '@/shared/lib/utils'
import { useMediaQuery, breakpoints } from '@/shared/hooks'
import { SidebarToggle } from '@/features/sidebar-navigation'
import { PlayerSidebarProfile } from './PlayerSidebarProfile'
import { PlayerSidebarNav } from './PlayerSidebarNav'

function HockeyLogo({ className }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      className={className}
      stroke="currentColor"
      strokeWidth="1.5"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <circle cx="12" cy="12" r="3" fill="currentColor" />
      <path d="M2 12c0-4 3-8 10-8s10 4 10 8" />
      <path d="M4 18l3-6" />
      <path d="M20 18l-3-6" />
      <path d="M4 18h16" strokeWidth="2.5" />
    </svg>
  )
}

interface PlayerSidebarProps {
  isOpen: boolean
  isCollapsed: boolean
  onToggle: () => void
  onClose: () => void
}

export const PlayerSidebar = memo(function PlayerSidebar({
  isOpen,
  isCollapsed,
  onToggle,
  onClose,
}: PlayerSidebarProps) {
  const isDesktop = useMediaQuery(breakpoints.lg)

  return (
    <>
      {/* Mobile overlay */}
      <AnimatePresence>
        {!isDesktop && isOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm lg:hidden"
            aria-hidden="true"
          />
        )}
      </AnimatePresence>

      {/* Sidebar */}
      <AnimatePresence mode="wait">
        {(isDesktop || isOpen) && (
          <motion.aside
            initial={isDesktop ? false : { x: -280 }}
            animate={{ x: 0 }}
            exit={{ x: -280 }}
            transition={{ type: 'spring', damping: 25, stiffness: 200 }}
            className={cn(
              'fixed left-0 top-0 z-50 flex h-screen flex-col',
              'bg-[#0d1224]/95 backdrop-blur-xl',
              'border-r border-[#00d4ff]/10',
              'transition-[width] duration-300',
              isCollapsed ? 'w-16' : 'w-64',
              'lg:relative lg:z-auto'
            )}
          >
            {/* Header */}
            <div
              className={cn(
                'flex h-16 items-center border-b border-white/5',
                isCollapsed ? 'justify-center px-2' : 'justify-between px-4'
              )}
            >
              {!isCollapsed && (
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  className="flex items-center gap-2"
                >
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-[#00d4ff] to-[#8b5cf6]">
                    <HockeyLogo className="h-5 w-5 text-white" />
                  </div>
                  <span className="text-lg font-bold text-white">
                    Star<span className="text-[#00d4ff]">Rink</span>
                  </span>
                </motion.div>
              )}

              {isDesktop && <SidebarToggle isCollapsed={isCollapsed} onToggle={onToggle} />}
            </div>

            {/* Profile */}
            <PlayerSidebarProfile isCollapsed={isCollapsed} />

            {/* Divider */}
            <div className="mx-4 border-t border-white/5" />

            {/* Navigation */}
            <PlayerSidebarNav isCollapsed={isCollapsed} />

            {/* Footer */}
            <div className={cn('border-t border-white/5 p-4', isCollapsed && 'px-2')}>
              {!isCollapsed && <p className="text-xs text-gray-500">© 2025 StarRink</p>}
            </div>
          </motion.aside>
        )}
      </AnimatePresence>
    </>
  )
})
