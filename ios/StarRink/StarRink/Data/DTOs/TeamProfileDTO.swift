import Foundation

struct TeamProfileDTO: Decodable {
    let id: String
    let name: String
    let city: String
    let logoUrl: String?
    let tournaments: [String]
    let playersCount: Int
    let roster: [PlayerItemDTO]
    let stats: TeamStatsDTO
    let recentMatches: [MatchDTO]
}

struct TeamStatsDTO: Decodable {
    let wins: Int
    let losses: Int
    let draws: Int
    let goalsFor: Int
    let goalsAgainst: Int
}
