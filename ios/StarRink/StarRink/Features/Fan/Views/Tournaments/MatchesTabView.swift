import SwiftUI

struct MatchesTabView: View {
    let matches: [MatchDTO]

    private var finishedMatches: [MatchDTO] {
        matches.filter { $0.status == "finished" }
    }

    private var scheduledMatches: [MatchDTO] {
        matches.filter { $0.status != "finished" }
    }

    var body: some View {
        if matches.isEmpty {
            emptyView
        } else {
            ScrollView {
                VStack(spacing: AppSpacing.md) {
                    if !finishedMatches.isEmpty {
                        matchSection("Завершённые", matches: finishedMatches)
                    }
                    if !scheduledMatches.isEmpty {
                        matchSection("Предстоящие", matches: scheduledMatches)
                    }
                }
                .padding(.horizontal, AppSpacing.screenHorizontal)
                .padding(.bottom, 100)
            }
            .scrollContentBackground(.hidden)
        }
    }

    private func matchSection(_ title: String, matches: [MatchDTO]) -> some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            Text(title)
                .font(.system(size: 11, weight: .semibold))
                .foregroundColor(.srTextMuted)
                .textCase(.uppercase)
                .tracking(1)

            VStack(spacing: 0) {
                ForEach(matches) { match in
                    matchRow(match)
                    if match.id != matches.last?.id {
                        Divider().background(Color.srBorder.opacity(0.2))
                    }
                }
            }
            .glassCard(padding: 0)
        }
    }

    private func matchRow(_ match: MatchDTO) -> some View {
        HStack(spacing: AppSpacing.sm) {
            VStack(alignment: .trailing, spacing: 2) {
                Text(match.homeTeam)
                    .font(.system(size: 13, weight: .medium))
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
            }
            .frame(maxWidth: .infinity, alignment: .trailing)

            scoreView(match)

            VStack(alignment: .leading, spacing: 2) {
                Text(match.awayTeam)
                    .font(.system(size: 13, weight: .medium))
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
            }
            .frame(maxWidth: .infinity, alignment: .leading)
        }
        .padding(.horizontal, AppSpacing.sm)
        .padding(.vertical, 10)
    }

    private func scoreView(_ match: MatchDTO) -> some View {
        Group {
            if match.status == "finished", let h = match.homeScore, let a = match.awayScore {
                Text("\(h) : \(a)")
                    .font(.system(size: 15, weight: .bold))
                    .foregroundColor(.srCyan)
            } else {
                VStack(spacing: 0) {
                    Text(match.date)
                        .font(.system(size: 10))
                    Text(match.time)
                        .font(.system(size: 10, weight: .medium))
                }
                .foregroundColor(.srTextMuted)
            }
        }
        .frame(width: 60)
    }

    private var emptyView: some View {
        VStack(spacing: AppSpacing.sm) {
            Image(systemName: "sportscourt")
                .font(.system(size: 36))
                .foregroundColor(.srTextMuted)
            Text("Нет матчей")
                .font(.srBody)
                .foregroundColor(.srTextSecondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .padding(.top, AppSpacing.xxl)
    }
}
