import { motion } from 'framer-motion'
import { HoloBarChart } from './HoloBarChart'
import { HoloRingChart } from './HoloRingChart'
import { HoloLineChart } from './HoloLineChart'

const barData = [
  { label: 'Янв', value: 42 },
  { label: 'Фев', value: 58 },
  { label: 'Мар', value: 35 },
  { label: 'Апр', value: 67 },
  { label: 'Май', value: 89 },
  { label: 'Июн', value: 54 },
]

const ringData = [
  { label: 'Голы', value: 45, color: '#00d4ff' },
  { label: 'Ассисты', value: 32, color: '#8b5cf6' },
  { label: 'Штрафы', value: 23, color: '#ec4899' },
]

const lineData = [
  { label: 'П1', value: 120 },
  { label: 'П2', value: 180 },
  { label: 'П3', value: 150 },
  { label: 'П4', value: 220 },
  { label: 'П5', value: 190 },
  { label: 'П6', value: 250 },
]

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.2 },
  },
}

const itemVariants = {
  hidden: { opacity: 0, y: 30 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.6 } },
}

export function HolographicCharts() {
  return (
    <section className="relative py-24">
      {/* Background gradient */}
      <div className="gradient-mesh absolute inset-0" />

      <div className="relative mx-auto max-w-6xl px-4">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mb-12 text-center"
        >
          <h2 className="mb-4 text-4xl font-bold">
            <span className="text-gradient">Аналитика</span>
            <span className="text-white"> в реальном времени</span>
          </h2>
          <p className="text-gray-400">
            Современные графики и визуализация данных
          </p>
        </motion.div>

        <motion.div
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: '-100px' }}
          className="grid gap-6 md:grid-cols-2 lg:grid-cols-3"
        >
          <motion.div variants={itemVariants}>
            <HoloBarChart data={barData} title="Результативность по месяцам" />
          </motion.div>

          <motion.div variants={itemVariants}>
            <HoloRingChart
              data={ringData}
              title="Распределение статистики"
              centerValue="100"
              centerLabel="Всего"
            />
          </motion.div>

          <motion.div variants={itemVariants}>
            <HoloLineChart data={lineData} title="Динамика показателей" />
          </motion.div>
        </motion.div>

        {/* Decorative elements */}
        <div className="pointer-events-none absolute -left-20 top-1/2 h-64 w-64 -translate-y-1/2 rounded-full bg-[#00d4ff]/10 blur-3xl" />
        <div className="pointer-events-none absolute -right-20 top-1/3 h-64 w-64 rounded-full bg-[#8b5cf6]/10 blur-3xl" />
      </div>
    </section>
  )
}
