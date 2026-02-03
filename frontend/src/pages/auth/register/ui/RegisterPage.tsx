import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Mail, Lock, Eye, EyeOff, Loader2, UserPlus, CheckCircle2 } from 'lucide-react'
import { useAuthStore } from '@/shared/stores'
import { cn } from '@/shared/lib/utils'

export function RegisterPage() {
  const navigate = useNavigate()
  const { register, isLoading, error, clearError } = useAuthStore()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [localError, setLocalError] = useState<string | null>(null)

  const passwordRequirements = [
    { label: 'Минимум 6 символов', met: password.length >= 6 },
    { label: 'Содержит цифру', met: /\d/.test(password) },
    { label: 'Пароли совпадают', met: password === confirmPassword && password.length > 0 },
  ]

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setLocalError(null)

    if (password !== confirmPassword) {
      setLocalError('Пароли не совпадают')
      return
    }

    if (password.length < 6) {
      setLocalError('Пароль должен быть не менее 6 символов')
      return
    }

    const success = await register(email, password)
    if (success) {
      navigate('/select-role')
    }
  }

  const displayError = localError || error

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-[#0a0e1a] via-[#0f1629] to-[#0a0e1a] px-4 py-8">
      {/* Background effects */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 -left-20 w-96 h-96 bg-[#8b5cf6]/10 rounded-full blur-3xl" />
        <div className="absolute bottom-1/4 -right-20 w-96 h-96 bg-[#00d4ff]/10 rounded-full blur-3xl" />
      </div>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="relative w-full max-w-md"
      >
        {/* Card */}
        <div className="glass-card rounded-2xl p-8 backdrop-blur-xl border border-white/10">
          {/* Logo/Header */}
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-[#8b5cf6] to-[#ec4899] mb-4">
              <UserPlus className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-2xl font-bold text-white mb-2">Создать аккаунт</h1>
            <p className="text-gray-400 text-sm">
              Зарегистрируйтесь, чтобы получить доступ к своему профилю игрока
            </p>
          </div>

          {/* Error message */}
          {displayError && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              className="mb-6 p-4 rounded-lg bg-red-500/10 border border-red-500/30"
            >
              <p className="text-red-400 text-sm text-center">{displayError}</p>
            </motion.div>
          )}

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-5">
            {/* Email */}
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-300 mb-2">
                Email
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <Mail className="w-5 h-5 text-gray-500" />
                </div>
                <input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => {
                    setEmail(e.target.value)
                    if (error) clearError()
                    setLocalError(null)
                  }}
                  placeholder="player@example.com"
                  className={cn(
                    'w-full pl-12 pr-4 py-3 rounded-xl',
                    'bg-white/5 border border-white/10',
                    'text-white placeholder-gray-500',
                    'focus:outline-none focus:border-[#8b5cf6]/50 focus:ring-1 focus:ring-[#8b5cf6]/50',
                    'transition-all duration-200'
                  )}
                  required
                />
              </div>
            </div>

            {/* Password */}
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-300 mb-2">
                Пароль
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <Lock className="w-5 h-5 text-gray-500" />
                </div>
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => {
                    setPassword(e.target.value)
                    if (error) clearError()
                    setLocalError(null)
                  }}
                  placeholder="Придумайте пароль"
                  className={cn(
                    'w-full pl-12 pr-12 py-3 rounded-xl',
                    'bg-white/5 border border-white/10',
                    'text-white placeholder-gray-500',
                    'focus:outline-none focus:border-[#8b5cf6]/50 focus:ring-1 focus:ring-[#8b5cf6]/50',
                    'transition-all duration-200'
                  )}
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-500 hover:text-gray-300 transition-colors"
                >
                  {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                </button>
              </div>
            </div>

            {/* Confirm Password */}
            <div>
              <label
                htmlFor="confirmPassword"
                className="block text-sm font-medium text-gray-300 mb-2"
              >
                Подтвердите пароль
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                  <Lock className="w-5 h-5 text-gray-500" />
                </div>
                <input
                  id="confirmPassword"
                  type={showPassword ? 'text' : 'password'}
                  value={confirmPassword}
                  onChange={(e) => {
                    setConfirmPassword(e.target.value)
                    setLocalError(null)
                  }}
                  placeholder="Повторите пароль"
                  className={cn(
                    'w-full pl-12 pr-4 py-3 rounded-xl',
                    'bg-white/5 border border-white/10',
                    'text-white placeholder-gray-500',
                    'focus:outline-none focus:border-[#8b5cf6]/50 focus:ring-1 focus:ring-[#8b5cf6]/50',
                    'transition-all duration-200'
                  )}
                  required
                />
              </div>
            </div>

            {/* Password requirements */}
            {password.length > 0 && (
              <motion.div
                initial={{ opacity: 0, height: 0 }}
                animate={{ opacity: 1, height: 'auto' }}
                className="space-y-2"
              >
                {passwordRequirements.map((req) => (
                  <div key={req.label} className="flex items-center gap-2 text-sm">
                    <CheckCircle2
                      className={cn(
                        'w-4 h-4 transition-colors',
                        req.met ? 'text-green-400' : 'text-gray-600'
                      )}
                    />
                    <span className={cn(req.met ? 'text-green-400' : 'text-gray-500')}>
                      {req.label}
                    </span>
                  </div>
                ))}
              </motion.div>
            )}

            {/* Submit button */}
            <button
              type="submit"
              disabled={isLoading}
              className={cn(
                'w-full py-3 rounded-xl font-medium text-white',
                'bg-gradient-to-r from-[#8b5cf6] to-[#ec4899]',
                'hover:shadow-[0_0_30px_rgba(139,92,246,0.3)]',
                'focus:outline-none focus:ring-2 focus:ring-[#8b5cf6]/50 focus:ring-offset-2 focus:ring-offset-[#0a0e1a]',
                'transition-all duration-300',
                'disabled:opacity-50 disabled:cursor-not-allowed',
                'flex items-center justify-center gap-2'
              )}
            >
              {isLoading ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin" />
                  Регистрация...
                </>
              ) : (
                'Создать аккаунт'
              )}
            </button>
          </form>

          {/* Divider */}
          <div className="relative my-8">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-white/10"></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-4 bg-[#0f1629] text-gray-500">или</span>
            </div>
          </div>

          {/* Login link */}
          <p className="text-center text-gray-400">
            Уже есть аккаунт?{' '}
            <Link
              to="/login"
              className="text-[#8b5cf6] hover:text-[#8b5cf6]/80 font-medium transition-colors"
            >
              Войти
            </Link>
          </p>
        </div>

        {/* Footer note */}
        <p className="text-center text-gray-500 text-xs mt-6">
          Регистрируясь, вы соглашаетесь с{' '}
          <button className="text-gray-400 hover:text-white transition-colors underline">
            условиями использования
          </button>
        </p>
      </motion.div>
    </div>
  )
}
