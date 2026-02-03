import Foundation

struct LoginRequest: Encodable {
    let email: String
    let password: String
}

struct RegisterRequest: Encodable {
    let email: String
    let password: String
    let name: String
}

struct RefreshRequest: Encodable {
    let refreshToken: String
}

struct LinkPlayerRequest: Encodable {
    let playerId: String
}

struct AuthResponse: Decodable {
    let accessToken: String
    let refreshToken: String
    let expiresIn: Int
    let user: UserDTO
}

struct UserDTO: Decodable {
    let id: String
    let email: String
    let name: String
    let role: String?
    let subscriptionTier: String
    let playerId: String?
    let emailVerified: Bool
    let createdAt: Date?
}

typealias UserResponse = UserDTO
