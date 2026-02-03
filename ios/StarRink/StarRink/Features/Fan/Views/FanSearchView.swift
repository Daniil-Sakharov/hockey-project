import SwiftUI

struct FanSearchView: View {
    var body: some View {
        PlayersSearchView()
            .navigationDestination(for: PlayerRoute.self) { route in
                PlayerProfileView(playerId: route.playerId)
            }
    }
}
