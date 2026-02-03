import { memo, useEffect, useState } from 'react'
import { motion } from 'framer-motion'
import { Users, Building2, Trophy } from 'lucide-react'
import { cn } from '@/shared/lib/utils'

interface StatsData {
  totalPlayers: number
  totalTeams: number
  totalTournaments: number
}

interface KPICardProps {
  title: string
  value: number
  icon: React.ReactNode
  color: 'cyan' | 'purple' | 'pink'
  delay: number
}

const colorConfig = {
  cyan: {
    bg: 'from-[#00d4ff]/20 to-[#00d4ff]/5',
    border: 'border-[#00d4ff]/30',
    text: 'text-[#00d4ff]',
    glow: 'shadow-[0_0_30px_rgba(0,212,255,0.2)]',
    iconBg: 'bg-[#00d4ff]/20',
  },
  purple: {
    bg: 'from-[#8b5cf6]/20 to-[#8b5cf6]/5',
    border: 'border-[#8b5cf6]/30',
    text: 'text-[#8b5cf6]',
    glow: 'shadow-[0_0_30px_rgba(139,92,246,0.2)]',
    iconBg: 'bg-[#8b5cf6]/20',
  },
  pink: {
    bg: 'from-[#ec4899]/20 to-[#ec4899]/5',
    border: 'border-[#ec4899]/30',
    text: 'text-[#ec4899]',
    glow: 'shadow-[0_0_30px_rgba(236,72,153,0.2)]',
    iconBg: 'bg-[#ec4899]/20',
  },
}

const KPICard = memo(function KPICard({
  title,
  value,
  icon,
  color,
  delay,
}: KPICardProps) {
  const colors = colorConfig[color]
  const [displayValue, setDisplayValue] = useState(0)

  // Animated counter
  useEffect(() => {
    const duration = 1500
    const startTime = Date.now()
    const startValue = 0

    const animate = () => {
      const now = Date.now()
      const progress = Math.min((now - startTime) / duration, 1)
      const eased = 1 - Math.pow(1 - progress, 3) // easeOutCubic
      const current = Math.floor(startValue + (value - startValue) * eased)

      setDisplayValue(current)

      if (progress < 1) {
        requestAnimationFrame(animate)
      }
    }

    const timeout = setTimeout(() => {
      requestAnimationFrame(animate)
    }, delay * 1000)

    return () => clearTimeout(timeout)
  }, [value, delay])

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay, duration: 0.5 }}
      className={cn(
        'relative overflow-hidden rounded-xl border p-5',
        'bg-gradient-to-br',
        colors.bg,
        colors.border,
        colors.glow,
        'transition-all duration-300 hover:scale-[1.02]'
      )}
    >
      {/* Background glow effect */}
      <div
        className={cn(
          'absolute -right-8 -top-8 h-24 w-24 rounded-full opacity-30 blur-2xl',
          color === 'cyan' && 'bg-[#00d4ff]',
          color === 'purple' && 'bg-[#8b5cf6]',
          color === 'pink' && 'bg-[#ec4899]'
        )}
      />

      <div className="relative flex items-center gap-4">
        {/* Icon */}
        <div
          className={cn(
            'flex h-12 w-12 items-center justify-center rounded-xl',
            colors.iconBg
          )}
        >
          <span className={colors.text}>{icon}</span>
        </div>

        {/* Content */}
        <div>
          <p className="text-sm text-gray-400">{title}</p>
          <motion.p
            className={cn('text-3xl font-bold', colors.text)}
            key={displayValue}
          >
            {displayValue.toLocaleString('ru-RU')}
          </motion.p>
        </div>
      </div>
    </motion.div>
  )
})

interface PlatformStatsKPIProps {
  className?: string
}

export const PlatformStatsKPI = memo(function PlatformStatsKPI({
  className,
}: PlatformStatsKPIProps) {
  // TODO: Replace with actual API call using React Query
  // const { data, isLoading } = useQuery({
  //   queryKey: ['platform-stats'],
  //   queryFn: () => apiClient.get('/stats/overview'),
  // })

  // Mock data for now
  const stats: StatsData = {
    totalPlayers: 2847,
    totalTeams: 156,
    totalTournaments: 45,
  }

  return (
    <div className={cn('grid grid-cols-1 gap-4 sm:grid-cols-3', className)}>
      <KPICard
        title="Игроков в базе"
        value={stats.totalPlayers}
        icon={<Users size={24} />}
        color="cyan"
        delay={0}
      />
      <KPICard
        title="Команд"
        value={stats.totalTeams}
        icon={<Building2 size={24} />}
        color="purple"
        delay={0.1}
      />
      <KPICard
        title="Турниров"
        value={stats.totalTournaments}
        icon={<Trophy size={24} />}
        color="pink"
        delay={0.2}
      />
    </div>
  )
})
