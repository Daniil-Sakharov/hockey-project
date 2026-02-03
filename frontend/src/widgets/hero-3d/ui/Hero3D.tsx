import { motion } from 'framer-motion'
import { HeroScene } from './HeroScene'

export function Hero3D() {
  return (
    <section className="relative h-screen w-full overflow-hidden">
      {/* 3D Scene Background */}
      <HeroScene />

      {/* Gradient Overlay */}
      <div className="absolute inset-0 z-10 bg-gradient-to-t from-[#0a0e1a] via-transparent to-[#0a0e1a]/50" />

      {/* Content Overlay */}
      <div className="relative z-20 flex h-full flex-col items-center justify-center px-4">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, ease: 'easeOut' }}
          className="text-center"
        >
          <motion.h1
            className="mb-4 text-6xl font-bold tracking-tight md:text-8xl"
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.6, delay: 0.2 }}
          >
            <span className="text-gradient">Hockey</span>
            <span className="text-white">Stats</span>
          </motion.h1>

          <motion.p
            className="mx-auto mb-8 max-w-2xl text-lg text-gray-400 md:text-xl"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.6, delay: 0.4 }}
          >
            Аналитика российского хоккея нового поколения.
            <br />
            <span className="neon-text">Статистика • Рейтинги • Прогнозы</span>
          </motion.p>
        </motion.div>

        {/* Scroll Indicator */}
        <motion.div
          className="absolute bottom-8 left-1/2 -translate-x-1/2"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1.2 }}
        >
          <motion.div
            animate={{ y: [0, 10, 0] }}
            transition={{ duration: 1.5, repeat: Infinity }}
            className="flex flex-col items-center"
          >
            <span className="mb-2 text-sm text-gray-500">Scroll</span>
            <div className="h-12 w-6 rounded-full border-2 border-[#00d4ff]/30 p-1">
              <motion.div
                animate={{ y: [0, 16, 0] }}
                transition={{ duration: 1.5, repeat: Infinity }}
                className="h-2 w-2 rounded-full bg-[#00d4ff]"
              />
            </div>
          </motion.div>
        </motion.div>
      </div>
    </section>
  )
}
