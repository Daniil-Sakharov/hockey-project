import SwiftUI

@main
struct StarRinkApp: App {
    @StateObject private var authViewModel = AuthViewModel()

    init() {
        // Make all UI backgrounds transparent to show aurora
        UINavigationBar.appearance().backgroundColor = .clear
        UIScrollView.appearance().backgroundColor = .clear
        UITableView.appearance().backgroundColor = .clear
        UICollectionView.appearance().backgroundColor = .clear

        // Configure navigation bar appearance
        let appearance = UINavigationBarAppearance()
        appearance.configureWithTransparentBackground()
        appearance.backgroundColor = .clear
        appearance.shadowColor = .clear
        UINavigationBar.appearance().standardAppearance = appearance
        UINavigationBar.appearance().scrollEdgeAppearance = appearance
        UINavigationBar.appearance().compactAppearance = appearance
    }

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(authViewModel)
                .preferredColorScheme(.dark)
        }
    }
}

struct ContentView: View {
    @EnvironmentObject var authViewModel: AuthViewModel
    @State private var showSplash = true

    var body: some View {
        ZStack {
            // Global animated background
            AnimatedBackgroundView()
                .ignoresSafeArea()

            Group {
                if authViewModel.isAuthenticated {
                    MainContainerView()
                } else {
                    LoginView()
                }
            }
            .animation(.easeInOut(duration: 0.3), value: authViewModel.isAuthenticated)

            if showSplash {
                SplashView {
                    withAnimation(.easeOut(duration: 0.5)) {
                        showSplash = false
                    }
                }
                .transition(.opacity)
                .zIndex(100)
            }
        }
    }
}
