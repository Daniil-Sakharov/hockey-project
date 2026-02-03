import Foundation

struct TournamentsResponse: Decodable {
    let tournaments: [TournamentItemDTO]
}

struct TournamentItemDTO: Decodable, Identifiable {
    let id: String
    let name: String
    let domain: String
    let season: String
    let source: String
    let birthYearGroups: [String: [GroupStatsDTO]]?
    let teamsCount: Int
    let matchesCount: Int
    let isEnded: Bool
}

struct GroupStatsDTO: Decodable {
    let name: String
    let teamsCount: Int
    let matchesCount: Int
}
