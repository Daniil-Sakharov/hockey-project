import Foundation

@MainActor
final class HomeViewModel: ObservableObject {
    @Published var isLoading = false
}
