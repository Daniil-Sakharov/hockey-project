import { motion } from 'framer-motion'
import { AuthButtons } from '@/features/auth-buttons'

export function CTASection() {
  return (
    <section className="relative overflow-hidden py-24">
      {/* Background gradient */}
      <div className="absolute inset-0 bg-gradient-to-t from-[#0a0e1a] via-[#0d1224] to-transparent" />

      {/* Animated background lines */}
      <div className="absolute inset-0 opacity-20">
        {[...Array(5)].map((_, i) => (
          <motion.div
            key={i}
            className="absolute h-px w-full bg-gradient-to-r from-transparent via-[#00d4ff] to-transparent"
            style={{ top: `${20 + i * 15}%` }}
            animate={{ x: ['-100%', '100%'] }}
            transition={{
              duration: 8 + i * 2,
              repeat: Infinity,
              ease: 'linear',
              delay: i * 0.5,
            }}
          />
        ))}
      </div>

      <div className="relative mx-auto max-w-4xl px-4 text-center">
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          className="rounded-3xl border border-[#00d4ff]/20 bg-[#0d1224]/80 p-12 backdrop-blur-xl"
        >
          <motion.h2
            className="mb-4 text-4xl font-bold md:text-5xl"
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
          >
            <span className="text-white">Готовы начать?</span>
          </motion.h2>

          <motion.p
            className="mx-auto mb-8 max-w-xl text-lg text-gray-400"
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            viewport={{ once: true }}
            transition={{ delay: 0.2 }}
          >
            Присоединяйтесь к платформе аналитики хоккея нового поколения.
            <br />
            <span className="text-[#00d4ff]">Бесплатная регистрация.</span>
          </motion.p>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.4 }}
            className="flex justify-center"
          >
            <AuthButtons />
          </motion.div>

          {/* Decorative corner accents */}
          <div className="absolute left-4 top-4 h-8 w-8 border-l-2 border-t-2 border-[#00d4ff]/30" />
          <div className="absolute right-4 top-4 h-8 w-8 border-r-2 border-t-2 border-[#00d4ff]/30" />
          <div className="absolute bottom-4 left-4 h-8 w-8 border-b-2 border-l-2 border-[#00d4ff]/30" />
          <div className="absolute bottom-4 right-4 h-8 w-8 border-b-2 border-r-2 border-[#00d4ff]/30" />
        </motion.div>
      </div>
    </section>
  )
}
