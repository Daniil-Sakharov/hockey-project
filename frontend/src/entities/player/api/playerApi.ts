import { apiClient } from '@/shared/api'
import type { TopScorer } from '../model/types'

interface TopScorersResponse {
  players: TopScorer[]
}

export async function getTopScorers(limit = 5): Promise<TopScorer[]> {
  const response = await apiClient.get<TopScorersResponse>(
    `/rankings/scorers?limit=${limit}`
  )
  return response.data.players
}
