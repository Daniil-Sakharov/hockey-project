import { useState, useRef, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { User, TrendingUp } from 'lucide-react'
import type { PlayerItem } from '@/shared/api/exploreTypes'
import { cn } from '@/shared/lib/utils'

const POSITION_LABELS: Record<string, string> = {
  forward: 'Нападающий',
  defender: 'Защитник',
  goalie: 'Вратарь',
}

interface Props {
  player: PlayerItem
  index: number
  linkParams?: string
}

export function PlayerSearchCard({ player, index, linkParams }: Props) {
  const cardRef = useRef<HTMLDivElement>(null)
  const [tilt, setTilt] = useState({ x: 0, y: 0 })
  const [isHovered, setIsHovered] = useState(false)
  const [photoError, setPhotoError] = useState(false)

  const handleMouseMove = useCallback((e: React.MouseEvent<HTMLDivElement>) => {
    if (!cardRef.current) return
    const rect = cardRef.current.getBoundingClientRect()
    const x = (e.clientX - rect.left) / rect.width - 0.5
    const y = (e.clientY - rect.top) / rect.height - 0.5
    setTilt({ x: y * -8, y: x * 8 })
  }, [])

  const handleMouseLeave = useCallback(() => {
    setTilt({ x: 0, y: 0 })
    setIsHovered(false)
  }, [])

  const positionColor = getPositionColor(player.position)
  const profileUrl = `/explore/players/${player.id}${linkParams ? `?${linkParams}` : ''}`

  return (
    <Link to={profileUrl}>
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
            background: `linear-gradient(${(index * 25) % 360}deg, ${positionColor}, #8b5cf6, ${positionColor})`,
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
          {player.teamLogoUrl && (
            <div className="absolute inset-0 flex items-center justify-center overflow-hidden">
              <img
                src={player.teamLogoUrl}
                alt=""
                className="w-40 h-40 object-contain opacity-[0.08] group-hover:opacity-[0.15] transition-opacity duration-500 scale-125 drop-shadow-[0_0_12px_rgba(255,255,255,0.6)]"
              />
            </div>
          )}

          {/* Small team logo badge */}
          {player.teamLogoUrl && (
            <div className="absolute top-3 right-3 w-10 h-10 rounded-xl bg-white/5 p-1.5 backdrop-blur-sm border border-white/10">
              <img src={player.teamLogoUrl} alt="" className="w-full h-full object-contain drop-shadow-[0_0_6px_rgba(255,255,255,0.5)]" />
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
          <div className="relative p-5 pt-14">
            <div className="flex items-start gap-4">
              {/* Photo */}
              <motion.div
                className="relative flex-shrink-0"
                animate={isHovered ? { scale: 1.05 } : { scale: 1 }}
                transition={{ type: 'spring', stiffness: 300, damping: 20 }}
              >
                <div
                  className="w-16 h-16 rounded-xl overflow-hidden border-2"
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
                      <User size={28} className="text-gray-500" />
                    </div>
                  )}
                </div>
                {/* Photo glow */}
                <div
                  className="absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-500 -z-10 blur-lg"
                  style={{ background: positionColor }}
                />
              </motion.div>

              {/* Info */}
              <div className="flex-1 min-w-0">
                <h3 className="text-base font-bold text-white truncate group-hover:text-[#00d4ff] transition-colors">
                  {player.name}
                </h3>
                <p className="text-sm text-gray-400 truncate mt-0.5">{player.team}</p>
                <div className="flex items-center gap-2 mt-2 flex-wrap">
                  <span
                    className="text-xs font-medium px-2 py-0.5 rounded-full"
                    style={{
                      background: `${positionColor}20`,
                      color: positionColor,
                    }}
                  >
                    {POSITION_LABELS[player.position] ?? player.position}
                  </span>
                  <span className="text-xs text-gray-500">{player.birthYear} г.р.</span>
                </div>
              </div>
            </div>

            {/* Stats */}
            {player.stats && (player.stats.games > 0 || player.stats.points > 0) && (
              <div className="mt-4 flex items-center justify-between rounded-xl bg-white/5 px-4 py-3 border border-white/5">
                <StatCell label="И" value={player.stats.games} />
                <StatCell label="Г" value={player.stats.goals} color="#10b981" />
                <StatCell label="П" value={player.stats.assists} color="#8b5cf6" />
                <StatCell label="О" value={player.stats.points} color="#00d4ff" highlight />
                <TrendingUp size={16} className="text-gray-600 group-hover:text-[#00d4ff] transition-colors" />
              </div>
            )}
          </div>

          {/* Spotlight effect */}
          {isHovered && (
            <div
              className="absolute inset-0 pointer-events-none"
              style={{
                background: `radial-gradient(circle at ${(tilt.y / 8 + 0.5) * 100}% ${(tilt.x / -8 + 0.5) * 100}%, ${positionColor}08 0%, transparent 60%)`,
              }}
            />
          )}
        </div>
      </div>
    </Link>
  )
}

function StatCell({ label, value, color, highlight }: { label: string; value: number; color?: string; highlight?: boolean }) {
  return (
    <div className="text-center">
      <p className={cn('text-sm font-semibold', highlight ? 'text-[#00d4ff]' : 'text-white')} style={color ? { color } : undefined}>
        {value}
      </p>
      <p className="text-[10px] text-gray-500">{label}</p>
    </div>
  )
}

function getPositionColor(position: string): string {
  switch (position) {
    case 'goalie':
      return '#10b981' // green
    case 'defender':
      return '#8b5cf6' // purple
    case 'forward':
      return '#00d4ff' // cyan
    default:
      return '#00d4ff'
  }
}
