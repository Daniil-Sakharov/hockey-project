import SwiftUI

struct TournamentDetailView: View {
    let route: TournamentRoute
    @StateObject private var viewModel = TournamentDetailViewModel()

    var body: some View {
        ZStack {
            AnimatedBackgroundView()
                .ignoresSafeArea()

            content
        }
        .navigationTitle(route.name)
        .navigationBarTitleDisplayMode(.inline)
        .toolbarBackground(.hidden, for: .navigationBar)
        .task {
            await viewModel.loadAll(
                tournamentId: route.tournamentId,
                birthYear: route.birthYear,
                group: route.groupName
            )
        }
    }

    @ViewBuilder
    private var content: some View {
        if viewModel.isLoading {
            loadingView
        } else if let error = viewModel.errorMessage {
            errorView(error)
        } else {
            detailContent
        }
    }

    private var detailContent: some View {
        VStack(spacing: 0) {
            tabPicker
            tabContent
        }
    }

    private var tabPicker: some View {
        Picker("", selection: $viewModel.selectedTab) {
            ForEach(TournamentTab.allCases, id: \.self) { tab in
                Text(tab.rawValue).tag(tab)
            }
        }
        .pickerStyle(.segmented)
        .padding(.horizontal, AppSpacing.screenHorizontal)
        .padding(.vertical, AppSpacing.sm)
    }

    @ViewBuilder
    private var tabContent: some View {
        switch viewModel.selectedTab {
        case .standings:
            StandingsTabView(standings: viewModel.standings)
        case .matches:
            MatchesTabView(matches: viewModel.matches)
        case .scorers:
            ScorersTabView(scorers: viewModel.scorers)
        }
    }

    private var loadingView: some View {
        VStack(spacing: AppSpacing.md) {
            ProgressView().tint(.srCyan)
            Text("Загрузка...")
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)
        }
    }

    private func errorView(_ message: String) -> some View {
        VStack(spacing: AppSpacing.md) {
            Image(systemName: "wifi.slash")
                .font(.system(size: 36))
                .foregroundColor(.srTextMuted)
            Text(message)
                .font(.srBody)
                .foregroundColor(.srTextSecondary)
            Button("Повторить") {
                Task {
                    await viewModel.loadAll(
                        tournamentId: route.tournamentId,
                        birthYear: route.birthYear,
                        group: route.groupName
                    )
                }
            }
            .font(.srBodyMedium)
            .foregroundColor(.srCyan)
        }
        .padding(AppSpacing.xl)
    }
}
