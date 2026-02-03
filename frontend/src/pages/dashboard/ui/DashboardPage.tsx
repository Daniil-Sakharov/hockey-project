import { motion } from 'framer-motion'
import { useScoutStore } from '@/shared/stores'
import { PlatformStatsKPI } from '@/widgets/platform-stats'
import { QuickSearchBar } from '@/widgets/quick-search'
import { TopPlayersWidget } from '@/widgets/top-players-widget'
import { WatchlistWidget } from '@/widgets/watchlist-widget'
import { DashboardLayout } from './DashboardLayout'

export function DashboardPage() {
  const scoutName = useScoutStore((state) => state.profile.name)

  return (
    <DashboardLayout>
      {/* Page header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="mb-8"
      >
        <h1 className="text-gradient text-3xl font-bold">Scout Dashboard</h1>
        <p className="mt-1 text-gray-400">
          Добро пожаловать, {scoutName}! Найдите талантливых игроков.
        </p>
      </motion.div>

      {/* Platform Stats KPI */}
      <PlatformStatsKPI className="mb-6" />

      {/* Quick Search */}
      <QuickSearchBar className="mb-8" />

      {/* Main content grid */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Top Players Widget */}
        <TopPlayersWidget />

        {/* Watchlist Widget */}
        <WatchlistWidget />
      </div>
    </DashboardLayout>
  )
}
