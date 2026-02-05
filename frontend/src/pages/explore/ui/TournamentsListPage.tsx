import { memo, useMemo } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { MapPin, Loader2, Lock } from 'lucide-react'
import { useTournaments } from '@/shared/api/useExploreQueries'

// Все доступные регионы с маппингом домен → название
const ALL_REGIONS: { key: string; domain: string; name: string; description: string }[] = [
  { key: 'junior', domain: 'junior.fhr.ru', name: 'Юниорская лига', description: 'Всероссийские юниорские соревнования' },
  { key: 'szfo', domain: 'szfo.fhr.ru', name: 'СЗФО', description: 'Северо-Западный федеральный округ' },
  { key: 'ufo', domain: 'ufo.fhr.ru', name: 'УФО', description: 'Уральский федеральный округ' },
  { key: 'cfo', domain: 'cfo.fhr.ru', name: 'ЦФО', description: 'Центральный федеральный округ' },
  { key: 'dfo', domain: 'dfo.fhr.ru', name: 'ДФО', description: 'Дальневосточный федеральный округ' },
  { key: 'pfo', domain: 'pfo.fhr.ru', name: 'ПФО', description: 'Приволжский федеральный округ' },
  { key: 'sfo', domain: 'sfo.fhr.ru', name: 'СФО', description: 'Сибирский федеральный округ' },
  { key: 'yfo', domain: 'yfo.fhr.ru', name: 'ЮФО', description: 'Южный федеральный округ' },
  { key: 'spb', domain: 'spb.fhr.ru', name: 'Санкт-Петербург', description: 'Санкт-Петербург' },
  { key: 'len', domain: 'len.fhr.ru', name: 'Ленинградская обл.', description: 'Ленинградская область' },
  { key: 'nsk', domain: 'nsk.fhr.ru', name: 'Новосибирск', description: 'Новосибирская область' },
  { key: 'sam', domain: 'sam.fhr.ru', name: 'Самара', description: 'Самарская область' },
  { key: 'vrn', domain: 'vrn.fhr.ru', name: 'Воронеж', description: 'Воронежская область' },
  { key: 'komi', domain: 'komi.fhr.ru', name: 'Коми', description: 'Республика Коми' },
  { key: 'kuzbass', domain: 'kuzbass.fhr.ru', name: 'Кузбасс', description: 'Кузбасс (Кемеровская область)' },
]

// Градиенты для карточек (циклически применяются к регионам)
const GRADIENTS = [
  { gradient: 'from-cyan-500 via-blue-400 to-cyan-600', bg: 'rgba(6,182,212,0.15)' },
  { gradient: 'from-purple-500 via-indigo-400 to-purple-600', bg: 'rgba(139,92,246,0.15)' },
  { gradient: 'from-emerald-500 via-teal-400 to-emerald-600', bg: 'rgba(16,185,129,0.15)' },
  { gradient: 'from-orange-500 via-amber-400 to-orange-600', bg: 'rgba(249,115,22,0.15)' },
  { gradient: 'from-pink-500 via-rose-400 to-pink-600', bg: 'rgba(236,72,153,0.15)' },
  { gradient: 'from-blue-500 via-sky-400 to-blue-600', bg: 'rgba(59,130,246,0.15)' },
]

interface RegionData {
  key: string
  domain: string
  name: string
  description: string
  tournamentCount: number
  gradient: string
  bgPattern: string
  available: boolean
}

export const TournamentsListPage = memo(function TournamentsListPage() {
  // Загружаем все турниры
  const { data: tournaments, isLoading } = useTournaments()

  // Создаём данные для всех регионов с реальным количеством турниров
  const regions = useMemo<RegionData[]>(() => {
    // Считаем турниры по доменам
    const domainCounts = new Map<string, number>()
    if (tournaments?.length) {
      for (const t of tournaments) {
        // Извлекаем домен из URL (например, "https://pfo.fhr.ru" → "pfo.fhr.ru")
        const domain = t.domain?.replace(/^https?:\/\//, '') || ''
        if (domain) {
          domainCounts.set(domain, (domainCounts.get(domain) || 0) + 1)
        }
      }
    }

    // Создаём массив всех регионов
    return ALL_REGIONS.map((region, i) => {
      const count = domainCounts.get(region.domain) || 0
      const style = GRADIENTS[i % GRADIENTS.length]
      return {
        key: region.key,
        domain: region.domain,
        name: region.name,
        description: region.description,
        tournamentCount: count,
        gradient: style.gradient,
        bgPattern: style.bg,
        available: count > 0,
      }
    }).sort((a, b) => {
      // Сначала доступные, потом недоступные; внутри - по количеству турниров
      if (a.available !== b.available) return a.available ? -1 : 1
      return b.tournamentCount - a.tournamentCount
    })
  }, [tournaments])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="h-8 w-8 animate-spin text-cyan-400" />
      </div>
    )
  }

  const availableCount = regions.filter(r => r.available).length

  return (
    <div className="space-y-8">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <h1 className="text-3xl font-bold text-white">Турниры</h1>
        <p className="text-gray-400 mt-1">
          {availableCount > 0
            ? `${availableCount} из ${regions.length} регионов с турнирами`
            : 'Выберите регион для просмотра турниров'}
        </p>
      </motion.div>

      {/* Region cards grid */}
      <div className="grid gap-6 sm:grid-cols-2">
        {regions.map((region, i) => (
          <motion.div
            key={region.key}
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 + i * 0.05, duration: 0.4 }}
          >
            {region.available ? (
              <Link to={`/explore/tournaments/${region.key}`}>
                <RegionCard region={region} />
              </Link>
            ) : (
              <RegionCard region={region} disabled />
            )}
          </motion.div>
        ))}
      </div>
    </div>
  )
})

interface RegionCardProps {
  region: RegionData
  disabled?: boolean
}

function RegionCard({ region, disabled }: RegionCardProps) {
  return (
    <div
      className={`
        relative overflow-hidden rounded-2xl border transition-all duration-300 group
        ${disabled
          ? 'border-white/5 opacity-50 cursor-default'
          : 'border-white/10 hover:border-white/20 hover:scale-[1.02] hover:shadow-xl hover:shadow-black/20 cursor-pointer'
        }
      `}
      style={{ background: region.bgPattern }}
    >
      {/* Gradient top bar */}
      <div className={`h-1.5 bg-gradient-to-r ${region.gradient}`} />

      <div className="p-6 bg-[#0a0e1a]/80 backdrop-blur-sm">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-4">
            {/* Icon with gradient background */}
            <div
              className={`
                flex h-14 w-14 items-center justify-center rounded-xl
                bg-gradient-to-br ${region.gradient} shadow-lg
                ${!disabled ? 'group-hover:scale-110 transition-transform duration-300' : ''}
              `}
            >
              <MapPin size={26} className="text-white" />
            </div>

            <div>
              <h3 className="text-xl font-bold text-white">{region.name}</h3>
              <p className="text-sm text-gray-400 mt-0.5">{region.description}</p>
            </div>
          </div>

          {disabled && (
            <div className="flex items-center gap-1.5 rounded-full bg-white/5 px-3 py-1">
              <Lock size={12} className="text-gray-500" />
              <span className="text-xs text-gray-500">Скоро</span>
            </div>
          )}
        </div>

        {/* Stats row */}
        {!disabled && region.tournamentCount > 0 && (
          <div className="mt-5 flex items-center gap-3">
            <span className={`text-sm font-medium bg-gradient-to-r ${region.gradient} bg-clip-text text-transparent`}>
              {region.tournamentCount} {pluralizeTournaments(region.tournamentCount)}
            </span>
            <span className="text-gray-600 text-xs">
              Нажмите для просмотра
            </span>
          </div>
        )}
      </div>
    </div>
  )
}

function pluralizeTournaments(n: number): string {
  const mod10 = n % 10
  const mod100 = n % 100
  if (mod10 === 1 && mod100 !== 11) return 'турнир'
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) return 'турнира'
  return 'турниров'
}
