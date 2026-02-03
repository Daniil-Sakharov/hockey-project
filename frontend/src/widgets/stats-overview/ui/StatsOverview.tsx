import { useQuery } from '@tanstack/react-query'
import { getStatsOverview } from '@/entities/stats'
import { StatCard } from './StatCard'

export function StatsOverview() {
  const { data, isLoading } = useQuery({
    queryKey: ['stats-overview'],
    queryFn: getStatsOverview,
  })

  return (
    <section className="py-8">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">
        Статистика базы данных
      </h2>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <StatCard
          title="Игроков"
          value={data?.players ?? 0}
          icon="👤"
          isLoading={isLoading}
        />
        <StatCard
          title="Команд"
          value={data?.teams ?? 0}
          icon="🏆"
          isLoading={isLoading}
        />
        <StatCard
          title="Турниров"
          value={data?.tournaments ?? 0}
          icon="🎯"
          isLoading={isLoading}
        />
      </div>
    </section>
  )
}
