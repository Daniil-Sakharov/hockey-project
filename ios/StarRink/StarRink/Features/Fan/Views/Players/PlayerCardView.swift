import SwiftUI

struct PlayerCardView: View {
    let player: PlayerItemDTO

    private var positionShort: String {
        switch player.position {
        case "forward": return "НАП"
        case "defender": return "ЗАЩ"
        case "goalie": return "ВРТ"
        default: return player.position.prefix(3).uppercased()
        }
    }

    private var positionColor: Color {
        switch player.position {
        case "forward": return .srCyan
        case "defender": return .srAmber
        case "goalie": return .srPurple
        default: return .srTextSecondary
        }
    }

    var body: some View {
        HStack(spacing: AppSpacing.md) {
            avatar
            playerInfo
            Spacer()
            statsColumn
        }
        .glassCard()
    }

    private var avatar: some View {
        CachedAsyncImage(
            url: player.photoUrl.flatMap { URL(string: $0) }
        ) {
            avatarPlaceholder
        }
        .frame(width: 48, height: 48)
        .clipShape(Circle())
    }

    private var avatarPlaceholder: some View {
        ZStack {
            Circle().fill(positionColor.opacity(0.15))
            Text(String(player.name.prefix(1)))
                .font(.srBodyMedium)
                .foregroundColor(positionColor)
        }
    }

    private var playerInfo: some View {
        VStack(alignment: .leading, spacing: 3) {
            Text(player.name)
                .font(.srBodyMedium)
                .foregroundColor(.srTextPrimary)
                .lineLimit(1)
            Text(player.team)
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)
                .lineLimit(1)
            HStack(spacing: AppSpacing.sm) {
                Text(positionShort)
                    .font(.system(size: 9, weight: .bold))
                    .foregroundColor(positionColor)
                    .padding(.horizontal, 6)
                    .padding(.vertical, 2)
                    .background(positionColor.opacity(0.15))
                    .clipShape(Capsule())
                Text(String(player.birthYear))
                    .font(.system(size: 10, weight: .semibold))
                    .foregroundColor(.srAmber)
                    .padding(.horizontal, 6)
                    .padding(.vertical, 2)
                    .background(Color.srAmber.opacity(0.15))
                    .clipShape(Capsule())
            }
        }
    }

    @ViewBuilder
    private var statsColumn: some View {
        if let stats = player.stats {
            VStack(alignment: .trailing, spacing: 4) {
                Text("\(stats.points)")
                    .font(.system(size: 18, weight: .bold, design: .rounded))
                    .foregroundColor(.srTextPrimary)
                + Text(" очк")
                    .font(.system(size: 10, weight: .medium))
                    .foregroundColor(.srTextSecondary)
                HStack(spacing: 10) {
                    HStack(spacing: 2) {
                        Text("\(stats.goals)")
                            .font(.system(size: 12, weight: .bold))
                            .foregroundColor(.srCyan)
                        Text("гол.")
                            .font(.system(size: 9, weight: .medium))
                            .foregroundColor(.srCyan)
                    }
                    HStack(spacing: 2) {
                        Text("\(stats.assists)")
                            .font(.system(size: 12, weight: .bold))
                            .foregroundColor(.srPurple)
                        Text("пер.")
                            .font(.system(size: 9, weight: .medium))
                            .foregroundColor(.srPurple)
                    }
                }
                Text("\(stats.games) игр")
                    .font(.system(size: 9, weight: .medium))
                    .foregroundColor(.srTextSecondary)
            }
        }
    }
}
