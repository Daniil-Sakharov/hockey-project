import { memo, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  ArrowLeft,
  Users2,
  MapPin,
  Trophy,
  UserRound,
  User,
  Shield,
  ChevronRight,
  Loader2,
} from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import { useTeamProfile } from '@/shared/api/useExploreQueries'

function PlayerPhoto({ url, name }: { url?: string; name: string }) {
  const [hasError, setHasError] = useState(false)
  if (!url || hasError) {
    return (
      <div className="w-10 h-10 rounded-full bg-white/10 flex items-center justify-center flex-shrink-0">
        <User size={20} className="text-gray-600" />
      </div>
    )
  }
  return (
    <img
      src={url}
      alt={name}
      className="w-10 h-10 rounded-full object-cover flex-shrink-0"
      onError={() => setHasError(true)}
    />
  )
}

function TeamLogoLarge({ url, name }: { url?: string; name: string }) {
  const [hasError, setHasError] = useState(false)
  if (!url || hasError) {
    return (
      <div className="flex h-24 w-24 items-center justify-center rounded-2xl bg-gradient-to-br from-[#8b5cf6] to-[#ec4899] flex-shrink-0">
        <Shield size={48} className="text-white" />
      </div>
    )
  }
  return (
    <img
      src={url}
      alt={name}
      className="h-24 w-24 object-contain flex-shrink-0"
      onError={() => setHasError(true)}
    />
  )
}

export const TeamProfilePage = memo(function TeamProfilePage() {
  const { id } = useParams<{ id: string }>()
  const { data: team, isLoading } = useTeamProfile(id ?? '')

  if (isLoading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 size={32} className="animate-spin text-gray-500" />
      </div>
    )
  }

  if (!team) {
    return (
      <div className="flex flex-col items-center justify-center py-20 text-center">
        <Users2 size={48} className="text-gray-700 mb-4" />
        <p className="text-gray-400 text-lg">Команда не найдена</p>
        <Link to="/explore/players" className="mt-4 text-sm text-[#00d4ff] hover:underline">
          Вернуться к поиску
        </Link>
      </div>
    )
  }

  const { stats } = team
  const totalGames = stats.wins + stats.losses + stats.draws

  return (
    <div className="space-y-6">
      <Link
        to="/explore/players"
        className="inline-flex items-center gap-2 text-sm text-gray-400 hover:text-white transition-colors"
      >
        <ArrowLeft size={16} />
        Назад
      </Link>

      {/* Team hero */}
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        <GlassCard className="p-6" glowColor="purple">
          <div className="flex items-center gap-6">
            <TeamLogoLarge url={team.logoUrl} name={team.name} />
            <div className="flex-1">
              <h1 className="text-2xl md:text-3xl font-bold text-white">{team.name}</h1>
              <div className="flex items-center gap-4 mt-2">
                {team.city && (
                  <span className="inline-flex items-center gap-1 text-sm text-gray-400">
                    <MapPin size={14} /> {team.city}
                  </span>
                )}
                <span className="inline-flex items-center gap-1 text-sm text-gray-400">
                  <UserRound size={14} /> {team.playersCount} игроков
                </span>
              </div>
              {team.tournaments.length > 0 && (
                <div className="flex flex-wrap gap-2 mt-4">
                  {team.tournaments.map((t) => (
                    <span
                      key={t}
                      className="px-2 py-1 rounded-lg bg-[#8b5cf6]/20 text-[#8b5cf6] text-xs font-medium"
                    >
                      <Trophy size={12} className="inline mr-1" />
                      {t}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>
        </GlassCard>
      </motion.div>

      {/* Stats row */}
      <div className="grid grid-cols-3 gap-4 md:grid-cols-6">
        {[
          { label: 'Игры', value: totalGames },
          { label: 'Победы', value: stats.wins, color: 'text-[#10b981]' },
          { label: 'Ничьи', value: stats.draws, color: 'text-[#f59e0b]' },
          { label: 'Поражения', value: stats.losses, color: 'text-[#ef4444]' },
          { label: 'Забито', value: stats.goalsFor, color: 'text-[#00d4ff]' },
          { label: 'Пропущено', value: stats.goalsAgainst, color: 'text-[#ec4899]' },
        ].map((s, i) => (
          <motion.div
            key={s.label}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 + i * 0.04 }}
          >
            <GlassCard className="p-3 text-center">
              <p className={cn('text-xl font-bold', s.color || 'text-white')}>{s.value}</p>
              <p className="text-[10px] text-gray-500 mt-0.5">{s.label}</p>
            </GlassCard>
          </motion.div>
        ))}
      </div>

      {/* Roster */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
      >
        <GlassCard className="p-6" glowColor="cyan">
          <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Users2 size={20} className="text-[#00d4ff]" />
            Состав ({team.roster.length})
          </h3>
          {team.roster.length > 0 ? (
            <div className="space-y-2">
              {team.roster.map((player) => (
                <Link
                  key={player.id}
                  to={`/explore/players/${player.id}`}
                  className="flex items-center gap-3 rounded-lg bg-white/5 p-3 hover:bg-white/[0.08] transition-colors"
                >
                  <PlayerPhoto url={player.photoUrl} name={player.name} />
                  {player.jerseyNumber > 0 && (
                    <span className="text-base font-bold text-gray-400 w-10 text-center">
                      #{player.jerseyNumber}
                    </span>
                  )}
                  <div className="flex-1 min-w-0">
                    <span className="text-base font-medium text-white block truncate">{player.name}</span>
                    <span className="text-xs text-gray-500">{player.birthYear} г.р.</span>
                  </div>
                  <PositionBadge position={player.position} />
                  {player.stats && (
                    <div className="text-right">
                      <span className="text-sm font-bold text-[#00d4ff] block">
                        {player.stats.points}
                      </span>
                      <span className="text-[10px] text-gray-500">очков</span>
                    </div>
                  )}
                  <ChevronRight size={16} className="text-gray-600 flex-shrink-0" />
                </Link>
              ))}
            </div>
          ) : (
            <p className="text-sm text-gray-500">Состав команды не найден</p>
          )}
        </GlassCard>
      </motion.div>

      {/* Recent matches */}
      {team.recentMatches && team.recentMatches.length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.35 }}
        >
          <GlassCard className="p-6">
            <h3 className="text-lg font-semibold text-white mb-4">Последние матчи</h3>
            <div className="space-y-2">
              {team.recentMatches.map((match) => (
                <Link
                  key={match.id}
                  to={`/explore/matches/${match.id}`}
                  className="flex items-center justify-between rounded-lg bg-white/5 p-3 hover:bg-white/[0.08] transition-colors"
                >
                  <span className="text-sm text-white flex-1 truncate">{match.homeTeam}</span>
                  <span className="text-lg font-bold text-[#00d4ff] mx-3">
                    {match.homeScore}:{match.awayScore}
                  </span>
                  <span className="text-sm text-white flex-1 truncate text-right">{match.awayTeam}</span>
                  <ChevronRight size={16} className="text-gray-600 ml-2 flex-shrink-0" />
                </Link>
              ))}
            </div>
          </GlassCard>
        </motion.div>
      )}
    </div>
  )
})

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
