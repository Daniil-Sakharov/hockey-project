import SwiftUI

struct LockedFeatureView: View {
    let title: String
    let requiredTier: SubscriptionTier

    var body: some View {
        VStack(spacing: AppSpacing.lg) {
            Image(systemName: "lock.fill")
                .font(.system(size: 60))
                .foregroundColor(.srTextMuted)

            Text(title)
                .font(.srHeading3)
                .foregroundColor(.srTextPrimary)

            Text("Доступно в \(requiredTier.displayName)")
                .font(.srBody)
                .foregroundColor(.srTextSecondary)

            PrimaryButton("Улучшить подписку") {
                // Navigate to subscription
            }
            .padding(.horizontal, AppSpacing.xl)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(Color.srBackground)
    }
}
