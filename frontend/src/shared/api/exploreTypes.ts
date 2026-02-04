// Explore API response types (match backend DTOs)

export interface ExploreOverview {
  players: number
  teams: number
  tournaments: number
  matches: number
}

export interface GroupStats {
  name: string
  teamsCount: number
  matchesCount: number
}

export interface TeamItem {
  id: string
  name: string
  city?: string
  logoUrl?: string
  playersCount: number
  groupName?: string
  birthYear?: number
}

export interface TournamentItem {
  id: string
  name: string
  domain: string
  season: string
  source: string
  birthYearGroups?: Record<string, GroupStats[]>
  teamsCount: number
  matchesCount: number
  isEnded: boolean
}

export interface Standing {
  position: number
  team: string
  teamId: string
  logoUrl?: string
  games: number
  wins: number
  winsOt: number
  losses: number
  lossesOt: number
  draws: number
  goalsFor: number
  goalsAgainst: number
  points: number
  groupName?: string
}

export interface Scorer {
  position: number
  playerId: string
  name: string
  photoUrl?: string
  team: string
  teamId: string
  logoUrl?: string
  games: number
  goals: number
  assists: number
  points: number
}

export interface MatchItem {
  id: string
  homeTeam: string
  awayTeam: string
  homeTeamId: string
  awayTeamId: string
  homeLogoUrl?: string
  awayLogoUrl?: string
  homeScore: number | null
  awayScore: number | null
  resultType?: string
  date: string
  time: string
  tournament: string
  venue?: string
  status: string
}

export interface PlayerStats {
  games: number
  goals: number
  assists: number
  points: number
  plusMinus: number
  penaltyMinutes: number
}

export interface PlayerItem {
  id: string
  name: string
  position: string
  birthDate: string
  birthYear: number
  team: string
  teamId: string
  jerseyNumber: number
  photoUrl?: string
  stats?: PlayerStats
}

export interface PlayerProfile {
  id: string
  name: string
  position: string
  birthDate: string
  birthYear: number
  team: string
  teamId: string
  jerseyNumber: number
  height?: number
  weight?: number
  handedness?: string
  city?: string
  photoUrl?: string
  stats?: PlayerStats
}

export interface TeamStats {
  wins: number
  losses: number
  draws: number
  goalsFor: number
  goalsAgainst: number
}

export interface TeamProfile {
  id: string
  name: string
  city: string
  logoUrl?: string
  tournaments: string[]
  playersCount: number
  roster: PlayerItem[]
  stats: TeamStats
  recentMatches: MatchItem[]
}

export interface PlayerStatEntry {
  season: string
  tournamentId: string
  tournamentName: string
  groupName: string
  birthYear: number
  games: number
  goals: number
  assists: number
  points: number
  plusMinus: number
  penaltyMinutes: number
}

export interface RankedPlayer {
  rank: number
  id: string
  name: string
  photoUrl?: string
  position: string
  birthYear: number
  team: string
  teamId: string
  games: number
  goals: number
  assists: number
  points: number
  plusMinus: number
  penaltyMinutes: number
}

export interface RankingsData {
  season: string
  players: RankedPlayer[]
}

export interface DomainOption {
  domain: string
  label: string
}

export interface TournamentOption {
  id: string
  name: string
  domain: string
  birthYears?: number[]
}

export interface GroupOption {
  name: string
  tournamentId: string
}

export interface RankingsFilters {
  birthYears: number[]
  domains: DomainOption[]
  tournaments: TournamentOption[]
  groups: GroupOption[]
}

// Match Detail Types
export interface MatchTeam {
  id: string
  name: string
  city?: string
  logoUrl?: string
}

export interface ScoreByPeriod {
  homeP1?: number
  awayP1?: number
  homeP2?: number
  awayP2?: number
  homeP3?: number
  awayP3?: number
  homeOt?: number
  awayOt?: number
}

export interface TournamentInfo {
  id: string
  name: string
}

export interface MatchEvent {
  type: string
  period?: number
  time?: string
  isHome: boolean
  teamName?: string
  teamLogoUrl?: string
  playerId?: string
  playerName?: string
  playerPhoto?: string
  assist1Id?: string
  assist1Name?: string
  assist2Id?: string
  assist2Name?: string
  goalType?: string
  penaltyMins?: number
  penaltyText?: string
}

export interface LineupPlayer {
  playerId: string
  playerName: string
  playerPhoto?: string
  jerseyNumber?: number
  position?: string
  goals: number
  assists: number
  points: number
  penaltyMinutes: number
  plusMinus: number
  saves?: number
  goalsAgainst?: number
}

export interface MatchDetail {
  id: string
  externalId: string
  homeTeam: MatchTeam
  awayTeam: MatchTeam
  homeScore: number | null
  awayScore: number | null
  scoreByPeriod?: ScoreByPeriod
  resultType?: string
  date: string
  time: string
  tournament: TournamentInfo
  venue?: string
  status: string
  groupName?: string
  birthYear?: number
  matchNumber?: number
  events: MatchEvent[]
  homeLineup: LineupPlayer[]
  awayLineup: LineupPlayer[]
}
