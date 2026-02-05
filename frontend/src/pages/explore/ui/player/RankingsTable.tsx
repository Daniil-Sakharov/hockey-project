import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Medal, UserRound } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import type { RankedPlayer } from '@/shared/api/exploreTypes'

type SortKey = 'points' | 'goals' | 'assists' | 'penaltyMinutes' | 'plusMinus'

const MEDAL_COLORS = ['text-[#f59e0b]', 'text-gray-300', 'text-[#cd7f32]']
const LOGOS_PER_SET = 20
const GRID_COLS = 'grid-cols-[48px_48px_1fr_minmax(100px,200px)_60px_40px_40px_40px_40px_48px_48px]'

function LogoMarquee({ src, index }: { src: string; index: number }) {
  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none z-0">
      <div className="absolute inset-0 flex items-center opacity-[0.07] group-hover:opacity-[0.14] transition-opacity duration-500">
        <div
          className="logo-marquee"
          style={{ '--scroll-duration': `${18 + (index % 5) * 3}s` } as React.CSSProperties}
        >
          {[0, 1].map((set) =>
            Array.from({ length: LOGOS_PER_SET }).map((_, j) => (
              <img
                key={`${set}-${j}`}
                src={src}
                alt=""
                className="h-10 w-10 object-contain shrink-0 drop-shadow-[0_0_6px_rgba(255,255,255,0.5)]"
              />
            ))
          )}
        </div>
      </div>
    </div>
  )
}

export function RankingsTable({ ranked, sortBy }: { ranked: RankedPlayer[]; sortBy: SortKey }) {
  return (
    <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.15 }}>
      <GlassCard className="p-0 overflow-hidden">
        {/* Header */}
        <div className={cn('grid items-center text-gray-500 text-sm border-b border-white/10 bg-white/[0.02]', GRID_COLS)}>
          <div className="py-3.5 pl-4 pr-2">#</div>
          <div className="py-3.5 px-1" />
          <div className="py-3.5 px-3">Игрок</div>
          <div className="py-3.5 px-3 hidden md:block">Команда</div>
          <div className="py-3.5 px-3 hidden lg:block">Поз.</div>
          <div className="py-3.5 px-2 text-center">И</div>
          <div className={cn('py-3.5 px-2 text-center', sortBy === 'goals' && 'text-[#00d4ff]')}>Г</div>
          <div className={cn('py-3.5 px-2 text-center', sortBy === 'assists' && 'text-[#00d4ff]')}>П</div>
          <div className={cn('py-3.5 px-2 text-center font-semibold', sortBy === 'points' && 'text-[#00d4ff]')}>О</div>
          <div className={cn('py-3.5 px-2 text-center hidden md:block', sortBy === 'plusMinus' && 'text-[#00d4ff]')}>+/-</div>
          <div className={cn('py-3.5 pr-4 px-2 text-center hidden md:block', sortBy === 'penaltyMinutes' && 'text-[#00d4ff]')}>ШМ</div>
        </div>

        {/* Rows */}
        {ranked.map((player, i) => (
          <motion.div
            key={player.id}
            initial={{ opacity: 0, x: -10 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.02 * i }}
            className={cn(
              'relative grid items-center text-sm border-b border-white/5 hover:bg-white/[0.03] transition-colors group',
              GRID_COLS,
            )}
          >
            {/* Animated logo marquee background */}
            {player.teamLogoUrl && <LogoMarquee src={player.teamLogoUrl} index={i} />}

            <div className="py-5 pl-4 pr-2 relative z-10">
              {i < 3 ? (
                <Medal size={22} className={MEDAL_COLORS[i]} />
              ) : (
                <span className="text-gray-500 font-medium">{player.rank}</span>
              )}
            </div>
            <div className="py-4 px-1 relative z-10">
              {player.photoUrl ? (
                <img src={player.photoUrl} alt={player.name} className="h-11 w-11 rounded-full object-cover ring-1 ring-white/10" />
              ) : (
                <div className="flex h-11 w-11 items-center justify-center rounded-full bg-white/10">
                  <UserRound size={18} className="text-gray-500" />
                </div>
              )}
            </div>
            <div className="py-5 px-3 relative z-10 min-w-0">
              <Link to={`/explore/players/${player.id}`} className="font-medium text-white hover:text-[#00d4ff] transition-colors text-[13px] truncate block">
                {player.name}
              </Link>
            </div>
            <div className="py-5 px-3 hidden md:block relative z-10 min-w-0">
              <Link to={`/explore/teams/${player.teamId}`} className="hover:text-[#8b5cf6] transition-colors flex items-center gap-2">
                {player.teamLogoUrl && (
                  <img src={player.teamLogoUrl} alt="" className="h-5 w-5 object-contain drop-shadow-[0_0_6px_rgba(255,255,255,0.5)] shrink-0" />
                )}
                <div className="min-w-0">
                  <span className="text-gray-400 text-[13px] truncate block">{player.team}</span>
                  {player.teamCity && (
                    <span className="text-gray-600 text-[11px] truncate block">{player.teamCity}</span>
                  )}
                </div>
              </Link>
            </div>
            <div className="py-5 px-3 hidden lg:block relative z-10">
              <PositionBadge position={player.position} />
            </div>
            <div className="py-5 px-2 text-center text-gray-400 relative z-10">{player.games}</div>
            <div className={cn('py-5 px-2 text-center relative z-10', sortBy === 'goals' ? 'text-[#00d4ff] font-semibold' : 'text-gray-300')}>{player.goals}</div>
            <div className={cn('py-5 px-2 text-center relative z-10', sortBy === 'assists' ? 'text-[#00d4ff] font-semibold' : 'text-gray-300')}>{player.assists}</div>
            <div className={cn('py-5 px-2 text-center font-semibold relative z-10', sortBy === 'points' ? 'text-[#00d4ff]' : 'text-white')}>{player.points}</div>
            <div className={cn('py-5 px-2 text-center hidden md:block relative z-10', sortBy === 'plusMinus' ? 'text-[#00d4ff] font-semibold' : player.plusMinus >= 0 ? 'text-[#10b981]' : 'text-[#ef4444]')}>
              {player.plusMinus > 0 ? `+${player.plusMinus}` : player.plusMinus}
            </div>
            <div className={cn('py-5 pr-4 px-2 text-center hidden md:block relative z-10', sortBy === 'penaltyMinutes' ? 'text-[#00d4ff] font-semibold' : 'text-gray-400')}>
              {player.penaltyMinutes}
            </div>
          </motion.div>
        ))}
      </GlassCard>
    </motion.div>
  )
}

function PositionBadge({ position }: { position: string }) {
  const labels: Record<string, string> = { forward: 'Нап', defender: 'Защ', goalie: 'Вр' }
  return (
    <span className={cn(
      'text-[10px] px-1.5 py-0.5 rounded font-medium',
      position === 'forward' ? 'bg-[#00d4ff]/20 text-[#00d4ff]' :
      position === 'defender' ? 'bg-[#8b5cf6]/20 text-[#8b5cf6]' :
      'bg-[#f59e0b]/20 text-[#f59e0b]'
    )}>
      {labels[position] || position}
    </span>
  )
}
