import { StatsOverview } from '@/widgets/stats-overview'
import { TopPlayers } from '@/widgets/top-players'
import { QuickSearch } from '@/features/search-quick'

export function HomePage() {
  return (
    <div className="container mx-auto px-4 py-8">
      {/* Hero section */}
      <section className="text-center py-12">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          Статистика юношеского хоккея России
        </h1>
        <p className="text-lg text-gray-600 mb-8 max-w-2xl mx-auto">
          Поиск игроков, статистика турниров, рейтинги бомбардиров
        </p>
        <div className="flex justify-center">
          <QuickSearch />
        </div>
      </section>

      {/* Stats overview */}
      <StatsOverview />

      {/* Top players */}
      <TopPlayers />
    </div>
  )
}
