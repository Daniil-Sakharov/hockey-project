import { memo } from 'react'
import { useParams, useSearchParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  ArrowLeft,
  UserRound,
  Ruler,
  Weight,
  Hand,
  MapPin,
  Calendar,
  Target,
  Shield,
  TrendingUp,
  Loader2,
} from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import { useCountUp } from '@/shared/hooks'
import { usePlayerProfile } from '@/shared/api/useExploreQueries'
import { PlayerStatsHistory } from './PlayerStatsHistory'
import { PlayerChartsSection } from './charts/PlayerChartsSection'

const POSITION_LABELS: Record<string, string> = {
  forward: 'Нападающий',
  defender: 'Защитник',
  goalie: 'Вратарь',
}

const HANDEDNESS_LABELS: Record<string, string> = {
  left: 'Левый',
  right: 'Правый',
}

export const PlayerProfilePage = memo(function PlayerProfilePage() {
  const { id } = useParams<{ id: string }>()
  const [searchParams] = useSearchParams()
  const season = searchParams.get('season') ?? undefined
  const fromQuery = searchParams.get('from') ?? ''
  const backUrl = `/explore/players${fromQuery ? `?${fromQuery}` : ''}`
  const { data: player, isLoading } = usePlayerProfile(id ?? '', season)

  if (isLoading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 size={32} className="animate-spin text-gray-500" />
      </div>
    )
  }

  if (!player) {
    return (
      <div className="flex flex-col items-center justify-center py-20 text-center">
        <UserRound size={48} className="text-gray-700 mb-4" />
        <p className="text-gray-400 text-lg">Игрок не найден</p>
        <Link to={backUrl} className="mt-4 text-sm text-[#00d4ff] hover:underline">
          Вернуться к поиску
        </Link>
      </div>
    )
  }

  const stats = player.stats

  return (
    <div className="space-y-6">
      {/* Back button */}
      <Link
        to={backUrl}
        className="inline-flex items-center gap-2 text-sm text-gray-400 hover:text-white transition-colors"
      >
        <ArrowLeft size={16} />
        Назад к поиску
      </Link>

      {/* Hero card */}
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        <GlassCard className="relative overflow-hidden" glowColor="cyan">
          {/* Team logo background */}
          {player.teamLogoUrl && (
            <div className="absolute inset-0 flex items-center justify-end overflow-hidden pointer-events-none">
              <img
                src={player.teamLogoUrl}
                alt=""
                className="w-80 h-80 object-contain opacity-[0.08] translate-x-20 drop-shadow-[0_0_16px_rgba(255,255,255,0.6)]"
              />
            </div>
          )}

          <div className="relative p-6">
            <div className="flex flex-col gap-6 md:flex-row md:items-start">
              {/* Avatar with team logo badge */}
              <div className="relative flex-shrink-0">
                {player.photoUrl ? (
                  <img
                    src={player.photoUrl}
                    alt={player.name}
                    className="h-28 w-28 rounded-2xl object-cover border-2 border-white/10"
                  />
                ) : (
                  <div className="flex h-28 w-28 items-center justify-center rounded-2xl bg-gradient-to-br from-[#00d4ff] to-[#8b5cf6]">
                    <UserRound size={56} className="text-white" />
                  </div>
                )}
                {/* Team logo badge */}
                {player.teamLogoUrl && (
                  <div className="absolute -bottom-2 -right-2 w-12 h-12 rounded-xl bg-[#0a0a0f]/90 p-1.5 border border-white/10 shadow-lg">
                    <img src={player.teamLogoUrl} alt="" className="w-full h-full object-contain drop-shadow-[0_0_6px_rgba(255,255,255,0.5)]" />
                  </div>
                )}
              </div>

              {/* Info */}
              <div className="flex-1">
                <div className="flex items-center gap-3 mb-2 flex-wrap">
                  <h1 className="text-2xl font-bold text-white">{player.name}</h1>
                  {player.jerseyNumber > 0 && (
                    <span className="text-xl font-bold text-[#00d4ff]">#{player.jerseyNumber}</span>
                  )}
                  {season && (
                    <span className="text-sm text-[#8b5cf6] bg-[#8b5cf6]/10 px-2 py-0.5 rounded-lg">
                      {season}
                    </span>
                  )}
                </div>

                <Link
                  to={`/explore/teams/${player.teamId}`}
                  className="inline-flex items-center gap-2 text-[#00d4ff] hover:underline text-sm"
                >
                  {player.teamLogoUrl && (
                    <img src={player.teamLogoUrl} alt="" className="w-5 h-5 object-contain drop-shadow-[0_0_4px_rgba(255,255,255,0.5)]" />
                  )}
                  {player.team}
                </Link>

                <div className="flex flex-wrap items-center gap-4 mt-4">
                  <span className={cn(
                    'px-3 py-1.5 rounded-lg text-xs font-semibold',
                    player.position === 'forward' ? 'bg-[#00d4ff]/20 text-[#00d4ff] border border-[#00d4ff]/30' :
                    player.position === 'defender' ? 'bg-[#8b5cf6]/20 text-[#8b5cf6] border border-[#8b5cf6]/30' :
                    'bg-[#10b981]/20 text-[#10b981] border border-[#10b981]/30'
                  )}>
                    {POSITION_LABELS[player.position] ?? player.position}
                  </span>

                  <InfoChip icon={<Calendar size={14} />} text={player.birthDate} />
                  {player.city && <InfoChip icon={<MapPin size={14} />} text={player.city} />}
                  {player.height && <InfoChip icon={<Ruler size={14} />} text={`${player.height} см`} />}
                  {player.weight && <InfoChip icon={<Weight size={14} />} text={`${player.weight} кг`} />}
                  {player.handedness && (
                    <InfoChip
                      icon={<Hand size={14} />}
                      text={HANDEDNESS_LABELS[player.handedness] ?? player.handedness}
                    />
                  )}
                </div>
              </div>
            </div>
          </div>
        </GlassCard>
      </motion.div>

      {/* Stats cards */}
      {stats && (
        <>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-6">
            {[
              { label: 'Игры', value: stats.games, icon: <Calendar size={20} />, color: 'gray' },
              { label: 'Голы', value: stats.goals, icon: <Target size={20} />, color: 'cyan' },
              { label: 'Передачи', value: stats.assists, icon: <TrendingUp size={20} />, color: 'purple' },
              { label: 'Очки', value: stats.points, icon: <TrendingUp size={20} />, color: 'pink' },
              { label: '+/-', value: stats.plusMinus, icon: <Shield size={20} />, color: 'green' },
              { label: 'Штраф', value: stats.penaltyMinutes, icon: <Shield size={20} />, color: 'red' },
            ].map((s, i) => (
              <motion.div
                key={s.label}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.1 + i * 0.05 }}
              >
                <GlassCard className="p-4 text-center">
                  <AnimatedStatValue
                    value={s.value}
                    prefix={s.value > 0 && s.label === '+/-' ? '+' : ''}
                    className={cn(
                      'text-2xl font-bold',
                      s.color === 'cyan' ? 'text-[#00d4ff]' :
                      s.color === 'purple' ? 'text-[#8b5cf6]' :
                      s.color === 'pink' ? 'text-[#ec4899]' :
                      s.color === 'green' ? 'text-[#10b981]' :
                      s.color === 'red' ? 'text-[#ef4444]' : 'text-white'
                    )}
                  />
                  <p className="text-xs text-gray-500 mt-1">{s.label}</p>
                </GlassCard>
              </motion.div>
            ))}
          </div>

          {/* Per game averages */}
          {stats.games > 0 && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.4 }}
            >
              <GlassCard className="p-6" glowColor="purple">
                <h3 className="text-lg font-semibold text-white mb-4">Средние за игру</h3>
                <div className="grid grid-cols-3 gap-4">
                  <AvgStat label="Голов" value={stats.goals / stats.games} />
                  <AvgStat label="Передач" value={stats.assists / stats.games} />
                  <AvgStat label="Очков" value={stats.points / stats.games} highlight />
                </div>
              </GlassCard>
            </motion.div>
          )}
        </>
      )}

      {/* Charts */}
      <PlayerChartsSection playerId={id ?? ''} />

      {/* Detailed stats history */}
      <PlayerStatsHistory playerId={id ?? ''} />
    </div>
  )
})

function InfoChip({ icon, text }: { icon: React.ReactNode; text: string }) {
  return (
    <span className="inline-flex items-center gap-1.5 text-xs text-gray-400">
      {icon}
      {text}
    </span>
  )
}

function AnimatedStatValue({ value, prefix = '', className }: { value: number; prefix?: string; className?: string }) {
  const animated = useCountUp(value)
  return <p className={className}>{prefix}{animated}</p>
}

function AvgStat({ label, value, highlight }: { label: string; value: number; highlight?: boolean }) {
  return (
    <div className="text-center">
      <p className={cn('text-xl font-bold', highlight ? 'text-[#00d4ff]' : 'text-white')}>
        {value.toFixed(2)}
      </p>
      <p className="text-xs text-gray-500 mt-1">{label}</p>
    </div>
  )
}
