import SwiftUI

struct CalendarView: View {
    @StateObject private var viewModel = CalendarViewModel()
    @State private var selectedTournament = "Все"

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.md) {
                headerSection
                tournamentFilter
                matchesList
            }
            .padding(.top, AppSpacing.md)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
        .overlay {
            if viewModel.isLoading && viewModel.upcomingMatches.isEmpty {
                ProgressView().tint(.srCyan)
            }
        }
        .task { await viewModel.load() }
        .onChange(of: selectedTournament) { _, newValue in
            viewModel.selectedTournament = newValue == "Все" ? nil : newValue
            Task { await viewModel.load() }
        }
    }

    private var headerSection: some View {
        HStack {
            Text("Календарь матчей")
                .font(.srHeading3)
                .foregroundColor(.srTextPrimary)
            Spacer()
            Image(systemName: "calendar")
                .font(.system(size: 24))
                .foregroundColor(.srCyan)
        }
        .padding(.horizontal, AppSpacing.screenHorizontal)
    }

    private var tournamentFilter: some View {
        Group {
            if !viewModel.tournamentNames.isEmpty {
                HorizontalChipPicker(
                    items: ["Все"] + viewModel.tournamentNames,
                    selectedItem: $selectedTournament
                )
            }
        }
    }

    private var matchesList: some View {
        LazyVStack(spacing: AppSpacing.lg) {
            if viewModel.groupedByDate.isEmpty && !viewModel.isLoading {
                Text("Нет предстоящих матчей")
                    .font(.srBody)
                    .foregroundColor(.srTextMuted)
                    .padding(.top, AppSpacing.xl)
            }

            ForEach(viewModel.groupedByDate, id: \.date) { group in
                VStack(alignment: .leading, spacing: AppSpacing.sm) {
                    Text(formatSectionDate(group.date))
                        .font(.srBodyMedium)
                        .foregroundColor(.srTextSecondary)
                        .padding(.horizontal, AppSpacing.screenHorizontal)

                    VStack(spacing: 0) {
                        ForEach(group.matches) { match in
                            NavigationLink(value: MatchRoute(matchId: match.id)) {
                                CalendarMatchRow(match: match)
                            }
                            Divider().background(Color.srBorder.opacity(0.3))
                        }
                    }
                    .glassCard(padding: 0)
                    .padding(.horizontal, AppSpacing.screenHorizontal)
                }
            }
        }
    }

    private func formatSectionDate(_ dateStr: String) -> String {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd"
        guard let date = formatter.date(from: dateStr) else { return dateStr }
        let display = DateFormatter()
        display.locale = Locale(identifier: "ru_RU")
        display.dateFormat = "d MMMM, EEEE"
        return display.string(from: date).capitalized
    }
}

// MARK: - Row

private struct CalendarMatchRow: View {
    let match: MatchDTO

    var body: some View {
        HStack(spacing: AppSpacing.md) {
            Text(match.time)
                .font(.system(size: 14, weight: .semibold, design: .monospaced))
                .foregroundColor(.srCyan)
                .frame(width: 50)

            VStack(alignment: .leading, spacing: 2) {
                Text(match.homeTeam)
                    .font(.srBodyMedium)
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
                Text(match.awayTeam)
                    .font(.srBodyMedium)
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
            }

            Spacer()

            if let tournament = match.tournament {
                Text(tournament)
                    .font(.system(size: 10))
                    .foregroundColor(.srTextMuted)
                    .lineLimit(1)
                    .frame(maxWidth: 80, alignment: .trailing)
            }
        }
        .padding(.horizontal, AppSpacing.md)
        .padding(.vertical, AppSpacing.sm)
    }
}
