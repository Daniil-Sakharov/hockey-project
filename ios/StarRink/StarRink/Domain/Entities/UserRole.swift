import Foundation

enum UserRole: String, Codable, CaseIterable {
    case fan
    case player
    case scout
    case coach

    var displayName: String {
        switch self {
        case .fan: return "Пользователь"
        case .player: return "Игрок"
        case .scout: return "Скаут"
        case .coach: return "Тренер"
        }
    }

    var subtitle: String {
        switch self {
        case .fan: return "Слежу за хоккеем"
        case .player: return "Играю в хоккей"
        case .scout: return "Ищу таланты"
        case .coach: return "Тренирую команду"
        }
    }

    var icon: String {
        switch self {
        case .fan: return "person.fill"
        case .player: return "figure.hockey"
        case .scout: return "binoculars.fill"
        case .coach: return "clipboard.fill"
        }
    }
}
