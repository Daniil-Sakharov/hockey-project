import SwiftUI

struct MainContainerView: View {
    @EnvironmentObject var authViewModel: AuthViewModel
    @State private var selectedTab: TabItem = .home
    @State private var selectedSidebarItem: NavigationItem = .fanHome
    @State private var isSidebarOpen = false

    private var userRole: UserRole {
        authViewModel.currentUser?.role ?? .fan
    }

    private var availableTabs: [TabItem] {
        TabItem.tabs(for: userRole)
    }

    private var sidebarWidth: CGFloat {
        UIScreen.main.bounds.width * 0.72
    }

    var body: some View {
        GeometryReader { geometry in
            ZStack {
                SidebarView(
                    selectedTab: $selectedTab,
                    selectedItem: $selectedSidebarItem,
                    isOpen: $isSidebarOpen,
                    onItemSelected: { item in
                        handleSidebarSelection(item)
                    }
                )
                .frame(width: geometry.size.width, height: geometry.size.height)
                .zIndex(0)

                mainContent
                    .frame(width: geometry.size.width, height: geometry.size.height)
                    .clipShape(RoundedRectangle(cornerRadius: isSidebarOpen ? 24 : 0))
                    .scaleEffect(isSidebarOpen ? 0.85 : 1.0, anchor: .trailing)
                    .offset(x: isSidebarOpen ? sidebarWidth : 0)
                    .shadow(color: .black.opacity(isSidebarOpen ? 0.4 : 0), radius: 20, x: -10, y: 0)
                    .zIndex(1)
                    .allowsHitTesting(!isSidebarOpen)

                if isSidebarOpen {
                    Color.black.opacity(0.001)
                        .contentShape(Rectangle())
                        .offset(x: sidebarWidth)
                        .onTapGesture { closeSidebar() }
                        .zIndex(2)
                }
            }
        }
        .background(Color.clear)
        .ignoresSafeArea()
        .animation(.spring(response: 0.4, dampingFraction: 0.82), value: isSidebarOpen)
        .gesture(sidebarDragGesture)
        .onChange(of: userRole) { _, newRole in
            selectedTab = TabItem.defaultTab(for: newRole)
        }
        .onChange(of: selectedTab) { _, newTab in
            selectedSidebarItem = NavigationItem.defaultItem(for: newTab, role: userRole)
        }
    }

    private var mainContent: some View {
        NavigationStack {
            ZStack {
                AnimatedBackgroundView()
                    .ignoresSafeArea()

                tabContentView
            }
            .background(Color.clear)
            .scrollContentBackground(.hidden)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    MenuButton(isSidebarOpen: isSidebarOpen) {
                        toggleSidebar()
                    }
                }
                ToolbarItem(placement: .principal) {
                    HStack(spacing: AppSpacing.xs) {
                        Image(systemName: "star.fill")
                            .foregroundColor(.srCyan)
                        Text("StarRink")
                            .font(.srHeading4)
                            .foregroundColor(.srTextPrimary)
                    }
                }
            }
            .toolbarBackground(.hidden, for: .navigationBar)
            .navigationDestination(for: MatchRoute.self) { route in
                MatchDetailView(matchId: route.matchId)
            }
            .navigationDestination(for: TeamRoute.self) { route in
                TeamProfileView(teamId: route.teamId, teamName: route.teamName)
            }
            .navigationDestination(for: TeamRosterRoute.self) { route in
                TeamRosterView(teamId: route.teamId, tournamentId: route.tournamentId, teamName: route.teamName)
            }
            .navigationDestination(for: PlayerRoute.self) { route in
                PlayerProfileView(playerId: route.playerId)
            }
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
        .safeAreaInset(edge: .bottom) {
            RoleTabBar(
                tabs: availableTabs,
                selectedTab: $selectedTab
            )
        }
    }

    @ViewBuilder
    private var tabContentView: some View {
        switch selectedTab {
        // Fan tabs
        case .home:
            if userRole == .fan {
                fanHomeContent
            } else {
                HomeView()
            }
        case .tournaments:
            FanTournamentsView()
        case .search:
            FanSearchView()
        // Player tabs
        case .calendar:
            CalendarView()
        case .profile:
            ProfileView()
        // Shared
        case .settings:
            SettingsView()
        }
    }

    @ViewBuilder
    private var fanHomeContent: some View {
        switch selectedSidebarItem {
        case .rankings:
            RankingsView()
        case .matchResults:
            MatchResultsView()
        case .matchCalendar:
            CalendarView()
        default:
            FanHomeView()
        }
    }

    // MARK: - Sidebar

    private var sidebarDragGesture: some Gesture {
        DragGesture()
            .onEnded { gesture in
                let threshold: CGFloat = 50
                if gesture.translation.width > threshold && !isSidebarOpen {
                    openSidebar()
                } else if gesture.translation.width < -threshold && isSidebarOpen {
                    closeSidebar()
                }
            }
    }

    private func toggleSidebar() {
        withAnimation(.spring(response: 0.4, dampingFraction: 0.82)) {
            isSidebarOpen.toggle()
        }
    }

    private func openSidebar() {
        withAnimation(.spring(response: 0.4, dampingFraction: 0.82)) {
            isSidebarOpen = true
        }
    }

    private func closeSidebar() {
        withAnimation(.spring(response: 0.4, dampingFraction: 0.82)) {
            isSidebarOpen = false
        }
    }

    private func handleSidebarSelection(_ item: NavigationItem) {
        selectedSidebarItem = item
        switch item {
        case .fanHome: selectedTab = .home
        case .tournaments: selectedTab = .tournaments
        case .playerSearch: selectedTab = .search
        case .profile: selectedTab = .profile
        case .calendar: selectedTab = .calendar
        case .settings: selectedTab = .settings
        case .rankings, .matchResults, .matchCalendar:
            selectedTab = .home
        default: break
        }
        closeSidebar()
    }
}
