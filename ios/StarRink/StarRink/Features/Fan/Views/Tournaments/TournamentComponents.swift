import SwiftUI

// MARK: - Navigation Routes

struct TournamentRoute: Hashable {
    let tournamentId: String
    let name: String
    let birthYear: Int?
    let groupName: String?
}

struct PlayerRoute: Hashable {
    let playerId: String
}

// MARK: - Region Model

struct Region: Identifiable, Hashable {
    let id: String
    let name: String
    let icon: String
    let domain: String
    let source: String
}

extension Region {
    static let allRegions: [Region] = [
        Region(id: "volga", name: "Приволжский ФО", icon: "building.2.fill", domain: "volga.juniorhl.ru", source: "junior"),
        Region(id: "moscow", name: "Москва", icon: "star.fill", domain: "moscow.fhr.ru", source: "fhmoscow"),
        Region(id: "spb", name: "Санкт-Петербург", icon: "snowflake", domain: "spb.fhr.ru", source: "fhspb"),
        Region(id: "ural", name: "Урал", icon: "mountain.2.fill", domain: "ural.fhr.ru", source: "ural"),
        Region(id: "cfo", name: "ЦФО", icon: "mappin.circle.fill", domain: "cfo.fhr.ru", source: "cfo"),
        Region(id: "siberia", name: "Сибирь", icon: "wind.snow", domain: "siberia.fhr.ru", source: "siberia"),
    ]
}

// MARK: - Region Card

struct RegionCard: View {
    let region: Region

    var body: some View {
        VStack(spacing: AppSpacing.sm) {
            ZStack {
                Circle()
                    .fill(Color.srCyan.opacity(0.12))
                    .frame(width: 48, height: 48)
                Image(systemName: region.icon)
                    .font(.system(size: 20))
                    .foregroundColor(.srCyan)
            }

            Text(region.name)
                .font(.srCaption)
                .foregroundColor(.srTextPrimary)
                .multilineTextAlignment(.center)
                .lineLimit(2)
        }
        .frame(maxWidth: .infinity)
        .glassCard(padding: AppSpacing.md)
    }
}

// MARK: - Tournament Name Helpers

enum TournamentNameHelper {
    static func cleanName(_ name: String) -> String {
        var result = name
        // "До 11 Лет" must be checked BEFORE the number-only pattern
        if let range = result.range(of: #"\s*[Дд]о\s*\d+\s*[Лл]ет"#, options: .regularExpression) {
            result = String(result[result.startIndex..<range.lowerBound])
        }
        // "14/13/12 Лет", "14 лет" etc.
        if let range = result.range(of: #"\s*\d+(/\d+)*\s*[Лл]ет"#, options: .regularExpression) {
            result = String(result[result.startIndex..<range.lowerBound])
        }
        return result.trimmingCharacters(in: .whitespaces)
    }
}

// MARK: - Tournament Row

struct TournamentRow: View {
    let name: String
    let groupName: String?
    let teamsCount: Int
    let matchesCount: Int
    let isEnded: Bool

    private var statusText: String {
        isEnded ? "Завершён" : "Идёт"
    }

    private var statusColor: Color {
        isEnded ? .srTextMuted : .srSuccess
    }

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            HStack(spacing: AppSpacing.md) {
                ZStack {
                    RoundedRectangle(cornerRadius: 10)
                        .fill(Color.srCyan.opacity(0.15))
                        .frame(width: 44, height: 44)
                    Image(systemName: "trophy.fill")
                        .foregroundColor(.srCyan)
                }

                VStack(alignment: .leading, spacing: 2) {
                    Text(name)
                        .font(.srBodyMedium)
                        .foregroundColor(.srTextPrimary)
                        .lineLimit(2)
                    Text("\(teamsCount) команд · \(matchesCount) матчей")
                        .font(.srCaption)
                        .foregroundColor(.srTextSecondary)
                }

                Spacer()

                statusBadge
            }

            if let groupName, !groupName.isEmpty {
                groupBadge(groupName)
            }
        }
        .glassCard()
    }

    private var statusBadge: some View {
        Text(statusText)
            .font(.system(size: 11, weight: .medium))
            .foregroundColor(statusColor)
            .padding(.horizontal, 8)
            .padding(.vertical, 4)
            .background(
                Capsule().fill(statusColor.opacity(0.15))
            )
    }

    private func groupBadge(_ group: String) -> some View {
        HStack(spacing: 4) {
            Image(systemName: "square.grid.2x2")
                .font(.system(size: 9))
            Text(group)
                .font(.system(size: 11, weight: .medium))
        }
        .foregroundColor(.srPurple)
        .padding(.horizontal, 10)
        .padding(.vertical, 4)
        .background(
            Capsule().fill(Color.srPurple.opacity(0.15))
        )
    }
}
