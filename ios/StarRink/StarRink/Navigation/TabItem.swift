import Foundation

enum TabItem: String, CaseIterable, Hashable {
    // Fan tabs
    case home
    case tournaments
    case search
    // Player tabs
    case calendar
    case profile
    // Shared
    case settings

    var title: String {
        switch self {
        case .home: return "Главная"
        case .tournaments: return "Турниры"
        case .search: return "Поиск"
        case .calendar: return "Календарь"
        case .profile: return "Профиль"
        case .settings: return "Настройки"
        }
    }

    var icon: String {
        switch self {
        case .home: return "house"
        case .tournaments: return "trophy"
        case .search: return "magnifyingglass"
        case .calendar: return "calendar"
        case .profile: return "person"
        case .settings: return "gearshape"
        }
    }

    var selectedIcon: String {
        switch self {
        case .home: return "house.fill"
        case .tournaments: return "trophy.fill"
        case .search: return "magnifyingglass"
        case .calendar: return "calendar"
        case .profile: return "person.fill"
        case .settings: return "gearshape.fill"
        }
    }

    static func tabs(for role: UserRole) -> [TabItem] {
        switch role {
        case .fan:
            return [.home, .tournaments, .search, .settings]
        case .player:
            return [.home, .calendar, .profile, .settings]
        case .scout:
            return [.home, .search, .settings]
        case .coach:
            return [.home, .calendar, .settings]
        }
    }

    static func defaultTab(for role: UserRole) -> TabItem {
        .home
    }
}
