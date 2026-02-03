import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'

export function Logo() {
  return (
    <Link to="/" className="group flex items-center gap-3">
      {/* Animated icon */}
      <motion.div
        className="relative flex h-10 w-10 items-center justify-center"
        whileHover={{ scale: 1.1 }}
      >
        <div className="absolute inset-0 rounded-lg bg-gradient-to-br from-[#00d4ff] to-[#0066ff] opacity-20" />
        <div className="absolute inset-0 rounded-lg border border-[#00d4ff]/30" />
        <span className="relative text-xl">🏒</span>
      </motion.div>

      {/* Text */}
      <div className="flex flex-col">
        <span className="text-lg font-bold leading-tight">
          <span className="text-gradient">Hockey</span>
          <span className="text-white">Stats</span>
        </span>
        <span className="text-[10px] uppercase tracking-widest text-gray-500">
          Analytics
        </span>
      </div>
    </Link>
  )
}
