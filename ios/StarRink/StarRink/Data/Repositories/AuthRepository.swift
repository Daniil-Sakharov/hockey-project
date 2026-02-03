import Foundation

protocol AuthRepositoryProtocol {
    func register(email: String, password: String, name: String) async throws -> User
    func login(email: String, password: String) async throws -> User
    func refreshTokens() async throws -> User
    func getCurrentUser() async throws -> User
    func linkPlayer(playerId: String) async throws -> User
    func logout() async throws
}

final class AuthRepository: AuthRepositoryProtocol {
    private let apiClient: APIClient
    private let keychain: KeychainManager

    init(
        apiClient: APIClient = .shared,
        keychain: KeychainManager = .shared
    ) {
        self.apiClient = apiClient
        self.keychain = keychain
    }

    func register(email: String, password: String, name: String) async throws -> User {
        let request = RegisterRequest(email: email, password: password, name: name)
        let response: AuthResponse = try await apiClient.request(
            endpoint: .register,
            body: request
        )
        saveTokens(from: response)
        return User(from: response.user)
    }

    func login(email: String, password: String) async throws -> User {
        let request = LoginRequest(email: email, password: password)
        let response: AuthResponse = try await apiClient.request(
            endpoint: .login,
            body: request
        )
        saveTokens(from: response)
        return User(from: response.user)
    }

    func refreshTokens() async throws -> User {
        guard let refreshToken = keychain.refreshToken else {
            throw APIError.unauthorized
        }
        let request = RefreshRequest(refreshToken: refreshToken)
        let response: AuthResponse = try await apiClient.request(
            endpoint: .refresh,
            body: request
        )
        saveTokens(from: response)
        return User(from: response.user)
    }

    func getCurrentUser() async throws -> User {
        let response: UserResponse = try await apiClient.request(endpoint: .me)
        return User(from: response)
    }

    func linkPlayer(playerId: String) async throws -> User {
        let request = LinkPlayerRequest(playerId: playerId)
        let response: UserResponse = try await apiClient.request(
            endpoint: .linkPlayer,
            body: request
        )
        return User(from: response)
    }

    func logout() async throws {
        try await apiClient.requestVoid(endpoint: .logout)
        keychain.clearTokens()
    }

    private func saveTokens(from response: AuthResponse) {
        keychain.saveTokens(
            access: response.accessToken,
            refresh: response.refreshToken
        )
    }
}
