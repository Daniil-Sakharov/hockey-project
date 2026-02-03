import SwiftUI

struct CachedAsyncImage: View {
    let url: URL?
    let placeholder: AnyView

    @StateObject private var loader = ImageLoader()

    init(url: URL?, @ViewBuilder placeholder: () -> some View) {
        self.url = url
        self.placeholder = AnyView(placeholder())
    }

    var body: some View {
        Group {
            if let image = loader.image {
                Image(uiImage: image)
                    .resizable()
                    .scaledToFill()
            } else if loader.isFailed {
                placeholder
            } else if loader.isLoading {
                ProgressView().tint(.srCyan)
            } else {
                placeholder
            }
        }
        .onAppear { loader.load(url: url) }
        .onChange(of: url) { _, newURL in loader.load(url: newURL) }
    }
}

@MainActor
private final class ImageLoader: ObservableObject {
    @Published var image: UIImage?
    @Published var isLoading = false
    @Published var isFailed = false

    private static let cache = NSCache<NSURL, UIImage>()
    private static let session: URLSession = {
        let config = URLSessionConfiguration.default
        config.timeoutIntervalForRequest = 15
        config.timeoutIntervalForResource = 30
        return URLSession(configuration: config)
    }()

    private var currentURL: URL?

    func load(url: URL?) {
        guard let url else {
            isFailed = true
            return
        }
        guard url != currentURL else { return }
        currentURL = url

        if let cached = Self.cache.object(forKey: url as NSURL) {
            image = cached
            return
        }

        isLoading = true
        isFailed = false
        image = nil

        let proxyURL = Self.proxyURL(for: url)

        Task {
            do {
                let (data, response) = try await Self.session.data(from: proxyURL)
                guard let httpResponse = response as? HTTPURLResponse,
                      httpResponse.statusCode == 200,
                      let uiImage = UIImage(data: data) else {
                    if currentURL == url {
                        isFailed = true
                        isLoading = false
                    }
                    return
                }
                Self.cache.setObject(uiImage, forKey: url as NSURL)
                if currentURL == url {
                    image = uiImage
                    isLoading = false
                }
            } catch {
                if currentURL == url {
                    isFailed = true
                    isLoading = false
                }
            }
        }
    }

    private static func proxyURL(for originalURL: URL) -> URL {
        let encoded = originalURL.absoluteString.addingPercentEncoding(
            withAllowedCharacters: .urlQueryAllowed
        ) ?? originalURL.absoluteString
        let proxy = "\(APIEndpoint.baseURL)/api/v1/proxy/image?url=\(encoded)"
        return URL(string: proxy) ?? originalURL
    }
}
