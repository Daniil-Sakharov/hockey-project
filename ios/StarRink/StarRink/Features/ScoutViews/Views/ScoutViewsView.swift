import SwiftUI

struct ScoutViewsView: View {
    var body: some View {
        LockedFeatureView(title: "Просмотры скаутов", requiredTier: .pro)
    }
}
