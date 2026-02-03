import SwiftUI

struct GradientText: View {
    let text: String
    let gradient: LinearGradient
    let font: Font

    init(_ text: String, gradient: LinearGradient = .srGradientPrimary, font: Font = .srHeading2) {
        self.text = text
        self.gradient = gradient
        self.font = font
    }

    var body: some View {
        Text(text)
            .font(font)
            .foregroundStyle(gradient)
    }
}

extension LinearGradient {
    static let srGradientPrimary = LinearGradient(
        colors: [Color.srCyan, Color.srPurple],
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
}
