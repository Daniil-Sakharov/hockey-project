import SwiftUI

struct MatchLineupSection: View {
    let title: String
    let players: [LineupPlayerDTO]

    private var skaters: [LineupPlayerDTO] {
        players.filter { $0.position != "goalie" }
    }

    private var goalies: [LineupPlayerDTO] {
        players.filter { $0.position == "goalie" }
    }

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            Text(title)
                .font(.system(size: 11, weight: .semibold))
                .foregroundColor(.srTextMuted)
                .textCase(.uppercase)
                .tracking(1)

            if !skaters.isEmpty {
                skatersTable
            }
            if !goalies.isEmpty {
                goaliesTable
            }
        }
    }

    // MARK: - Skaters

    private var skatersTable: some View {
        VStack(spacing: 0) {
            skaterHeader
            ForEach(skaters) { player in
                skaterRow(player)
                if player.id != skaters.last?.id {
                    Divider().background(Color.srBorder.opacity(0.2))
                }
            }
        }
        .glassCard(padding: 0)
    }

    private var skaterHeader: some View {
        HStack(spacing: 0) {
            Text("#").frame(width: 28, alignment: .center)
            Text("Игрок").frame(maxWidth: .infinity, alignment: .leading)
            statHeader("Г")
            statHeader("П")
            statHeader("+/-")
            statHeader("Шт")
        }
        .font(.system(size: 10, weight: .semibold))
        .foregroundColor(.srTextMuted)
        .padding(.horizontal, AppSpacing.sm)
        .padding(.vertical, AppSpacing.xs)
    }

    private func skaterRow(_ p: LineupPlayerDTO) -> some View {
        NavigationLink(value: PlayerRoute(playerId: p.playerId)) {
            HStack(spacing: 0) {
                Text(p.jerseyNumber != nil ? String(p.jerseyNumber!) : "-")
                    .font(.system(size: 11, weight: .bold))
                    .foregroundColor(.srTextMuted)
                    .frame(width: 28, alignment: .center)

                Text(p.playerName)
                    .font(.system(size: 12, weight: .medium))
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
                    .frame(maxWidth: .infinity, alignment: .leading)

                statCell(p.goals, highlight: p.goals > 0)
                statCell(p.assists, highlight: p.assists > 0)
                plusMinusCell(p.plusMinus)
                statCell(p.penaltyMinutes, highlight: p.penaltyMinutes > 0, color: .srAmber)
            }
            .padding(.horizontal, AppSpacing.sm)
            .padding(.vertical, 6)
        }
    }

    // MARK: - Goalies

    private var goaliesTable: some View {
        VStack(spacing: 0) {
            goalieHeader
            ForEach(goalies) { player in
                goalieRow(player)
            }
        }
        .glassCard(padding: 0)
    }

    private var goalieHeader: some View {
        HStack(spacing: 0) {
            Text("#").frame(width: 28, alignment: .center)
            Text("Вратарь").frame(maxWidth: .infinity, alignment: .leading)
            statHeader("СВ")
            statHeader("ПГ")
        }
        .font(.system(size: 10, weight: .semibold))
        .foregroundColor(.srTextMuted)
        .padding(.horizontal, AppSpacing.sm)
        .padding(.vertical, AppSpacing.xs)
    }

    private func goalieRow(_ p: LineupPlayerDTO) -> some View {
        NavigationLink(value: PlayerRoute(playerId: p.playerId)) {
            HStack(spacing: 0) {
                Text(p.jerseyNumber != nil ? String(p.jerseyNumber!) : "-")
                    .font(.system(size: 11, weight: .bold))
                    .foregroundColor(.srTextMuted)
                    .frame(width: 28, alignment: .center)

                Text(p.playerName)
                    .font(.system(size: 12, weight: .medium))
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
                    .frame(maxWidth: .infinity, alignment: .leading)

                statCell(p.saves ?? 0, highlight: false)
                statCell(p.goalsAgainst ?? 0, highlight: (p.goalsAgainst ?? 0) > 0, color: .srAmber)
            }
            .padding(.horizontal, AppSpacing.sm)
            .padding(.vertical, 6)
        }
    }

    // MARK: - Helpers

    private func statHeader(_ title: String) -> some View {
        Text(title).frame(width: 32, alignment: .center)
    }

    private func statCell(_ value: Int, highlight: Bool, color: Color = .srCyan) -> some View {
        Text(String(value))
            .font(.system(size: 12, weight: highlight ? .bold : .regular))
            .foregroundColor(highlight ? color : .srTextSecondary)
            .frame(width: 32, alignment: .center)
    }

    private func plusMinusCell(_ value: Int) -> some View {
        let text = value > 0 ? "+\(value)" : String(value)
        let color: Color = value > 0 ? .srSuccess : value < 0 ? .srAmber : .srTextSecondary
        return Text(text)
            .font(.system(size: 12, weight: value != 0 ? .bold : .regular))
            .foregroundColor(color)
            .frame(width: 32, alignment: .center)
    }
}
