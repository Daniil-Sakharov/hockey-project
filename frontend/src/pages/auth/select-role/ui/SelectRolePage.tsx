import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { UserRound, Search, ClipboardList, ArrowRight, Loader2 } from 'lucide-react'
import { useAuthStore } from '@/shared/stores'
import { cn } from '@/shared/lib/utils'
import type { UserRole } from '@/entities/user'

interface RoleOption {
  role: UserRole
  label: string
  description: string
  icon: React.ReactNode
  color: string
  gradient: string
  redirectTo: string
}

const ROLE_OPTIONS: RoleOption[] = [
  {
    role: 'player',
    label: 'Игрок',
    description: 'Отслеживайте свою статистику, достижения и прогресс',
    icon: <UserRound size={32} />,
    color: '#00d4ff',
    gradient: 'from-[#00d4ff] to-[#0ea5e9]',
    redirectTo: '/link-player',
  },
  {
    role: 'scout',
    label: 'Скаут',
    description: 'Находите талантливых игроков, ведите watchlist',
    icon: <Search size={32} />,
    color: '#8b5cf6',
    gradient: 'from-[#8b5cf6] to-[#7c3aed]',
    redirectTo: '/dashboard',
  },
  {
    role: 'coach',
    label: 'Тренер',
    description: 'Управляйте командой и анализируйте игроков',
    icon: <ClipboardList size={32} />,
    color: '#10b981',
    gradient: 'from-[#10b981] to-[#059669]',
    redirectTo: '/dashboard',
  },
]

export function SelectRolePage() {
  const navigate = useNavigate()
  const { updateRole } = useAuthStore()
  const [selected, setSelected] = useState<UserRole | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const handleSelect = async (option: RoleOption) => {
    setSelected(option.role)
    setIsLoading(true)

    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 600))

    updateRole(option.role)
    navigate(option.redirectTo)
  }

  const handleSkip = () => {
    navigate('/explore')
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-[#0a0e1a] via-[#0f1629] to-[#0a0e1a] px-4 py-8">
      {/* Background effects */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/3 -left-20 w-96 h-96 bg-[#00d4ff]/10 rounded-full blur-3xl" />
        <div className="absolute bottom-1/3 -right-20 w-96 h-96 bg-[#8b5cf6]/10 rounded-full blur-3xl" />
      </div>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="relative w-full max-w-2xl"
      >
        {/* Header */}
        <div className="text-center mb-10">
          <motion.h1
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="text-3xl font-bold text-white mb-3"
          >
            Кто вы в хоккее?
          </motion.h1>
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.3 }}
            className="text-gray-400"
          >
            Выберите роль для персонализированного опыта
          </motion.p>
        </div>

        {/* Role cards */}
        <div className="grid gap-4 md:grid-cols-3">
          {ROLE_OPTIONS.map((option, index) => (
            <motion.button
              key={option.role}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 + index * 0.1 }}
              onClick={() => handleSelect(option)}
              disabled={isLoading}
              className={cn(
                'group relative flex flex-col items-center gap-4 rounded-2xl p-6',
                'border backdrop-blur-xl transition-all duration-300',
                'disabled:opacity-50 disabled:cursor-not-allowed',
                selected === option.role
                  ? 'border-white/30 bg-white/10 scale-[1.02]'
                  : 'border-white/10 bg-white/5 hover:border-white/20 hover:bg-white/[0.07]'
              )}
            >
              {/* Icon */}
              <div
                className={cn(
                  'flex h-16 w-16 items-center justify-center rounded-2xl',
                  'bg-gradient-to-br transition-transform duration-300',
                  'group-hover:scale-110',
                  option.gradient
                )}
              >
                <span className="text-white">{option.icon}</span>
              </div>

              {/* Text */}
              <div className="text-center">
                <h3 className="text-lg font-semibold text-white mb-1">{option.label}</h3>
                <p className="text-sm text-gray-400 leading-relaxed">{option.description}</p>
              </div>

              {/* Loading indicator */}
              {selected === option.role && isLoading && (
                <Loader2 className="absolute top-4 right-4 w-5 h-5 text-white animate-spin" />
              )}
            </motion.button>
          ))}
        </div>

        {/* Skip button */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.7 }}
          className="mt-8 text-center"
        >
          <button
            onClick={handleSkip}
            disabled={isLoading}
            className={cn(
              'inline-flex items-center gap-2 text-gray-400',
              'hover:text-white transition-colors duration-200',
              'disabled:opacity-50 disabled:cursor-not-allowed'
            )}
          >
            Пропустить и продолжить как пользователь
            <ArrowRight size={16} />
          </button>
          <p className="mt-2 text-xs text-gray-600">
            Роль можно изменить позже в настройках
          </p>
        </motion.div>
      </motion.div>
    </div>
  )
}
