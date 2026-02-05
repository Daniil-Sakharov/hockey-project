import { useState, useRef, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { User, Ruler, Weight, Hand } from 'lucide-react'
import type { RosterPlayer } from '@/shared/api/exploreTypes'

interface Props {
  player: RosterPlayer
  teamLogoUrl?: string
  index: number
}

export function RosterPlayerCard({ player, teamLogoUrl, index }: Props) {
  const cardRef = useRef<HTMLDivElement>(null)
  const [tilt, setTilt] = useState({ x: 0, y: 0 })
  const [isHovered, setIsHovered] = useState(false)
  const [photoError, setPhotoError] = useState(false)

  const handleMouseMove = useCallback((e: React.MouseEvent<HTMLDivElement>) => {
    if (!cardRef.current) return
    const rect = cardRef.current.getBoundingClientRect()
    const x = (e.clientX - rect.left) / rect.width - 0.5
    const y = (e.clientY - rect.top) / rect.height - 0.5
    setTilt({ x: y * -10, y: x * 10 })
  }, [])

  const handleMouseLeave = useCallback(() => {
    setTilt({ x: 0, y: 0 })
    setIsHovered(false)
  }, [])

  const positionColor = getPositionColor(player.position)
  const birthYear = player.birthDate ? new Date(player.birthDate).getFullYear() : player.birthYear

  return (
    <Link to={`/explore/players/${player.id}`}>
      <div
        ref={cardRef}
        className="relative group h-full"
        style={{ perspective: '1000px' }}
        onMouseMove={handleMouseMove}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={handleMouseLeave}
      >
        {/* Animated border glow */}
        <motion.div
          className="absolute -inset-[1px] rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"
          style={{
            background: `linear-gradient(${(index * 30) % 360}deg, ${positionColor}, #8b5cf6, ${positionColor})`,
            backgroundSize: '200% 200%',
            filter: 'blur(2px)',
          }}
          animate={isHovered ? { backgroundPosition: ['0% 50%', '100% 50%', '0% 50%'] } : {}}
          transition={{ duration: 3, repeat: Infinity, ease: 'linear' }}
        />

        {/* Card */}
        <div
          className="relative h-full rounded-2xl bg-[#0a0a0f]/95 border border-white/[0.08] backdrop-blur-xl overflow-hidden transition-all duration-300 group-hover:border-white/[0.15]"
          style={{
            transform: `rotateX(${tilt.x}deg) rotateY(${tilt.y}deg)`,
            transition: isHovered ? 'transform 0.1s ease-out' : 'transform 0.4s ease-out',
          }}
        >
          {/* Team logo background */}
          {teamLogoUrl && (
            <div className="absolute inset-0 flex items-center justify-center overflow-hidden">
              <img
                src={teamLogoUrl}
                alt=""
                className="w-48 h-48 object-contain opacity-[0.12] group-hover:opacity-[0.20] transition-opacity duration-500 scale-150 drop-shadow-[0_0_12px_rgba(255,255,255,0.6)]"
              />
            </div>
          )}

          {/* Small team logo badge */}
          {teamLogoUrl && (
            <div className="absolute top-3 right-3 w-8 h-8 rounded-lg bg-white/5 p-1 backdrop-blur-sm border border-white/10">
              <img src={teamLogoUrl} alt="" className="w-full h-full object-contain drop-shadow-[0_0_6px_rgba(255,255,255,0.5)]" />
            </div>
          )}

          {/* Jersey number */}
          {player.jerseyNumber > 0 && (
            <div
              className="absolute top-3 left-3 px-2.5 py-1 rounded-lg font-bold text-lg"
              style={{
                background: `linear-gradient(135deg, ${positionColor}30, ${positionColor}10)`,
                color: positionColor,
                boxShadow: `0 0 12px ${positionColor}20`,
              }}
            >
              #{player.jerseyNumber}
            </div>
          )}

          {/* Content */}
          <div className="relative flex flex-col items-center px-4 pt-14 pb-5">
            {/* Photo */}
            <motion.div
              className="relative mb-4"
              animate={isHovered ? { y: -4, scale: 1.05 } : { y: 0, scale: 1 }}
              transition={{ type: 'spring', stiffness: 300, damping: 20 }}
            >
              <div
                className="w-24 h-24 rounded-full overflow-hidden border-2"
                style={{ borderColor: `${positionColor}40` }}
              >
                {player.photoUrl && !photoError ? (
                  <img
                    src={player.photoUrl}
                    alt={player.name}
                    className="w-full h-full object-cover"
                    onError={() => setPhotoError(true)}
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-white/10 to-white/5">
                    <User size={40} className="text-gray-500" />
                  </div>
                )}
              </div>
              {/* Photo glow */}
              <div
                className="absolute inset-0 rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-500 -z-10 blur-xl"
                style={{ background: positionColor }}
              />
            </motion.div>

            {/* Name */}
            <h3 className="text-base font-bold text-white text-center leading-tight mb-1 line-clamp-2 min-h-[2.5rem] group-hover:text-[#00d4ff] transition-colors">
              {formatName(player.name)}
            </h3>

            {/* Position */}
            {player.position && (
              <span
                className="text-xs font-medium px-2 py-0.5 rounded-full mb-3"
                style={{
                  background: `${positionColor}20`,
                  color: positionColor,
                }}
              >
                {player.position}
              </span>
            )}

            {/* Divider */}
            <div className="w-12 h-px bg-gradient-to-r from-transparent via-white/20 to-transparent mb-3" />

            {/* Stats row */}
            <div className="flex items-center justify-center gap-3 text-xs text-gray-400 flex-wrap">
              {birthYear && (
                <span>{birthYear} г.р.</span>
              )}
              {player.height && player.height > 0 && (
                <span className="flex items-center gap-1">
                  <Ruler size={12} className="text-[#00d4ff]" />
                  {player.height}
                </span>
              )}
              {player.weight && player.weight > 0 && (
                <span className="flex items-center gap-1">
                  <Weight size={12} className="text-[#8b5cf6]" />
                  {player.weight}
                </span>
              )}
              {player.handedness && (
                <span className="flex items-center gap-1">
                  <Hand size={12} className="text-[#f59e0b]" />
                  {player.handedness}
                </span>
              )}
            </div>
          </div>

          {/* Spotlight effect */}
          {isHovered && (
            <div
              className="absolute inset-0 pointer-events-none"
              style={{
                background: `radial-gradient(circle at ${(tilt.y / 10 + 0.5) * 100}% ${(tilt.x / -10 + 0.5) * 100}%, ${positionColor}08 0%, transparent 60%)`,
              }}
            />
          )}
        </div>
      </div>
    </Link>
  )
}

function getPositionColor(position?: string): string {
  if (!position) return '#00d4ff'
  const pos = position.toLowerCase()
  if (pos.includes('врат')) return '#10b981' // green for goalies
  if (pos.includes('защит')) return '#8b5cf6' // purple for defensemen
  if (pos.includes('напад')) return '#00d4ff' // cyan for forwards
  return '#00d4ff'
}

function formatName(name: string): string {
  // "Фамилия Имя Отчество" -> "Имя Фамилия"
  const parts = name.split(' ')
  if (parts.length >= 2) {
    return `${parts[1]} ${parts[0]}`
  }
  return name
}
