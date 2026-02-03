import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Input } from '@/shared/ui'

export function QuickSearch() {
  const [query, setQuery] = useState('')
  const navigate = useNavigate()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (query.trim()) {
      navigate(`/search?q=${encodeURIComponent(query.trim())}`)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="w-full max-w-md">
      <Input
        type="text"
        placeholder="Поиск игроков..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        icon={<span className="text-gray-400">🔍</span>}
      />
    </form>
  )
}
