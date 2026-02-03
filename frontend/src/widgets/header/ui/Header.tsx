import { motion } from 'framer-motion'
import { Logo } from './Logo'
import { AuthButtons } from '@/features/auth-buttons'

export function Header() {
  return (
    <motion.header
      className="fixed left-0 right-0 top-0 z-50"
      initial={{ y: -100, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.5 }}
    >
      <div className="mx-auto max-w-7xl px-4 py-4">
        <div className="flex items-center justify-between rounded-2xl border border-white/5 bg-[#0a0e1a]/80 px-6 py-3 backdrop-blur-xl">
          <Logo />
          <AuthButtons />
        </div>
      </div>
    </motion.header>
  )
}
