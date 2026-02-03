import Foundation
import SwiftUI

@MainActor
final class AuthViewModel: ObservableObject {
    @Published var isAuthenticated = false
    @Published var currentUser: User?
    @Published var isLoading = false
    @Published var errorMessage: String?

    @Published var email = ""
    @Published var password = ""
    @Published var name = ""
    @Published var confirmPassword = ""

    private let repository: AuthRepositoryProtocol
    private let keychain: KeychainManager

    init(repository: AuthRepositoryProtocol = AuthRepository(), keychain: KeychainManager = .shared) {
        self.repository = repository
        self.keychain = keychain

        if keychain.hasTokens {
            Task { await checkAuthStatus() }
        }
    }

    func login() async {
        guard validateLoginForm() else { return }

        isLoading = true
        errorMessage = nil

        do {
            let user = try await repository.login(email: email, password: password)
            currentUser = user
            isAuthenticated = true
            clearForm()
        } catch let error as APIError {
            errorMessage = error.errorDescription
        } catch {
            errorMessage = "Не удалось войти. Попробуйте ещё раз."
        }

        isLoading = false
    }

    func register() async {
        guard validateRegisterForm() else { return }

        isLoading = true
        errorMessage = nil

        do {
            let user = try await repository.register(email: email, password: password, name: name)
            currentUser = user
            isAuthenticated = true
            clearForm()
        } catch let error as APIError {
            errorMessage = error.errorDescription
        } catch {
            errorMessage = "Не удалось зарегистрироваться. Попробуйте ещё раз."
        }

        isLoading = false
    }

    func logout() async {
        if keychain.hasTokens {
            do {
                try await repository.logout()
            } catch {}
        }
        keychain.clearTokens()
        currentUser = nil
        isAuthenticated = false
    }

    #if DEBUG
    func mockLogin(role: UserRole = .fan) {
        let name: String
        switch role {
        case .fan: name = "Тестовый Пользователь"
        case .player: name = "Тестовый Игрок"
        case .scout: name = "Тестовый Скаут"
        case .coach: name = "Тестовый Тренер"
        }

        currentUser = User(
            id: "test-123",
            email: "test@starrink.ru",
            displayName: name,
            role: role,
            subscriptionTier: .free
        )
        isAuthenticated = true
    }
    #endif

    private func checkAuthStatus() async {
        do {
            let user = try await repository.getCurrentUser()
            currentUser = user
            isAuthenticated = true
        } catch {
            keychain.clearTokens()
            isAuthenticated = false
        }
    }

    private func validateLoginForm() -> Bool {
        if email.isEmpty {
            errorMessage = "Введите email"
            return false
        }
        if password.isEmpty {
            errorMessage = "Введите пароль"
            return false
        }
        return true
    }

    private func validateRegisterForm() -> Bool {
        if email.isEmpty {
            errorMessage = "Введите email"
            return false
        }
        if password.isEmpty {
            errorMessage = "Введите пароль"
            return false
        }
        if password.count < 6 {
            errorMessage = "Пароль должен быть не менее 6 символов"
            return false
        }
        if password != confirmPassword {
            errorMessage = "Пароли не совпадают"
            return false
        }
        return true
    }

    private func clearForm() {
        email = ""
        password = ""
        name = ""
        confirmPassword = ""
        errorMessage = nil
    }
}
