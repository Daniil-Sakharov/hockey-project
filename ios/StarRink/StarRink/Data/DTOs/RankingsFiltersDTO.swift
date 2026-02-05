import Foundation

struct RankingsFiltersResponse: Decodable {
    let birthYears: [Int]
    let domains: [DomainOption]
    let tournaments: [TournamentOption]
    let groups: [GroupOption]
}

struct DomainOption: Decodable, Identifiable {
    var id: String { domain }
    let domain: String
    let label: String
}

struct TournamentOption: Decodable, Identifiable {
    let id: String
    let name: String
    let domain: String
    let birthYears: [Int]?
}

struct GroupOption: Decodable, Identifiable {
    var id: String { "\(tournamentId)-\(name)" }
    let name: String
    let tournamentId: String
}
