import { memo, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { cn } from '@/shared/lib/utils'

interface BarData {
  label: string
  value: number
  maxValue?: number
}

interface HoloBarChartProps {
  data: BarData[]
  title: string
}

export const HoloBarChart = memo(function HoloBarChart({ data, title }: HoloBarChartProps) {
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null)
  const maxValue = Math.max(...data.map((d) => d.maxValue || d.value))

  return (
    <div className="glass-card rounded-xl p-6">
      <h3 className="mb-4 text-sm font-medium uppercase tracking-wider text-gray-400">
        {title}
      </h3>

      <div className="relative flex items-end justify-between gap-2" style={{ height: 160 }}>
        {data.map((item, index) => {
          const height = (item.value / maxValue) * 100
          const isHovered = hoveredIndex === index

          return (
            <div
              key={item.label}
              className="relative flex flex-1 flex-col items-center"
              onMouseEnter={() => setHoveredIndex(index)}
              onMouseLeave={() => setHoveredIndex(null)}
            >
              {/* Tooltip */}
              <AnimatePresence>
                {isHovered && (
                  <motion.div
                    initial={{ opacity: 0, y: 10, scale: 0.9 }}
                    animate={{ opacity: 1, y: 0, scale: 1 }}
                    exit={{ opacity: 0, y: 10, scale: 0.9 }}
                    className={cn(
                      'absolute -top-16 left-1/2 z-10 -translate-x-1/2',
                      'rounded-lg bg-[#0d1224] px-3 py-2',
                      'border border-[#00d4ff]/30',
                      'shadow-[0_0_20px_rgba(0,212,255,0.3)]'
                    )}
                  >
                    <div className="text-center">
                      <div className="text-lg font-bold text-[#00d4ff]">
                        {item.value}
                      </div>
                      <div className="text-xs text-gray-400">{item.label}</div>
                    </div>
                    {/* Arrow */}
                    <div className="absolute -bottom-1.5 left-1/2 h-3 w-3 -translate-x-1/2 rotate-45 border-b border-r border-[#00d4ff]/30 bg-[#0d1224]" />
                  </motion.div>
                )}
              </AnimatePresence>

              {/* Bar container */}
              <motion.div
                className="relative w-full max-w-[50px] cursor-pointer"
                style={{ height: `${height}%` }}
                initial={{ height: 0 }}
                animate={{ height: `${height}%` }}
                transition={{ duration: 0.8, delay: index * 0.1, ease: 'easeOut' }}
                whileHover={{ scale: 1.05 }}
              >
                {/* Bar background glow */}
                <motion.div
                  className="absolute inset-0 rounded-t-lg"
                  animate={{
                    boxShadow: isHovered
                      ? '0 0 30px rgba(0, 212, 255, 0.6), 0 0 60px rgba(0, 212, 255, 0.3)'
                      : '0 0 20px rgba(0, 212, 255, 0.3)',
                  }}
                  transition={{ duration: 0.3 }}
                />

                {/* Bar */}
                <motion.div
                  className="absolute inset-0 rounded-t-lg"
                  style={{
                    background: isHovered
                      ? 'linear-gradient(180deg, #00ffff 0%, #00d4ff 50%, #0066ff 100%)'
                      : 'linear-gradient(180deg, #00d4ff 0%, #0066ff 100%)',
                  }}
                  animate={{
                    filter: isHovered ? 'brightness(1.2)' : 'brightness(1)',
                  }}
                  transition={{ duration: 0.3 }}
                />

                {/* Glow overlay */}
                <motion.div
                  className="absolute inset-0 rounded-t-lg"
                  animate={{
                    opacity: isHovered ? [0.5, 0.8, 0.5] : [0.3, 0.5, 0.3],
                  }}
                  transition={{ duration: 1.5, repeat: Infinity }}
                  style={{
                    background:
                      'linear-gradient(180deg, rgba(255,255,255,0.4) 0%, transparent 50%)',
                  }}
                />

                {/* Scanline effect */}
                {isHovered && (
                  <motion.div
                    className="absolute inset-x-0 h-1 bg-white/50"
                    initial={{ top: '100%' }}
                    animate={{ top: '0%' }}
                    transition={{ duration: 0.8, repeat: Infinity }}
                  />
                )}
              </motion.div>

              {/* Label */}
              <motion.span
                className={cn(
                  'mt-3 text-xs transition-colors duration-200',
                  isHovered ? 'text-[#00d4ff]' : 'text-gray-500'
                )}
              >
                {item.label}
              </motion.span>
            </div>
          )
        })}
      </div>

      {/* Bottom glow line */}
      <div className="mt-4 h-px bg-gradient-to-r from-transparent via-[#00d4ff]/50 to-transparent" />
    </div>
  )
})
