import { memo, useEffect } from 'react'
import type { ReactNode } from 'react'
import { Outlet, useNavigate } from 'react-router-dom'
import { PlayerSidebar } from '@/features/player-sidebar-navigation'
import { MobileMenuButton, useSidebar } from '@/features/sidebar-navigation'
import { useAuthStore, usePlayerDashboardStore } from '@/shared/stores'
import { MOCK_PLAYER_PROFILE, MOCK_TEAM_MATCHES, MOCK_ACHIEVEMENTS, MOCK_SCOUT_NOTIFICATIONS } from '@/shared/mocks'
import { cn } from '@/shared/lib/utils'

interface PlayerDashboardLayoutProps {
  children?: ReactNode
}

export const PlayerDashboardLayout = memo(function PlayerDashboardLayout({
  children,
}: PlayerDashboardLayoutProps) {
  const navigate = useNavigate()
  const { isOpen, isCollapsed, toggle, close } = useSidebar()
  const { isAuthenticated, user } = useAuthStore()
  const { linkedPlayer, setLinkedPlayer } = usePlayerDashboardStore()

  // Auth check
  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login')
    }
  }, [isAuthenticated, navigate])

  // Load mock data if player is linked but data not loaded
  useEffect(() => {
    if (user?.linkedPlayerId && !linkedPlayer) {
      // В реальности здесь будет API запрос
      setLinkedPlayer(MOCK_PLAYER_PROFILE)

      // Also load other data
      const store = usePlayerDashboardStore.getState()
      if (store.teamMatches.length === 0) {
        usePlayerDashboardStore.setState({
          teamMatches: MOCK_TEAM_MATCHES,
          achievements: MOCK_ACHIEVEMENTS,
          scoutNotifications: MOCK_SCOUT_NOTIFICATIONS,
        })
      }
    }
  }, [user?.linkedPlayerId, linkedPlayer, setLinkedPlayer])

  if (!isAuthenticated) {
    return null
  }

  return (
    <div className="flex min-h-screen bg-[#0a0e1a]">
      {/* Mobile menu button */}
      <MobileMenuButton isOpen={isOpen} onToggle={toggle} />

      {/* Sidebar */}
      <PlayerSidebar
        isOpen={isOpen}
        isCollapsed={isCollapsed}
        onToggle={toggle}
        onClose={close}
      />

      {/* Main content */}
      <main
        className={cn(
          'flex-1 transition-all duration-300',
          'p-6 pt-20 lg:pt-6'
        )}
      >
        <div className="w-full">
          {children || <Outlet />}
        </div>
      </main>
    </div>
  )
})
