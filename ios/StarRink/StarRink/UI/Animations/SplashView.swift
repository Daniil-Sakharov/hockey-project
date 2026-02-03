import SwiftUI

struct SplashView: View {
    let onComplete: () -> Void
    @State private var scale: CGFloat = 0.5
    @State private var opacity: Double = 0
    @State private var rotation: Double = 0

    var body: some View {
        ZStack {
            Color.srBackground.ignoresSafeArea()

            VStack(spacing: AppSpacing.md) {
                Image(systemName: "star.fill")
                    .font(.system(size: 80))
                    .foregroundStyle(Color.srGradientPrimary)
                    .scaleEffect(scale)
                    .rotationEffect(.degrees(rotation))

                Text("StarRink")
                    .font(.system(size: 36, weight: .bold))
                    .foregroundStyle(Color.srGradientPrimary)
                    .opacity(opacity)
            }
        }
        .onAppear {
            withAnimation(.spring(response: 0.6, dampingFraction: 0.6)) {
                scale = 1.0
                rotation = 360
            }
            withAnimation(.easeIn(duration: 0.4).delay(0.3)) {
                opacity = 1.0
            }
            DispatchQueue.main.asyncAfter(deadline: .now() + 2.0) {
                onComplete()
            }
        }
    }
}
