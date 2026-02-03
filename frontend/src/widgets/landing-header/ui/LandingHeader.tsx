import { memo, useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { LogIn, UserPlus } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useAuthStore } from '@/shared/stores'

export const LandingHeader = memo(function LandingHeader() {
  const [scrolled, setScrolled] = useState(false)
  const { isAuthenticated, user } = useAuthStore()

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 50)
    }

    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  return (
    <motion.header
      initial={{ y: -100, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.6, delay: 0.2 }}
      className={cn(
        'fixed top-0 left-0 right-0 z-50 transition-all duration-300',
        scrolled
          ? 'bg-[#0a0e1a]/90 backdrop-blur-xl border-b border-white/5 shadow-lg'
          : 'bg-transparent'
      )}
    >
      <div className="max-w-7xl mx-auto px-6 py-4">
        <div className="flex items-center justify-between">
          {/* Logo */}
          <Link to="/" className="flex items-center gap-3 group">
            <div className="relative">
              <div className="absolute inset-0 bg-[#00d4ff] blur-lg opacity-30 group-hover:opacity-50 transition-opacity rounded-full" />
              <img
                src="/logo.png"
                alt="StarRink"
                className="relative h-10 w-10 object-cover rounded-full"
              />
            </div>
            <div className="flex flex-col">
              <span className="text-xl font-bold text-white tracking-tight">
                Star<span className="text-[#00d4ff]">Rink</span>
              </span>
              <span className="text-[10px] text-gray-500 uppercase tracking-widest -mt-1">
                Hockey Platform
              </span>
            </div>
          </Link>

          {/* Navigation */}
          <nav className="hidden md:flex items-center gap-8">
            <a
              href="#features"
              className="text-sm text-gray-400 hover:text-white transition-colors"
            >
              Возможности
            </a>
            <a
              href="#pricing"
              className="text-sm text-gray-400 hover:text-white transition-colors"
            >
              Тарифы
            </a>
            <a
              href="#about"
              className="text-sm text-gray-400 hover:text-white transition-colors"
            >
              О нас
            </a>
          </nav>

          {/* Auth Buttons */}
          <div className="flex items-center gap-3">
            {isAuthenticated ? (
              <Link
                to="/player"
                className={cn(
                  'flex items-center gap-2 px-5 py-2.5 rounded-xl',
                  'bg-gradient-to-r from-[#00d4ff] to-[#8b5cf6]',
                  'text-white font-medium text-sm',
                  'hover:shadow-[0_0_30px_rgba(0,212,255,0.4)]',
                  'transition-all duration-300'
                )}
              >
                <span>Мой профиль</span>
              </Link>
            ) : (
              <>
                <Link
                  to="/login"
                  className={cn(
                    'flex items-center gap-2 px-4 py-2.5 rounded-xl',
                    'text-gray-300 hover:text-white',
                    'border border-white/10 hover:border-white/20',
                    'hover:bg-white/5',
                    'transition-all duration-200',
                    'text-sm font-medium'
                  )}
                >
                  <LogIn size={16} />
                  <span className="hidden sm:inline">Войти</span>
                </Link>
                <Link
                  to="/register"
                  className={cn(
                    'flex items-center gap-2 px-5 py-2.5 rounded-xl',
                    'bg-gradient-to-r from-[#00d4ff] to-[#8b5cf6]',
                    'text-white font-medium text-sm',
                    'hover:shadow-[0_0_30px_rgba(0,212,255,0.4)]',
                    'transition-all duration-300'
                  )}
                >
                  <UserPlus size={16} />
                  <span className="hidden sm:inline">Регистрация</span>
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </motion.header>
  )
})
