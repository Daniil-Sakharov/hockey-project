import SwiftUI

struct GoalCelebrationView: View {
    @Binding var isShowing: Bool
    @State private var scale: CGFloat = 0.1

    var body: some View {
        if isShowing {
            ZStack {
                Color.black.opacity(0.8).ignoresSafeArea()

                VStack(spacing: AppSpacing.lg) {
                    Text("GOAL!")
                        .font(.system(size: 60, weight: .black))
                        .foregroundStyle(Color.srGradientPrimary)
                        .scaleEffect(scale)
                }
            }
            .onAppear {
                withAnimation(.spring(response: 0.5, dampingFraction: 0.6)) {
                    scale = 1.0
                }
                DispatchQueue.main.asyncAfter(deadline: .now() + 2.0) {
                    withAnimation {
                        isShowing = false
                    }
                }
            }
        }
    }
}
