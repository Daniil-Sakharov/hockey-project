import Foundation

struct TeamRosterResponse: Decodable {
    let team: TeamInfoDTO
    let players: [RosterPlayerDTO]
}

struct TeamInfoDTO: Decodable {
    let id: String
    let name: String
    let city: String?
    let logoUrl: String?
}

struct RosterPlayerDTO: Decodable, Identifiable {
    let id: String
    let name: String
    let photoUrl: String?
    let birthDate: String?
    let position: String?
    let jerseyNumber: Int
    let height: Int?
    let weight: Int?
    let birthYear: Int?
    let groupName: String?
    let handedness: String?
}
