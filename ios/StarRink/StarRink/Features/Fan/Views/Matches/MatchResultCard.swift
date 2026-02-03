import SwiftUI

struct MatchResultCard: View {
    let match: MatchDTO

    private var homeWins: Bool { (match.homeScore ?? 0) > (match.awayScore ?? 0) }
    private var awayWins: Bool { (match.awayScore ?? 0) > (match.homeScore ?? 0) }

    var body: some View {
        VStack(spacing: 0) {
            HStack {
                Text(match.homeTeam)
                    .font(.srBodyMedium)
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
                Spacer()
                Text("\(match.homeScore ?? 0)")
                    .font(.srHeading4)
                    .foregroundColor(homeWins ? .srCyan : .srTextSecondary)
            }

            HStack {
                Text(match.awayTeam)
                    .font(.srBodyMedium)
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
                Spacer()
                Text("\(match.awayScore ?? 0)")
                    .font(.srHeading4)
                    .foregroundColor(awayWins ? .srCyan : .srTextSecondary)
            }

            if let tournament = match.tournament {
                HStack {
                    Text(tournament)
                        .font(.system(size: 10))
                        .foregroundColor(.srTextMuted)
                        .lineLimit(1)
                    Spacer()
                    if let venue = match.venue, !venue.isEmpty {
                        Text(venue)
                            .font(.system(size: 10))
                            .foregroundColor(.srTextMuted)
                            .lineLimit(1)
                    }
                }
                .padding(.top, AppSpacing.xxs)
            }
        }
        .glassCard()
    }
}
