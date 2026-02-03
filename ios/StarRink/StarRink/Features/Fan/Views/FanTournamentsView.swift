import SwiftUI

struct FanTournamentsView: View {
    var body: some View {
        RegionListView()
            .navigationDestination(for: Region.self) { region in
                RegionTournamentsView(region: region)
            }
            .navigationDestination(for: TournamentRoute.self) { route in
                TournamentDetailView(route: route)
            }
            .navigationDestination(for: PlayerRoute.self) { route in
                PlayerProfileView(playerId: route.playerId)
            }
    }
}
