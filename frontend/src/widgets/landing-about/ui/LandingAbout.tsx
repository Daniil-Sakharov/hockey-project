import { memo, useEffect, useState } from 'react'
import { motion, useInView } from 'framer-motion'
import { useRef } from 'react'
import { Users, Building2, MapPin, UserCheck, Mail, Send } from 'lucide-react'
import { cn } from '@/shared/lib/utils'

const stats = [
  { icon: Users, value: 15000, suffix: '+', label: 'Игроков' },
  { icon: Building2, value: 500, suffix: '+', label: 'Команд' },
  { icon: MapPin, value: 50, suffix: '+', label: 'Регионов' },
  { icon: UserCheck, value: 200, suffix: '+', label: 'Скаутов' },
]

function AnimatedCounter({
  value,
  suffix = '',
  duration = 2000
}: {
  value: number
  suffix?: string
  duration?: number
}) {
  const [count, setCount] = useState(0)
  const ref = useRef<HTMLSpanElement>(null)
  const isInView = useInView(ref, { once: true })

  useEffect(() => {
    if (!isInView) return

    let startTime: number
    let animationFrame: number

    const animate = (timestamp: number) => {
      if (!startTime) startTime = timestamp
      const progress = Math.min((timestamp - startTime) / duration, 1)

      // Easing function
      const eased = 1 - Math.pow(1 - progress, 3)
      setCount(Math.floor(eased * value))

      if (progress < 1) {
        animationFrame = requestAnimationFrame(animate)
      }
    }

    animationFrame = requestAnimationFrame(animate)

    return () => cancelAnimationFrame(animationFrame)
  }, [isInView, value, duration])

  return (
    <span ref={ref}>
      {count.toLocaleString('ru-RU')}{suffix}
    </span>
  )
}

export const LandingAbout = memo(function LandingAbout() {
  return (
    <section
      id="about"
      className="relative py-24 px-4 overflow-hidden"
    >
      {/* Background */}
      <div className="absolute inset-0 bg-gradient-to-b from-[#050810] via-[#0a0e1a] to-[#050810]" />

      {/* Decorative elements */}
      <div className="absolute bottom-0 left-0 w-96 h-96 bg-[#00d4ff]/5 rounded-full blur-[120px]" />
      <div className="absolute top-0 right-0 w-96 h-96 bg-[#8b5cf6]/5 rounded-full blur-[120px]" />

      <div className="relative max-w-6xl mx-auto">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-16 items-center">
          {/* Left column - Text */}
          <motion.div
            initial={{ opacity: 0, x: -30 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true, margin: '-100px' }}
            transition={{ duration: 0.6 }}
          >
            <div className="flex items-center gap-4 mb-6">
              <img
                src="/logo.png"
                alt="StarRink"
                className="h-16 w-16 object-cover rounded-full"
              />
              <h2 className="text-4xl md:text-5xl font-bold text-white">
                О{' '}
                <span className="bg-gradient-to-r from-[#00d4ff] to-[#8b5cf6] bg-clip-text text-transparent">
                  StarRink
                </span>
              </h2>
            </div>

            <p className="text-gray-400 text-lg leading-relaxed mb-6">
              Мы помогаем молодым хоккеистам показать свой талант миру.
              StarRink объединяет игроков, тренеров и скаутов на одной платформе.
            </p>

            <p className="text-gray-500 leading-relaxed mb-8">
              Наша миссия — сделать путь от дворового катка до профессионального
              хоккея прозрачным и доступным. Каждый игрок заслуживает быть замеченным.
            </p>

            {/* Contacts */}
            <div className="space-y-4">
              <h3 className="text-white font-semibold mb-3">Контакты</h3>

              <a
                href="mailto:support@starrink.ru"
                className="flex items-center gap-3 text-gray-400 hover:text-[#00d4ff] transition-colors group"
              >
                <div className="w-10 h-10 rounded-lg bg-white/5 flex items-center justify-center group-hover:bg-[#00d4ff]/10 transition-colors">
                  <Mail className="w-5 h-5" />
                </div>
                <span>support@starrink.ru</span>
              </a>

              {/* Social links */}
              <div className="flex items-center gap-3 pt-2">
                <a
                  href="https://t.me/starrink"
                  target="_blank"
                  rel="noopener noreferrer"
                  className={cn(
                    'w-10 h-10 rounded-lg flex items-center justify-center',
                    'bg-white/5 text-gray-400',
                    'hover:bg-[#00d4ff]/10 hover:text-[#00d4ff]',
                    'transition-all duration-200'
                  )}
                  aria-label="Telegram"
                >
                  <Send className="w-5 h-5" />
                </a>
                <a
                  href="https://vk.com/starrink"
                  target="_blank"
                  rel="noopener noreferrer"
                  className={cn(
                    'w-10 h-10 rounded-lg flex items-center justify-center',
                    'bg-white/5 text-gray-400',
                    'hover:bg-[#8b5cf6]/10 hover:text-[#8b5cf6]',
                    'transition-all duration-200'
                  )}
                  aria-label="VKontakte"
                >
                  <svg className="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12.785 16.241c-4.932 0-7.748-3.377-7.867-8.991h2.473c.082 4.126 1.9 5.875 3.34 6.234V7.25h2.327v3.559c1.422-.154 2.916-1.787 3.422-3.559h2.327c-.39 2.189-2.015 3.822-3.172 4.489 1.157.541 2.994 1.986 3.696 4.502h-2.57c-.548-1.71-1.916-3.03-3.703-3.208v3.208h-.273z" />
                  </svg>
                </a>
              </div>
            </div>
          </motion.div>

          {/* Right column - Stats */}
          <motion.div
            initial={{ opacity: 0, x: 30 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true, margin: '-100px' }}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="grid grid-cols-2 gap-4"
          >
            {stats.map((stat, index) => {
              const Icon = stat.icon
              return (
                <motion.div
                  key={stat.label}
                  initial={{ opacity: 0, scale: 0.9 }}
                  whileInView={{ opacity: 1, scale: 1 }}
                  viewport={{ once: true }}
                  transition={{ delay: 0.3 + index * 0.1 }}
                  className={cn(
                    'p-6 rounded-2xl text-center',
                    'bg-white/[0.02] backdrop-blur-sm',
                    'border border-white/5',
                    'hover:border-[#00d4ff]/20 transition-colors'
                  )}
                >
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#00d4ff]/20 to-[#8b5cf6]/20 flex items-center justify-center mx-auto mb-4">
                    <Icon className="w-6 h-6 text-[#00d4ff]" />
                  </div>
                  <div className="text-3xl md:text-4xl font-bold text-white mb-1">
                    <AnimatedCounter value={stat.value} suffix={stat.suffix} />
                  </div>
                  <div className="text-gray-500 text-sm">{stat.label}</div>
                </motion.div>
              )
            })}
          </motion.div>
        </div>
      </div>
    </section>
  )
})
