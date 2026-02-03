import SwiftUI

struct AnimatedBackgroundView: View {
    var body: some View {
        ZStack {
            // Colorful base - NOT black
            AuroraGradientView()
                .ignoresSafeArea()

            // Animated waves on top
            AuroraWavesView()
                .ignoresSafeArea()

            // Subtle particles
            BackgroundParticlesView()
                .ignoresSafeArea()
        }
    }
}

// Static colorful gradient that covers entire screen
struct AuroraGradientView: View {
    var body: some View {
        ZStack {
            // Main gradient - covers full screen with color
            LinearGradient(
                colors: [
                    Color(red: 0.05, green: 0.15, blue: 0.25),  // Teal-blue top
                    Color(red: 0.08, green: 0.12, blue: 0.22),  // Blue-purple
                    Color(red: 0.10, green: 0.10, blue: 0.20),  // Purple center
                    Color(red: 0.08, green: 0.12, blue: 0.22),  // Blue-purple
                    Color(red: 0.06, green: 0.10, blue: 0.18)   // Dark teal bottom
                ],
                startPoint: .top,
                endPoint: .bottom
            )

            // Diagonal color overlay
            LinearGradient(
                colors: [
                    Color.srCyan.opacity(0.15),
                    Color.clear,
                    Color.srPurple.opacity(0.12)
                ],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )

            // Second diagonal for richness
            LinearGradient(
                colors: [
                    Color.srPurple.opacity(0.1),
                    Color.clear,
                    Color.srCyan.opacity(0.1)
                ],
                startPoint: .topTrailing,
                endPoint: .bottomLeading
            )
        }
    }
}

// Animated aurora waves
struct AuroraWavesView: View {
    var body: some View {
        TimelineView(.animation(minimumInterval: 1/30)) { timeline in
            Canvas { context, size in
                let time = timeline.date.timeIntervalSinceReferenceDate

                // Draw multiple overlapping waves across FULL screen
                // Top waves
                drawWave(context: context, size: size, time: time,
                        yPosition: 0.15, amplitude: 0.2,
                        color: .srCyan, opacity: 0.3)

                drawWave(context: context, size: size, time: time,
                        yPosition: 0.25, amplitude: 0.25,
                        color: .srPurple, opacity: 0.25)

                // Middle waves - IMPORTANT for center coverage
                drawWave(context: context, size: size, time: time,
                        yPosition: 0.4, amplitude: 0.3,
                        color: Color(red: 0.2, green: 0.4, blue: 0.8), opacity: 0.2)

                drawWave(context: context, size: size, time: time,
                        yPosition: 0.55, amplitude: 0.25,
                        color: .srCyan, opacity: 0.2)

                // Bottom waves
                drawWave(context: context, size: size, time: time,
                        yPosition: 0.7, amplitude: 0.3,
                        color: .srPurple, opacity: 0.25)

                drawWave(context: context, size: size, time: time,
                        yPosition: 0.85, amplitude: 0.2,
                        color: .srCyan, opacity: 0.3)
            }
        }
        .blendMode(.plusLighter)
    }

    private func drawWave(context: GraphicsContext, size: CGSize, time: TimeInterval,
                         yPosition: CGFloat, amplitude: CGFloat, color: Color, opacity: Double) {
        var path = Path()
        let baseY = size.height * yPosition
        let waveHeight = size.height * amplitude
        let timeVal = Double(time)

        // Start from left edge, full height
        path.move(to: CGPoint(x: 0, y: 0))

        // Draw wave top edge
        for x in stride(from: 0, through: size.width, by: 3) {
            let nx = Double(x / size.width)
            let wave1 = sin((nx * 2.5 + timeVal * 0.15 + Double(yPosition)) * .pi)
            let wave2 = sin((nx * 4.0 + timeVal * 0.1) * .pi) * 0.5
            let y = baseY + waveHeight * CGFloat(wave1 + wave2) * 0.5
            path.addLine(to: CGPoint(x: x, y: y))
        }

        // Close to bottom right
        path.addLine(to: CGPoint(x: size.width, y: size.height))
        path.addLine(to: CGPoint(x: 0, y: size.height))
        path.closeSubpath()

        // Fill with gradient from wave position
        let gradient = Gradient(colors: [
            color.opacity(opacity),
            color.opacity(opacity * 0.3),
            color.opacity(0)
        ])

        context.fill(path, with: .linearGradient(
            gradient,
            startPoint: CGPoint(x: size.width / 2, y: baseY - waveHeight),
            endPoint: CGPoint(x: size.width / 2, y: baseY + size.height * 0.4)
        ))
    }
}

struct BackgroundParticlesView: View {
    @State private var particles: [BGParticle] = (0..<15).map { _ in
        BGParticle(
            x: CGFloat.random(in: 0...1),
            y: CGFloat.random(in: 0...1),
            size: CGFloat.random(in: 2...8),
            opacity: Double.random(in: 0.1...0.3),
            speed: CGFloat.random(in: 0.02...0.05),
            phase: CGFloat.random(in: 0...(.pi * 2))
        )
    }

    var body: some View {
        TimelineView(.animation(minimumInterval: 1/30)) { timeline in
            Canvas { context, size in
                let time = CGFloat(timeline.date.timeIntervalSinceReferenceDate)
                for p in particles {
                    let x = (p.x * size.width + sin(time * p.speed + p.phase) * 30)
                        .truncatingRemainder(dividingBy: size.width)
                    let y = (p.y * size.height + cos(time * p.speed * 0.7 + p.phase) * 20)
                        .truncatingRemainder(dividingBy: size.height)

                    let rect = CGRect(x: x - p.size/2, y: y - p.size/2, width: p.size, height: p.size)
                    context.fill(
                        Path(ellipseIn: rect),
                        with: .color(Color.white.opacity(p.opacity))
                    )
                }
            }
        }
    }
}

struct BGParticle {
    let x, y, size: CGFloat
    let opacity: Double
    let speed, phase: CGFloat
}
