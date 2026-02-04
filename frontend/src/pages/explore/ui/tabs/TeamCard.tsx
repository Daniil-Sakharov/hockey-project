import { useRef, useState, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Shield, Users, MapPin } from 'lucide-react'
import type { TeamItem } from '@/shared/api/exploreTypes'

function TeamLogo({ url, name }: { url?: string; name: string }) {
  const [hasError, setHasError] = useState(false)

  if (!url || hasError) {
    return (
      <div className="flex h-20 w-20 items-center justify-center rounded-2xl bg-gradient-to-br from-[#00d4ff]/10 to-[#8b5cf6]/10 border border-white/10">
        <Shield size={36} className="text-gray-500" />
      </div>
    )
  }

  return (
    <img
      src={url}
      alt={name}
      className="h-20 w-20 object-contain drop-shadow-[0_0_12px_rgba(0,212,255,0.3)]"
      onError={() => setHasError(true)}
    />
  )
}

interface Props {
  team: TeamItem
  index: number
}

export function TeamCard({ team, index }: Props) {
  const cardRef = useRef<HTMLDivElement>(null)
  const [tilt, setTilt] = useState({ x: 0, y: 0 })
  const [isHovered, setIsHovered] = useState(false)

  const handleMouseMove = useCallback((e: React.MouseEvent<HTMLDivElement>) => {
    if (!cardRef.current) return
    const rect = cardRef.current.getBoundingClientRect()
    const x = (e.clientX - rect.left) / rect.width - 0.5
    const y = (e.clientY - rect.top) / rect.height - 0.5
    setTilt({ x: y * -12, y: x * 12 })
  }, [])

  const handleMouseLeave = useCallback(() => {
    setTilt({ x: 0, y: 0 })
    setIsHovered(false)
  }, [])

  const gradientAngle = (index * 40) % 360

  return (
    <Link to={`/explore/teams/${team.id}`}>
      <div
        ref={cardRef}
        className="relative group"
        style={{ perspective: '800px' }}
        onMouseMove={handleMouseMove}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={handleMouseLeave}
      >
        {/* Animated gradient border */}
        <div
          className="absolute -inset-[1px] rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500 blur-[1px]"
          style={{
            background: `linear-gradient(${gradientAngle}deg, #00d4ff, #8b5cf6, #00d4ff)`,
            backgroundSize: '200% 200%',
            animation: isHovered ? 'gradientShift 3s ease infinite' : 'none',
          }}
        />

        <div
          className="relative rounded-2xl bg-[#0a0a0f]/90 border border-white/[0.08] backdrop-blur-xl overflow-hidden transition-shadow duration-300 group-hover:shadow-[0_0_40px_rgba(0,212,255,0.1)]"
          style={{
            transform: `rotateX(${tilt.x}deg) rotateY(${tilt.y}deg)`,
            transition: isHovered ? 'transform 0.1s ease-out' : 'transform 0.4s ease-out',
          }}
        >
          {/* Top gradient accent */}
          <div
            className="h-1 w-full"
            style={{ background: `linear-gradient(90deg, #00d4ff, #8b5cf6, #00d4ff)` }}
          />

          <div className="flex flex-col items-center px-6 pt-7 pb-6">
            {/* Logo */}
            <motion.div
              animate={isHovered ? { y: -4, scale: 1.08 } : { y: 0, scale: 1 }}
              transition={{ type: 'spring' as const, stiffness: 300, damping: 20 }}
              className="mb-4"
            >
              <TeamLogo url={team.logoUrl} name={team.name} />
            </motion.div>

            {/* Team name */}
            <h3 className="text-lg font-bold text-white text-center leading-tight group-hover:text-[#00d4ff] transition-colors duration-300 line-clamp-2 min-h-[3.5rem]">
              {team.name}
            </h3>

            {/* City */}
            {team.city && (
              <div className="flex items-center gap-1.5 mt-2 text-gray-400 text-sm">
                <MapPin size={14} className="shrink-0 text-[#8b5cf6]" />
                <span>{team.city}</span>
              </div>
            )}

            {/* Divider */}
            <div className="w-12 h-px bg-gradient-to-r from-transparent via-white/20 to-transparent mt-4 mb-4" />

            {/* Players count */}
            <div className="flex items-center gap-2 px-4 py-1.5 rounded-full bg-[#00d4ff]/10 border border-[#00d4ff]/20">
              <Users size={16} className="text-[#00d4ff]" />
              <span className="text-sm font-semibold text-[#00d4ff]">{team.playersCount}</span>
              <span className="text-xs text-gray-400">игроков</span>
            </div>
          </div>

          {/* Hover spotlight */}
          {isHovered && (
            <div
              className="absolute inset-0 pointer-events-none"
              style={{
                background: `radial-gradient(circle at ${(tilt.y / 12 + 0.5) * 100}% ${(tilt.x / -12 + 0.5) * 100}%, rgba(0,212,255,0.06) 0%, transparent 60%)`,
              }}
            />
          )}
        </div>
      </div>
    </Link>
  )
}
