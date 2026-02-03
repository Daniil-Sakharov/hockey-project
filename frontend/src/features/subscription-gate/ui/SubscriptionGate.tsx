import { memo, type ReactNode } from 'react'
import type { FeatureKey } from '@/entities/user'
import { useAuthStore } from '@/shared/stores'
import { UpgradePrompt } from './UpgradePrompt'

interface SubscriptionGateProps {
  feature: FeatureKey
  children: ReactNode
  fallback?: ReactNode
  showUpgrade?: boolean
}

export const SubscriptionGate = memo(function SubscriptionGate({
  feature,
  children,
  fallback,
  showUpgrade = true,
}: SubscriptionGateProps) {
  const hasAccess = useAuthStore((state) => state.hasFeature(feature))

  // Если есть доступ - показываем контент
  if (hasAccess) {
    return <>{children}</>
  }

  // Если есть fallback - показываем его
  if (fallback) {
    return <>{fallback}</>
  }

  // Если нужно показать промо апгрейда
  if (showUpgrade) {
    return <UpgradePrompt feature={feature} />
  }

  // Иначе ничего не показываем
  return null
})
