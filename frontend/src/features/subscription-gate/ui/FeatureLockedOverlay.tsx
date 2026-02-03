import { memo, type ReactNode } from 'react'
import { motion } from 'framer-motion'
import type { FeatureKey } from '@/entities/user'
import { cn } from '@/shared/lib/utils'
import { UpgradePrompt } from './UpgradePrompt'

interface FeatureLockedOverlayProps {
  feature: FeatureKey
  previewContent: ReactNode
  className?: string
  blurAmount?: 'sm' | 'md' | 'lg'
}

const BLUR_CLASSES = {
  sm: 'blur-sm',
  md: 'blur-md',
  lg: 'blur-lg',
}

export const FeatureLockedOverlay = memo(function FeatureLockedOverlay({
  feature,
  previewContent,
  className,
  blurAmount = 'md',
}: FeatureLockedOverlayProps) {
  return (
    <div className={cn('relative overflow-hidden rounded-xl', className)}>
      {/* Blurred preview content */}
      <div
        className={cn(
          'pointer-events-none select-none',
          BLUR_CLASSES[blurAmount]
        )}
        aria-hidden="true"
      >
        {previewContent}
      </div>

      {/* Overlay */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="absolute inset-0 flex items-center justify-center bg-[#0a0e1a]/60 backdrop-blur-sm"
      >
        <UpgradePrompt feature={feature} />
      </motion.div>
    </div>
  )
})
