import Foundation

enum NavigationItem: String, CaseIterable, Identifiable, Hashable {
    // Fan/Explore items
    case fanHome
    case tournaments
    case playerSearch
    case rankings
    case matchResults
    case matchCalendar
    case predictions
    case favorites
    // Player items
    case profile
    case stats
    case calendar
    case achievements
    // Premium items
    case compare
    case scoutViews
    case videoHighlights
    case aiRecommendations
    // Shared
    case subscription
    case settings

    var id: String { rawValue }

    var title: String {
        switch self {
        case .fanHome: return "Обзор"
        case .tournaments: return "Турниры"
        case .playerSearch: return "Игроки"
        case .rankings: return "Рейтинг"
        case .matchResults: return "Результаты"
        case .matchCalendar: return "Календарь"
        case .predictions: return "Прогнозы"
        case .favorites: return "Избранное"
        case .profile: return "Мой профиль"
        case .stats: return "Статистика"
        case .calendar: return "Календарь"
        case .achievements: return "Достижения"
        case .compare: return "Сравнение"
        case .scoutViews: return "Просмотры"
        case .videoHighlights: return "Видео"
        case .aiRecommendations: return "AI советы"
        case .subscription: return "Подписка"
        case .settings: return "Настройки"
        }
    }

    var icon: String {
        switch self {
        case .fanHome: return "house.fill"
        case .tournaments: return "trophy.fill"
        case .playerSearch: return "magnifyingglass"
        case .rankings: return "medal.fill"
        case .matchResults: return "clipboard.fill"
        case .matchCalendar: return "calendar"
        case .predictions: return "bolt.fill"
        case .favorites: return "star.fill"
        case .profile: return "person.fill"
        case .stats: return "chart.bar.fill"
        case .calendar: return "calendar"
        case .achievements: return "trophy.fill"
        case .compare: return "person.2.fill"
        case .scoutViews: return "eye.fill"
        case .videoHighlights: return "play.rectangle.fill"
        case .aiRecommendations: return "brain.head.profile"
        case .subscription: return "creditcard.fill"
        case .settings: return "gearshape.fill"
        }
    }

    var tier: SubscriptionTier? {
        switch self {
        case .compare, .scoutViews, .favorites: return .pro
        case .videoHighlights, .aiRecommendations: return .ultra
        default: return nil
        }
    }

    // MARK: - Fan/Explore groups

    static func mainItems(for role: UserRole) -> [NavigationItem] {
        switch role {
        case .fan:
            return [.fanHome, .tournaments, .playerSearch, .rankings]
        case .player:
            return [.profile, .stats, .calendar, .achievements]
        case .scout:
            return [.fanHome, .playerSearch, .rankings, .tournaments]
        case .coach:
            return [.fanHome, .calendar, .rankings, .tournaments]
        }
    }

    static func secondaryItems(for role: UserRole) -> [NavigationItem] {
        switch role {
        case .fan:
            return [.matchResults, .matchCalendar, .predictions]
        default:
            return []
        }
    }

    static func proItems(for role: UserRole) -> [NavigationItem] {
        switch role {
        case .fan: return [.favorites]
        case .player: return [.compare, .scoutViews]
        case .scout: return [.compare]
        default: return []
        }
    }

    static func ultraItems(for role: UserRole) -> [NavigationItem] {
        switch role {
        case .player: return [.videoHighlights, .aiRecommendations]
        default: return []
        }
    }

    static var bottomItems: [NavigationItem] {
        [.settings]
    }

    static func defaultItem(for tab: TabItem, role: UserRole) -> NavigationItem {
        switch tab {
        case .home:
            return role == .player ? .profile : .fanHome
        case .tournaments:
            return .tournaments
        case .search:
            return .playerSearch
        case .calendar:
            return .calendar
        case .profile:
            return .profile
        case .settings:
            return .settings
        }
    }
}
