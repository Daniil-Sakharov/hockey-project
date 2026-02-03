import SwiftUI

struct RoleTabBar: View {
    let tabs: [TabItem]
    @Binding var selectedTab: TabItem

    var body: some View {
        HStack(spacing: 0) {
            ForEach(tabs, id: \.self) { tab in
                tabButton(for: tab)
            }
        }
        .padding(.horizontal, AppSpacing.sm)
        .padding(.vertical, AppSpacing.xs)
        .background(tabBarBackground)
        .padding(.horizontal, AppSpacing.lg)
        .padding(.bottom, 16)
    }

    private func tabButton(for tab: TabItem) -> some View {
        Button {
            withAnimation(.spring(response: 0.3, dampingFraction: 0.7)) {
                selectedTab = tab
            }
        } label: {
            VStack(spacing: 3) {
                ZStack {
                    if selectedTab == tab {
                        Circle()
                            .fill(
                                RadialGradient(
                                    colors: [Color.srCyan.opacity(0.4), Color.clear],
                                    center: .center,
                                    startRadius: 0,
                                    endRadius: 25
                                )
                            )
                            .frame(width: 50, height: 50)
                            .blur(radius: 8)
                    }

                    Image(systemName: selectedTab == tab ? tab.selectedIcon : tab.icon)
                        .font(.system(size: 22, weight: selectedTab == tab ? .semibold : .regular))
                        .foregroundStyle(
                            selectedTab == tab ?
                            LinearGradient(colors: [.srCyan, .srPurple], startPoint: .topLeading, endPoint: .bottomTrailing) :
                            LinearGradient(colors: [.srTextMuted, .srTextMuted], startPoint: .top, endPoint: .bottom)
                        )
                        .scaleEffect(selectedTab == tab ? 1.1 : 1)
                }

                Text(tab.title)
                    .font(.system(size: 9, weight: selectedTab == tab ? .semibold : .regular))
                    .foregroundColor(selectedTab == tab ? .srCyan : .srTextMuted.opacity(0.7))
            }
            .frame(maxWidth: .infinity)
            .padding(.vertical, AppSpacing.xs)
        }
        .buttonStyle(.plain)
    }

    private var tabBarBackground: some View {
        Capsule()
            .fill(Color.srBackground.opacity(0.6))
            .background(
                Capsule().fill(.ultraThinMaterial.opacity(0.3))
            )
            .overlay(
                Capsule()
                    .stroke(
                        LinearGradient(
                            colors: [Color.srCyan.opacity(0.3), Color.srPurple.opacity(0.2)],
                            startPoint: .leading,
                            endPoint: .trailing
                        ),
                        lineWidth: 0.5
                    )
            )
    }
}
