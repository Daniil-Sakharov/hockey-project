import Foundation

// MARK: - Rankings Response

struct RankingsResponse: Decodable {
    let season: String
    let players: [RankedPlayerDTO]
}

struct RankedPlayerDTO: Decodable, Identifiable {
    let rank: Int
    let id: String
    let name: String
    let photoUrl: String?
    let position: String
    let birthYear: Int
    let team: String
    let teamId: String
    let games: Int
    let goals: Int
    let assists: Int
    let points: Int
    let plusMinus: Int
    let penaltyMinutes: Int
}
