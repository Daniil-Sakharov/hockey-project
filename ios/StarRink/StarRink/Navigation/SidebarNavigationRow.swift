import SwiftUI

struct SidebarNavigationRow: View {
    let item: NavigationItem
    let isSelected: Bool
    let userTier: SubscriptionTier
    let action: () -> Void

    @State private var isPressed = false

    private var isLocked: Bool {
        guard let requiredTier = item.tier else { return false }
        let tierOrder: [SubscriptionTier] = [.free, .pro, .ultra]
        guard let userIndex = tierOrder.firstIndex(of: userTier),
              let requiredIndex = tierOrder.firstIndex(of: requiredTier) else { return true }
        return userIndex < requiredIndex
    }

    var body: some View {
        Button {
            let impactFeedback = UIImpactFeedbackGenerator(style: .light)
            impactFeedback.impactOccurred()
            action()
        } label: {
            HStack(spacing: AppSpacing.md) {
                iconView
                titleView
                Spacer()
                if isLocked {
                    Image(systemName: "lock.fill")
                        .font(.system(size: 12))
                        .foregroundColor(.srTextMuted.opacity(0.5))
                }
            }
            .padding(.horizontal, AppSpacing.md)
            .padding(.vertical, AppSpacing.sm)
            .background(rowBackground)
            .scaleEffect(isPressed ? 0.98 : 1.0)
        }
        .buttonStyle(SidebarButtonStyle(isPressed: $isPressed))
    }

    private var iconView: some View {
        ZStack {
            if isSelected {
                RoundedRectangle(cornerRadius: 12)
                    .fill(Color.srCyan.opacity(0.25))
                    .frame(width: 42, height: 42)
            }
            Image(systemName: item.icon)
                .font(.system(size: 18, weight: isSelected ? .semibold : .regular))
                .foregroundColor(
                    isSelected ? .srCyan :
                    isLocked ? .srTextMuted.opacity(0.5) :
                    .srTextSecondary
                )
                .frame(width: 42, height: 42)
        }
    }

    private var titleView: some View {
        Text(item.title)
            .font(.system(size: 16, weight: isSelected ? .semibold : .regular))
            .foregroundColor(
                isSelected ? .srTextPrimary :
                isLocked ? .srTextMuted.opacity(0.5) :
                .srTextSecondary
            )
    }

    private var rowBackground: some View {
        RoundedRectangle(cornerRadius: 14)
            .fill(
                isSelected ? Color.srCyan.opacity(0.15) :
                isPressed ? Color.srCyan.opacity(0.08) :
                Color.clear
            )
            .overlay(
                RoundedRectangle(cornerRadius: 14)
                    .stroke(
                        isSelected ? Color.srCyan.opacity(0.3) : Color.clear,
                        lineWidth: 1
                    )
            )
    }
}

struct SidebarButtonStyle: ButtonStyle {
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
