import Foundation

// MARK: - Search

struct PlayersSearchResponse: Decodable {
    let players: [PlayerItemDTO]
    let total: Int
}

struct PlayerItemDTO: Decodable, Identifiable {
    let id: String
    let name: String
    let position: String
    let birthDate: String?
    let birthYear: Int
    let team: String
    let teamId: String
    let jerseyNumber: Int?
    let photoUrl: String?
    let stats: PlayerStatsDTO?

    enum CodingKeys: String, CodingKey {
        case id, name, position, birthDate, birthYear, team, teamId, jerseyNumber, stats
        case photoUrl
    }
}

struct PlayerStatsDTO: Decodable {
    let games: Int
    let goals: Int
    let assists: Int
    let points: Int
    let plusMinus: Int
    let penaltyMinutes: Int
}

// MARK: - Profile

struct PlayerProfileDTO: Decodable {
    let id: String
    let name: String
    let position: String
    let birthDate: String?
    let birthYear: Int
    let team: String
    let teamId: String
    let jerseyNumber: Int?
    let height: Int?
    let weight: Int?
    let handedness: String?
    let city: String?
    let photoUrl: String?
    let stats: PlayerStatsDTO?

    enum CodingKeys: String, CodingKey {
        case id, name, position, birthDate, birthYear, team, teamId
        case jerseyNumber, height, weight, handedness, city, stats
        case photoUrl
    }
}

// MARK: - Stats History

struct PlayerStatsHistoryResponse: Decodable {
    let stats: [PlayerStatEntryDTO]
}

struct PlayerStatEntryDTO: Decodable, Identifiable {
    var id: String { "\(tournamentId)_\(groupName)_\(birthYear)" }
    let season: String
    let tournamentId: String
    let tournamentName: String
    let groupName: String
    let birthYear: Int
    let games: Int
    let goals: Int
    let assists: Int
    let points: Int
    let plusMinus: Int
    let penaltyMinutes: Int
}

// MARK: - Seasons

struct SeasonsResponse: Decodable {
    let seasons: [String]
}
