import { memo } from 'react'
import type { ReactNode } from 'react'
import { Sidebar, MobileMenuButton, useSidebar } from '@/features/sidebar-navigation'
import { cn } from '@/shared/lib/utils'

interface DashboardLayoutProps {
  children: ReactNode
}

export const DashboardLayout = memo(function DashboardLayout({
  children,
}: DashboardLayoutProps) {
  const { isOpen, isCollapsed, toggle, close } = useSidebar()

  return (
    <div className="flex min-h-screen bg-[#0a0e1a]">
      {/* Mobile menu button */}
      <MobileMenuButton isOpen={isOpen} onToggle={toggle} />

      {/* Sidebar */}
      <Sidebar
        isOpen={isOpen}
        isCollapsed={isCollapsed}
        onToggle={toggle}
        onClose={close}
      />

      {/* Main content */}
      <main
        className={cn(
          'flex-1 transition-all duration-300',
          'p-6 pt-20 lg:pt-6',
          isCollapsed ? 'lg:ml-16' : 'lg:ml-64'
        )}
      >
        <div className="mx-auto max-w-7xl">
          {children}
        </div>
      </main>
    </div>
  )
})
