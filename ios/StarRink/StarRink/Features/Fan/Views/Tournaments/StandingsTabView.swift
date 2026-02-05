import SwiftUI

struct StandingsTabView: View {
    let standings: [StandingDTO]

    var body: some View {
        if standings.isEmpty {
            emptyView
        } else {
            ScrollView {
                VStack(spacing: 0) {
                    headerRow
                    ForEach(standings) { row in
                        standingRow(row)
                        if row.id != standings.last?.id {
                            Divider().background(Color.srBorder.opacity(0.2))
                        }
                    }
                }
                .glassCard(padding: 0)
                .padding(.horizontal, AppSpacing.screenHorizontal)
                .padding(.bottom, 100)
            }
            .scrollContentBackground(.hidden)
        }
    }

    private var headerRow: some View {
        HStack(spacing: 0) {
            Text("#").frame(width: 28, alignment: .center)
            Text("Команда").frame(maxWidth: .infinity, alignment: .leading)
            statHeader("И")
            statHeader("В")
            statHeader("Н")
            statHeader("П")
            statHeader("ШЗ")
            statHeader("ШП")
            statHeader("О")
        }
        .font(.system(size: 10, weight: .semibold))
        .foregroundColor(.srTextMuted)
        .textCase(.uppercase)
        .padding(.horizontal, AppSpacing.sm)
        .padding(.vertical, AppSpacing.xs)
    }

    private func standingRow(_ s: StandingDTO) -> some View {
        NavigationLink(value: TeamRoute(teamId: s.teamId, teamName: s.team)) {
            HStack(spacing: 0) {
                positionBadge(s.position)
                    .frame(width: 28)
                Text(s.team)
                    .font(.system(size: 12, weight: .medium))
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
                    .frame(maxWidth: .infinity, alignment: .leading)
                statCell(s.games)
                statCell(s.wins)
                statCell(s.draws)
                statCell(s.losses)
                statCell(s.goalsFor)
                statCell(s.goalsAgainst)
                Text("\(s.points)")
                    .font(.system(size: 12, weight: .bold))
                    .foregroundColor(.srCyan)
                    .frame(width: 28, alignment: .center)
            }
            .padding(.horizontal, AppSpacing.sm)
            .padding(.vertical, 8)
        }
    }

    private func positionBadge(_ pos: Int) -> some View {
        Group {
            if pos <= 3 {
                Text("\(pos)")
                    .font(.system(size: 11, weight: .bold))
                    .foregroundColor(medalColor(pos))
            } else {
                Text("\(pos)")
                    .font(.system(size: 11))
                    .foregroundColor(.srTextSecondary)
            }
        }
    }

    private func medalColor(_ pos: Int) -> Color {
        switch pos {
        case 1: return .srAmber
        case 2: return .srTextSecondary
        case 3: return Color(red: 0.8, green: 0.5, blue: 0.2)
        default: return .srTextSecondary
        }
    }

    private func statHeader(_ title: String) -> some View {
        Text(title).frame(width: 28, alignment: .center)
    }

    private func statCell(_ value: Int) -> some View {
        Text("\(value)")
            .font(.system(size: 12))
            .foregroundColor(.srTextSecondary)
            .frame(width: 28, alignment: .center)
    }

    private var emptyView: some View {
        VStack(spacing: AppSpacing.sm) {
            Image(systemName: "table")
                .font(.system(size: 36))
                .foregroundColor(.srTextMuted)
            Text("Нет данных таблицы")
                .font(.srBody)
                .foregroundColor(.srTextSecondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .padding(.top, AppSpacing.xxl)
    }
}
