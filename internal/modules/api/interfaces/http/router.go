package http

import (
	"net/http"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/interfaces/http/handlers"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/interfaces/http/middleware"
)

// Router represents the HTTP router with all handlers.
type Router struct {
	mux                   *http.ServeMux
	healthHandler         *handlers.HealthHandler
	statsHandler          *handlers.StatsHandler
	rankingHandler        *handlers.RankingHandler
	authHandler           *handlers.AuthHandler
	exploreHandler        *handlers.ExploreHandler
	explorePlayersHandler *handlers.ExplorePlayersHandler
	exploreMatchesHandler *handlers.ExploreMatchesHandler
	imageProxyHandler     *handlers.ImageProxyHandler
	authMiddleware        *middleware.AuthMiddleware
	allowedOrigins        []string
}

// NewRouter creates a new HTTP router.
func NewRouter(
	healthHandler *handlers.HealthHandler,
	statsHandler *handlers.StatsHandler,
	rankingHandler *handlers.RankingHandler,
	authHandler *handlers.AuthHandler,
	exploreHandler *handlers.ExploreHandler,
	explorePlayersHandler *handlers.ExplorePlayersHandler,
	exploreMatchesHandler *handlers.ExploreMatchesHandler,
	imageProxyHandler *handlers.ImageProxyHandler,
	authMiddleware *middleware.AuthMiddleware,
	allowedOrigins []string,
) *Router {
	return &Router{
		mux:                   http.NewServeMux(),
		healthHandler:         healthHandler,
		statsHandler:          statsHandler,
		rankingHandler:        rankingHandler,
		authHandler:           authHandler,
		exploreHandler:        exploreHandler,
		explorePlayersHandler: explorePlayersHandler,
		exploreMatchesHandler: exploreMatchesHandler,
		imageProxyHandler:     imageProxyHandler,
		authMiddleware:        authMiddleware,
		allowedOrigins:        allowedOrigins,
	}
}

// Setup registers all routes.
func (r *Router) Setup() http.Handler {
	// Health check (public)
	r.mux.HandleFunc("GET /api/v1/health", r.healthHandler.Health)

	// Auth routes (public)
	r.mux.HandleFunc("POST /api/v1/auth/register", r.authHandler.Register)
	r.mux.HandleFunc("POST /api/v1/auth/login", r.authHandler.Login)
	r.mux.HandleFunc("POST /api/v1/auth/refresh", r.authHandler.Refresh)

	// Auth routes (protected)
	r.mux.Handle("GET /api/v1/auth/me", r.authMiddleware.RequireAuth(http.HandlerFunc(r.authHandler.Me)))
	r.mux.Handle("POST /api/v1/auth/link-player", r.authMiddleware.RequireAuth(http.HandlerFunc(r.authHandler.LinkPlayer)))
	r.mux.Handle("POST /api/v1/auth/logout", r.authMiddleware.RequireAuth(http.HandlerFunc(r.authHandler.Logout)))

	// Stats routes (public)
	r.mux.HandleFunc("GET /api/v1/stats/overview", r.statsHandler.Overview)

	// Rankings routes (public)
	r.mux.HandleFunc("GET /api/v1/rankings/scorers", r.rankingHandler.TopScorers)

	// Explore routes (public)
	r.mux.HandleFunc("GET /api/v1/explore/overview", r.exploreHandler.Overview)
	r.mux.HandleFunc("GET /api/v1/explore/seasons", r.exploreHandler.Seasons)
	r.mux.HandleFunc("GET /api/v1/explore/tournaments", r.exploreHandler.Tournaments)
	r.mux.HandleFunc("GET /api/v1/explore/tournaments/{id}/standings", r.exploreHandler.Standings)
	r.mux.HandleFunc("GET /api/v1/explore/tournaments/{id}/matches", r.exploreHandler.TournamentMatches)
	r.mux.HandleFunc("GET /api/v1/explore/tournaments/{id}/scorers", r.exploreHandler.Scorers)
	r.mux.HandleFunc("GET /api/v1/explore/tournaments/{id}/teams", r.exploreHandler.TournamentTeams)
	r.mux.HandleFunc("GET /api/v1/explore/players/{id}/stats", r.explorePlayersHandler.PlayerStats)
	r.mux.HandleFunc("GET /api/v1/explore/players/{id}", r.explorePlayersHandler.PlayerProfile)
	r.mux.HandleFunc("GET /api/v1/explore/players", r.explorePlayersHandler.SearchPlayers)
	r.mux.HandleFunc("GET /api/v1/explore/teams/{teamId}/roster/{tournamentId}", r.exploreHandler.TeamRoster)
	r.mux.HandleFunc("GET /api/v1/explore/teams/{id}", r.explorePlayersHandler.TeamProfile)
	r.mux.HandleFunc("GET /api/v1/explore/results", r.exploreMatchesHandler.RecentResults)
	r.mux.HandleFunc("GET /api/v1/explore/calendar", r.exploreMatchesHandler.UpcomingMatches)
	r.mux.HandleFunc("GET /api/v1/explore/rankings", r.exploreMatchesHandler.Rankings)
	r.mux.HandleFunc("GET /api/v1/explore/rankings/filters", r.exploreMatchesHandler.RankingsFilters)
	r.mux.HandleFunc("GET /api/v1/explore/matches/{id}", r.exploreMatchesHandler.MatchDetail)

	// Image proxy (public)
	r.mux.HandleFunc("GET /api/v1/proxy/image", r.imageProxyHandler.ProxyImage)

	// Apply middleware chain
	handler := r.applyMiddleware(r.mux)

	return handler
}

func (r *Router) applyMiddleware(handler http.Handler) http.Handler {
	// Apply in reverse order (last applied = first executed)
	handler = middleware.Logging()(handler)
	handler = middleware.CORS(r.allowedOrigins)(handler)

	return handler
}
