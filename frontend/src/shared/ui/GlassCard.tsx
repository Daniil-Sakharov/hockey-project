import { memo } from 'react'
import type { ReactNode } from 'react'
import { cn } from '@/shared/lib/utils'

type GlowColor = 'blue' | 'purple' | 'pink' | 'cyan'

interface GlassCardProps {
  children: ReactNode
  className?: string
  title?: string
  glowColor?: GlowColor
  neonBorder?: boolean
  interactive?: boolean
}

const glowStyles: Record<GlowColor, string> = {
  blue: 'hover:shadow-[0_0_30px_rgba(0,212,255,0.15),0_0_60px_rgba(0,212,255,0.08)] hover:border-[rgba(0,212,255,0.35)]',
  purple: 'hover:shadow-[0_0_30px_rgba(139,92,246,0.15),0_0_60px_rgba(139,92,246,0.08)] hover:border-[rgba(139,92,246,0.35)]',
  pink: 'hover:shadow-[0_0_30px_rgba(236,72,153,0.15),0_0_60px_rgba(236,72,153,0.08)] hover:border-[rgba(236,72,153,0.35)]',
  cyan: 'hover:shadow-[0_0_30px_rgba(0,255,255,0.15),0_0_60px_rgba(0,255,255,0.08)] hover:border-[rgba(0,255,255,0.35)]',
}

export const GlassCard = memo(function GlassCard({
  children,
  className,
  title,
  glowColor = 'blue',
  neonBorder = false,
  interactive = true,
}: GlassCardProps) {
  return (
    <div
      className={cn(
        'glass-card rounded-xl p-6 transition-all duration-300',
        interactive && glowStyles[glowColor],
        neonBorder && 'neon-border',
        className
      )}
    >
      {title && (
        <h3 className="mb-4 text-sm font-medium uppercase tracking-wider text-gray-400">
          {title}
        </h3>
      )}
      {children}
    </div>
  )
})
