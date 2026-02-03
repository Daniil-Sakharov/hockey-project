import { memo, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { cn } from '@/shared/lib/utils'

interface RingData {
  label: string
  value: number
  color: string
}

interface HoloRingChartProps {
  data: RingData[]
  title: string
  centerValue?: string
  centerLabel?: string
}

export const HoloRingChart = memo(function HoloRingChart({
  data,
  title,
  centerValue,
  centerLabel,
}: HoloRingChartProps) {
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null)

  const total = data.reduce((sum, d) => sum + d.value, 0)
  const size = 180
  const strokeWidth = 16
  const hoverStrokeWidth = 20
  const radius = (size - hoverStrokeWidth) / 2
  const circumference = 2 * Math.PI * radius

  // Calculate segment positions
  const segments = data.map((item, index) => {
    const percentage = item.value / total
    const previousPercentage = data
      .slice(0, index)
      .reduce((sum, d) => sum + d.value / total, 0)

    return {
      ...item,
      percentage,
      dashLength: circumference * percentage,
      dashOffset: circumference * previousPercentage,
    }
  })

  return (
    <div className="glass-card rounded-xl p-6">
      <h3 className="mb-4 text-sm font-medium uppercase tracking-wider text-gray-400">
        {title}
      </h3>

      <div className="relative mx-auto" style={{ width: size, height: size }}>
        <svg width={size} height={size} className="-rotate-90">
          <defs>
            {data.map((item, index) => (
              <filter key={`glow-${index}`} id={`segmentGlow-${index}`}>
                <feGaussianBlur stdDeviation="4" result="blur" />
                <feFlood floodColor={item.color} floodOpacity="0.6" />
                <feComposite in2="blur" operator="in" />
                <feMerge>
                  <feMergeNode />
                  <feMergeNode in="SourceGraphic" />
                </feMerge>
              </filter>
            ))}
          </defs>

          {/* Background ring */}
          <circle
            cx={size / 2}
            cy={size / 2}
            r={radius}
            fill="none"
            stroke="rgba(0, 212, 255, 0.1)"
            strokeWidth={strokeWidth}
          />

          {/* Data segments */}
          {segments.map((segment, index) => {
            const isHovered = hoveredIndex === index

            return (
              <motion.circle
                key={segment.label}
                cx={size / 2}
                cy={size / 2}
                r={radius}
                fill="none"
                stroke={segment.color}
                strokeWidth={isHovered ? hoverStrokeWidth : strokeWidth}
                strokeDasharray={`${segment.dashLength} ${circumference}`}
                strokeDashoffset={-segment.dashOffset}
                strokeLinecap="round"
                initial={{ strokeDasharray: `0 ${circumference}` }}
                animate={{
                  strokeDasharray: `${segment.dashLength} ${circumference}`,
                  strokeWidth: isHovered ? hoverStrokeWidth : strokeWidth,
                  filter: isHovered
                    ? `drop-shadow(0 0 12px ${segment.color})`
                    : `drop-shadow(0 0 6px ${segment.color})`,
                }}
                transition={{
                  strokeDasharray: { duration: 1, delay: index * 0.2, ease: 'easeOut' },
                  strokeWidth: { duration: 0.2 },
                  filter: { duration: 0.2 },
                }}
                style={{ cursor: 'pointer' }}
                onMouseEnter={() => setHoveredIndex(index)}
                onMouseLeave={() => setHoveredIndex(null)}
              />
            )
          })}
        </svg>

        {/* Center content */}
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <AnimatePresence mode="wait">
            {hoveredIndex !== null ? (
              <motion.div
                key="hovered"
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.8 }}
                className="text-center"
              >
                <motion.span
                  className="block text-3xl font-bold"
                  style={{ color: segments[hoveredIndex].color }}
                >
                  {segments[hoveredIndex].value}
                </motion.span>
                <span className="text-xs text-gray-400">
                  {segments[hoveredIndex].label}
                </span>
                <span className="mt-1 block text-xs text-gray-500">
                  {((segments[hoveredIndex].percentage) * 100).toFixed(0)}%
                </span>
              </motion.div>
            ) : (
              <motion.div
                key="default"
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.8 }}
                className="text-center"
              >
                {centerValue && (
                  <span className="block text-3xl font-bold text-white">
                    {centerValue}
                  </span>
                )}
                {centerLabel && (
                  <span className="text-xs text-gray-500">{centerLabel}</span>
                )}
              </motion.div>
            )}
          </AnimatePresence>
        </div>

        {/* Rotating glow ring */}
        <motion.div
          className="pointer-events-none absolute inset-3 rounded-full"
          animate={{ rotate: 360 }}
          transition={{ duration: 20, repeat: Infinity, ease: 'linear' }}
          style={{
            background:
              'conic-gradient(from 0deg, transparent 0%, rgba(0, 212, 255, 0.2) 10%, transparent 20%)',
          }}
        />

        {/* Pulse effect on hover */}
        <AnimatePresence>
          {hoveredIndex !== null && (
            <motion.div
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: [0.5, 0, 0.5], scale: [0.9, 1.1, 0.9] }}
              exit={{ opacity: 0 }}
              transition={{ duration: 1.5, repeat: Infinity }}
              className="pointer-events-none absolute inset-0 rounded-full"
              style={{
                border: `2px solid ${segments[hoveredIndex].color}`,
              }}
            />
          )}
        </AnimatePresence>
      </div>

      {/* Legend */}
      <div className="mt-6 flex flex-wrap justify-center gap-4">
        {data.map((item, index) => {
          const isHovered = hoveredIndex === index

          return (
            <motion.div
              key={item.label}
              className={cn(
                'flex cursor-pointer items-center gap-2 rounded-lg px-2 py-1 transition-colors',
                isHovered && 'bg-white/5'
              )}
              onMouseEnter={() => setHoveredIndex(index)}
              onMouseLeave={() => setHoveredIndex(null)}
              whileHover={{ scale: 1.05 }}
            >
              <motion.div
                className="h-3 w-3 rounded-full"
                style={{
                  backgroundColor: item.color,
                }}
                animate={{
                  boxShadow: isHovered
                    ? `0 0 12px ${item.color}`
                    : `0 0 6px ${item.color}`,
                }}
              />
              <span
                className={cn(
                  'text-xs transition-colors',
                  isHovered ? 'text-white' : 'text-gray-400'
                )}
              >
                {item.label}
              </span>
              <span
                className={cn(
                  'text-xs font-medium transition-colors',
                  isHovered ? 'text-white' : 'text-gray-500'
                )}
              >
                {item.value}
              </span>
            </motion.div>
          )
        })}
      </div>

      {/* Bottom glow line */}
      <div className="mt-4 h-px bg-gradient-to-r from-transparent via-[#ec4899]/50 to-transparent" />
    </div>
  )
})
