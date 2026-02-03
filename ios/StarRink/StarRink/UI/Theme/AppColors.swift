import SwiftUI

extension Color {
    static let srBackground = Color(hex: "0D0D0F")
    static let srCard = Color(hex: "1A1A1F")
    static let srBorder = Color(hex: "2A2A30")

    static let srCyan = Color(hex: "00D4FF")
    static let srPurple = Color(hex: "8B5CF6")
    static let srAmber = Color(hex: "F59E0B")

    static let srSuccess = Color(hex: "10B981")
    static let srError = Color(hex: "EF4444")
    static let srWarning = Color(hex: "F59E0B")

    static let srTextPrimary = Color.white
    static let srTextSecondary = Color(hex: "9CA3AF")
    static let srTextMuted = Color(hex: "6B7280")

    static let srGradientPrimary = LinearGradient(
        colors: [Color(hex: "00D4FF"), Color(hex: "8B5CF6")],
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )

    static let srGradientAccent = LinearGradient(
        colors: [Color(hex: "8B5CF6"), Color(hex: "EC4899")],
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )

    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 3:
            (a, r, g, b) = (255, (int >> 8) * 17, (int >> 4 & 0xF) * 17, (int & 0xF) * 17)
        case 6:
            (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8:
            (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default:
            (a, r, g, b) = (255, 0, 0, 0)
        }
        self.init(
            .sRGB,
            red: Double(r) / 255,
            green: Double(g) / 255,
            blue: Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}
