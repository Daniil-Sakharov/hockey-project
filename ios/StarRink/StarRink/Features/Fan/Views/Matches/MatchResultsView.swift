import SwiftUI

struct MatchResultsView: View {
    @StateObject private var viewModel = MatchResultsViewModel()
    @State private var selectedTournament = "Все"

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.md) {
                headerSection
                tournamentFilter
                resultsList
            }
            .padding(.top, AppSpacing.md)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
        .overlay {
            if viewModel.isLoading && viewModel.results.isEmpty {
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
            Text("Результаты матчей")
                .font(.srHeading3)
                .foregroundColor(.srTextPrimary)
            Spacer()
            Image(systemName: "clipboard.fill")
                .font(.system(size: 24))
                .foregroundColor(.srPurple)
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

    private var resultsList: some View {
        LazyVStack(spacing: AppSpacing.lg) {
            if viewModel.groupedByDate.isEmpty && !viewModel.isLoading {
                Text("Нет результатов")
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

                    VStack(spacing: AppSpacing.sm) {
                        ForEach(group.matches) { match in
                            MatchResultCard(match: match)
                        }
                    }
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
