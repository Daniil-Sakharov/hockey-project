import { memo } from 'react'
import type { ReactNode } from 'react'
import { motion } from 'framer-motion'
import {
  Target,
  Crosshair,
  Calendar,
  Star,
  TrendingUp,
  Clock,
} from 'lucide-react'
import type { PlayerDetailedStats } from '@/entities/player'
import { KPICard } from './KPICard'

interface PlayerKPIGridProps {
  stats: PlayerDetailedStats | null
  isLoading?: boolean
}

interface KPIConfig {
  key: keyof PlayerDetailedStats
  title: string
  icon: ReactNode
  color: 'blue' | 'purple' | 'pink' | 'cyan' | 'green' | 'red' | 'orange'
}

const kpiConfigs: KPIConfig[] = [
  { key: 'goals', title: 'Голы', icon: <Target size={24} strokeWidth={2} />, color: 'blue' },
  { key: 'assists', title: 'Передачи', icon: <Crosshair size={24} strokeWidth={2} />, color: 'purple' },
  { key: 'games', title: 'Игры', icon: <Calendar size={24} strokeWidth={2} />, color: 'cyan' },
  { key: 'points', title: 'Очки', icon: <Star size={24} strokeWidth={2} />, color: 'pink' },
  { key: 'plusMinus', title: '+/-', icon: <TrendingUp size={24} strokeWidth={2} />, color: 'green' },
  { key: 'penaltyMinutes', title: 'Штраф. мин.', icon: <Clock size={24} strokeWidth={2} />, color: 'orange' },
]

export const PlayerKPIGrid = memo(function PlayerKPIGrid({
  stats,
  isLoading = false,
}: PlayerKPIGridProps) {
  return (
    <motion.section
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="mb-8"
    >
      <h2 className="mb-4 text-lg font-semibold text-white">
        Статистика сезона
      </h2>

      <div className="grid grid-cols-2 gap-4 md:grid-cols-3 xl:grid-cols-6">
        {kpiConfigs.map((config, index) => {
          const value = stats?.[config.key] ?? 0
          const color =
            config.key === 'plusMinus'
              ? value >= 0
                ? 'green'
                : 'red'
              : config.color

          return (
            <KPICard
              key={config.key}
              title={config.title}
              value={value}
              icon={config.icon}
              color={color}
              isLoading={isLoading}
              delay={index * 0.1}
              trend={config.key === 'goals' ? 'up' : undefined}
              trendValue={config.key === 'goals' ? '+5' : undefined}
            />
          )
        })}
      </div>
    </motion.section>
  )
})
