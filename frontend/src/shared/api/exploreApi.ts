import { apiClient } from './client'
import type {
  ExploreOverview,
  TournamentItem,
  Standing,
  Scorer,
  MatchItem,
  PlayerItem,
  PlayerProfile,
  PlayerStatEntry,
  TeamProfile,
  TeamItem,
  RankingsData,
  RankingsFilters,
  MatchDetail,
} from './exploreTypes'

export async function getExploreOverview(): Promise<ExploreOverview> {
  const { data } = await apiClient.get('/explore/overview')
  return data
}

export async function getTournaments(source?: string, domain?: string): Promise<TournamentItem[]> {
  const params: Record<string, string> = {}
  if (source && source !== 'all') params.source = source
  if (domain) params.domain = domain
  const { data } = await apiClient.get('/explore/tournaments', { params })
  return data.tournaments
}

export async function getTournamentStandings(id: string, birthYear?: number, groupName?: string): Promise<Standing[]> {
  const params: Record<string, string | number> = {}
  if (birthYear) params.birthYear = birthYear
  if (groupName) params.group = groupName
  const { data } = await apiClient.get(`/explore/tournaments/${id}/standings`, { params })
  return data.standings
}

export async function getTournamentMatches(id: string, limit?: number, birthYear?: number, groupName?: string): Promise<MatchItem[]> {
  const params: Record<string, string | number> = {}
  if (limit) params.limit = limit
  if (birthYear) params.birthYear = birthYear
  if (groupName) params.group = groupName
  const { data } = await apiClient.get(`/explore/tournaments/${id}/matches`, { params })
  return data.matches
}

export async function getTournamentScorers(id: string, limit?: number, birthYear?: number, groupName?: string): Promise<Scorer[]> {
  const params: Record<string, string | number> = {}
  if (limit) params.limit = limit
  if (birthYear) params.birthYear = birthYear
  if (groupName) params.group = groupName
  const { data } = await apiClient.get(`/explore/tournaments/${id}/scorers`, { params })
  return data.scorers
}

export async function getTournamentTeams(id: string, birthYear?: number, groupName?: string): Promise<TeamItem[]> {
  const params: Record<string, string | number> = {}
  if (birthYear) params.birthYear = birthYear
  if (groupName) params.group = groupName
  const { data } = await apiClient.get(`/explore/tournaments/${id}/teams`, { params })
  return data.teams
}

export async function getSeasons(): Promise<string[]> {
  const { data } = await apiClient.get('/explore/seasons')
  return data.seasons
}

export async function searchPlayers(params: {
  q?: string
  position?: string
  season?: string
  birthYear?: number
  limit?: number
  offset?: number
}): Promise<{ players: PlayerItem[]; total: number }> {
  const query: Record<string, string | number> = {}
  if (params.q) query.q = params.q
  if (params.position && params.position !== 'all') query.position = params.position
  if (params.season) query.season = params.season
  if (params.birthYear && params.birthYear > 0) query.birthYear = params.birthYear
  if (params.limit) query.limit = params.limit
  if (params.offset) query.offset = params.offset
  const { data } = await apiClient.get('/explore/players', { params: query })
  return data
}

export async function getPlayerProfile(id: string, season?: string): Promise<PlayerProfile> {
  const params: Record<string, string> = {}
  if (season) params.season = season
  const { data } = await apiClient.get(`/explore/players/${id}`, { params })
  return data
}

export async function getPlayerStats(id: string): Promise<PlayerStatEntry[]> {
  const { data } = await apiClient.get(`/explore/players/${id}/stats`)
  return data.stats
}

export async function getTeamProfile(id: string): Promise<TeamProfile> {
  const { data } = await apiClient.get(`/explore/teams/${id}`)
  return data
}

export async function getRecentResults(tournament?: string, limit?: number): Promise<MatchItem[]> {
  const params: Record<string, string | number> = {}
  if (tournament && tournament !== 'all') params.tournament = tournament
  if (limit) params.limit = limit
  const { data } = await apiClient.get('/explore/results', { params })
  return data.matches
}

export async function getUpcomingMatches(tournament?: string, limit?: number): Promise<MatchItem[]> {
  const params: Record<string, string | number> = {}
  if (tournament && tournament !== 'all') params.tournament = tournament
  if (limit) params.limit = limit
  const { data } = await apiClient.get('/explore/calendar', { params })
  return data.matches
}

export interface RankingsParams {
  sort?: string
  limit?: number
  birthYear?: number
  domain?: string
  tournamentId?: string
  groupName?: string
}

export async function getRankings(params: RankingsParams = {}): Promise<RankingsData> {
  const query: Record<string, string | number> = {}
  if (params.sort) query.sort = params.sort
  if (params.limit) query.limit = params.limit
  if (params.birthYear) query.birthYear = params.birthYear
  if (params.domain) query.domain = params.domain
  if (params.tournamentId) query.tournamentId = params.tournamentId
  if (params.groupName) query.groupName = params.groupName
  const { data } = await apiClient.get('/explore/rankings', { params: query })
  return { season: data.season, players: data.players }
}

export async function getRankingsFilters(): Promise<RankingsFilters> {
  const { data } = await apiClient.get('/explore/rankings/filters')
  return data
}

export async function getMatchDetail(id: string): Promise<MatchDetail> {
  const { data } = await apiClient.get(`/explore/matches/${id}`)
  return {
    ...data,
    events: data.events ?? [],
    homeLineup: data.homeLineup ?? [],
    awayLineup: data.awayLineup ?? [],
  }
}
