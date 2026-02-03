import { memo, useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Search, ArrowRight } from 'lucide-react'
import { cn } from '@/shared/lib/utils'

interface QuickSearchBarProps {
  className?: string
}

export const QuickSearchBar = memo(function QuickSearchBar({
  className,
}: QuickSearchBarProps) {
  const [query, setQuery] = useState('')
  const [isFocused, setIsFocused] = useState(false)
  const navigate = useNavigate()

  const handleSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault()
      if (query.trim()) {
        navigate(`/dashboard/search?q=${encodeURIComponent(query.trim())}`)
      }
    },
    [query, navigate]
  )

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (e.key === 'Enter' && query.trim()) {
        navigate(`/dashboard/search?q=${encodeURIComponent(query.trim())}`)
      }
    },
    [query, navigate]
  )

  return (
    <motion.form
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.3 }}
      onSubmit={handleSubmit}
      className={cn('relative', className)}
    >
      <div
        className={cn(
          'relative flex items-center rounded-xl border bg-[#0d1224]/80',
          'transition-all duration-300',
          isFocused
            ? 'border-[#00d4ff]/50 shadow-[0_0_30px_rgba(0,212,255,0.15)]'
            : 'border-white/10 hover:border-white/20'
        )}
      >
        {/* Search icon */}
        <div className="flex items-center pl-4">
          <Search
            size={20}
            className={cn(
              'transition-colors duration-200',
              isFocused ? 'text-[#00d4ff]' : 'text-gray-500'
            )}
          />
        </div>

        {/* Input */}
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          onKeyDown={handleKeyDown}
          placeholder="Найти игрока по имени, команде или городу..."
          className={cn(
            'flex-1 bg-transparent px-4 py-4 text-white',
            'placeholder:text-gray-500',
            'focus:outline-none'
          )}
        />

        {/* Submit button */}
        <motion.button
          type="submit"
          className={cn(
            'mr-2 flex items-center gap-2 rounded-lg px-4 py-2',
            'bg-[#00d4ff]/20 text-[#00d4ff]',
            'transition-all duration-200',
            'hover:bg-[#00d4ff]/30',
            'disabled:cursor-not-allowed disabled:opacity-50'
          )}
          disabled={!query.trim()}
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          <span className="hidden text-sm font-medium sm:inline">Найти</span>
          <ArrowRight size={16} />
        </motion.button>
      </div>

      {/* Glow effect */}
      {isFocused && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="absolute inset-0 -z-10 rounded-xl bg-[#00d4ff]/5 blur-xl"
        />
      )}

      {/* Keyboard shortcut hint */}
      <div className="mt-2 flex items-center gap-2 text-xs text-gray-500">
        <span>Нажмите</span>
        <kbd className="rounded border border-white/10 bg-white/5 px-1.5 py-0.5 font-mono">
          Enter
        </kbd>
        <span>для поиска</span>
      </div>
    </motion.form>
  )
})
