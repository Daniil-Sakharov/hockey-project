import { memo } from 'react'
import { useParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { ArrowLeft, Loader2, Calendar } from 'lucide-react'
import { useMatchDetail } from '@/shared/api/useExploreQueries'
import { MatchHero } from './match/MatchHero'
import { MatchEvents } from './match/MatchEvents'
import { MatchLineups } from './match/MatchLineups'

export const MatchDetailPage = memo(function MatchDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { data: match, isLoading, error } = useMatchDetail(id ?? '')

  if (isLoading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 size={32} className="animate-spin text-gray-500" />
      </div>
    )
  }

  if (error || !match) {
    return (
      <div className="flex flex-col items-center justify-center py-20 text-center">
        <Calendar size={48} className="text-gray-700 mb-4" />
        <p className="text-gray-400 text-lg">Матч не найден</p>
        <Link to="/explore/results" className="mt-4 text-sm text-[#00d4ff] hover:underline">
          Вернуться к результатам
        </Link>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <Link
        to="/explore/results"
        className="inline-flex items-center gap-2 text-sm text-gray-400 hover:text-white transition-colors"
      >
        <ArrowLeft size={16} />
        Назад к матчам
      </Link>

      {/* Hero section with teams and score */}
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
        <MatchHero match={match} />
      </motion.div>

      {/* Events */}
      {match.status === 'finished' && match.events.length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
        >
          <MatchEvents match={match} />
        </motion.div>
      )}

      {/* Lineups */}
      {(match.homeLineup.length > 0 || match.awayLineup.length > 0) && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
        >
          <h2 className="text-xl font-bold text-white mb-4">Составы команд</h2>
          <MatchLineups
            homeLineup={match.homeLineup}
            awayLineup={match.awayLineup}
            homeTeam={match.homeTeam}
            awayTeam={match.awayTeam}
          />
        </motion.div>
      )}

      {/* Empty state for events */}
      {match.status === 'finished' && match.events.length === 0 && match.homeLineup.length === 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="text-center py-12 text-gray-500"
        >
          <p>Детальная статистика матча пока недоступна</p>
        </motion.div>
      )}
    </div>
  )
})
