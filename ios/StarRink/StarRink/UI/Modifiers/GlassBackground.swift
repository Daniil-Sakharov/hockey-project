import SwiftUI

struct GlassBackground: ViewModifier {
    var padding: CGFloat = AppSpacing.md
    var cornerRadius: CGFloat = AppSpacing.radiusMedium
    var material: Material = .ultraThinMaterial
    var borderOpacity: Double = 0.3

    func body(content: Content) -> some View {
        content
            .padding(padding)
            .background(
                RoundedRectangle(cornerRadius: cornerRadius)
                    .fill(material)
                    .overlay(
                        RoundedRectangle(cornerRadius: cornerRadius)
                            .fill(
                                LinearGradient(
                                    colors: [
                                        Color.white.opacity(0.1),
                                        Color.clear
                                    ],
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                    )
            )
            .overlay(
                RoundedRectangle(cornerRadius: cornerRadius)
                    .stroke(
                        LinearGradient(
                            colors: [
                                Color.white.opacity(borderOpacity),
                                Color.srBorder.opacity(borderOpacity * 0.5)
                            ],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        ),
                        lineWidth: 0.5
                    )
            )
    }
}

extension View {
    func glassCard(padding: CGFloat = AppSpacing.md) -> some View {
        modifier(GlassBackground(padding: padding))
    }

    func glassCard(
        padding: CGFloat = AppSpacing.md,
        cornerRadius: CGFloat = AppSpacing.radiusMedium,
        material: Material = .ultraThinMaterial
    ) -> some View {
        modifier(GlassBackground(
            padding: padding,
            cornerRadius: cornerRadius,
            material: material
        ))
    }

    func thickGlassCard(padding: CGFloat = AppSpacing.md) -> some View {
        modifier(GlassBackground(
            padding: padding,
            material: .thinMaterial,
            borderOpacity: 0.4
        ))
    }
}
