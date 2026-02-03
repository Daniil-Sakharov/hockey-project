import SwiftUI

struct SettingsView: View {
    @EnvironmentObject var authViewModel: AuthViewModel
    @State private var showLogoutConfirmation = false

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.lg) {
                profileSection
                logoutButton
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.vertical, AppSpacing.md)
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
        .alert("Выход", isPresented: $showLogoutConfirmation) {
            Button("Отмена", role: .cancel) {}
            Button("Выйти", role: .destructive) {
                Task { await authViewModel.logout() }
            }
        } message: {
            Text("Вы уверены?")
        }
    }

    private var profileSection: some View {
        VStack(spacing: AppSpacing.md) {
            ZStack {
                Circle()
                    .fill(Color.srGradientPrimary)
                    .frame(width: 80, height: 80)
                Text(authViewModel.currentUser?.displayName.prefix(1).uppercased() ?? "?")
                    .font(.srHeading1)
                    .foregroundColor(.white)
            }
            Text(authViewModel.currentUser?.displayName ?? "Player")
                .font(.srHeading3)
                .foregroundColor(.srTextPrimary)
            Text(authViewModel.currentUser?.email ?? "")
                .font(.srBodySmall)
                .foregroundColor(.srTextSecondary)
        }
        .glassCard()
    }

    private var logoutButton: some View {
        Button { showLogoutConfirmation = true } label: {
            HStack {
                Image(systemName: "rectangle.portrait.and.arrow.right")
                Text("Выйти")
            }
            .font(.srButton)
            .foregroundColor(.srError)
            .frame(maxWidth: .infinity)
            .frame(height: 50)
            .background(Color.srCard)
            .clipShape(RoundedRectangle(cornerRadius: AppSpacing.radiusMedium))
        }
    }
}
