import SwiftUI

struct HorizontalChipPicker: View {
    let items: [String]
    @Binding var selectedItem: String

    var body: some View {
        ScrollViewReader { proxy in
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: AppSpacing.xs) {
                    ForEach(items, id: \.self) { item in
                        chipButton(item)
                            .id(item)
                    }
                }
                .padding(.horizontal, AppSpacing.screenHorizontal)
                .padding(.vertical, AppSpacing.xxs)
            }
            .onChange(of: selectedItem) { _, newValue in
                withAnimation(.spring(response: 0.3, dampingFraction: 0.8)) {
                    proxy.scrollTo(newValue, anchor: .center)
                }
            }
        }
    }

    private func chipButton(_ item: String) -> some View {
        let isSelected = item == selectedItem

        return Button {
            withAnimation(.spring(response: 0.3, dampingFraction: 0.8)) {
                selectedItem = item
            }
        } label: {
            Text(item)
                .font(.system(size: 13, weight: isSelected ? .bold : .medium))
                .foregroundColor(isSelected ? .white : .srTextSecondary)
                .padding(.horizontal, AppSpacing.sm)
                .padding(.vertical, AppSpacing.xs)
                .background(
                    Capsule()
                        .fill(isSelected ? Color.srCyan.opacity(0.25) : Color.clear)
                )
                .overlay(
                    Capsule()
                        .stroke(
                            isSelected ? Color.srCyan.opacity(0.6) : Color.srBorder.opacity(0.3),
                            lineWidth: isSelected ? 1.0 : 0.5
                        )
                )
                .shadow(
                    color: isSelected ? Color.srCyan.opacity(0.3) : .clear,
                    radius: 8, x: 0, y: 0
                )
        }
        .buttonStyle(.plain)
    }
}
