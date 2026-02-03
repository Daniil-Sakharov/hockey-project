import { Link } from 'react-router-dom'
import { Button } from '@/shared/ui'

export function NotFoundPage() {
  return (
    <div className="container mx-auto px-4 py-16 text-center">
      <h1 className="text-6xl font-bold text-gray-300 mb-4">404</h1>
      <p className="text-xl text-gray-600 mb-8">Страница не найдена</p>
      <Link to="/">
        <Button>Вернуться на главную</Button>
      </Link>
    </div>
  )
}
