import { memo } from 'react'
import { motion } from 'framer-motion'
import {
  BarChart3,
  Trophy,
  Award,
  Eye,
  GitCompare,
  Sparkles,
} from 'lucide-react'
import { cn } from '@/shared/lib/utils'

const features = [
  {
    icon: BarChart3,
    title: 'Полная статистика',
    description: 'Голы, передачи, +/-, все показатели за каждый сезон',
    color: 'cyan',
  },
  {
    icon: Trophy,
    title: 'Рейтинг региона',
    description: 'Узнай своё место среди игроков твоего региона',
    color: 'purple',
  },
  {
    icon: Award,
    title: 'Достижения',
    description: 'Зарабатывай бейджи за успехи на льду',
    color: 'amber',
  },
  {
    icon: Eye,
    title: 'Видимость для скаутов',
    description: 'PRO подписчики видны скаутам и тренерам',
    color: 'purple',
  },
  {
    icon: GitCompare,
    title: 'Сравнение игроков',
    description: 'Сравни свои показатели с любым игроком России',
    color: 'cyan',
  },
  {
    icon: Sparkles,
    title: 'AI рекомендации',
    description: 'Персональные советы по улучшению игры',
    color: 'amber',
  },
]

const colorStyles = {
  cyan: {
    gradient: 'from-[#00d4ff] to-[#0099cc]',
    glow: 'group-hover:shadow-[0_0_30px_rgba(0,212,255,0.3)]',
    border: 'group-hover:border-[#00d4ff]/30',
  },
  purple: {
    gradient: 'from-[#8b5cf6] to-[#6d28d9]',
    glow: 'group-hover:shadow-[0_0_30px_rgba(139,92,246,0.3)]',
    border: 'group-hover:border-[#8b5cf6]/30',
  },
  amber: {
    gradient: 'from-[#f59e0b] to-[#d97706]',
    glow: 'group-hover:shadow-[0_0_30px_rgba(245,158,11,0.3)]',
    border: 'group-hover:border-[#f59e0b]/30',
  },
}

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
    },
  },
}

const itemVariants = {
  hidden: { opacity: 0, y: 30 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      type: 'spring',
      damping: 20,
      stiffness: 100,
    },
  },
}

export const LandingFeatures = memo(function LandingFeatures() {
  return (
    <section
      id="features"
      className="relative py-24 px-4 overflow-hidden"
    >
      {/* Background gradient */}
      <div className="absolute inset-0 bg-gradient-to-b from-[#050810] via-[#0a0e1a] to-[#050810]" />

      {/* Decorative elements */}
      <div className="absolute top-1/4 left-0 w-72 h-72 bg-[#00d4ff]/5 rounded-full blur-[100px]" />
      <div className="absolute bottom-1/4 right-0 w-72 h-72 bg-[#8b5cf6]/5 rounded-full blur-[100px]" />

      <div className="relative max-w-6xl mx-auto">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-100px' }}
          transition={{ duration: 0.6 }}
          className="text-center mb-16"
        >
          <h2 className="text-4xl md:text-5xl font-bold text-white mb-4">
            Возможности для{' '}
            <span className="bg-gradient-to-r from-[#00d4ff] to-[#8b5cf6] bg-clip-text text-transparent">
              игроков
            </span>
          </h2>
          <p className="text-gray-400 text-lg max-w-2xl mx-auto">
            Все инструменты для отслеживания прогресса и развития карьеры в одном месте
          </p>
        </motion.div>

        {/* Features grid */}
        <motion.div
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: '-50px' }}
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
        >
          {features.map((feature) => {
            const styles = colorStyles[feature.color as keyof typeof colorStyles]
            const Icon = feature.icon

            return (
              <motion.div
                key={feature.title}
                variants={itemVariants}
                className={cn(
                  'group relative p-6 rounded-2xl',
                  'bg-white/[0.02] backdrop-blur-sm',
                  'border border-white/5',
                  'transition-all duration-300',
                  styles.glow,
                  styles.border
                )}
              >
                {/* Icon */}
                <div
                  className={cn(
                    'w-12 h-12 rounded-xl flex items-center justify-center mb-4',
                    'bg-gradient-to-br',
                    styles.gradient
                  )}
                >
                  <Icon className="w-6 h-6 text-white" />
                </div>

                {/* Content */}
                <h3 className="text-xl font-semibold text-white mb-2">
                  {feature.title}
                </h3>
                <p className="text-gray-400 text-sm leading-relaxed">
                  {feature.description}
                </p>

                {/* Hover gradient overlay */}
                <div
                  className={cn(
                    'absolute inset-0 rounded-2xl opacity-0 group-hover:opacity-100',
                    'bg-gradient-to-br pointer-events-none transition-opacity duration-300',
                    styles.gradient,
                    'opacity-0 group-hover:opacity-[0.03]'
                  )}
                />
              </motion.div>
            )
          })}
        </motion.div>
      </div>
    </section>
  )
})
