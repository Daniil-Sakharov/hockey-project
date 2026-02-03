import SwiftUI

struct PlayerStatsHistoryView: View {
    let groupedHistory: [(season: String, entries: [PlayerStatEntryDTO])]

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.md) {
            Text("История по сезонам")
                .font(.srHeading4)
                .foregroundColor(.srTextPrimary)

            ForEach(groupedHistory, id: \.season) { group in
                seasonSection(group.season, entries: group.entries)
            }
        }
    }

    private func seasonSection(_ season: String, entries: [PlayerStatEntryDTO]) -> some View {
        VStack(alignment: .leading, spacing: AppSpacing.xs) {
            Text(season)
                .font(.system(size: 11, weight: .semibold))
                .foregroundColor(.srCyan)
                .textCase(.uppercase)

            VStack(spacing: 0) {
                headerRow
                ForEach(entries) { entry in
                    entryRow(entry)
                    if entry.id != entries.last?.id {
                        Divider().background(Color.srBorder.opacity(0.2))
                    }
                }
            }
            .glassCard(padding: 0)
        }
    }

    private var headerRow: some View {
        HStack(spacing: 0) {
            Text("Турнир").frame(maxWidth: .infinity, alignment: .leading)
            statHeader("И")
            statHeader("Г")
            statHeader("П")
            statHeader("О")
            statHeader("+/-")
        }
        .font(.system(size: 9, weight: .semibold))
        .foregroundColor(.srTextMuted)
        .textCase(.uppercase)
        .padding(.horizontal, AppSpacing.sm)
        .padding(.vertical, AppSpacing.xs)
    }

    private func entryRow(_ e: PlayerStatEntryDTO) -> some View {
        HStack(spacing: 0) {
            VStack(alignment: .leading, spacing: 1) {
                Text(e.tournamentName)
                    .font(.system(size: 11, weight: .medium))
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
                if !e.groupName.isEmpty {
                    Text(e.groupName)
                        .font(.system(size: 9))
                        .foregroundColor(.srTextMuted)
                        .lineLimit(1)
                }
            }
            .frame(maxWidth: .infinity, alignment: .leading)

            statCell(e.games)
            statCell(e.goals)
            statCell(e.assists)
            Text("\(e.points)")
                .font(.system(size: 11, weight: .bold))
                .foregroundColor(.srPurple)
                .frame(width: 28, alignment: .center)
            Text("\(e.plusMinus)")
                .font(.system(size: 11))
                .foregroundColor(e.plusMinus >= 0 ? .srSuccess : .srError)
                .frame(width: 28, alignment: .center)
        }
        .padding(.horizontal, AppSpacing.sm)
        .padding(.vertical, 6)
    }

    private func statHeader(_ title: String) -> some View {
        Text(title).frame(width: 28, alignment: .center)
    }

    private func statCell(_ value: Int) -> some View {
        Text("\(value)")
            .font(.system(size: 11))
            .foregroundColor(.srTextSecondary)
            .frame(width: 28, alignment: .center)
    }
}
