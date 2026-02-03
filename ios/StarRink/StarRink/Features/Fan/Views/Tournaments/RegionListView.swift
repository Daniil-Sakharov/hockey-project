import SwiftUI

struct RegionListView: View {
    let regions: [Region] = Region.allRegions

    private let columns = [
        GridItem(.flexible(), spacing: AppSpacing.sm),
        GridItem(.flexible(), spacing: AppSpacing.sm),
    ]

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: AppSpacing.lg) {
                Text("Выберите регион")
                    .font(.srHeading3)
                    .foregroundColor(.srTextPrimary)

                LazyVGrid(columns: columns, spacing: AppSpacing.sm) {
                    ForEach(regions) { region in
                        NavigationLink(value: region) {
                            RegionCard(region: region)
                        }
                        .buttonStyle(.plain)
                    }
                }
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.top, AppSpacing.md)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
    }
}
