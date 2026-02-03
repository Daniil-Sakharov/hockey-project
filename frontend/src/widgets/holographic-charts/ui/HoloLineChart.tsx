import { memo, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { cn } from '@/shared/lib/utils'

interface DataPoint {
  label: string
  value: number
}

interface HoloLineChartProps {
  data: DataPoint[]
  title: string
}

export const HoloLineChart = memo(function HoloLineChart({ data, title }: HoloLineChartProps) {
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null)

  const maxValue = Math.max(...data.map((d) => d.value))
  const minValue = Math.min(...data.map((d) => d.value))
  const range = maxValue - minValue || 1

  const width = 280
  const height = 140
  const padding = 25

  const points = data.map((d, i) => ({
    x: padding + (i / (data.length - 1)) * (width - padding * 2),
    y: height - padding - ((d.value - minValue) / range) * (height - padding * 2),
    value: d.value,
    label: d.label,
  }))

  const pathD = points.reduce((path, point, i) => {
    return path + (i === 0 ? `M ${point.x},${point.y}` : ` L ${point.x},${point.y}`)
  }, '')

  const areaD = pathD + ` L ${points[points.length - 1].x},${height - padding} L ${padding},${height - padding} Z`

  return (
    <div className="glass-card rounded-xl p-6">
      <h3 className="mb-4 text-sm font-medium uppercase tracking-wider text-gray-400">
        {title}
      </h3>

      <div className="relative">
        <svg width={width} height={height} className="overflow-visible">
          <defs>
            <linearGradient id="lineGradient" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stopColor="#00d4ff" />
              <stop offset="50%" stopColor="#8b5cf6" />
              <stop offset="100%" stopColor="#ec4899" />
            </linearGradient>
            <linearGradient id="areaGradient" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stopColor="rgba(0, 212, 255, 0.4)" />
              <stop offset="50%" stopColor="rgba(139, 92, 246, 0.2)" />
              <stop offset="100%" stopColor="rgba(0, 212, 255, 0)" />
            </linearGradient>
            <filter id="lineGlow">
              <feGaussianBlur stdDeviation="4" result="coloredBlur" />
              <feMerge>
                <feMergeNode in="coloredBlur" />
                <feMergeNode in="SourceGraphic" />
              </feMerge>
            </filter>
            <filter id="pointGlow">
              <feGaussianBlur stdDeviation="3" result="blur" />
              <feMerge>
                <feMergeNode in="blur" />
                <feMergeNode in="SourceGraphic" />
              </feMerge>
            </filter>
          </defs>

          {/* Grid lines */}
          {[0, 1, 2, 3, 4].map((i) => (
            <motion.line
              key={i}
              x1={padding}
              y1={padding + (i * (height - padding * 2)) / 4}
              x2={width - padding}
              y2={padding + (i * (height - padding * 2)) / 4}
              stroke="rgba(0, 212, 255, 0.1)"
              strokeDasharray="4 4"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: i * 0.1 }}
            />
          ))}

          {/* Vertical hover line */}
          <AnimatePresence>
            {hoveredIndex !== null && (
              <motion.line
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                x1={points[hoveredIndex].x}
                y1={padding}
                x2={points[hoveredIndex].x}
                y2={height - padding}
                stroke="#00d4ff"
                strokeWidth={1}
                strokeDasharray="4 4"
                style={{ filter: 'drop-shadow(0 0 4px #00d4ff)' }}
              />
            )}
          </AnimatePresence>

          {/* Area fill */}
          <motion.path
            d={areaD}
            fill="url(#areaGradient)"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 1 }}
          />

          {/* Line */}
          <motion.path
            d={pathD}
            fill="none"
            stroke="url(#lineGradient)"
            strokeWidth={3}
            strokeLinecap="round"
            strokeLinejoin="round"
            filter="url(#lineGlow)"
            initial={{ pathLength: 0 }}
            animate={{ pathLength: 1 }}
            transition={{ duration: 1.5, ease: 'easeInOut' }}
          />

          {/* Data points with hover areas */}
          {points.map((point, i) => {
            const isHovered = hoveredIndex === i

            return (
              <g key={i}>
                {/* Invisible hover area */}
                <circle
                  cx={point.x}
                  cy={point.y}
                  r={20}
                  fill="transparent"
                  style={{ cursor: 'pointer' }}
                  onMouseEnter={() => setHoveredIndex(i)}
                  onMouseLeave={() => setHoveredIndex(null)}
                />

                {/* Outer glow ring */}
                <motion.circle
                  cx={point.x}
                  cy={point.y}
                  fill="rgba(0, 212, 255, 0.15)"
                  initial={{ r: 0 }}
                  animate={{
                    r: isHovered ? 16 : 8,
                    scale: isHovered ? 1 : [1, 1.3, 1],
                  }}
                  transition={{
                    r: { duration: 0.2 },
                    scale: { duration: 2, repeat: Infinity, delay: i * 0.2 }
                  }}
                />

                {/* Middle ring */}
                <motion.circle
                  cx={point.x}
                  cy={point.y}
                  fill="rgba(0, 212, 255, 0.3)"
                  initial={{ r: 0 }}
                  animate={{ r: isHovered ? 10 : 5 }}
                  transition={{ duration: 0.2 }}
                />

                {/* Inner dot */}
                <motion.circle
                  cx={point.x}
                  cy={point.y}
                  fill={isHovered ? '#00ffff' : '#00d4ff'}
                  filter="url(#pointGlow)"
                  initial={{ r: 0 }}
                  animate={{ r: isHovered ? 6 : 4 }}
                  transition={{ duration: 0.2 }}
                />
              </g>
            )
          })}
        </svg>

        {/* Tooltip */}
        <AnimatePresence>
          {hoveredIndex !== null && (
            <motion.div
              initial={{ opacity: 0, y: 10, scale: 0.9 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: 10, scale: 0.9 }}
              className={cn(
                'absolute z-10 -translate-x-1/2',
                'rounded-lg bg-[#0d1224] px-3 py-2',
                'border border-[#00d4ff]/30',
                'shadow-[0_0_20px_rgba(0,212,255,0.3)]',
                'pointer-events-none'
              )}
              style={{
                left: points[hoveredIndex].x,
                top: points[hoveredIndex].y - 50,
              }}
            >
              <div className="text-center whitespace-nowrap">
                <div className="text-lg font-bold text-[#00d4ff]">
                  {points[hoveredIndex].value}
                </div>
                <div className="text-xs text-gray-400">
                  {points[hoveredIndex].label}
                </div>
              </div>
              <div className="absolute -bottom-1.5 left-1/2 h-3 w-3 -translate-x-1/2 rotate-45 border-b border-r border-[#00d4ff]/30 bg-[#0d1224]" />
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* X-axis labels */}
      <div className="mt-3 flex justify-between px-5">
        {data.map((d, i) => (
          <motion.span
            key={i}
            className={cn(
              'text-xs transition-colors duration-200',
              hoveredIndex === i ? 'text-[#00d4ff]' : 'text-gray-500'
            )}
          >
            {d.label}
          </motion.span>
        ))}
      </div>

      {/* Bottom glow line */}
      <div className="mt-3 h-px bg-gradient-to-r from-transparent via-[#8b5cf6]/50 to-transparent" />
    </div>
  )
})
