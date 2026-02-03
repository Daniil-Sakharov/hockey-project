import { memo } from 'react'
import type { ReactNode } from 'react'
import { motion } from 'framer-motion'
import { TrendingUp, TrendingDown, Minus } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { Skeleton } from '@/shared/ui'

type KPIColor = 'blue' | 'purple' | 'pink' | 'cyan' | 'green' | 'red' | 'orange'

interface KPICardProps {
  title: string
  value: number
  icon: ReactNode
  color?: KPIColor
  trend?: 'up' | 'down' | 'neutral'
  trendValue?: string
  isLoading?: boolean
  delay?: number
}

const colorStyles: Record<KPIColor, { text: string; glow: string; bg: string; iconColor: string }> = {
  blue: {
    text: 'text-[#00d4ff]',
    glow: 'shadow-[0_0_20px_rgba(0,212,255,0.3)]',
    bg: 'bg-[#00d4ff]/10',
    iconColor: '#00d4ff',
  },
  purple: {
    text: 'text-[#8b5cf6]',
    glow: 'shadow-[0_0_20px_rgba(139,92,246,0.3)]',
    bg: 'bg-[#8b5cf6]/10',
    iconColor: '#8b5cf6',
  },
  pink: {
    text: 'text-[#ec4899]',
    glow: 'shadow-[0_0_20px_rgba(236,72,153,0.3)]',
    bg: 'bg-[#ec4899]/10',
    iconColor: '#ec4899',
  },
  cyan: {
    text: 'text-[#00ffff]',
    glow: 'shadow-[0_0_20px_rgba(0,255,255,0.3)]',
    bg: 'bg-[#00ffff]/10',
    iconColor: '#00ffff',
  },
  green: {
    text: 'text-[#22c55e]',
    glow: 'shadow-[0_0_20px_rgba(34,197,94,0.3)]',
    bg: 'bg-[#22c55e]/10',
    iconColor: '#22c55e',
  },
  red: {
    text: 'text-[#ef4444]',
    glow: 'shadow-[0_0_20px_rgba(239,68,68,0.3)]',
    bg: 'bg-[#ef4444]/10',
    iconColor: '#ef4444',
  },
  orange: {
    text: 'text-[#f97316]',
    glow: 'shadow-[0_0_20px_rgba(249,115,22,0.3)]',
    bg: 'bg-[#f97316]/10',
    iconColor: '#f97316',
  },
}

const TrendIcons = {
  up: TrendingUp,
  down: TrendingDown,
  neutral: Minus,
}

const trendColors = {
  up: 'text-green-500',
  down: 'text-red-500',
  neutral: 'text-gray-500',
}

export const KPICard = memo(function KPICard({
  title,
  value,
  icon,
  color = 'blue',
  trend,
  trendValue,
  isLoading = false,
  delay = 0,
}: KPICardProps) {
  const styles = colorStyles[color]

  if (isLoading) {
    return (
      <div className="glass-card rounded-xl p-4">
        <div className="flex items-start justify-between">
          <Skeleton className="h-10 w-10 rounded-lg" />
          <Skeleton className="h-4 w-12" />
        </div>
        <Skeleton className="mt-3 h-8 w-20" />
        <Skeleton className="mt-1 h-4 w-16" />
      </div>
    )
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, delay }}
      whileHover={{ scale: 1.02 }}
      className={cn(
        'glass-card rounded-xl p-4 transition-shadow duration-300',
        'hover:shadow-[0_0_30px_rgba(0,212,255,0.2)]'
      )}
    >
      <div className="flex items-start justify-between">
        {/* Icon */}
        <motion.div
          whileHover={{ scale: 1.1, rotate: 5 }}
          className={cn(
            'flex h-12 w-12 items-center justify-center rounded-xl',
            styles.bg,
            styles.glow
          )}
          style={{ color: styles.iconColor }}
        >
          {icon}
        </motion.div>

        {/* Trend */}
        {trend && (
          <div className={cn('flex items-center gap-1 text-xs font-medium', trendColors[trend])}>
            {(() => {
              const TrendIcon = TrendIcons[trend]
              return <TrendIcon size={14} />
            })()}
            {trendValue && <span>{trendValue}</span>}
          </div>
        )}
      </div>

      {/* Value */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: delay + 0.2 }}
        className={cn('mt-3 text-3xl font-bold', styles.text)}
      >
        {value.toLocaleString('ru-RU')}
      </motion.div>

      {/* Title */}
      <div className="mt-1 text-sm text-gray-400">{title}</div>
    </motion.div>
  )
})
