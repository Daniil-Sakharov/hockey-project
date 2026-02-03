import { motion } from 'framer-motion'

interface StatsRevealProps {
  revealProgress: number
  stats: {
    goals: number
    assists: number
    points: number
    games: number
  }
  teamColor: string
}

export function StatsReveal({ revealProgress, stats, teamColor }: StatsRevealProps) {
  return (
    <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
      {/* Left panel */}
      <motion.div
        className="absolute left-[5%] top-1/2 -translate-y-1/2"
        initial={{ opacity: 0, x: -50 }}
        animate={{
          opacity: revealProgress,
          x: revealProgress > 0.2 ? 0 : -50,
        }}
        transition={{ duration: 0.5, ease: 'easeOut' }}
      >
        <StatPanel
          title="ГОЛЫ"
          value={stats.goals}
          subtitle="за сезон"
          color={teamColor}
          delay={0}
          progress={revealProgress}
        />
      </motion.div>

      {/* Right panel */}
      <motion.div
        className="absolute right-[5%] top-1/2 -translate-y-1/2"
        initial={{ opacity: 0, x: 50 }}
        animate={{
          opacity: revealProgress,
          x: revealProgress > 0.2 ? 0 : 50,
        }}
        transition={{ duration: 0.5, ease: 'easeOut', delay: 0.1 }}
      >
        <StatPanel
          title="ПЕРЕДАЧИ"
          value={stats.assists}
          subtitle="за сезон"
          color={teamColor}
          delay={0.1}
          progress={revealProgress}
        />
      </motion.div>

      {/* Bottom panel */}
      <motion.div
        className="absolute bottom-[10%] left-1/2 -translate-x-1/2"
        initial={{ opacity: 0, y: 50 }}
        animate={{
          opacity: revealProgress > 0.3 ? revealProgress : 0,
          y: revealProgress > 0.3 ? 0 : 50,
        }}
        transition={{ duration: 0.5, ease: 'easeOut', delay: 0.2 }}
      >
        <div className="flex gap-8">
          <MiniStat label="ИГРЫ" value={stats.games} />
          <MiniStat label="ОЧКИ" value={stats.points} highlight />
          <MiniStat label="+/-" value="+32" />
        </div>
      </motion.div>

      {/* Top title */}
      <motion.div
        className="absolute top-[10%] left-1/2 -translate-x-1/2"
        initial={{ opacity: 0, y: -30 }}
        animate={{
          opacity: revealProgress > 0.4 ? 1 : 0,
          y: revealProgress > 0.4 ? 0 : -30,
        }}
        transition={{ duration: 0.5, ease: 'easeOut', delay: 0.3 }}
      >
        <h2 className="text-center text-2xl font-bold text-white">
          СТАТИСТИКА СЕЗОНА
        </h2>
        <p className="mt-2 text-center text-gray-400">2023-2024 NHL Regular Season</p>
      </motion.div>
    </div>
  )
}

interface StatPanelProps {
  title: string
  value: number
  subtitle: string
  color: string
  delay: number
  progress: number
}

function StatPanel({ title, value, subtitle, color, delay, progress }: StatPanelProps) {
  return (
    <div
      className="rounded-xl border border-white/10 bg-black/60 p-6 backdrop-blur-md"
      style={{
        boxShadow: `0 0 30px ${color}33`,
      }}
    >
      <p className="text-sm font-medium text-gray-400">{title}</p>
      <motion.p
        className="mt-2 text-5xl font-bold"
        style={{ color }}
        initial={{ scale: 0.5 }}
        animate={{ scale: progress > 0.3 ? 1 : 0.5 }}
        transition={{ duration: 0.4, delay: delay + 0.2, type: 'spring' }}
      >
        {value}
      </motion.p>
      <p className="mt-1 text-sm text-gray-500">{subtitle}</p>
    </div>
  )
}

interface MiniStatProps {
  label: string
  value: number | string
  highlight?: boolean
}

function MiniStat({ label, value, highlight }: MiniStatProps) {
  return (
    <div className="text-center">
      <p className="text-xs font-medium text-gray-500">{label}</p>
      <p
        className={`mt-1 text-2xl font-bold ${
          highlight ? 'text-yellow-400' : 'text-white'
        }`}
      >
        {value}
      </p>
    </div>
  )
}
