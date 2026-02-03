import { memo, useMemo, useCallback } from 'react'
import { useParams, useSearchParams, Link, Navigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { ArrowLeft, Trophy, Loader2, MapPin, Building2, Landmark, Mountain } from 'lucide-react'
import { useTournaments } from '@/shared/api/useExploreQueries'
import type { TournamentItem } from '@/shared/api/exploreTypes'
import { TournamentFilters } from './TournamentFilters'
import { TournamentGroupCard } from './TournamentGroupCard'

const REGIONS_MAP: Record<string, {
  name: string
  domain?: string
  source?: string
  description: string
  gradient: string
  icon: React.ComponentType<{ size?: number; className?: string }>
}> = {
  pfo: { name: 'Приволжский', domain: 'pfo', description: 'ПФО', gradient: 'from-blue-500 to-cyan-400', icon: MapPin },
  cfo: { name: 'Центральный', domain: 'cfo', description: 'ЦФО', gradient: 'from-red-500 to-rose-400', icon: Landmark },
  sfo: { name: 'Сибирский', domain: 'sfo', description: 'СФО', gradient: 'from-sky-500 to-blue-400', icon: Mountain },
  ufo: { name: 'Уральский', domain: 'ufo', description: 'УрФО', gradient: 'from-emerald-500 to-teal-400', icon: Mountain },
  dfo: { name: 'Дальневосточный', domain: 'dfo', description: 'ДФО', gradient: 'from-orange-500 to-amber-400', icon: Mountain },
  szfo: { name: 'Северо-Западный', domain: 'szfo', description: 'СЗФО', gradient: 'from-cyan-500 to-blue-400', icon: MapPin },
  yfo: { name: 'Южный', domain: 'yfo', description: 'ЮФО', gradient: 'from-yellow-500 to-orange-400', icon: MapPin },
  spb: { name: 'Санкт-Петербург', domain: 'spb', description: 'СПб', gradient: 'from-purple-500 to-indigo-400', icon: Landmark },
  moscow: { name: 'Москва', source: 'fhmoscow', description: 'Москва и МО', gradient: 'from-red-500 to-orange-400', icon: Building2 },
  junior: { name: 'Всероссийские', domain: 'junior', description: 'Юниор', gradient: 'from-violet-500 to-purple-400', icon: Trophy },
}

interface ExpandedCard {
  tournament: TournamentItem
  groupName: string
  birthYear: number
  teamsCount: number
  matchesCount: number
}

function getTournamentBirthYears(t: TournamentItem): number[] {
  if (!t.birthYearGroups) return []
  return Object.keys(t.birthYearGroups)
    .map(Number)
    .filter((y) => !isNaN(y) && (t.birthYearGroups?.[String(y)]?.length ?? 0) > 0)
}

export const RegionTournamentsPage = memo(function RegionTournamentsPage() {
  const { region } = useParams<{ region: string }>()
  const [searchParams, setSearchParams] = useSearchParams()
  const regionData = region ? REGIONS_MAP[region] : undefined
  const source = regionData?.source ?? ''
  const domain = regionData?.domain ?? ''

  const { data: tournaments, isLoading } = useTournaments(source, domain)

  const seasons = useMemo(() => {
    if (!tournaments) return []
    const set = new Set(tournaments.map((t) => t.season))
    return Array.from(set).sort().reverse()
  }, [tournaments])

  const currentSeason = seasons[0] ?? ''
  const seasonParam = searchParams.get('season')
  const activeSeason = seasonParam ?? currentSeason

  const seasonTournaments = useMemo(() => {
    if (!tournaments || !activeSeason) return []
    return tournaments.filter((t) => t.season === activeSeason)
  }, [tournaments, activeSeason])

  const birthYears = useMemo(() => {
    const allYears = new Set<number>()
    for (const t of seasonTournaments) {
      for (const y of getTournamentBirthYears(t)) allYears.add(y)
    }
    return Array.from(allYears).sort()
  }, [seasonTournaments])

  const oldestBirthYear = birthYears[0] ?? null
  const birthYearParam = searchParams.get('birthYear')
  const activeBirthYear = birthYearParam ? Number(birthYearParam) : oldestBirthYear

  const updateParam = useCallback((key: string, value: string | null) => {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev)
      if (!value) {
        next.delete(key)
      } else {
        next.set(key, value)
      }
      return next
    }, { replace: true })
  }, [setSearchParams])

  const handleSeasonChange = useCallback((s: string) => {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev)
      next.set('season', s)
      next.delete('birthYear')
      return next
    }, { replace: true })
  }, [setSearchParams])

  const handleBirthYearChange = useCallback((y: number) => {
    updateParam('birthYear', String(y))
  }, [updateParam])

  // Развёрнутые карточки: турнир × группа
  const expandedCards = useMemo<ExpandedCard[]>(() => {
    if (!activeBirthYear) return []
    const cards: ExpandedCard[] = []
    for (const t of seasonTournaments) {
      const groups = t.birthYearGroups?.[String(activeBirthYear)] ?? []
      if (groups.length === 0) continue
      for (const g of groups) {
        cards.push({
          tournament: t, groupName: g.name, birthYear: activeBirthYear,
          teamsCount: g.teamsCount, matchesCount: g.matchesCount,
        })
      }
    }
    return cards
  }, [seasonTournaments, activeBirthYear])

  if (!regionData) return <Navigate to="/explore/tournaments" replace />

  const Icon = regionData.icon

  return (
    <div className="space-y-6">
      <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }}>
        <Link to="/explore/tournaments" className="inline-flex items-center gap-1 text-sm text-gray-400 hover:text-white transition-colors mb-4">
          <ArrowLeft size={16} />
          Все регионы
        </Link>
        <div className="flex items-center gap-4">
          <div className={`flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br ${regionData.gradient} shadow-lg`}>
            <Icon size={22} className="text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-white">{regionData.name}</h1>
            <p className="text-gray-400 text-sm">{regionData.description}</p>
          </div>
        </div>
      </motion.div>

      <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 0.1 }}>
        <TournamentFilters
          seasons={seasons}
          activeSeason={activeSeason}
          onSeasonChange={handleSeasonChange}
          birthYears={birthYears}
          activeBirthYear={activeBirthYear}
          onBirthYearChange={handleBirthYearChange}
        />
      </motion.div>

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Loader2 size={32} className="animate-spin text-gray-500" />
        </div>
      ) : expandedCards.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-16 text-center">
          <Trophy size={48} className="text-gray-600 mb-4" />
          <p className="text-gray-400">Турниров пока нет</p>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {expandedCards.map((card, i) => (
            <TournamentGroupCard
              key={`${card.tournament.id}-${card.groupName}`}
              tournament={card.tournament}
              groupName={card.groupName}
              birthYear={card.birthYear}
              teamsCount={card.teamsCount}
              matchesCount={card.matchesCount}
              index={i}
              region={region ?? ''}
            />
          ))}
        </div>
      )}
    </div>
  )
})
