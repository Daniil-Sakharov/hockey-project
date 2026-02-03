import { memo } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { MapPin, Building2, Landmark, Mountain, Lock } from 'lucide-react'
import { useTournaments } from '@/shared/api/useExploreQueries'

const REGIONS = [
  {
    id: 'pfo',
    name: 'Приволжский',
    source: 'junior',
    description: 'Приволжский Федеральный Округ',
    icon: MapPin,
    gradient: 'from-blue-500 via-cyan-400 to-blue-600',
    bgPattern: 'radial-gradient(circle at 80% 20%, rgba(56,189,248,0.15) 0%, transparent 50%)',
    available: true,
  },
  {
    id: 'moscow',
    name: 'Москва',
    source: 'fhmoscow',
    description: 'Москва и Московская область',
    icon: Building2,
    gradient: 'from-red-500 via-orange-400 to-red-600',
    bgPattern: 'radial-gradient(circle at 80% 20%, rgba(239,68,68,0.15) 0%, transparent 50%)',
    available: false,
  },
  {
    id: 'spb',
    name: 'Санкт-Петербург',
    source: 'fhspb',
    description: 'Санкт-Петербург и Ленинградская область',
    icon: Landmark,
    gradient: 'from-purple-500 via-indigo-400 to-purple-600',
    bgPattern: 'radial-gradient(circle at 80% 20%, rgba(139,92,246,0.15) 0%, transparent 50%)',
    available: false,
  },
  {
    id: 'ural',
    name: 'Уральский',
    source: 'ural',
    description: 'Уральский Федеральный Округ',
    icon: Mountain,
    gradient: 'from-emerald-500 via-teal-400 to-emerald-600',
    bgPattern: 'radial-gradient(circle at 80% 20%, rgba(16,185,129,0.15) 0%, transparent 50%)',
    available: false,
  },
] as const

export const TournamentsListPage = memo(function TournamentsListPage() {
  const { data: juniorTournaments } = useTournaments('junior')

  return (
    <div className="space-y-8">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <h1 className="text-3xl font-bold text-white">Турниры</h1>
        <p className="text-gray-400 mt-1">Выберите регион для просмотра турниров</p>
      </motion.div>

      {/* Region cards grid */}
      <div className="grid gap-6 sm:grid-cols-2">
        {REGIONS.map((region, i) => {
          const Icon = region.icon
          const tournamentCount = region.source === 'junior'
            ? (juniorTournaments?.length ?? 0)
            : 0

          return (
            <motion.div
              key={region.id}
              initial={{ opacity: 0, y: 30 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 + i * 0.08, duration: 0.4 }}
            >
              {region.available ? (
                <Link to={`/explore/tournaments/${region.id}`}>
                  <RegionCard
                    region={region}
                    icon={Icon}
                    tournamentCount={tournamentCount}
                  />
                </Link>
              ) : (
                <RegionCard
                  region={region}
                  icon={Icon}
                  tournamentCount={0}
                  disabled
                />
              )}
            </motion.div>
          )
        })}
      </div>
    </div>
  )
})

interface RegionCardProps {
  region: (typeof REGIONS)[number]
  icon: React.ComponentType<{ size?: number; className?: string }>
  tournamentCount: number
  disabled?: boolean
}

function RegionCard({ region, icon: Icon, tournamentCount, disabled }: RegionCardProps) {
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
              <Icon size={26} className="text-white" />
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
        {!disabled && tournamentCount > 0 && (
          <div className="mt-5 flex items-center gap-3">
            <span className={`text-sm font-medium bg-gradient-to-r ${region.gradient} bg-clip-text text-transparent`}>
              {tournamentCount} {pluralizeTournaments(tournamentCount)}
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
