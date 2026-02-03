import SwiftUI

struct MenuButton: View {
    let isSidebarOpen: Bool
    let action: () -> Void
    @State private var isPressed = false

    var body: some View {
        Button(action: {
            let impactFeedback = UIImpactFeedbackGenerator(style: .light)
            impactFeedback.impactOccurred()
            action()
        }) {
            ZStack {
                Circle()
                    .fill(Color.srCyan.opacity(isPressed || isSidebarOpen ? 0.2 : 0))
                    .frame(width: 40, height: 40)

                Image(systemName: "line.3.horizontal")
                    .font(.system(size: 20, weight: .medium))
                    .foregroundColor(.srCyan)
                    .scaleEffect(isPressed ? 0.9 : 1.0)
            }
        }
        .buttonStyle(PressableButtonStyle(isPressed: $isPressed))
    }
}

struct PressableButtonStyle: ButtonStyle {
    @Binding var isPressed: Bool

    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .onChange(of: configuration.isPressed) { _, newValue in
                withAnimation(.easeInOut(duration: 0.1)) {
                    isPressed = newValue
                }
            }
    }
}
