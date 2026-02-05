import Foundation

struct TeamsResponse: Decodable {
    let teams: [TeamItemDTO]
}

struct TeamItemDTO: Decodable, Identifiable {
    let id: String
    let name: String
    let city: String?
    let logoUrl: String?
    let playersCount: Int
    let groupName: String?
    let birthYear: Int?
}
