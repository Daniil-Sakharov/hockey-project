import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Mail, Lock, Eye, EyeOff, Loader2, UserCircle2 } from 'lucide-react'
import { useAuthStore } from '@/shared/stores'
import { cn } from '@/shared/lib/utils'

export function LoginPage() {
  const navigate = useNavigate()
  const { login, isLoading, error, clearError } = useAuthStore()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    const success = await login(email, password)
    if (success) {
      const user = useAuthStore.getState().user
      if (!user) return

      // Роутинг по роли
      switch (user.role) {
        case 'player':
          navigate(user.linkedPlayerId ? '/player' : '/link-player')
          break
        case 'scout':
        case 'coach':
          navigate('/dashboard')
          break
        default:
          navigate('/explore')
      }
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-[#0a0e1a] via-[#0f1629] to-[#0a0e1a] px-4">
      {/* Background effects */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 -left-20 w-96 h-96 bg-[#00d4ff]/10 rounded-full blur-3xl" />
        <div className="absolute bottom-1/4 -right-20 w-96 h-96 bg-[#8b5cf6]/10 rounded-full blur-3xl" />
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
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-[#00d4ff] to-[#8b5cf6] mb-4">
              <UserCircle2 className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-2xl font-bold text-white mb-2">Добро пожаловать</h1>
            <p className="text-gray-400 text-sm">Войдите в свой аккаунт игрока</p>
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
                  }}
                  placeholder="player@example.com"
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
                  }}
                  placeholder="Введите пароль"
                  className={cn(
                    'w-full pl-12 pr-12 py-3 rounded-xl',
                    'bg-white/5 border border-white/10',
                    'text-white placeholder-gray-500',
                    'focus:outline-none focus:border-[#00d4ff]/50 focus:ring-1 focus:ring-[#00d4ff]/50',
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

            {/* Forgot password link */}
            <div className="flex justify-end">
              <button
                type="button"
                className="text-sm text-[#00d4ff] hover:text-[#00d4ff]/80 transition-colors"
              >
                Забыли пароль?
              </button>
            </div>

            {/* Submit button */}
            <button
              type="submit"
              disabled={isLoading}
              className={cn(
                'w-full py-3 rounded-xl font-medium text-white',
                'bg-gradient-to-r from-[#00d4ff] to-[#8b5cf6]',
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
                  Вход...
                </>
              ) : (
                'Войти'
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

          {/* Register link */}
          <p className="text-center text-gray-400">
            Нет аккаунта?{' '}
            <Link
              to="/register"
              className="text-[#00d4ff] hover:text-[#00d4ff]/80 font-medium transition-colors"
            >
              Зарегистрироваться
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
