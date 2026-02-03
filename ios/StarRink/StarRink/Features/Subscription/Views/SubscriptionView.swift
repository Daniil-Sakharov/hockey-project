import SwiftUI

struct SubscriptionView: View {
    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.lg) {
                Text("Подписка")
                    .font(.srHeading2)
                    .foregroundColor(.srTextPrimary)
            }
            .padding()
        }
        .background(Color.clear)
    }
}
