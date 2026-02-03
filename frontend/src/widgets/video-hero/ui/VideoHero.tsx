import { useRef, useEffect } from 'react'
import { motion } from 'framer-motion'
import { AuthButtons } from '@/features/auth-buttons'

export function VideoHero() {
  const videoRef = useRef<HTMLVideoElement>(null)

  useEffect(() => {
    if (videoRef.current) {
      videoRef.current.playbackRate = 0.75 // Slight slow motion
    }
  }, [])

  return (
    <section className="relative h-screen w-full overflow-hidden">
      {/* Video Background */}
      <div className="absolute inset-0 z-0">
        <video
          ref={videoRef}
          autoPlay
          loop
          muted
          playsInline
          className="h-full w-full object-cover"
          poster="/videos/hockey-poster.jpg"
        >
          <source src="/videos/hockey-background.mp4" type="video/mp4" />
        </video>
      </div>

      {/* Dark Overlay */}
      <div className="absolute inset-0 z-10 bg-black/60" />

      {/* Gradient Overlays */}
      <div className="absolute inset-0 z-10 bg-gradient-to-t from-[#0a0e1a] via-transparent to-[#0a0e1a]/70" />
      <div className="absolute inset-0 z-10 bg-gradient-to-r from-[#0a0e1a]/50 via-transparent to-[#0a0e1a]/50" />

      {/* Scan Lines Effect */}
      <div
        className="pointer-events-none absolute inset-0 z-20 opacity-30"
        style={{
          background:
            'repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(0, 212, 255, 0.03) 2px, rgba(0, 212, 255, 0.03) 4px)',
        }}
      />

      {/* Animated Grid Overlay */}
      <div
        className="pointer-events-none absolute inset-0 z-20 opacity-10"
        style={{
          backgroundImage: `
            linear-gradient(rgba(0, 212, 255, 0.1) 1px, transparent 1px),
            linear-gradient(90deg, rgba(0, 212, 255, 0.1) 1px, transparent 1px)
          `,
          backgroundSize: '50px 50px',
        }}
      />

      {/* Content */}
      <div className="relative z-30 flex h-full flex-col items-center justify-center px-4">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, ease: 'easeOut' }}
          className="text-center"
        >
          {/* Badge */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="mb-6 inline-block"
          >
            <span className="rounded-full border border-[#00d4ff]/30 bg-[#00d4ff]/10 px-4 py-2 text-sm text-[#00d4ff]">
              🏒 Платформа аналитики хоккея
            </span>
          </motion.div>

          {/* Title */}
          <motion.h1
            className="mb-6 text-5xl font-bold tracking-tight md:text-7xl lg:text-8xl"
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.6, delay: 0.3 }}
          >
            <span className="text-gradient">Hockey</span>
            <span className="text-white">Stats</span>
          </motion.h1>

          {/* Subtitle */}
          <motion.p
            className="mx-auto mb-8 max-w-2xl text-lg text-gray-300 md:text-xl"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.6, delay: 0.5 }}
          >
            Аналитика российского хоккея нового поколения
          </motion.p>

          {/* Stats Row */}
          <motion.div
            className="mb-10 flex flex-wrap justify-center gap-8"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.7 }}
          >
            {[
              { value: '2,847', label: 'Игроков' },
              { value: '156', label: 'Команд' },
              { value: '12K+', label: 'Матчей' },
            ].map((stat) => (
              <div key={stat.label} className="text-center">
                <div className="text-3xl font-bold text-[#00d4ff]">{stat.value}</div>
                <div className="text-sm text-gray-500">{stat.label}</div>
              </div>
            ))}
          </motion.div>

          {/* Auth Buttons */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.9 }}
          >
            <AuthButtons />
          </motion.div>
        </motion.div>

        {/* Scroll Indicator */}
        <motion.div
          className="absolute bottom-8 left-1/2 -translate-x-1/2"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1.5 }}
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

      {/* Corner Decorations */}
      <div className="absolute left-4 top-4 z-30 h-16 w-16 border-l-2 border-t-2 border-[#00d4ff]/30" />
      <div className="absolute right-4 top-4 z-30 h-16 w-16 border-r-2 border-t-2 border-[#00d4ff]/30" />
      <div className="absolute bottom-4 left-4 z-30 h-16 w-16 border-b-2 border-l-2 border-[#00d4ff]/30" />
      <div className="absolute bottom-4 right-4 z-30 h-16 w-16 border-b-2 border-r-2 border-[#00d4ff]/30" />
    </section>
  )
}
