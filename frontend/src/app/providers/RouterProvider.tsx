import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { NotFoundPage } from '@/pages/not-found'
import { DemoPage } from '@/pages/demo'
import { DashboardPage } from '@/pages/dashboard'
import { WatchlistPage } from '@/pages/dashboard/watchlist'
import { DashboardLayout } from '@/pages/dashboard/ui/DashboardLayout'

// Auth pages
import { LoginPage, RegisterPage, LinkPlayerPage, SelectRolePage } from '@/pages/auth'

// Player Dashboard
import { PlayerDashboardLayout, PlayerDashboardPage } from '@/pages/player-dashboard'
import { TeamCalendarPage } from '@/pages/player-dashboard/calendar'
import { AchievementsPage } from '@/pages/player-dashboard/achievements'

// Explore (Fan Dashboard)
import {
  ExploreLayout,
  ExploreDashboardPage,
  TournamentsListPage,
  RegionTournamentsPage,
  TournamentDetailPage,
  PlayersSearchPage,
  PlayerProfilePage,
  TeamProfilePage,
  MatchResultsPage,
  MatchCalendarPage,
  MatchDetailPage,
  RankingsPage,
  PredictionsPage,
} from '@/pages/explore'

// Route guards
import { ProtectedRoute, GuestOnlyRoute } from '@/shared/lib/hoc'

export function RouterProvider() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Demo page without header/footer */}
        <Route path="/demo" element={<DemoPage />} />

        {/* Auth pages (guest only) */}
        <Route
          path="/login"
          element={
            <GuestOnlyRoute>
              <LoginPage />
            </GuestOnlyRoute>
          }
        />
        <Route
          path="/register"
          element={
            <GuestOnlyRoute>
              <RegisterPage />
            </GuestOnlyRoute>
          }
        />

        {/* Role selection (after registration) */}
        <Route
          path="/select-role"
          element={
            <ProtectedRoute>
              <SelectRolePage />
            </ProtectedRoute>
          }
        />

        {/* Link player page (requires auth, but not linked player) */}
        <Route
          path="/link-player"
          element={
            <ProtectedRoute>
              <LinkPlayerPage />
            </ProtectedRoute>
          }
        />

        {/* Explore (Fan) Dashboard routes */}
        <Route
          path="/explore"
          element={
            <ProtectedRoute>
              <ExploreLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<ExploreDashboardPage />} />
          <Route path="tournaments" element={<TournamentsListPage />} />
          <Route path="tournaments/:region" element={<RegionTournamentsPage />} />
          <Route path="tournaments/detail/:id" element={<TournamentDetailPage />} />
          <Route path="players" element={<PlayersSearchPage />} />
          <Route path="players/:id" element={<PlayerProfilePage />} />
          <Route path="teams/:id" element={<TeamProfilePage />} />
          <Route path="results" element={<MatchResultsPage />} />
          <Route path="matches/:id" element={<MatchDetailPage />} />
          <Route path="calendar" element={<MatchCalendarPage />} />
          <Route path="rankings" element={<RankingsPage />} />
          <Route path="predictions" element={<PredictionsPage />} />
          <Route path="favorites" element={<ExplorePlaceholder title="Избранное" description="Отслеживание игроков (PRO)" />} />
          <Route path="settings" element={<ExplorePlaceholder title="Настройки" description="Настройки профиля" />} />
        </Route>

        {/* Player Dashboard routes */}
        <Route
          path="/player"
          element={
            <ProtectedRoute>
              <PlayerDashboardLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<PlayerDashboardPage />} />
          <Route path="stats" element={<PlayerStatsPlaceholder />} />
          <Route path="calendar" element={<TeamCalendarPage />} />
          <Route path="achievements" element={<AchievementsPage />} />
          <Route path="compare" element={<PlayerComparePlaceholder />} />
          <Route path="notifications" element={<PlayerNotificationsPlaceholder />} />
          <Route path="highlights" element={<PlayerHighlightsPlaceholder />} />
          <Route path="recommendations" element={<PlayerRecommendationsPlaceholder />} />
          <Route path="subscription" element={<PlayerSubscriptionPlaceholder />} />
          <Route path="settings" element={<PlayerSettingsPlaceholder />} />
        </Route>

        {/* Scout Dashboard pages with own layout */}
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/dashboard/watchlist" element={<WatchlistPage />} />
        <Route path="/dashboard/search" element={<SearchPagePlaceholder />} />
        <Route path="/dashboard/compare" element={<ComparePagePlaceholder />} />
        <Route path="/dashboard/notifications" element={<NotificationsPagePlaceholder />} />
        <Route path="/dashboard/settings" element={<SettingsPagePlaceholder />} />

        {/* Main landing page */}
        <Route path="/" element={<DemoPage />} />

        {/* 404 */}
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </BrowserRouter>
  )
}

// Explore placeholder
function ExplorePlaceholder({ title, description }: { title: string; description: string }) {
  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center text-center">
      <h1 className="text-gradient mb-4 text-3xl font-bold">{title}</h1>
      <p className="text-gray-400">{description}</p>
    </div>
  )
}

// Player Dashboard placeholder components
function PlayerStatsPlaceholder() {
  return <PlayerPlaceholderPage title="Статистика" description="Расширенная статистика (скоро)" />
}

function PlayerComparePlaceholder() {
  return <PlayerPlaceholderPage title="Сравнение" description="Сравнение с другими игроками (PRO)" />
}

function PlayerNotificationsPlaceholder() {
  return <PlayerPlaceholderPage title="Просмотры скаутов" description="Кто просматривал ваш профиль (PRO)" />
}

function PlayerHighlightsPlaceholder() {
  return <PlayerPlaceholderPage title="Видео хайлайты" description="Ваши лучшие моменты (ULTRA)" />
}

function PlayerRecommendationsPlaceholder() {
  return <PlayerPlaceholderPage title="AI рекомендации" description="Советы по улучшению игры (ULTRA)" />
}

function PlayerSubscriptionPlaceholder() {
  return <PlayerPlaceholderPage title="Подписка" description="Управление подпиской" />
}

function PlayerSettingsPlaceholder() {
  return <PlayerPlaceholderPage title="Настройки" description="Настройки профиля и уведомлений" />
}

function PlayerPlaceholderPage({ title, description }: { title: string; description: string }) {
  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center text-center">
      <h1 className="text-gradient mb-4 text-3xl font-bold">{title}</h1>
      <p className="text-gray-400">{description}</p>
    </div>
  )
}

// Scout Dashboard placeholder components
function SearchPagePlaceholder() {
  return <PlaceholderPage title="Поиск игроков" description="Страница расширенного поиска (скоро)" />
}

function ComparePagePlaceholder() {
  return <PlaceholderPage title="Сравнение игроков" description="Страница сравнения игроков (скоро)" />
}

function NotificationsPagePlaceholder() {
  return <PlaceholderPage title="Уведомления" description="Страница уведомлений (скоро)" />
}

function SettingsPagePlaceholder() {
  return <PlaceholderPage title="Настройки" description="Страница настроек (скоро)" />
}

function PlaceholderPage({ title, description }: { title: string; description: string }) {
  return (
    <DashboardLayout>
      <div className="flex min-h-[60vh] flex-col items-center justify-center text-center">
        <h1 className="text-gradient mb-4 text-3xl font-bold">{title}</h1>
        <p className="text-gray-400">{description}</p>
      </div>
    </DashboardLayout>
  )
}
