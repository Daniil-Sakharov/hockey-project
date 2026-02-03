import { memo, useState } from 'react'
import { motion } from 'framer-motion'
import { Zap, Trophy, Calendar, Check, X } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import { cn } from '@/shared/lib/utils'
import {
  MOCK_PREDICTION_MATCHES,
  MOCK_PREDICTORS_RANKING,
  type MockPredictionMatch,
} from '@/shared/mocks/explorePlayers'

export const PredictionsPage = memo(function PredictionsPage() {
  const [predictions, setPredictions] = useState<
    Record<string, { home: string; away: string }>
  >({})
  const [submitted, setSubmitted] = useState<Set<string>>(new Set())

  const handleChange = (matchId: string, side: 'home' | 'away', value: string) => {
    if (submitted.has(matchId)) return
    setPredictions((prev) => ({
      ...prev,
      [matchId]: { ...prev[matchId], [side]: value },
    }))
  }

  const handleSubmit = (matchId: string) => {
    const pred = predictions[matchId]
    if (!pred?.home || !pred?.away) return
    setSubmitted((prev) => new Set(prev).add(matchId))
  }

  return (
    <div className="space-y-6">
      <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }}>
        <h1 className="text-2xl font-bold text-white flex items-center gap-2">
          <Zap size={24} className="text-[#f59e0b]" />
          Прогнозы
        </h1>
        <p className="text-gray-400">Угадайте счёт и зарабатывайте очки</p>
      </motion.div>

      {/* Rules */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        <GlassCard className="p-4">
          <div className="flex flex-wrap gap-4 text-xs text-gray-400">
            <RuleItem icon={<Check size={14} />} text="Точный счёт = 5 очков" color="text-[#10b981]" />
            <RuleItem icon={<Check size={14} />} text="Угадал победителя = 2 очка" color="text-[#00d4ff]" />
            <RuleItem icon={<X size={14} />} text="Не угадал = 0 очков" color="text-gray-600" />
          </div>
        </GlassCard>
      </motion.div>

      {/* Upcoming matches for prediction */}
      <div className="space-y-3">
        <h2 className="text-lg font-semibold text-white">Ближайшие матчи</h2>
        {MOCK_PREDICTION_MATCHES.map((match, i) => (
          <PredictionCard
            key={match.id}
            match={match}
            prediction={predictions[match.id]}
            isSubmitted={submitted.has(match.id)}
            onChange={(side, val) => handleChange(match.id, side, val)}
            onSubmit={() => handleSubmit(match.id)}
            delay={0.15 + i * 0.05}
          />
        ))}
      </div>

      {/* Leaderboard */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
      >
        <GlassCard className="p-6" glowColor="purple">
          <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <Trophy size={20} className="text-[#f59e0b]" />
            Таблица лидеров
          </h3>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="text-gray-500 border-b border-white/5">
                  <th className="text-left py-2 pr-3">#</th>
                  <th className="text-left py-2 pr-3">Прогнозист</th>
                  <th className="text-center py-2 px-2">Точных</th>
                  <th className="text-center py-2 px-2">Исходов</th>
                  <th className="text-center py-2 px-2">Всего</th>
                  <th className="text-center py-2 pl-2 font-semibold text-[#f59e0b]">Очки</th>
                </tr>
              </thead>
              <tbody>
                {MOCK_PREDICTORS_RANKING.map((pred) => (
                  <tr key={pred.position} className="border-b border-white/5">
                    <td className="py-2.5 pr-3">
                      {pred.position <= 3 ? (
                        <Trophy
                          size={16}
                          className={
                            pred.position === 1 ? 'text-[#f59e0b]' :
                            pred.position === 2 ? 'text-gray-300' :
                            'text-[#cd7f32]'
                          }
                        />
                      ) : (
                        <span className="text-gray-500">{pred.position}</span>
                      )}
                    </td>
                    <td className="py-2.5 pr-3 font-medium text-white">{pred.name}</td>
                    <td className="py-2.5 px-2 text-center text-[#10b981]">{pred.correctScores}</td>
                    <td className="py-2.5 px-2 text-center text-[#00d4ff]">{pred.correctOutcomes}</td>
                    <td className="py-2.5 px-2 text-center text-gray-400">{pred.totalPredictions}</td>
                    <td className="py-2.5 pl-2 text-center font-bold text-[#f59e0b]">{pred.points}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </GlassCard>
      </motion.div>
    </div>
  )
})

function RuleItem({ icon, text, color }: { icon: React.ReactNode; text: string; color: string }) {
  return (
    <span className={cn('inline-flex items-center gap-1.5', color)}>
      {icon}
      {text}
    </span>
  )
}

interface PredictionCardProps {
  match: MockPredictionMatch
  prediction?: { home: string; away: string }
  isSubmitted: boolean
  onChange: (side: 'home' | 'away', val: string) => void
  onSubmit: () => void
  delay: number
}

function PredictionCard({ match, prediction, isSubmitted, onChange, onSubmit, delay }: PredictionCardProps) {
  const canSubmit = prediction?.home !== undefined && prediction?.away !== undefined &&
    prediction.home !== '' && prediction.away !== ''

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay }}
    >
      <GlassCard className={cn('p-4', isSubmitted && 'ring-1 ring-[#10b981]/30')}>
        <div className="flex items-center justify-between mb-2">
          <span className="text-xs text-gray-500 flex items-center gap-1">
            <Calendar size={12} />
            {new Date(match.date).toLocaleDateString('ru-RU', {
              day: 'numeric',
              month: 'short',
            })}{' '}
            {match.time}
          </span>
          <span className="text-xs text-gray-600">{match.tournament}</span>
        </div>

        <div className="flex items-center gap-3">
          <span className="text-sm font-medium text-white flex-1 truncate text-right">
            {match.homeTeam}
          </span>

          <div className="flex items-center gap-2">
            <input
              type="number"
              min={0}
              max={99}
              value={prediction?.home ?? ''}
              onChange={(e) => onChange('home', e.target.value)}
              disabled={isSubmitted}
              className={cn(
                'w-10 h-10 rounded-lg text-center text-lg font-bold',
                'bg-white/5 border border-white/10 text-white',
                'focus:outline-none focus:border-[#00d4ff]/50',
                'disabled:opacity-50',
                '[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none'
              )}
            />
            <span className="text-gray-600 font-bold">:</span>
            <input
              type="number"
              min={0}
              max={99}
              value={prediction?.away ?? ''}
              onChange={(e) => onChange('away', e.target.value)}
              disabled={isSubmitted}
              className={cn(
                'w-10 h-10 rounded-lg text-center text-lg font-bold',
                'bg-white/5 border border-white/10 text-white',
                'focus:outline-none focus:border-[#00d4ff]/50',
                'disabled:opacity-50',
                '[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none'
              )}
            />
          </div>

          <span className="text-sm font-medium text-white flex-1 truncate">
            {match.awayTeam}
          </span>
        </div>

        <div className="flex justify-center mt-3">
          {isSubmitted ? (
            <span className="inline-flex items-center gap-1 text-xs text-[#10b981]">
              <Check size={14} />
              Прогноз принят
            </span>
          ) : (
            <button
              onClick={onSubmit}
              disabled={!canSubmit}
              className={cn(
                'px-4 py-1.5 rounded-lg text-xs font-medium transition-all',
                canSubmit
                  ? 'bg-[#f59e0b]/20 text-[#f59e0b] hover:bg-[#f59e0b]/30 border border-[#f59e0b]/30'
                  : 'bg-white/5 text-gray-600 cursor-not-allowed'
              )}
            >
              Отправить прогноз
            </button>
          )}
        </div>
      </GlassCard>
    </motion.div>
  )
}
