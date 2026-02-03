import { motion } from 'framer-motion'

interface AuthButtonsProps {
  onLogin?: () => void
  onRegister?: () => void
}

export function AuthButtons({ onLogin, onRegister }: AuthButtonsProps) {
  return (
    <motion.div
      className="flex items-center gap-4"
      initial={{ opacity: 0, y: -20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, delay: 0.6 }}
    >
      <button
        onClick={onLogin}
        className="group relative overflow-hidden rounded-lg border border-[#00d4ff]/30 bg-transparent px-6 py-2.5 text-sm font-medium text-[#00d4ff] transition-all hover:border-[#00d4ff] hover:shadow-[0_0_20px_rgba(0,212,255,0.3)]"
      >
        <span className="relative z-10">Войти</span>
        <div className="absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-[#00d4ff]/10 to-transparent transition-transform group-hover:translate-x-full" />
      </button>

      <button
        onClick={onRegister}
        className="group relative overflow-hidden rounded-lg bg-gradient-to-r from-[#00d4ff] to-[#0066ff] px-6 py-2.5 text-sm font-medium text-white shadow-[0_0_20px_rgba(0,212,255,0.3)] transition-all hover:shadow-[0_0_30px_rgba(0,212,255,0.5)]"
      >
        <span className="relative z-10">Регистрация</span>
        <motion.div
          className="absolute inset-0 bg-gradient-to-r from-[#00d4ff] to-[#8b5cf6]"
          initial={{ opacity: 0 }}
          whileHover={{ opacity: 1 }}
          transition={{ duration: 0.3 }}
        />
      </button>
    </motion.div>
  )
}
