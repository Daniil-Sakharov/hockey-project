import SwiftUI

struct MatchDetailView: View {
    let matchId: String
    @StateObject private var viewModel = MatchDetailViewModel()

    var body: some View {
        ZStack {
            AnimatedBackgroundView().ignoresSafeArea()
            content
        }
        .navigationTitle("Матч")
        .navigationBarTitleDisplayMode(.inline)
        .toolbarBackground(.hidden, for: .navigationBar)
        .task { await viewModel.load(matchId: matchId) }
    }

    @ViewBuilder
    private var content: some View {
        if viewModel.isLoading {
            ProgressView().tint(.srCyan)
        } else if let error = viewModel.errorMessage {
            errorView(error)
        } else if let match = viewModel.match {
            matchContent(match)
        }
    }

    private func matchContent(_ match: MatchDetailResponse) -> some View {
        ScrollView {
            VStack(spacing: AppSpacing.md) {
                heroCard(match)
                if let sbp = match.scoreByPeriod {
                    periodScoresCard(sbp, resultType: match.resultType)
                }
                if !viewModel.eventsByPeriod.isEmpty {
                    eventsSection
                }
                if !match.homeLineup.isEmpty {
                    MatchLineupSection(
                        title: match.homeTeam.name,
                        players: match.homeLineup
                    )
                }
                if !match.awayLineup.isEmpty {
                    MatchLineupSection(
                        title: match.awayTeam.name,
                        players: match.awayLineup
                    )
                }
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.top, AppSpacing.sm)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
    }

    // MARK: - Hero

    private func heroCard(_ m: MatchDetailResponse) -> some View {
        VStack(spacing: AppSpacing.sm) {
            HStack {
                teamColumn(m.homeTeam)
                Spacer()
                scoreColumn(m)
                Spacer()
                teamColumn(m.awayTeam)
            }

            HStack(spacing: AppSpacing.xs) {
                Text(m.tournament.name)
                    .font(.system(size: 11))
                    .foregroundColor(.srTextMuted)
                    .lineLimit(1)
            }

            HStack(spacing: AppSpacing.sm) {
                Label(m.date, systemImage: "calendar")
                Label(m.time, systemImage: "clock")
                if let venue = m.venue, !venue.isEmpty {
                    Label(venue, systemImage: "mappin")
                }
            }
            .font(.system(size: 10))
            .foregroundColor(.srTextMuted)
        }
        .glassCard()
    }

    private func teamColumn(_ team: MatchTeamDTO) -> some View {
        NavigationLink(value: TeamRoute(teamId: team.id, teamName: team.name)) {
            VStack(spacing: AppSpacing.xs) {
                teamLogo(team.logoUrl)
                Text(team.name)
                    .font(.system(size: 12, weight: .medium))
                    .foregroundColor(.srTextPrimary)
                    .multilineTextAlignment(.center)
                    .lineLimit(2)
            }
            .frame(width: 100)
        }
    }

    private func teamLogo(_ urlString: String?) -> some View {
        Group {
            if let urlString, let url = URL(string: urlString) {
                CachedAsyncImage(url: url) {
                    teamLogoPlaceholder
                }
                .frame(width: 48, height: 48)
                .clipShape(Circle())
            } else {
                teamLogoPlaceholder
            }
        }
    }

    private var teamLogoPlaceholder: some View {
        Circle()
            .fill(Color.srCard)
            .frame(width: 48, height: 48)
            .overlay(
                Image(systemName: "shield.fill")
                    .foregroundColor(.srTextMuted)
            )
    }

    private func scoreColumn(_ m: MatchDetailResponse) -> some View {
        VStack(spacing: 2) {
            if m.status == "finished", let h = m.homeScore, let a = m.awayScore {
                Text("\(h) : \(a)")
                    .font(.system(size: 28, weight: .bold))
                    .foregroundColor(.srCyan)
                if let rt = m.resultType, !rt.isEmpty {
                    Text(rt.uppercased())
                        .font(.system(size: 10, weight: .semibold))
                        .foregroundColor(.srTextMuted)
                }
            } else {
                Text("vs")
                    .font(.system(size: 20, weight: .bold))
                    .foregroundColor(.srTextMuted)
            }
        }
    }

    // MARK: - Period Scores

    private func periodScoresCard(_ sbp: ScoreByPeriodDTO, resultType: String?) -> some View {
        VStack(spacing: AppSpacing.xs) {
            Text("Счёт по периодам")
                .font(.system(size: 11, weight: .semibold))
                .foregroundColor(.srTextMuted)
                .textCase(.uppercase)
                .tracking(1)

            HStack(spacing: 0) {
                Text("").frame(width: 60, alignment: .leading)
                periodHeader("П1")
                periodHeader("П2")
                periodHeader("П3")
                if sbp.homeOt != nil || sbp.awayOt != nil {
                    periodHeader("ОТ")
                }
            }
            .font(.system(size: 10, weight: .semibold))
            .foregroundColor(.srTextMuted)

            periodRow("Дом", values: [sbp.homeP1, sbp.homeP2, sbp.homeP3, sbp.homeOt])
            periodRow("Гости", values: [sbp.awayP1, sbp.awayP2, sbp.awayP3, sbp.awayOt])
        }
        .glassCard()
    }

    private func periodHeader(_ title: String) -> some View {
        Text(title).frame(maxWidth: .infinity)
    }

    private func periodRow(_ label: String, values: [Int?]) -> some View {
        HStack(spacing: 0) {
            Text(label)
                .font(.system(size: 12, weight: .medium))
                .foregroundColor(.srTextPrimary)
                .frame(width: 60, alignment: .leading)
            ForEach(Array(values.enumerated()), id: \.offset) { _, val in
                if let v = val {
                    Text(String(v))
                        .font(.system(size: 13, weight: .bold))
                        .foregroundColor(.srCyan)
                        .frame(maxWidth: .infinity)
                } else {
                    Text("-")
                        .font(.system(size: 13))
                        .foregroundColor(.srTextMuted)
                        .frame(maxWidth: .infinity)
                }
            }
        }
    }

    // MARK: - Events

    private var eventsSection: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            Text("События")
                .font(.system(size: 11, weight: .semibold))
                .foregroundColor(.srTextMuted)
                .textCase(.uppercase)
                .tracking(1)

            ForEach(Array(viewModel.eventsByPeriod.enumerated()), id: \.offset) { _, group in
                VStack(alignment: .leading, spacing: 0) {
                    Text(group.period)
                        .font(.system(size: 11, weight: .semibold))
                        .foregroundColor(.srAmber)
                        .padding(.horizontal, AppSpacing.sm)
                        .padding(.vertical, AppSpacing.xs)

                    ForEach(group.events) { event in
                        MatchEventRow(event: event)
                        if event.id != group.events.last?.id {
                            Divider().background(Color.srBorder.opacity(0.2))
                        }
                    }
                }
                .glassCard(padding: 0)
            }
        }
    }

    // MARK: - Error

    private func errorView(_ message: String) -> some View {
        VStack(spacing: AppSpacing.sm) {
            Image(systemName: "exclamationmark.triangle")
                .font(.system(size: 36))
                .foregroundColor(.srAmber)
            Text(message)
                .font(.srBody)
                .foregroundColor(.srTextSecondary)
            Button("Повторить") {
                Task { await viewModel.load(matchId: matchId) }
            }
            .font(.srBodyMedium)
            .foregroundColor(.srCyan)
        }
    }
}
