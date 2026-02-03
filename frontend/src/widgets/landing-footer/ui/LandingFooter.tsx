import { memo } from 'react'
import { Link } from 'react-router-dom'
import { Send } from 'lucide-react'
import { cn } from '@/shared/lib/utils'

const footerLinks = {
  platform: {
    title: 'Платформа',
    links: [
      { label: 'Для игроков', href: '/#features' },
      { label: 'Для скаутов', href: '/scouts' },
      { label: 'Тарифы', href: '/#pricing' },
      { label: 'FAQ', href: '/faq' },
    ],
  },
  company: {
    title: 'Компания',
    links: [
      { label: 'О нас', href: '/#about' },
      { label: 'Контакты', href: '/contacts' },
      { label: 'Карьера', href: '/careers' },
      { label: 'Блог', href: '/blog' },
    ],
  },
  legal: {
    title: 'Правовое',
    links: [
      { label: 'Политика конфиденциальности', href: '/privacy' },
      { label: 'Пользовательское соглашение', href: '/terms' },
      { label: 'Обработка данных', href: '/data-processing' },
    ],
  },
}

export const LandingFooter = memo(function LandingFooter() {
  return (
    <footer className="relative bg-[#0a0e1a] border-t border-white/5">
      {/* Top gradient line */}
      <div className="absolute top-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-[#00d4ff]/30 to-transparent" />

      <div className="max-w-6xl mx-auto px-4 py-12 lg:py-16">
        {/* Main footer content */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 lg:gap-12 mb-12">
          {/* Logo column */}
          <div className="col-span-2 md:col-span-1">
            <Link to="/" className="flex items-center gap-2 mb-4">
              <img
                src="/logo.png"
                alt="StarRink"
                className="h-8 w-8 object-cover rounded-full"
              />
              <span className="text-lg font-bold text-white">
                Star<span className="text-[#00d4ff]">Rink</span>
              </span>
            </Link>

            <p className="text-gray-500 text-sm mb-6 max-w-xs">
              Платформа для хоккеистов. Покажи свой талант миру.
            </p>

            {/* Social links */}
            <div className="flex items-center gap-3">
              <a
                href="https://t.me/starrink"
                target="_blank"
                rel="noopener noreferrer"
                className={cn(
                  'w-9 h-9 rounded-lg flex items-center justify-center',
                  'bg-white/5 text-gray-500',
                  'hover:bg-[#00d4ff]/10 hover:text-[#00d4ff]',
                  'transition-all duration-200'
                )}
                aria-label="Telegram"
              >
                <Send className="w-4 h-4" />
              </a>
              <a
                href="https://vk.com/starrink"
                target="_blank"
                rel="noopener noreferrer"
                className={cn(
                  'w-9 h-9 rounded-lg flex items-center justify-center',
                  'bg-white/5 text-gray-500',
                  'hover:bg-[#8b5cf6]/10 hover:text-[#8b5cf6]',
                  'transition-all duration-200'
                )}
                aria-label="VKontakte"
              >
                <svg className="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12.785 16.241c-4.932 0-7.748-3.377-7.867-8.991h2.473c.082 4.126 1.9 5.875 3.34 6.234V7.25h2.327v3.559c1.422-.154 2.916-1.787 3.422-3.559h2.327c-.39 2.189-2.015 3.822-3.172 4.489 1.157.541 2.994 1.986 3.696 4.502h-2.57c-.548-1.71-1.916-3.03-3.703-3.208v3.208h-.273z" />
                </svg>
              </a>
            </div>
          </div>

          {/* Platform links */}
          <div>
            <h4 className="text-white font-semibold mb-4">
              {footerLinks.platform.title}
            </h4>
            <ul className="space-y-3">
              {footerLinks.platform.links.map((link) => (
                <li key={link.label}>
                  <Link
                    to={link.href}
                    className="text-gray-500 text-sm hover:text-[#00d4ff] transition-colors"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Company links */}
          <div>
            <h4 className="text-white font-semibold mb-4">
              {footerLinks.company.title}
            </h4>
            <ul className="space-y-3">
              {footerLinks.company.links.map((link) => (
                <li key={link.label}>
                  <Link
                    to={link.href}
                    className="text-gray-500 text-sm hover:text-[#00d4ff] transition-colors"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Legal links */}
          <div>
            <h4 className="text-white font-semibold mb-4">
              {footerLinks.legal.title}
            </h4>
            <ul className="space-y-3">
              {footerLinks.legal.links.map((link) => (
                <li key={link.label}>
                  <Link
                    to={link.href}
                    className="text-gray-500 text-sm hover:text-[#00d4ff] transition-colors"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Bottom bar */}
        <div className="pt-8 border-t border-white/5">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            <p className="text-gray-600 text-sm">
              © 2025 StarRink. Все права защищены.
            </p>
            <p className="text-gray-700 text-xs">
              Сделано с любовью к хоккею
            </p>
          </div>
        </div>
      </div>
    </footer>
  )
})
