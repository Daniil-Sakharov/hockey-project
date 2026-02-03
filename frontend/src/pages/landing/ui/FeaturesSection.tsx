import { motion } from 'framer-motion'

const features = [
  {
    icon: '📊',
    title: 'Детальная статистика',
    description: 'Полная аналитика по всем игрокам, командам и турнирам российского хоккея',
  },
  {
    icon: '🏆',
    title: 'Рейтинги игроков',
    description: 'Автоматические рейтинги на основе реальных показателей и результатов матчей',
  },
  {
    icon: '📈',
    title: 'Динамика прогресса',
    description: 'Отслеживайте развитие карьеры игроков с интерактивными графиками',
  },
  {
    icon: '🔔',
    title: 'Уведомления',
    description: 'Получайте оповещения о новых матчах и изменениях в рейтингах',
  },
]

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.15 },
  },
}

const itemVariants = {
  hidden: { opacity: 0, y: 40 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.6 } },
}

export function FeaturesSection() {
  return (
    <section className="relative py-24">
      <div className="mx-auto max-w-6xl px-4">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mb-16 text-center"
        >
          <h2 className="mb-4 text-4xl font-bold text-white">
            Возможности <span className="text-gradient">платформы</span>
          </h2>
          <p className="mx-auto max-w-2xl text-gray-400">
            Все инструменты для анализа хоккейной статистики в одном месте
          </p>
        </motion.div>

        <motion.div
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: '-50px' }}
          className="grid gap-6 md:grid-cols-2"
        >
          {features.map((feature) => (
            <motion.div
              key={feature.title}
              variants={itemVariants}
              className="group glass-card rounded-2xl p-8 transition-all hover:border-[#00d4ff]/40"
            >
              <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-xl bg-gradient-to-br from-[#00d4ff]/20 to-[#8b5cf6]/20 text-2xl">
                {feature.icon}
              </div>
              <h3 className="mb-2 text-xl font-semibold text-white">
                {feature.title}
              </h3>
              <p className="text-gray-400">{feature.description}</p>

              {/* Hover glow effect */}
              <div className="absolute inset-0 -z-10 rounded-2xl opacity-0 transition-opacity group-hover:opacity-100">
                <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-[#00d4ff]/5 to-[#8b5cf6]/5" />
              </div>
            </motion.div>
          ))}
        </motion.div>
      </div>

      {/* Background decoration */}
      <div className="pointer-events-none absolute left-1/2 top-0 h-px w-1/2 bg-gradient-to-r from-transparent via-[#00d4ff]/50 to-transparent" />
    </section>
  )
}
