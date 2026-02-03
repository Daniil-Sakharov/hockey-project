import SwiftUI

struct FanStatCard: View {
    let value: String
    let label: String
    let icon: String
    let color: Color

    var body: some View {
        VStack(spacing: AppSpacing.xs) {
            Image(systemName: icon)
                .font(.system(size: 20))
                .foregroundColor(color)

            Text(value)
                .font(.srHeading3)
                .foregroundColor(.srTextPrimary)

            Text(label)
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)
        }
        .frame(maxWidth: .infinity)
        .glassCard(padding: AppSpacing.md)
    }
}

struct FanScorerRow: View {
    let scorer: RankedPlayerDTO

    private var rankColor: Color {
        switch scorer.rank {
        case 1: return .srAmber
        case 2: return .srTextSecondary
        case 3: return Color(red: 0.72, green: 0.45, blue: 0.2)
        default: return .srCyan
        }
    }

    var body: some View {
        HStack(spacing: AppSpacing.md) {
            Text("\(scorer.rank)")
                .font(.srHeading4)
                .foregroundColor(rankColor)
                .frame(width: 28)

            VStack(alignment: .leading, spacing: 2) {
                Text(scorer.name)
                    .font(.srBodyMedium)
                    .foregroundColor(.srTextPrimary)
                Text(scorer.team)
                    .font(.srCaption)
                    .foregroundColor(.srTextSecondary)
            }

            Spacer()

            HStack(spacing: AppSpacing.md) {
                StatColumn(value: scorer.goals, label: "Г", color: .srCyan)
                StatColumn(value: scorer.assists, label: "П", color: .srPurple)
                StatColumn(value: scorer.points, label: "О", color: .srAmber)
            }
        }
        .padding(AppSpacing.md)
    }
}

struct FanMatchCard: View {
    let match: MatchDTO

    private var homeWins: Bool { (match.homeScore ?? 0) > (match.awayScore ?? 0) }
    private var awayWins: Bool { (match.awayScore ?? 0) > (match.homeScore ?? 0) }

    var body: some View {
        HStack {
            VStack(alignment: .leading, spacing: 4) {
                Text(match.homeTeam)
                    .font(.srBodyMedium)
                    .foregroundColor(homeWins ? .srTextPrimary : .srTextSecondary)
                Text(match.awayTeam)
                    .font(.srBodyMedium)
                    .foregroundColor(awayWins ? .srTextPrimary : .srTextSecondary)
            }

            Spacer()

            if let hs = match.homeScore, let as_ = match.awayScore {
                VStack(spacing: 4) {
                    Text("\(hs)")
                        .font(.srHeading4)
                        .foregroundColor(homeWins ? .srCyan : .srTextSecondary)
                    Text("\(as_)")
                        .font(.srHeading4)
                        .foregroundColor(awayWins ? .srCyan : .srTextSecondary)
                }
            }

            Text(formatDate(match.date))
                .font(.srCaption)
                .foregroundColor(.srTextMuted)
                .frame(width: 70, alignment: .trailing)
        }
        .glassCard()
    }

    private func formatDate(_ dateStr: String) -> String {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd"
        guard let date = formatter.date(from: dateStr) else { return dateStr }

        let calendar = Calendar.current
        if calendar.isDateInToday(date) { return "Сегодня" }
        if calendar.isDateInYesterday(date) { return "Вчера" }

        let display = DateFormatter()
        display.dateFormat = "dd.MM"
        return display.string(from: date)
    }
}

// MARK: - Helpers

private struct StatColumn: View {
    let value: Int
    let label: String
    let color: Color

    var body: some View {
        VStack(spacing: 2) {
            Text("\(value)")
                .font(.srBodyMedium)
                .foregroundColor(color)
            Text(label)
                .font(.system(size: 9))
                .foregroundColor(.srTextMuted)
        }
    }
}
