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

struct MatchRoute: Hashable {
    let matchId: String
}

struct TeamRoute: Hashable {
    let teamId: String
    let teamName: String
}

struct TeamRosterRoute: Hashable {
    let teamId: String
    let tournamentId: String
    let teamName: String
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
        Region(id: "pfo", name: "ПФО", icon: "building.2.fill", domain: "https://pfo.fhr.ru", source: "junior"),
        Region(id: "ufo", name: "УрФО", icon: "mountain.2.fill", domain: "https://ufo.fhr.ru", source: "junior"),
        Region(id: "cfo", name: "ЦФО", icon: "mappin.circle.fill", domain: "https://cfo.fhr.ru", source: "junior"),
        Region(id: "sfo", name: "СФО", icon: "wind.snow", domain: "https://sfo.fhr.ru", source: "junior"),
        Region(id: "szfo", name: "СЗФО", icon: "snowflake", domain: "https://szfo.fhr.ru", source: "junior"),
        Region(id: "yfo", name: "ЮФО", icon: "sun.max.fill", domain: "https://yfo.fhr.ru", source: "junior"),
        Region(id: "dfo", name: "ДФО", icon: "globe.asia.australia.fill", domain: "https://dfo.fhr.ru", source: "junior"),
        Region(id: "junior", name: "Юниор", icon: "star.fill", domain: "https://junior.fhr.ru", source: "junior"),
        Region(id: "spb", name: "СПб", icon: "building.columns.fill", domain: "https://spb.fhr.ru", source: "fhspb"),
        Region(id: "komi", name: "Коми", icon: "tree.fill", domain: "https://komi.fhr.ru", source: "junior"),
        Region(id: "sam", name: "Самара", icon: "leaf.fill", domain: "https://sam.fhr.ru", source: "junior"),
        Region(id: "nsk", name: "Новосибирск", icon: "building.fill", domain: "https://nsk.fhr.ru", source: "junior"),
        Region(id: "vrn", name: "Воронеж", icon: "house.fill", domain: "https://vrn.fhr.ru", source: "junior"),
        Region(id: "len", name: "Ленобласть", icon: "map.fill", domain: "https://len.fhr.ru", source: "junior"),
    ]

    static func fromDomain(_ domain: String) -> Region? {
        allRegions.first { $0.domain == domain }
    }

    static func labelForDomain(_ domain: String) -> String {
        fromDomain(domain)?.name ?? domain
            .replacingOccurrences(of: "https://", with: "")
            .replacingOccurrences(of: ".fhr.ru", with: "")
            .uppercased()
    }
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
