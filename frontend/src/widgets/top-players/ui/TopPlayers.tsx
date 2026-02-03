import { useQuery } from '@tanstack/react-query'
import { getTopScorers } from '@/entities/player'
import { Card, CardTitle, Skeleton } from '@/shared/ui'
import { PlayerRow } from './PlayerRow'

export function TopPlayers() {
  const { data: players, isLoading } = useQuery({
    queryKey: ['top-scorers'],
    queryFn: () => getTopScorers(5),
  })

  return (
    <section className="py-8">
      <Card>
        <CardTitle className="mb-4">Топ-5 бомбардиров</CardTitle>
        {isLoading ? (
          <div className="space-y-3">
            {[...Array(5)].map((_, i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        ) : players && players.length > 0 ? (
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 text-sm text-gray-500">
                <th className="py-2 px-4 text-center w-12">#</th>
                <th className="py-2 px-4 text-left">Игрок</th>
                <th className="py-2 px-4 text-center">Голы</th>
                <th className="py-2 px-4 text-center">Передачи</th>
                <th className="py-2 px-4 text-center">Игры</th>
              </tr>
            </thead>
            <tbody>
              {players.map((player, index) => (
                <PlayerRow key={player.id} player={player} rank={index + 1} />
              ))}
            </tbody>
          </table>
        ) : (
          <p className="text-gray-500 text-center py-8">
            Данные пока отсутствуют
          </p>
        )}
      </Card>
    </section>
  )
}
