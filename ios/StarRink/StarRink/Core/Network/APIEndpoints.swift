import Foundation

enum APIEndpoint {
    #if DEBUG
    static let baseURL = "http://localhost:8080"
    #else
    static let baseURL = "https://api.starrink.ru"
    #endif

    case register
    case login
    case refresh
    case me
    case linkPlayer
    case logout
    case statsOverview
    case topScorers
    case player(id: String)
    case playerStats(id: String)
    case searchPlayers(query: String)
    case health
    case exploreTournaments
    case tournamentStandings(id: String)
    case tournamentMatches(id: String)
    case tournamentScorers(id: String)
    case explorePlayers
    case playerProfile(id: String)
    case playerStatsHistory(id: String)
    case seasons
    case exploreOverview
    case exploreRankings
    case exploreResults
    case exploreCalendar
    case teamProfile(id: String)

    var path: String {
        switch self {
        case .register: return "/api/v1/auth/register"
        case .login: return "/api/v1/auth/login"
        case .refresh: return "/api/v1/auth/refresh"
        case .me: return "/api/v1/auth/me"
        case .linkPlayer: return "/api/v1/auth/link-player"
        case .logout: return "/api/v1/auth/logout"
        case .statsOverview: return "/api/v1/stats/overview"
        case .topScorers: return "/api/v1/rankings/scorers"
        case .player(let id): return "/api/v1/players/\(id)"
        case .playerStats(let id): return "/api/v1/players/\(id)/stats"
        case .searchPlayers: return "/api/v1/players"
        case .health: return "/api/v1/health"
        case .exploreTournaments: return "/api/v1/explore/tournaments"
        case .tournamentStandings(let id): return "/api/v1/explore/tournaments/\(id)/standings"
        case .tournamentMatches(let id): return "/api/v1/explore/tournaments/\(id)/matches"
        case .tournamentScorers(let id): return "/api/v1/explore/tournaments/\(id)/scorers"
        case .explorePlayers: return "/api/v1/explore/players"
        case .playerProfile(let id): return "/api/v1/explore/players/\(id)"
        case .playerStatsHistory(let id): return "/api/v1/explore/players/\(id)/stats"
        case .seasons: return "/api/v1/explore/seasons"
        case .exploreOverview: return "/api/v1/explore/overview"
        case .exploreRankings: return "/api/v1/explore/rankings"
        case .exploreResults: return "/api/v1/explore/results"
        case .exploreCalendar: return "/api/v1/explore/calendar"
        case .teamProfile(let id): return "/api/v1/explore/teams/\(id)"
        }
    }

    var method: HTTPMethod {
        switch self {
        case .register, .login, .refresh, .linkPlayer, .logout:
            return .post
        default:
            return .get
        }
    }

    var requiresAuth: Bool {
        switch self {
        case .me, .linkPlayer, .logout:
            return true
        default:
            return false
        }
    }
}

enum HTTPMethod: String {
    case get = "GET"
    case post = "POST"
    case put = "PUT"
    case patch = "PATCH"
    case delete = "DELETE"
}
