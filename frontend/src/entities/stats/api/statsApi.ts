import { apiClient } from '@/shared/api'
import type { StatsOverview } from '../model/types'

export async function getStatsOverview(): Promise<StatsOverview> {
  const response = await apiClient.get<StatsOverview>('/stats/overview')
  return response.data
}
