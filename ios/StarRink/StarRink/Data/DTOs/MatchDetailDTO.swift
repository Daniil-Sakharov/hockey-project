import Foundation

// MARK: - Match Detail

struct MatchDetailResponse: Decodable {
    let id: String
    let externalId: String?
    let homeTeam: MatchTeamDTO
    let awayTeam: MatchTeamDTO
    let homeScore: Int?
    let awayScore: Int?
    let scoreByPeriod: ScoreByPeriodDTO?
    let resultType: String?
    let date: String
    let time: String
    let tournament: TournamentInfoDTO
    let venue: String?
    let status: String
    let groupName: String?
    let birthYear: Int?
    let matchNumber: Int?
    let events: [MatchEventDTO]
    let homeLineup: [LineupPlayerDTO]
    let awayLineup: [LineupPlayerDTO]
}

struct MatchTeamDTO: Decodable {
    let id: String
    let name: String
    let city: String?
    let logoUrl: String?
}

struct ScoreByPeriodDTO: Decodable {
    let homeP1: Int?
    let awayP1: Int?
    let homeP2: Int?
    let awayP2: Int?
    let homeP3: Int?
    let awayP3: Int?
    let homeOt: Int?
    let awayOt: Int?
}

struct TournamentInfoDTO: Decodable {
    let id: String
    let name: String
}

struct MatchEventDTO: Decodable, Identifiable {
    var id: String { "\(type)-\(time)-\(playerName ?? "")" }
    let type: String
    let period: Int?
    let time: String
    let isHome: Bool
    let teamName: String?
    let teamLogoUrl: String?
    let playerId: String?
    let playerName: String?
    let playerPhoto: String?
    let assist1Id: String?
    let assist1Name: String?
    let assist2Id: String?
    let assist2Name: String?
    let goalType: String?
    let penaltyMins: Int?
    let penaltyText: String?
}

struct LineupPlayerDTO: Decodable, Identifiable {
    var id: String { playerId }
    let playerId: String
    let playerName: String
    let playerPhoto: String?
    let jerseyNumber: Int?
    let position: String?
    let goals: Int
    let assists: Int
    let points: Int
    let penaltyMinutes: Int
    let plusMinus: Int
    let saves: Int?
    let goalsAgainst: Int?
}
