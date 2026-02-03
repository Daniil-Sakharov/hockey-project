import Foundation

struct ExploreOverviewDTO: Decodable {
    let players: Int
    let teams: Int
    let tournaments: Int
    let matches: Int
}
