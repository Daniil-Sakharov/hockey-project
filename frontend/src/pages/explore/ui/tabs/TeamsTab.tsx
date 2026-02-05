import { motion } from 'framer-motion'
import { Loader2, Shield } from 'lucide-react'
import { GlassCard } from '@/shared/ui'
import type { TeamItem } from '@/shared/api/exploreTypes'
import { TeamCard } from './TeamCard'

const containerVariants = {
  hidden: { opacity: 0 },
  show: {
    opacity: 1,
    transition: { staggerChildren: 0.06 },
  },
} as const

const itemVariants = {
  hidden: { opacity: 0, y: 30, scale: 0.9 },
  show: {
    opacity: 1,
    y: 0,
    scale: 1,
    transition: { type: 'spring' as const, stiffness: 260, damping: 22 },
  },
}

interface Props {
  teams: TeamItem[]
  isLoading: boolean
  tournamentId: string
  birthYear?: number
  groupName?: string
}

export function TeamsTab({ teams, isLoading, tournamentId, birthYear, groupName }: Props) {
  if (isLoading) {
    return (
      <div className="flex justify-center py-16">
        <Loader2 size={28} className="animate-spin text-gray-500" />
      </div>
    )
  }

  if (teams.length === 0) {
    return (
      <GlassCard className="p-12 text-center" glowColor="cyan">
        <Shield size={48} className="mx-auto mb-4 text-gray-600" />
        <h3 className="text-lg font-semibold text-white mb-2">Команды не найдены</h3>
        <p className="text-gray-400 text-sm">В данной группе пока нет команд</p>
      </GlassCard>
    )
  }

  return (
    <>
      <style>{`
        @keyframes gradientShift {
          0% { background-position: 0% 50%; }
          50% { background-position: 100% 50%; }
          100% { background-position: 0% 50%; }
        }
      `}</style>
      <motion.div
        variants={containerVariants}
        initial="hidden"
        animate="show"
        className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5"
      >
        {teams.map((team, index) => (
          <motion.div key={team.id} variants={itemVariants}>
            <TeamCard
              team={team}
              index={index}
              tournamentId={tournamentId}
              birthYear={birthYear}
              groupName={groupName}
            />
          </motion.div>
        ))}
      </motion.div>
    </>
  )
}
