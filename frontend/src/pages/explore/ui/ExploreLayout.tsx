import { memo, useEffect } from 'react'
import type { ReactNode } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { motion } from 'framer-motion'
import { ExploreSidebar } from '@/features/explore-sidebar'
import { MobileMenuButton, useSidebar } from '@/features/sidebar-navigation'
import { useAuthStore } from '@/shared/stores'
import { cn } from '@/shared/lib/utils'

interface ExploreLayoutProps {
  children?: ReactNode
}

export const ExploreLayout = memo(function ExploreLayout({
  children,
}: ExploreLayoutProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const { isOpen, isCollapsed, toggle, close } = useSidebar()
  const { isAuthenticated } = useAuthStore()

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login')
    }
  }, [isAuthenticated, navigate])

  if (!isAuthenticated) {
    return null
  }

  return (
    <div className="flex min-h-screen">
      {/* Animated background */}
      <div className="bg-blob bg-blob--cyan" style={{ top: '10%', left: '15%' }} />
      <div className="bg-blob bg-blob--purple" style={{ top: '60%', right: '10%' }} />
      <div className="bg-blob bg-blob--pink" style={{ top: '30%', right: '30%' }} />
      <div className="bg-grid" />

      <MobileMenuButton isOpen={isOpen} onToggle={toggle} />

      <ExploreSidebar
        isOpen={isOpen}
        isCollapsed={isCollapsed}
        onToggle={toggle}
        onClose={close}
      />

      <main
        className={cn(
          'relative z-10 flex-1 transition-all duration-300',
          'p-6 pt-20 lg:pt-6'
        )}
      >
        <motion.div
          key={location.pathname}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.15 }}
          className="w-full"
        >
          {children || <Outlet />}
        </motion.div>
      </main>
    </div>
  )
})
