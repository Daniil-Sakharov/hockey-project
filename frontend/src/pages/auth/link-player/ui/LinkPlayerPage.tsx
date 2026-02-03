import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { User, Calendar, Search, Loader2, CheckCircle, AlertCircle, Link2 } from 'lucide-react'
import { useAuthStore, usePlayerDashboardStore } from '@/shared/stores'
import { MOCK_PLAYER_PROFILE } from '@/shared/mocks'
import { cn } from '@/shared/lib/utils'

export function LinkPlayerPage() {
  const navigate = useNavigate()
  const { user, linkPlayer, isLoading, error, clearError } = useAuthStore()
  const { setLinkedPlayer } = usePlayerDashboardStore()

  const [fullName, setFullName] = useState('')
  const [birthDate, setBirthDate] = useState('')
  const [searchResult, setSearchResult] = useState<'idle' | 'found' | 'not_found'>('idle')

  // Если пользователь не авторизован - редирект на логин
  if (!user) {
    navigate('/login')
    return null
  }

  // Если игрок уже привязан - редирект на dashboard
  if (user.linkedPlayerId) {
    navigate('/player')
    return null
  }

  const handleSearch = async (e: FormEvent) => {
    e.preventDefault()
    clearError()
    setSearchResult('idle')

    // Симуляция поиска
    await new Promise((resolve) => setTimeout(resolve, 1000))

    // Mock: проверяем ФИО (в реальности - запрос к API)
    // Для демо ищем по части имени из mock данных
    const mockName = MOCK_PLAYER_PROFILE.name.toLowerCase()
    const searchName = fullName.toLowerCase()

    if (mockName.includes(searchName) || searchName.includes('иванов')) {
      setSearchResult('found')
    } else {
      setSearchResult('not_found')
    }
  }

  const handleLink = async () => {
    const success = await linkPlayer(MOCK_PLAYER_PROFILE.id, fullName, birthDate)
    if (success) {
      // Загружаем данные игрока в store
      setLinkedPlayer(MOCK_PLAYER_PROFILE)
      navigate('/player')
    }
  }

  const handleSkip = () => {
    // Можно пропустить и привязать позже
    navigate('/player')
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-[#0a0e1a] via-[#0f1629] to-[#0a0e1a] px-4 py-8">
      {/* Background effects */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/3 -left-20 w-96 h-96 bg-[#00d4ff]/10 rounded-full blur-3xl" />
        <div className="absolute bottom-1/3 -right-20 w-96 h-96 bg-[#10b981]/10 rounded-full blur-3xl" />
      </div>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="relative w-full max-w-lg"
      >
        {/* Card */}
        <div className="glass-card rounded-2xl p-8 backdrop-blur-xl border border-white/10">
          {/* Header */}
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-[#00d4ff] to-[#10b981] mb-4">
              <Link2 className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-2xl font-bold text-white mb-2">Привяжите свой профиль</h1>
            <p className="text-gray-400 text-sm">
              Найдите себя в базе игроков по ФИО и дате рождения
            </p>
          </div>

          {/* Error message */}
          {error && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              className="mb-6 p-4 rounded-lg bg-red-500/10 border border-red-500/30"
            >
              <p className="text-red-400 text-sm text-center">{error}</p>
            </motion.div>
          )}

          {/* Search Form */}
          <form onSubmit={handleSearch} className="space-y-5">
            {/* Full Name */}
            <div>
              <label htmlFor="fullName" className="block text-sm font-medium text-gray-300 mb-2">
                ФИО игрока
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <User className="w-5 h-5 text-gray-500" />
                </div>
                <input
                  id="fullName"
                  type="text"
                  value={fullName}
                  onChange={(e) => {
                    setFullName(e.target.value)
                    setSearchResult('idle')
                    clearError()
                  }}
                  placeholder="Иванов Александр Сергеевич"
                  className={cn(
                    'w-full pl-12 pr-4 py-3 rounded-xl',
                    'bg-white/5 border border-white/10',
                    'text-white placeholder-gray-500',
                    'focus:outline-none focus:border-[#00d4ff]/50 focus:ring-1 focus:ring-[#00d4ff]/50',
                    'transition-all duration-200'
                  )}
                  required
                />
              </div>
            </div>

            {/* Birth Date */}
            <div>
              <label htmlFor="birthDate" className="block text-sm font-medium text-gray-300 mb-2">
                Дата рождения
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <Calendar className="w-5 h-5 text-gray-500" />
                </div>
                <input
                  id="birthDate"
                  type="date"
                  value={birthDate}
                  onChange={(e) => {
                    setBirthDate(e.target.value)
                    setSearchResult('idle')
                    clearError()
                  }}
                  className={cn(
                    'w-full pl-12 pr-4 py-3 rounded-xl',
                    'bg-white/5 border border-white/10',
                    'text-white placeholder-gray-500',
                    'focus:outline-none focus:border-[#00d4ff]/50 focus:ring-1 focus:ring-[#00d4ff]/50',
                    'transition-all duration-200',
                    '[color-scheme:dark]'
                  )}
                  required
                />
              </div>
            </div>

            {/* Search button */}
            <button
              type="submit"
              disabled={isLoading || !fullName || !birthDate}
              className={cn(
                'w-full py-3 rounded-xl font-medium text-white',
                'bg-gradient-to-r from-[#00d4ff] to-[#10b981]',
                'hover:shadow-[0_0_30px_rgba(0,212,255,0.3)]',
                'focus:outline-none focus:ring-2 focus:ring-[#00d4ff]/50 focus:ring-offset-2 focus:ring-offset-[#0a0e1a]',
                'transition-all duration-300',
                'disabled:opacity-50 disabled:cursor-not-allowed',
                'flex items-center justify-center gap-2'
              )}
            >
              {isLoading ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin" />
                  Поиск...
                </>
              ) : (
                <>
                  <Search className="w-5 h-5" />
                  Найти игрока
                </>
              )}
            </button>
          </form>

          {/* Search Results */}
          {searchResult === 'found' && (
            <motion.div
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className="mt-6 p-4 rounded-xl bg-[#10b981]/10 border border-[#10b981]/30"
            >
              <div className="flex items-start gap-4">
                <div className="flex-shrink-0 w-12 h-12 rounded-full bg-[#10b981]/20 flex items-center justify-center">
                  <CheckCircle className="w-6 h-6 text-[#10b981]" />
                </div>
                <div className="flex-1">
                  <h3 className="font-semibold text-white">{MOCK_PLAYER_PROFILE.name}</h3>
                  <p className="text-sm text-gray-400 mt-1">
                    {MOCK_PLAYER_PROFILE.team} • {MOCK_PLAYER_PROFILE.position === 'forward' ? 'Нападающий' : 'Защитник'}
                  </p>
                  <p className="text-sm text-gray-500 mt-1">
                    {MOCK_PLAYER_PROFILE.region} • #{MOCK_PLAYER_PROFILE.jerseyNumber}
                  </p>
                </div>
              </div>
              <button
                onClick={handleLink}
                disabled={isLoading}
                className={cn(
                  'w-full mt-4 py-2.5 rounded-lg font-medium text-white',
                  'bg-[#10b981] hover:bg-[#10b981]/90',
                  'transition-all duration-200',
                  'disabled:opacity-50 disabled:cursor-not-allowed',
                  'flex items-center justify-center gap-2'
                )}
              >
                {isLoading ? (
                  <Loader2 className="w-5 h-5 animate-spin" />
                ) : (
                  <>
                    <Link2 className="w-4 h-4" />
                    Это я! Привязать профиль
                  </>
                )}
              </button>
            </motion.div>
          )}

          {searchResult === 'not_found' && (
            <motion.div
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className="mt-6 p-4 rounded-xl bg-amber-500/10 border border-amber-500/30"
            >
              <div className="flex items-start gap-4">
                <div className="flex-shrink-0 w-12 h-12 rounded-full bg-amber-500/20 flex items-center justify-center">
                  <AlertCircle className="w-6 h-6 text-amber-500" />
                </div>
                <div>
                  <h3 className="font-semibold text-white">Игрок не найден</h3>
                  <p className="text-sm text-gray-400 mt-1">
                    Проверьте правильность ФИО и даты рождения. Данные должны совпадать с
                    официальными записями в реестре игроков.
                  </p>
                </div>
              </div>
            </motion.div>
          )}

          {/* Skip link */}
          <div className="mt-8 text-center">
            <button
              onClick={handleSkip}
              className="text-gray-500 hover:text-gray-300 text-sm transition-colors"
            >
              Пропустить и привязать позже
            </button>
          </div>
        </div>

        {/* Help text */}
        <div className="mt-6 p-4 rounded-xl bg-white/5 border border-white/10">
          <p className="text-gray-400 text-sm text-center">
            Не можете найти себя? Возможно, ваши данные ещё не загружены в систему.{' '}
            <button className="text-[#00d4ff] hover:underline">Написать в поддержку</button>
          </p>
        </div>
      </motion.div>
    </div>
  )
}
