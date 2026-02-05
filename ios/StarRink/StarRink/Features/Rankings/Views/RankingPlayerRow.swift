import SwiftUI

struct RankingPlayerRow: View {
    let player: RankedPlayerDTO

    private var rankColor: Color {
        switch player.rank {
        case 1: return .srAmber
        case 2: return .srTextSecondary
        case 3: return Color(red: 0.72, green: 0.45, blue: 0.2)
        default: return .srCyan
        }
    }

    var body: some View {
        NavigationLink(value: PlayerRoute(playerId: player.id)) {
            HStack(spacing: AppSpacing.sm) {
                Text("\(player.rank)")
                    .font(.srHeading4)
                    .foregroundColor(rankColor)
                    .frame(width: 30)

                CachedAsyncImage(url: URL(string: player.photoUrl ?? "")) {
                    Image(systemName: "person.circle.fill")
                        .resizable()
                        .foregroundColor(.srTextMuted)
                }
                .frame(width: 36, height: 36)
                .clipShape(Circle())

                VStack(alignment: .leading, spacing: 2) {
                    Text(player.name)
                        .font(.srBodyMedium)
                        .foregroundColor(.srTextPrimary)
                        .lineLimit(1)
                    Text(player.team)
                        .font(.srCaption)
                        .foregroundColor(.srTextSecondary)
                        .lineLimit(1)
                }

                Spacer()

                HStack(spacing: AppSpacing.sm) {
                    RankingStatCell(value: player.games, label: "И")
                    RankingStatCell(value: player.goals, label: "Г")
                    RankingStatCell(value: player.assists, label: "П")
                    RankingStatCell(value: player.points, label: "О")
                }
            }
            .padding(.horizontal, AppSpacing.md)
            .padding(.vertical, AppSpacing.sm)
        }
    }
}

struct RankingStatCell: View {
    let value: Int
    let label: String

    var body: some View {
        VStack(spacing: 1) {
            Text("\(value)")
                .font(.system(size: 13, weight: .semibold, design: .monospaced))
                .foregroundColor(.srTextPrimary)
            Text(label)
                .font(.system(size: 9))
                .foregroundColor(.srTextMuted)
        }
        .frame(width: 28)
    }
}
