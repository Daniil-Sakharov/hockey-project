import { useQuery } from '@tanstack/react-query'
import {
  getExploreOverview,
  getTournaments,
  getTournamentStandings,
  getTournamentMatches,
  getTournamentScorers,
  getSeasons,
  searchPlayers,
  getPlayerProfile,
  getPlayerStats,
  getTeamProfile,
  getRecentResults,
  getUpcomingMatches,
  getRankings,
  getRankingsFilters,
  getMatchDetail,
} from './exploreApi'
import type { RankingsParams } from './exploreApi'

export function useExploreOverview() {
  return useQuery({
    queryKey: ['explore', 'overview'],
    queryFn: getExploreOverview,
    staleTime: 60_000,
  })
}

export function useTournaments(source?: string, domain?: string) {
  return useQuery({
    queryKey: ['explore', 'tournaments', source, domain],
    queryFn: () => getTournaments(source, domain),
    staleTime: 60_000,
  })
}

export function useTournamentStandings(id: string, birthYear?: number, groupName?: string) {
  return useQuery({
    queryKey: ['explore', 'tournaments', id, 'standings', birthYear, groupName],
    queryFn: () => getTournamentStandings(id, birthYear, groupName),
    enabled: !!id,
  })
}

export function useTournamentMatches(id: string, limit?: number, birthYear?: number, groupName?: string) {
  return useQuery({
    queryKey: ['explore', 'tournaments', id, 'matches', limit, birthYear, groupName],
    queryFn: () => getTournamentMatches(id, limit, birthYear, groupName),
    enabled: !!id,
  })
}

export function useTournamentScorers(id: string, limit?: number, birthYear?: number, groupName?: string) {
  return useQuery({
    queryKey: ['explore', 'tournaments', id, 'scorers', limit, birthYear, groupName],
    queryFn: () => getTournamentScorers(id, limit, birthYear, groupName),
    enabled: !!id,
  })
}

export function useSeasons() {
  return useQuery({
    queryKey: ['explore', 'seasons'],
    queryFn: getSeasons,
    staleTime: 300_000,
  })
}

export function usePlayersSearch(
  q: string,
  position: string,
  season: string,
  birthYear: number,
  limit: number,
  offset: number,
) {
  return useQuery({
    queryKey: ['explore', 'players', q, position, season, birthYear, limit, offset],
    queryFn: () => searchPlayers({ q, position, season, birthYear, limit, offset }),
    staleTime: 30_000,
  })
}

export function usePlayerProfile(id: string, season?: string) {
  return useQuery({
    queryKey: ['explore', 'players', id, 'profile', season],
    queryFn: () => getPlayerProfile(id, season),
    enabled: !!id,
  })
}

export function usePlayerStats(id: string) {
  return useQuery({
    queryKey: ['explore', 'players', id, 'stats'],
    queryFn: () => getPlayerStats(id),
    enabled: !!id,
  })
}

export function useTeamProfile(id: string) {
  return useQuery({
    queryKey: ['explore', 'teams', id],
    queryFn: () => getTeamProfile(id),
    enabled: !!id,
  })
}

export function useRecentResults(tournament?: string, limit?: number) {
  return useQuery({
    queryKey: ['explore', 'results', tournament, limit],
    queryFn: () => getRecentResults(tournament, limit),
    staleTime: 30_000,
  })
}

export function useUpcomingMatches(tournament?: string, limit?: number) {
  return useQuery({
    queryKey: ['explore', 'calendar', tournament, limit],
    queryFn: () => getUpcomingMatches(tournament, limit),
    staleTime: 30_000,
  })
}

export function useRankings(params: RankingsParams = {}) {
  return useQuery({
    queryKey: ['explore', 'rankings', params],
    queryFn: () => getRankings(params),
    staleTime: 30_000,
  })
}

export function useRankingsFilters() {
  return useQuery({
    queryKey: ['explore', 'rankings', 'filters'],
    queryFn: getRankingsFilters,
    staleTime: 60_000,
  })
}

export function useMatchDetail(id: string) {
  return useQuery({
    queryKey: ['explore', 'matches', id],
    queryFn: () => getMatchDetail(id),
    enabled: !!id,
  })
}
