import { type ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/shared/stores'

interface ProtectedRouteProps {
  children: ReactNode
  requireLinkedPlayer?: boolean
}

export function ProtectedRoute({ children, requireLinkedPlayer = false }: ProtectedRouteProps) {
  const location = useLocation()
  const { isAuthenticated, user } = useAuthStore()

  // Если не авторизован - редирект на логин
  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }

  // Если требуется привязанный игрок, но его нет - редирект на привязку
  if (requireLinkedPlayer && !user?.linkedPlayerId) {
    return <Navigate to="/link-player" state={{ from: location }} replace />
  }

  return <>{children}</>
}

interface GuestOnlyRouteProps {
  children: ReactNode
}

export function GuestOnlyRoute({ children }: GuestOnlyRouteProps) {
  const { isAuthenticated, user } = useAuthStore()

  // Если авторизован - редирект по роли
  if (isAuthenticated && user) {
    switch (user.role) {
      case 'player':
        return <Navigate to={user.linkedPlayerId ? '/player' : '/link-player'} replace />
      case 'scout':
      case 'coach':
        return <Navigate to="/dashboard" replace />
      default:
        return <Navigate to="/explore" replace />
    }
  }

  return <>{children}</>
}
