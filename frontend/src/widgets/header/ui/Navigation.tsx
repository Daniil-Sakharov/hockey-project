import { NavLink } from 'react-router-dom'
import { cn } from '@/shared/lib/utils'

const navItems = [
  { to: '/', label: 'Главная' },
  { to: '/search', label: 'Поиск игроков' },
  { to: '/rankings', label: 'Рейтинги' },
]

export function Navigation() {
  return (
    <nav className="flex items-center gap-6">
      {navItems.map((item) => (
        <NavLink
          key={item.to}
          to={item.to}
          className={({ isActive }) =>
            cn(
              'text-sm font-medium transition-colors',
              isActive
                ? 'text-primary-600'
                : 'text-gray-600 hover:text-gray-900'
            )
          }
        >
          {item.label}
        </NavLink>
      ))}
    </nav>
  )
}
