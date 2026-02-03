import SwiftUI

struct SidebarView: View {
    @Binding var selectedTab: TabItem
    @Binding var selectedItem: NavigationItem
    @Binding var isOpen: Bool
    var onItemSelected: ((NavigationItem) -> Void)?
    @EnvironmentObject var authViewModel: AuthViewModel

    private var userRole: UserRole {
        authViewModel.currentUser?.role ?? .fan
    }

    private var userTier: SubscriptionTier {
        authViewModel.currentUser?.subscriptionTier ?? .free
    }

    private var sidebarContentWidth: CGFloat {
        UIScreen.main.bounds.width * 0.72
    }

    var body: some View {
        ZStack(alignment: .leading) {
            AnimatedBackgroundView()
                .ignoresSafeArea()

            Color.srBackground.opacity(0.4)
                .ignoresSafeArea()

            sidebarContent
                .frame(width: sidebarContentWidth)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .leading)
    }

    private var sidebarContent: some View {
        VStack(alignment: .leading, spacing: 0) {
            sidebarHeader
                .padding(.horizontal, AppSpacing.lg)
                .padding(.top, 70)
                .padding(.bottom, AppSpacing.lg)

            userInfoSection
                .padding(.horizontal, AppSpacing.lg)
                .padding(.bottom, AppSpacing.lg)

            sidebarDivider

            ScrollView(showsIndicators: false) {
                VStack(alignment: .leading, spacing: 4) {
                    navigationGroups
                }
                .padding(.vertical, AppSpacing.md)
                .padding(.horizontal, AppSpacing.sm)
            }

            Spacer()

            subscriptionBadge
                .padding(.horizontal, AppSpacing.lg)
                .padding(.bottom, AppSpacing.xl)
        }
    }

    // MARK: - Navigation Groups

    @ViewBuilder
    private var navigationGroups: some View {
        // Main items
        ForEach(NavigationItem.mainItems(for: userRole)) { item in
            sidebarRow(item)
        }

        // Secondary items (with divider)
        let secondary = NavigationItem.secondaryItems(for: userRole)
        if !secondary.isEmpty {
            sidebarDivider
            ForEach(secondary) { item in
                sidebarRow(item)
            }
        }

        // PRO items (with divider + label)
        let proItems = NavigationItem.proItems(for: userRole)
        if !proItems.isEmpty {
            tierSectionHeader("PRO", color: .srPurple)
            ForEach(proItems) { item in
                sidebarRow(item)
            }
        }

        // ULTRA items (with divider + label)
        let ultraItems = NavigationItem.ultraItems(for: userRole)
        if !ultraItems.isEmpty {
            tierSectionHeader("ULTRA", color: .srAmber)
            ForEach(ultraItems) { item in
                sidebarRow(item)
            }
        }

        // Bottom items (with divider)
        sidebarDivider
        ForEach(NavigationItem.bottomItems) { item in
            sidebarRow(item)
        }
    }

    private func sidebarRow(_ item: NavigationItem) -> some View {
        SidebarNavigationRow(
            item: item,
            isSelected: selectedItem == item,
            userTier: userTier
        ) {
            onItemSelected?(item)
        }
    }

    // MARK: - Dividers & Section Headers

    private var sidebarDivider: some View {
        Rectangle()
            .fill(Color.srBorder.opacity(0.3))
            .frame(height: 1)
            .padding(.horizontal, AppSpacing.lg)
            .padding(.vertical, AppSpacing.sm)
    }

    private func tierSectionHeader(_ title: String, color: Color) -> some View {
        HStack(spacing: AppSpacing.sm) {
            Text(title)
                .font(.system(size: 11, weight: .bold))
                .foregroundColor(color.opacity(0.8))
                .tracking(1.5)

            Rectangle()
                .fill(
                    LinearGradient(
                        colors: [color.opacity(0.3), Color.clear],
                        startPoint: .leading,
                        endPoint: .trailing
                    )
                )
                .frame(height: 1)
        }
        .padding(.horizontal, AppSpacing.md)
        .padding(.top, AppSpacing.lg)
        .padding(.bottom, AppSpacing.xs)
    }

    // MARK: - Header & User Info

    private var sidebarHeader: some View {
        HStack(spacing: AppSpacing.sm) {
            Image(systemName: "star.fill")
                .font(.title)
                .foregroundStyle(
                    LinearGradient(
                        colors: [.srCyan, .srPurple],
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    )
                )
            Text("StarRink")
                .font(.system(size: 26, weight: .bold))
                .foregroundColor(.srTextPrimary)
            Spacer()
        }
    }

    private var userInfoSection: some View {
        HStack(spacing: AppSpacing.md) {
            ZStack {
                Circle()
                    .fill(
                        LinearGradient(
                            colors: [Color.srCyan, Color.srPurple],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(width: 48, height: 48)
                Text(authViewModel.currentUser?.displayName.prefix(1).uppercased() ?? "?")
                    .font(.system(size: 20, weight: .bold))
                    .foregroundColor(.white)
            }

            VStack(alignment: .leading, spacing: 3) {
                Text(authViewModel.currentUser?.displayName ?? "Пользователь")
                    .font(.system(size: 15, weight: .semibold))
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)

                Text(userRole.displayName)
                    .font(.system(size: 12))
                    .foregroundColor(.srTextSecondary)
            }
            Spacer()
        }
        .padding(AppSpacing.md)
        .background(
            RoundedRectangle(cornerRadius: 16)
                .fill(.ultraThinMaterial.opacity(0.5))
                .overlay(
                    RoundedRectangle(cornerRadius: 16)
                        .stroke(Color.srBorder.opacity(0.3), lineWidth: 0.5)
                )
        )
    }

    // MARK: - Subscription Badge

    private var subscriptionBadge: some View {
        HStack(spacing: AppSpacing.sm) {
            Image(systemName: badgeIcon)
                .font(.system(size: 14))
                .foregroundColor(badgeColor)

            Text(badgeText)
                .font(.system(size: 12, weight: .semibold))
                .foregroundColor(badgeColor)

            Spacer()
        }
        .padding(.horizontal, AppSpacing.md)
        .padding(.vertical, AppSpacing.sm)
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(badgeColor.opacity(0.1))
                .overlay(
                    RoundedRectangle(cornerRadius: 12)
                        .stroke(badgeColor.opacity(0.2), lineWidth: 0.5)
                )
        )
    }

    private var badgeIcon: String {
        switch userTier {
        case .ultra: return "crown.fill"
        case .pro: return "sparkles"
        case .free: return "star"
        }
    }

    private var badgeText: String {
        switch userTier {
        case .ultra: return "ULTRA"
        case .pro: return "PRO"
        case .free: return "Бесплатный план"
        }
    }

    private var badgeColor: Color {
        switch userTier {
        case .ultra: return .srAmber
        case .pro: return .srPurple
        case .free: return .srTextMuted
        }
    }
}
