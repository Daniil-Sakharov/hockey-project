import SwiftUI

struct ProfileView: View {
    @EnvironmentObject var authViewModel: AuthViewModel

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.lg) {
                profileHeader
                quickStats
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.vertical, AppSpacing.md)
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
    }

    private var profileHeader: some View {
        VStack(spacing: AppSpacing.md) {
            ZStack {
                Circle()
                    .fill(Color.srGradientPrimary)
                    .frame(width: 100, height: 100)
                Text(authViewModel.currentUser?.displayName.prefix(1).uppercased() ?? "?")
                    .font(.system(size: 40, weight: .bold))
                    .foregroundColor(.white)
            }
            Text(authViewModel.currentUser?.displayName ?? "Игрок")
                .font(.srHeading2)
                .foregroundColor(.srTextPrimary)
        }
        .glassCard()
    }

    private var quickStats: some View {
        VStack(alignment: .leading, spacing: AppSpacing.md) {
            Text("Статистика сезона")
                .font(.srHeading4)
                .foregroundColor(.srTextPrimary)
            HStack(spacing: AppSpacing.md) {
                ProfileStatCard(value: "12", label: "Игры")
                ProfileStatCard(value: "8", label: "Голы")
                ProfileStatCard(value: "15", label: "Передачи")
            }
        }
    }
}

struct ProfileStatCard: View {
    let value: String
    let label: String

    var body: some View {
        VStack(spacing: AppSpacing.xs) {
            Text(value)
                .font(.srHeading3)
                .foregroundColor(.srTextPrimary)
            Text(label)
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, AppSpacing.md)
        .background(Color.srCard)
        .clipShape(RoundedRectangle(cornerRadius: AppSpacing.radiusSmall))
    }
}
