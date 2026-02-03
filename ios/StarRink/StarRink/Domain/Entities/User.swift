import Foundation

struct User: Identifiable, Codable {
    let id: String
    let email: String
    let displayName: String
    let role: UserRole
    let subscriptionTier: SubscriptionTier
    let playerId: String?
    let emailVerified: Bool
    let createdAt: Date?

    init(from dto: UserDTO) {
        self.id = dto.id
        self.email = dto.email
        self.displayName = dto.name.isEmpty ? email.components(separatedBy: "@").first ?? "User" : dto.name
        self.role = UserRole(rawValue: dto.role ?? "fan") ?? .fan
        self.subscriptionTier = SubscriptionTier(rawValue: dto.subscriptionTier) ?? .free
        self.playerId = dto.playerId
        self.emailVerified = dto.emailVerified
        self.createdAt = dto.createdAt
    }

    #if DEBUG
    init(id: String, email: String, displayName: String, role: UserRole, subscriptionTier: SubscriptionTier) {
        self.id = id
        self.email = email
        self.displayName = displayName
        self.role = role
        self.subscriptionTier = subscriptionTier
        self.playerId = nil
        self.emailVerified = true
        self.createdAt = Date()
    }
    #endif
}

enum SubscriptionTier: String, Codable {
    case free
    case pro
    case ultra

    var displayName: String {
        switch self {
        case .free: return "Free"
        case .pro: return "PRO"
        case .ultra: return "ULTRA"
        }
    }
}
