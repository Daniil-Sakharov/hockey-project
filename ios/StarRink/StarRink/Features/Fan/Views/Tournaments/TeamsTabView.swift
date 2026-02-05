import SwiftUI

struct TeamsTabView: View {
    let teams: [TeamItemDTO]
    let tournamentId: String

    var body: some View {
        if teams.isEmpty {
            emptyView
        } else {
            ScrollView {
                VStack(spacing: 0) {
                    ForEach(teams) { team in
                        NavigationLink(value: TeamRosterRoute(
                            teamId: team.id,
                            tournamentId: tournamentId,
                            teamName: team.name
                        )) {
                            teamRow(team)
                        }
                        if team.id != teams.last?.id {
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

    private func teamRow(_ team: TeamItemDTO) -> some View {
        HStack(spacing: AppSpacing.sm) {
            teamLogo(team.logoUrl)

            VStack(alignment: .leading, spacing: 2) {
                Text(team.name)
                    .font(.system(size: 14, weight: .medium))
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)

                HStack(spacing: AppSpacing.xs) {
                    if let city = team.city, !city.isEmpty {
                        Text(city)
                            .font(.system(size: 11))
                            .foregroundColor(.srTextSecondary)
                    }
                    if let group = team.groupName, !group.isEmpty {
                        Text(group)
                            .font(.system(size: 10, weight: .medium))
                            .foregroundColor(.srPurple)
                            .padding(.horizontal, 6)
                            .padding(.vertical, 2)
                            .background(Color.srPurple.opacity(0.15))
                            .clipShape(Capsule())
                    }
                }
            }

            Spacer()

            Text("\(team.playersCount)")
                .font(.system(size: 12, weight: .bold))
                .foregroundColor(.srCyan)
            Image(systemName: "person.2.fill")
                .font(.system(size: 10))
                .foregroundColor(.srTextMuted)
        }
        .padding(.horizontal, AppSpacing.sm)
        .padding(.vertical, 10)
    }

    private func teamLogo(_ urlString: String?) -> some View {
        Group {
            if let urlString, let url = URL(string: urlString) {
                CachedAsyncImage(url: url) {
                    logoPlaceholder
                }
                .frame(width: 36, height: 36)
                .clipShape(Circle())
            } else {
                logoPlaceholder
            }
        }
    }

    private var logoPlaceholder: some View {
        Circle()
            .fill(Color.srCard)
            .frame(width: 36, height: 36)
            .overlay(
                Image(systemName: "shield.fill")
                    .font(.system(size: 14))
                    .foregroundColor(.srTextMuted)
            )
    }

    private var emptyView: some View {
        VStack(spacing: AppSpacing.sm) {
            Image(systemName: "person.3")
                .font(.system(size: 36))
                .foregroundColor(.srTextMuted)
            Text("Нет команд")
                .font(.srBody)
                .foregroundColor(.srTextSecondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .padding(.top, AppSpacing.xxl)
    }
}
