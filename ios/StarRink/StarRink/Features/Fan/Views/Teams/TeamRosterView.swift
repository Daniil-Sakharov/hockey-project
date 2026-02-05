import SwiftUI

struct TeamRosterView: View {
    let teamId: String
    let tournamentId: String
    let teamName: String

    @StateObject private var viewModel = TeamRosterViewModel()

    var body: some View {
        ZStack {
            AnimatedBackgroundView().ignoresSafeArea()
            content
        }
        .navigationTitle(teamName)
        .navigationBarTitleDisplayMode(.inline)
        .toolbarBackground(.hidden, for: .navigationBar)
        .task { await viewModel.load(teamId: teamId, tournamentId: tournamentId) }
    }

    @ViewBuilder
    private var content: some View {
        if viewModel.isLoading {
            ProgressView().tint(.srCyan)
        } else if let error = viewModel.errorMessage {
            errorView(error)
        } else {
            rosterContent
        }
    }

    private var rosterContent: some View {
        ScrollView {
            VStack(spacing: AppSpacing.md) {
                if let team = viewModel.team {
                    heroCard(team)
                }
                if !viewModel.players.isEmpty {
                    playersSection
                }
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.top, AppSpacing.sm)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
    }

    // MARK: - Hero

    private func heroCard(_ team: TeamInfoDTO) -> some View {
        VStack(spacing: AppSpacing.sm) {
            teamLogo(team.logoUrl)
            Text(team.name)
                .font(.srHeading3)
                .foregroundColor(.srTextPrimary)
                .multilineTextAlignment(.center)
            if let city = team.city, !city.isEmpty {
                Text(city)
                    .font(.srBody)
                    .foregroundColor(.srTextSecondary)
            }
            Text("\(viewModel.players.count) игроков")
                .font(.srCaption)
                .foregroundColor(.srTextMuted)
        }
        .frame(maxWidth: .infinity)
        .glassCard()
    }

    private func teamLogo(_ urlString: String?) -> some View {
        Group {
            if let urlString, let url = URL(string: urlString) {
                CachedAsyncImage(url: url) {
                    logoPlaceholder
                }
                .frame(width: 64, height: 64)
                .clipShape(RoundedRectangle(cornerRadius: 12))
            } else {
                logoPlaceholder
            }
        }
    }

    private var logoPlaceholder: some View {
        RoundedRectangle(cornerRadius: 12)
            .fill(Color.srCard)
            .frame(width: 64, height: 64)
            .overlay(
                Image(systemName: "shield.fill")
                    .font(.system(size: 24))
                    .foregroundColor(.srTextMuted)
            )
    }

    // MARK: - Players

    private var playersSection: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            Text("Состав")
                .font(.system(size: 11, weight: .semibold))
                .foregroundColor(.srTextMuted)
                .textCase(.uppercase)
                .tracking(1)

            VStack(spacing: 0) {
                ForEach(viewModel.players) { player in
                    NavigationLink(value: PlayerRoute(playerId: player.id)) {
                        rosterRow(player)
                    }
                    if player.id != viewModel.players.last?.id {
                        Divider().background(Color.srBorder.opacity(0.2))
                    }
                }
            }
            .glassCard(padding: 0)
        }
    }

    private func rosterRow(_ p: RosterPlayerDTO) -> some View {
        HStack(spacing: AppSpacing.sm) {
            playerPhoto(p.photoUrl)

            VStack(alignment: .leading, spacing: 2) {
                Text(p.name)
                    .font(.system(size: 14, weight: .medium))
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)

                HStack(spacing: AppSpacing.xs) {
                    if let pos = p.position {
                        Text(positionShort(pos))
                            .font(.system(size: 10, weight: .semibold))
                            .foregroundColor(.srPurple)
                    }
                    if let birth = p.birthDate, !birth.isEmpty {
                        Text(birth)
                            .font(.system(size: 10))
                            .foregroundColor(.srTextMuted)
                    }
                }
            }

            Spacer()

            Text("#\(p.jerseyNumber)")
                .font(.system(size: 14, weight: .bold, design: .monospaced))
                .foregroundColor(.srCyan)
        }
        .padding(.horizontal, AppSpacing.sm)
        .padding(.vertical, 8)
    }

    private func playerPhoto(_ urlString: String?) -> some View {
        Group {
            if let urlString, let url = URL(string: urlString) {
                CachedAsyncImage(url: url) {
                    Image(systemName: "person.circle.fill")
                        .resizable()
                        .foregroundColor(.srTextMuted)
                }
                .frame(width: 36, height: 36)
                .clipShape(Circle())
            } else {
                Image(systemName: "person.circle.fill")
                    .resizable()
                    .foregroundColor(.srTextMuted)
                    .frame(width: 36, height: 36)
            }
        }
    }

    private func positionShort(_ pos: String) -> String {
        switch pos {
        case "forward": return "НАП"
        case "defender": return "ЗАЩ"
        case "goalie": return "ВРТ"
        default: return pos.prefix(3).uppercased()
        }
    }

    // MARK: - Error

    private func errorView(_ message: String) -> some View {
        VStack(spacing: AppSpacing.sm) {
            Image(systemName: "exclamationmark.triangle")
                .font(.system(size: 36))
                .foregroundColor(.srAmber)
            Text(message)
                .font(.srBody)
                .foregroundColor(.srTextSecondary)
            Button("Повторить") {
                Task { await viewModel.load(teamId: teamId, tournamentId: tournamentId) }
            }
            .font(.srBodyMedium)
            .foregroundColor(.srCyan)
        }
    }
}
