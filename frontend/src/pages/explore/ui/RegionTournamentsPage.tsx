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
  // Федеральные округа
  junior: { name: 'Юниорская лига', source: 'junior', domain: 'https://junior.fhr.ru', description: 'Всероссийские соревнования', gradient: 'from-violet-500 to-purple-400', icon: Trophy },
  szfo: { name: 'СЗФО', source: 'junior', domain: 'https://szfo.fhr.ru', description: 'Северо-Западный федеральный округ', gradient: 'from-cyan-500 to-blue-400', icon: MapPin },
  ufo: { name: 'УФО', source: 'junior', domain: 'https://ufo.fhr.ru', description: 'Уральский федеральный округ', gradient: 'from-emerald-500 to-teal-400', icon: Mountain },
  cfo: { name: 'ЦФО', source: 'junior', domain: 'https://cfo.fhr.ru', description: 'Центральный федеральный округ', gradient: 'from-red-500 to-rose-400', icon: Landmark },
  dfo: { name: 'ДФО', source: 'junior', domain: 'https://dfo.fhr.ru', description: 'Дальневосточный федеральный округ', gradient: 'from-orange-500 to-amber-400', icon: Mountain },
  pfo: { name: 'ПФО', source: 'junior', domain: 'https://pfo.fhr.ru', description: 'Приволжский федеральный округ', gradient: 'from-blue-500 to-cyan-400', icon: MapPin },
  sfo: { name: 'СФО', source: 'junior', domain: 'https://sfo.fhr.ru', description: 'Сибирский федеральный округ', gradient: 'from-sky-500 to-blue-400', icon: Mountain },
  yfo: { name: 'ЮФО', source: 'junior', domain: 'https://yfo.fhr.ru', description: 'Южный федеральный округ', gradient: 'from-yellow-500 to-orange-400', icon: MapPin },
  // Регионы
  spb: { name: 'Санкт-Петербург', source: 'junior', domain: 'https://spb.fhr.ru', description: 'Санкт-Петербург', gradient: 'from-purple-500 to-indigo-400', icon: Landmark },
  len: { name: 'Ленинградская обл.', source: 'junior', domain: 'https://len.fhr.ru', description: 'Ленинградская область', gradient: 'from-indigo-500 to-violet-400', icon: MapPin },
  nsk: { name: 'Новосибирск', source: 'junior', domain: 'https://nsk.fhr.ru', description: 'Новосибирская область', gradient: 'from-teal-500 to-emerald-400', icon: Building2 },
  sam: { name: 'Самара', source: 'junior', domain: 'https://sam.fhr.ru', description: 'Самарская область', gradient: 'from-pink-500 to-rose-400', icon: Building2 },
  vrn: { name: 'Воронеж', source: 'junior', domain: 'https://vrn.fhr.ru', description: 'Воронежская область', gradient: 'from-amber-500 to-yellow-400', icon: Building2 },
  komi: { name: 'Коми', source: 'junior', domain: 'https://komi.fhr.ru', description: 'Республика Коми', gradient: 'from-lime-500 to-green-400', icon: Mountain },
  kuzbass: { name: 'Кузбасс', source: 'junior', domain: 'https://kuzbass.fhr.ru', description: 'Кемеровская область', gradient: 'from-slate-500 to-gray-400', icon: Mountain },
  // Другие источники
  moscow: { name: 'Москва', source: 'fhmoscow', description: 'Москва и МО', gradient: 'from-red-500 to-orange-400', icon: Building2 },
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
