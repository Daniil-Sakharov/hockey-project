import Foundation

// MARK: - Standings

struct StandingsResponse: Decodable {
    let standings: [StandingDTO]
}

struct StandingDTO: Decodable, Identifiable {
    var id: String { teamId }
    let position: Int
    let team: String
    let teamId: String
    let games: Int
    let wins: Int
    let winsOt: Int
    let losses: Int
    let lossesOt: Int
    let draws: Int
    let goalsFor: Int
    let goalsAgainst: Int
    let points: Int
    let groupName: String?
}

// MARK: - Matches

struct MatchListResponse: Decodable {
    let matches: [MatchDTO]
}

struct MatchDTO: Decodable, Identifiable {
    let id: String
    let homeTeam: String
    let awayTeam: String
    let homeTeamId: String
    let awayTeamId: String
    let homeScore: Int?
    let awayScore: Int?
    let date: String
    let time: String
    let tournament: String?
    let venue: String?
    let status: String
}

// MARK: - Scorers

struct ScorersResponse: Decodable {
    let scorers: [ScorerDTO]
}

struct ScorerDTO: Decodable, Identifiable {
    var id: String { playerId }
    let position: Int
    let playerId: String
    let name: String
    let team: String
    let teamId: String
    let games: Int
    let goals: Int
    let assists: Int
    let points: Int
}
