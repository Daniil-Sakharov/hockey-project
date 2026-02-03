import SwiftUI

struct AchievementsView: View {
    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.lg) {
                Text("Достижения")
                    .font(.srHeading2)
                    .foregroundColor(.srTextPrimary)
            }
            .padding()
        }
        .background(Color.clear)
    }
}
