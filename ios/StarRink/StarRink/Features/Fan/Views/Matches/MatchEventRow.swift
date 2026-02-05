import SwiftUI

struct MatchEventRow: View {
    let event: MatchEventDTO

    var body: some View {
        HStack(spacing: AppSpacing.sm) {
            timeLabel
            eventIcon
            eventContent
            Spacer()
        }
        .padding(.horizontal, AppSpacing.sm)
        .padding(.vertical, 8)
    }

    private var timeLabel: some View {
        Text(event.time)
            .font(.system(size: 11, weight: .medium, design: .monospaced))
            .foregroundColor(.srTextMuted)
            .frame(width: 36, alignment: .center)
    }

    private var eventIcon: some View {
        Group {
            if event.type == "goal" {
                Image(systemName: "hockey.puck.fill")
                    .foregroundColor(.srCyan)
            } else {
                Image(systemName: "clock.badge.exclamationmark")
                    .foregroundColor(.srAmber)
            }
        }
        .font(.system(size: 14))
        .frame(width: 24)
    }

    @ViewBuilder
    private var eventContent: some View {
        if event.type == "goal" {
            goalContent
        } else {
            penaltyContent
        }
    }

    private var goalContent: some View {
        VStack(alignment: .leading, spacing: 2) {
            if let name = event.playerName {
                playerLink(id: event.playerId, name: name)
                    .font(.system(size: 13, weight: .semibold))
                    .foregroundColor(.srTextPrimary)
            }

            if let a1 = event.assist1Name {
                HStack(spacing: 2) {
                    Text("Пас:")
                        .foregroundColor(.srTextMuted)
                    playerLink(id: event.assist1Id, name: a1)
                        .foregroundColor(.srTextSecondary)
                    if let a2 = event.assist2Name {
                        Text(",")
                            .foregroundColor(.srTextMuted)
                        playerLink(id: event.assist2Id, name: a2)
                            .foregroundColor(.srTextSecondary)
                    }
                }
                .font(.system(size: 11))
            }

            if let goalType = event.goalType, !goalType.isEmpty {
                Text(goalType)
                    .font(.system(size: 10))
                    .foregroundColor(.srPurple)
            }
        }
    }

    private var penaltyContent: some View {
        VStack(alignment: .leading, spacing: 2) {
            if let name = event.playerName {
                playerLink(id: event.playerId, name: name)
                    .font(.system(size: 13, weight: .medium))
                    .foregroundColor(.srTextPrimary)
            }

            HStack(spacing: AppSpacing.xxs) {
                if let mins = event.penaltyMins {
                    Text("\(mins) мин")
                        .font(.system(size: 11, weight: .semibold))
                        .foregroundColor(.srAmber)
                }
                if let text = event.penaltyText, !text.isEmpty {
                    Text(text)
                        .font(.system(size: 11))
                        .foregroundColor(.srTextMuted)
                        .lineLimit(1)
                }
            }
        }
    }

    private func playerLink(id: String?, name: String) -> some View {
        Group {
            if let id, !id.isEmpty {
                NavigationLink(value: PlayerRoute(playerId: id)) {
                    Text(name)
                }
            } else {
                Text(name)
            }
        }
    }
}
