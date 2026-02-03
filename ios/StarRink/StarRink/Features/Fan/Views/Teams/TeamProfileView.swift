import SwiftUI

struct TeamProfileView: View {
    let teamId: String
    let teamName: String

    @StateObject private var viewModel = TeamProfileViewModel()

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.lg) {
                if let team = viewModel.team {
                    heroSection(team)
                    statsSection(team.stats)
                    tournamentsSection(team.tournaments)
                    rosterSection(team.roster)
                    recentMatchesSection(team.recentMatches)
                }
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.top, AppSpacing.md)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
        .overlay {
            if viewModel.isLoading {
                ProgressView().tint(.srCyan)
            }
        }
        .task { await viewModel.load(teamId: teamId) }
    }

    private func heroSection(_ team: TeamProfileDTO) -> some View {
        VStack(spacing: AppSpacing.md) {
            CachedAsyncImage(url: URL(string: team.logoUrl ?? "")) {
                Image(systemName: "shield.fill")
                    .resizable()
                    .foregroundColor(.srTextMuted)
            }
            .frame(width: 80, height: 80)
            .clipShape(RoundedRectangle(cornerRadius: 16))

            Text(team.name)
                .font(.srHeading2)
                .foregroundColor(.srTextPrimary)
                .multilineTextAlignment(.center)

            if !team.city.isEmpty {
                Text(team.city)
                    .font(.srBody)
                    .foregroundColor(.srTextSecondary)
            }

            Text("\(team.playersCount) игроков")
                .font(.srCaption)
                .foregroundColor(.srTextMuted)
        }
        .frame(maxWidth: .infinity)
        .glassCard()
    }

    private func statsSection(_ stats: TeamStatsDTO) -> some View {
        VStack(alignment: .leading, spacing: AppSpacing.md) {
            Text("Статистика")
                .font(.srHeading4)
                .foregroundColor(.srTextPrimary)

            HStack(spacing: AppSpacing.sm) {
                FanStatCard(value: "\(stats.wins)", label: "Побед", icon: "checkmark.circle.fill", color: .srCyan)
                FanStatCard(value: "\(stats.losses)", label: "Поражений", icon: "xmark.circle.fill", color: .srPurple)
                FanStatCard(value: "\(stats.draws)", label: "Ничьих", icon: "equal.circle.fill", color: .srAmber)
            }

            HStack(spacing: AppSpacing.sm) {
                TeamStatRow(label: "Забито", value: "\(stats.goalsFor)", color: .srCyan)
                TeamStatRow(label: "Пропущено", value: "\(stats.goalsAgainst)", color: .srPurple)
            }
        }
    }

    private func tournamentsSection(_ tournaments: [String]) -> some View {
        Group {
            if !tournaments.isEmpty {
                VStack(alignment: .leading, spacing: AppSpacing.sm) {
                    Text("Турниры")
                        .font(.srHeading4)
                        .foregroundColor(.srTextPrimary)

                    FlowLayout(spacing: AppSpacing.xs) {
                        ForEach(tournaments, id: \.self) { name in
                            Text(name)
                                .font(.system(size: 12, weight: .medium))
                                .foregroundColor(.srCyan)
                                .padding(.horizontal, AppSpacing.sm)
                                .padding(.vertical, AppSpacing.xxs)
                                .background(Color.srCyan.opacity(0.15))
                                .clipShape(Capsule())
                        }
                    }
                }
            }
        }
    }

    private func rosterSection(_ roster: [PlayerItemDTO]) -> some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            Text("Состав (\(roster.count))")
                .font(.srHeading4)
                .foregroundColor(.srTextPrimary)

            VStack(spacing: 0) {
                ForEach(roster) { player in
                    TeamRosterRow(player: player)
                    Divider().background(Color.srBorder.opacity(0.3))
                }
            }
            .glassCard(padding: 0)
        }
    }

    private func recentMatchesSection(_ matches: [MatchDTO]) -> some View {
        Group {
            if !matches.isEmpty {
                VStack(alignment: .leading, spacing: AppSpacing.sm) {
                    Text("Последние матчи")
                        .font(.srHeading4)
                        .foregroundColor(.srTextPrimary)

                    VStack(spacing: AppSpacing.sm) {
                        ForEach(matches) { match in
                            FanMatchCard(match: match)
                        }
                    }
                }
            }
        }
    }
}

// MARK: - Helpers

private struct TeamStatRow: View {
    let label: String
    let value: String
    let color: Color

    var body: some View {
        HStack {
            Text(label)
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)
            Spacer()
            Text(value)
                .font(.srHeading4)
                .foregroundColor(color)
        }
        .frame(maxWidth: .infinity)
        .glassCard(padding: AppSpacing.md)
    }
}

private struct FlowLayout: Layout {
    var spacing: CGFloat

    func sizeThatFits(proposal: ProposedViewSize, subviews: Subviews, cache: inout ()) -> CGSize {
        let result = arrangeSubviews(proposal: proposal, subviews: subviews)
        return result.size
    }

    func placeSubviews(in bounds: CGRect, proposal: ProposedViewSize, subviews: Subviews, cache: inout ()) {
        let result = arrangeSubviews(proposal: proposal, subviews: subviews)
        for (index, position) in result.positions.enumerated() {
            subviews[index].place(at: CGPoint(x: bounds.minX + position.x, y: bounds.minY + position.y), proposal: .unspecified)
        }
    }

    private func arrangeSubviews(proposal: ProposedViewSize, subviews: Subviews) -> (size: CGSize, positions: [CGPoint]) {
        let maxWidth = proposal.width ?? .infinity
        var positions: [CGPoint] = []
        var x: CGFloat = 0
        var y: CGFloat = 0
        var rowHeight: CGFloat = 0

        for subview in subviews {
            let size = subview.sizeThatFits(.unspecified)
            if x + size.width > maxWidth && x > 0 {
                x = 0
                y += rowHeight + spacing
                rowHeight = 0
            }
            positions.append(CGPoint(x: x, y: y))
            rowHeight = max(rowHeight, size.height)
            x += size.width + spacing
        }

        return (CGSize(width: maxWidth, height: y + rowHeight), positions)
    }
}
