import { memo, useState } from 'react'
import { useParams, useSearchParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { ArrowLeft, Trophy, Calendar, TrendingUp, Users, Loader2 } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { cleanTournamentName, formatGroupName } from '@/shared/lib/formatters'
import {
  useTournaments,
  useTournamentStandings,
  useTournamentMatches,
  useTournamentScorers,
  useTournamentTeams,
} from '@/shared/api/useExploreQueries'
import { StandingsTab } from './tabs/StandingsTab'
import { MatchesTab } from './tabs/MatchesTab'
import { ScorersTab } from './tabs/ScorersTab'
import { TeamsTab } from './tabs/TeamsTab'

type Tab = 'standings' | 'matches' | 'scorers' | 'teams'

export const TournamentDetailPage = memo(function TournamentDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [searchParams] = useSearchParams()
  const [activeTab, setActiveTab] = useState<Tab>('standings')

  const birthYear = Number(searchParams.get('birthYear')) || undefined
  const groupName = searchParams.get('group') ?? undefined
  const fromRegion = searchParams.get('from') ?? ''
  const backUrl = fromRegion ? `/explore/tournaments/${fromRegion}` : '/explore/tournaments'

  const { data: tournaments, isLoading: tournamentsLoading } = useTournaments()
  const tournament = tournaments?.find((t) => t.id === id)

  const { data: standings, isLoading: standingsLoading } = useTournamentStandings(id ?? '', birthYear, groupName)
  const { data: matches, isLoading: matchesLoading } = useTournamentMatches(id ?? '', undefined, birthYear, groupName)
  const { data: scorers, isLoading: scorersLoading } = useTournamentScorers(id ?? '', undefined, birthYear, groupName)
  const { data: teams, isLoading: teamsLoading } = useTournamentTeams(id ?? '', birthYear, groupName)

  if (tournamentsLoading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 size={32} className="animate-spin text-gray-500" />
      </div>
    )
  }

  if (!tournament) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center text-center">
        <h1 className="text-2xl font-bold text-white mb-2">Турнир не найден</h1>
        <Link to={backUrl} className="text-[#00d4ff] hover:underline">
          Назад к турнирам
        </Link>
      </div>
    )
  }

  const tabs: { key: Tab; label: string; icon: React.ReactNode }[] = [
    { key: 'standings', label: 'Таблица', icon: <Trophy size={16} /> },
    { key: 'teams', label: 'Команды', icon: <Users size={16} /> },
    { key: 'matches', label: 'Матчи', icon: <Calendar size={16} /> },
    { key: 'scorers', label: 'Бомбардиры', icon: <TrendingUp size={16} /> },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }}>
        <Link
          to={backUrl}
          className="inline-flex items-center gap-1 text-sm text-gray-400 hover:text-white transition-colors mb-4"
        >
          <ArrowLeft size={16} />
          Назад к турнирам
        </Link>
        <h1 className="text-2xl font-bold text-white">
          {cleanTournamentName(tournament.name)}
          {groupName && <span className="text-[#00d4ff]"> — {formatGroupName(groupName)}</span>}
        </h1>
        <div className="flex items-center gap-3 mt-1 flex-wrap">
          <span className="text-gray-400 text-sm">{tournament.domain}</span>
          <span className="text-gray-600">·</span>
          <span className="text-gray-400 text-sm">Сезон {tournament.season}</span>
          {birthYear && (
            <>
              <span className="text-gray-600">·</span>
              <span className="text-[#8b5cf6] text-sm">{birthYear} г.р.</span>
            </>
          )}
          <span
            className={cn(
              'text-xs px-2 py-0.5 rounded-full',
              !tournament.isEnded ? 'bg-green-500/20 text-green-400' : 'bg-gray-500/20 text-gray-400',
            )}
          >
            {!tournament.isEnded ? 'Активный' : 'Завершён'}
          </span>
        </div>
      </motion.div>

      {/* Tabs */}
      <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 0.1 }} className="flex gap-2">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={cn(
              'flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium transition-all duration-200',
              activeTab === tab.key
                ? 'bg-[#00d4ff]/20 text-[#00d4ff] border border-[#00d4ff]/30'
                : 'bg-white/5 text-gray-400 border border-white/10 hover:bg-white/[0.07]',
            )}
          >
            {tab.icon}
            {tab.label}
          </button>
        ))}
      </motion.div>

      {/* Tab content */}
      <motion.div key={activeTab} initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.2 }}>
        {activeTab === 'standings' && <StandingsTab standings={standings ?? []} isLoading={standingsLoading} />}
        {activeTab === 'teams' && (
          <TeamsTab
            teams={teams ?? []}
            isLoading={teamsLoading}
            tournamentId={id ?? ''}
            birthYear={birthYear}
            groupName={groupName}
          />
        )}
        {activeTab === 'matches' && <MatchesTab matches={matches ?? []} isLoading={matchesLoading} />}
        {activeTab === 'scorers' && <ScorersTab scorers={scorers ?? []} isLoading={scorersLoading} />}
      </motion.div>
    </div>
  )
})
