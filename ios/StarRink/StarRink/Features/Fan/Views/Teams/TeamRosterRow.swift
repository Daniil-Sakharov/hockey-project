import SwiftUI

struct TeamRosterRow: View {
    let player: PlayerItemDTO

    var body: some View {
        HStack(spacing: AppSpacing.sm) {
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
                Text(positionLocalized(player.position))
                    .font(.srCaption)
                    .foregroundColor(.srTextSecondary)
            }

            Spacer()

            if let number = player.jerseyNumber, number > 0 {
                Text("#\(number)")
                    .font(.system(size: 14, weight: .bold, design: .monospaced))
                    .foregroundColor(.srCyan)
            }

            if let stats = player.stats {
                HStack(spacing: AppSpacing.xs) {
                    Text("\(stats.goals)Г")
                        .font(.system(size: 11))
                        .foregroundColor(.srCyan)
                    Text("\(stats.assists)П")
                        .font(.system(size: 11))
                        .foregroundColor(.srPurple)
                }
            }
        }
        .padding(.horizontal, AppSpacing.md)
        .padding(.vertical, AppSpacing.sm)
    }

    private func positionLocalized(_ pos: String) -> String {
        switch pos {
        case "forward": return "Нападающий"
        case "defender": return "Защитник"
        case "goalie": return "Вратарь"
        default: return pos
        }
    }
}
