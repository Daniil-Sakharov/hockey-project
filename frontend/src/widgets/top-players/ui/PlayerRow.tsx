import type { TopScorer } from '@/entities/player'

interface PlayerRowProps {
  player: TopScorer
  rank: number
}

export function PlayerRow({ player, rank }: PlayerRowProps) {
  return (
    <tr className="hover:bg-gray-50">
      <td className="py-3 px-4 text-center font-medium text-gray-500">
        {rank}
      </td>
      <td className="py-3 px-4">
        <div>
          <p className="font-medium text-gray-900">{player.name}</p>
          <p className="text-sm text-gray-500">{player.team || '—'}</p>
        </div>
      </td>
      <td className="py-3 px-4 text-center font-bold text-primary-600">
        {player.goals}
      </td>
      <td className="py-3 px-4 text-center text-gray-600">
        {player.assists}
      </td>
      <td className="py-3 px-4 text-center text-gray-500">
        {player.games}
      </td>
    </tr>
  )
}
