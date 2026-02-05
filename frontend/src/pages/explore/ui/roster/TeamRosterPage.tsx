import { memo, useMemo } from 'react'
import { useParams, useSearchParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { ArrowLeft, Users, MapPin, Shield, Loader2 } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { useTeamRoster } from '@/shared/api/useExploreQueries'
import type { RosterPlayer } from '@/shared/api/exploreTypes'
import { RosterPlayerCard } from './RosterPlayerCard'

const containerVariants = {
  hidden: { opacity: 0 },
  show: { opacity: 1, transition: { staggerChildren: 0.04 } },
}

const itemVariants = {
  hidden: { opacity: 0, y: 20, scale: 0.95 },
  show: { opacity: 1, y: 0, scale: 1, transition: { type: 'spring' as const, stiffness: 300, damping: 25 } },
}

export const TeamRosterPage = memo(function TeamRosterPage() {
  const { teamId, tournamentId } = useParams<{ teamId: string; tournamentId: string }>()
  const [searchParams] = useSearchParams()
  const birthYear = searchParams.get('birthYear') ? Number(searchParams.get('birthYear')) : undefined
  const groupName = searchParams.get('group') || undefined

  const { data, isLoading, error } = useTeamRoster(teamId!, tournamentId!, birthYear, groupName)

  // Group players by position
  // eslint-disable-next-line react-hooks/preserve-manual-memoization
  const groupedPlayers = useMemo(() => {
    if (!data?.players) return null

    const goalies = data.players.filter((p: RosterPlayer) => p.position?.toLowerCase().includes('врат'))
    const defensemen = data.players.filter((p: RosterPlayer) => p.position?.toLowerCase().includes('защит'))
    const forwards = data.players.filter((p: RosterPlayer) => p.position?.toLowerCase().includes('напад'))
    const other = data.players.filter(
      (p: RosterPlayer) =>
        !p.position ||
        (!p.position.toLowerCase().includes('врат') &&
          !p.position.toLowerCase().includes('защит') &&
          !p.position.toLowerCase().includes('напад'))
    )

    return { goalies, defensemen, forwards, other }
  }, [data?.players])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 size={40} className="animate-spin text-[#00d4ff]" />
      </div>
    )
  }

  if (error || !data) {
    return (
      <GlassCard className="p-8 text-center">
        <p className="text-gray-400">Не удалось загрузить состав команды</p>
      </GlassCard>
    )
  }

  const backUrl = `/explore/tournaments/detail/${tournamentId}?birthYear=${birthYear || ''}&group=${groupName || ''}`

  return (
    <div className="space-y-6">
      {/* Back link */}
      <Link
        to={backUrl}
        className="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
      >
        <ArrowLeft size={18} />
        <span>Назад к турниру</span>
      </Link>

      {/* Team header */}
      <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }}>
        <GlassCard className="p-6">
          <div className="flex items-center gap-6">
            {/* Team logo */}
            <div className="relative">
              {data.team.logoUrl ? (
                <img
                  src={data.team.logoUrl}
                  alt={data.team.name}
                  className="w-24 h-24 object-contain drop-shadow-[0_0_8px_rgba(255,255,255,0.5)]"
                />
              ) : (
                <div className="w-24 h-24 rounded-xl bg-gradient-to-br from-[#00d4ff]/10 to-[#8b5cf6]/10 border border-white/10 flex items-center justify-center">
                  <Shield size={40} className="text-gray-500" />
                </div>
              )}
              {/* Glow under logo */}
              <div className="absolute inset-0 bg-[#00d4ff]/20 blur-2xl -z-10 opacity-50" />
            </div>

            {/* Team info */}
            <div className="flex-1">
              <h1 className="text-2xl font-bold text-gradient-animated mb-2">{data.team.name}</h1>
              {data.team.city && (
                <div className="flex items-center gap-2 text-gray-400 mb-3">
                  <MapPin size={16} className="text-[#8b5cf6]" />
                  <span>{data.team.city}</span>
                </div>
              )}
              <div className="flex items-center gap-4">
                <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-[#00d4ff]/10 border border-[#00d4ff]/20">
                  <Users size={16} className="text-[#00d4ff]" />
                  <span className="text-sm font-semibold text-[#00d4ff]">{data.players.length}</span>
                  <span className="text-xs text-gray-400">игроков</span>
                </div>
                {birthYear && (
                  <div className="px-3 py-1.5 rounded-lg bg-[#8b5cf6]/10 border border-[#8b5cf6]/20">
                    <span className="text-sm font-semibold text-[#8b5cf6]">{birthYear} г.р.</span>
                  </div>
                )}
              </div>
            </div>
          </div>
        </GlassCard>
      </motion.div>

      {/* Players by position */}
      {groupedPlayers && (
        <div className="space-y-8">
          {/* Goalies */}
          {groupedPlayers.goalies.length > 0 && (
            <PositionSection
              title="Вратари"
              players={groupedPlayers.goalies}
              teamLogoUrl={data.team.logoUrl}
              color="#10b981"
            />
          )}

          {/* Defensemen */}
          {groupedPlayers.defensemen.length > 0 && (
            <PositionSection
              title="Защитники"
              players={groupedPlayers.defensemen}
              teamLogoUrl={data.team.logoUrl}
              color="#8b5cf6"
            />
          )}

          {/* Forwards */}
          {groupedPlayers.forwards.length > 0 && (
            <PositionSection
              title="Нападающие"
              players={groupedPlayers.forwards}
              teamLogoUrl={data.team.logoUrl}
              color="#00d4ff"
            />
          )}

          {/* Other */}
          {groupedPlayers.other.length > 0 && (
            <PositionSection
              title="Другие"
              players={groupedPlayers.other}
              teamLogoUrl={data.team.logoUrl}
              color="#f59e0b"
            />
          )}
        </div>
      )}
    </div>
  )
})

interface PositionSectionProps {
  title: string
  players: Array<{
    id: string
    name: string
    photoUrl?: string
    birthDate?: string
    position?: string
    jerseyNumber: number
    height?: number
    weight?: number
    birthYear?: number
  }>
  teamLogoUrl?: string
  color: string
}

function PositionSection({ title, players, teamLogoUrl, color }: PositionSectionProps) {
  return (
    <motion.section initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
      <div className="flex items-center gap-3 mb-4">
        <div className="w-1 h-6 rounded-full" style={{ background: color }} />
        <h2 className="text-lg font-semibold text-white">{title}</h2>
        <span className="text-sm text-gray-500">({players.length})</span>
      </div>

      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="show"
        className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4"
      >
        {players.map((player, index) => (
          <motion.div key={player.id} variants={itemVariants}>
            <RosterPlayerCard player={player} teamLogoUrl={teamLogoUrl} index={index} />
          </motion.div>
        ))}
      </motion.div>
    </motion.section>
  )
}
