import Foundation
import Security

final class KeychainManager {
    static let shared = KeychainManager()

    private let accessTokenKey = "starrink.accessToken"
    private let refreshTokenKey = "starrink.refreshToken"

    private init() {}

    var accessToken: String? {
        get { getString(forKey: accessTokenKey) }
        set {
            if let value = newValue {
                setString(value, forKey: accessTokenKey)
            } else {
                delete(forKey: accessTokenKey)
            }
        }
    }

    var refreshToken: String? {
        get { getString(forKey: refreshTokenKey) }
        set {
            if let value = newValue {
                setString(value, forKey: refreshTokenKey)
            } else {
                delete(forKey: refreshTokenKey)
            }
        }
    }

    var hasTokens: Bool {
        accessToken != nil && refreshToken != nil
    }

    func saveTokens(access: String, refresh: String) {
        accessToken = access
        refreshToken = refresh
    }

    func clearTokens() {
        accessToken = nil
        refreshToken = nil
    }

    private func setString(_ value: String, forKey key: String) {
        guard let data = value.data(using: .utf8) else { return }

        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecValueData as String: data
        ]

        SecItemDelete(query as CFDictionary)
        SecItemAdd(query as CFDictionary, nil)
    }

    private func getString(forKey key: String) -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecReturnData as String: true
        ]

        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)

        guard status == errSecSuccess,
              let data = result as? Data,
              let string = String(data: data, encoding: .utf8) else {
            return nil
        }

        return string
    }

    private func delete(forKey key: String) {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key
        ]
        SecItemDelete(query as CFDictionary)
    }
}
