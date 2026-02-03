import { memo } from 'react'
import { motion } from 'framer-motion'
import { cn } from '@/shared/lib/utils'

interface SidebarToggleProps {
  isCollapsed: boolean
  onToggle: () => void
}

export const SidebarToggle = memo(function SidebarToggle({
  isCollapsed,
  onToggle,
}: SidebarToggleProps) {
  return (
    <button
      onClick={onToggle}
      className={cn(
        'flex h-8 w-8 items-center justify-center rounded-lg',
        'bg-white/5 text-gray-400 transition-all duration-200',
        'hover:bg-[#00d4ff]/20 hover:text-[#00d4ff]',
        'focus:outline-none focus:ring-2 focus:ring-[#00d4ff]/50'
      )}
      aria-label={isCollapsed ? 'Развернуть меню' : 'Свернуть меню'}
      aria-expanded={!isCollapsed}
    >
      <motion.svg
        width="16"
        height="16"
        viewBox="0 0 16 16"
        fill="none"
        animate={{ rotate: isCollapsed ? 180 : 0 }}
        transition={{ duration: 0.2 }}
      >
        <path
          d="M10 12L6 8L10 4"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </motion.svg>
    </button>
  )
})

interface MobileMenuButtonProps {
  isOpen: boolean
  onToggle: () => void
}

export const MobileMenuButton = memo(function MobileMenuButton({
  isOpen,
  onToggle,
}: MobileMenuButtonProps) {
  return (
    <button
      onClick={onToggle}
      className={cn(
        'fixed left-4 top-4 z-50 flex h-10 w-10 items-center justify-center rounded-lg',
        'bg-[#0d1224] text-gray-400 shadow-lg',
        'transition-all duration-200',
        'hover:bg-[#00d4ff]/20 hover:text-[#00d4ff]',
        'focus:outline-none focus:ring-2 focus:ring-[#00d4ff]/50',
        'lg:hidden'
      )}
      aria-label={isOpen ? 'Закрыть меню' : 'Открыть меню'}
      aria-expanded={isOpen}
    >
      <motion.div
        animate={isOpen ? 'open' : 'closed'}
        className="flex flex-col gap-1"
      >
        <motion.span
          variants={{
            closed: { rotate: 0, y: 0 },
            open: { rotate: 45, y: 6 },
          }}
          className="h-0.5 w-5 bg-current"
        />
        <motion.span
          variants={{
            closed: { opacity: 1 },
            open: { opacity: 0 },
          }}
          className="h-0.5 w-5 bg-current"
        />
        <motion.span
          variants={{
            closed: { rotate: 0, y: 0 },
            open: { rotate: -45, y: -6 },
          }}
          className="h-0.5 w-5 bg-current"
        />
      </motion.div>
    </button>
  )
})
