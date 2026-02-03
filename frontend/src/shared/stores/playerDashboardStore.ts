import { create } from 'zustand'
import type {
  PlayerProfileExtended,
  Achievement,
  TeamMatch,
  ScoutNotification,
} from '@/entities/player'

interface PlayerDashboardStore {
  // State
  linkedPlayer: PlayerProfileExtended | null
  teamMatches: TeamMatch[]
  achievements: Achievement[]
  scoutNotifications: ScoutNotification[]
  isLoading: boolean
  error: string | null

  // Actions
  setLinkedPlayer: (player: PlayerProfileExtended | null) => void
  fetchPlayerData: (playerId: string) => Promise<void>
  fetchTeamCalendar: (teamId: string) => Promise<void>
  fetchAchievements: (playerId: string) => Promise<void>
  fetchScoutNotifications: (playerId: string) => Promise<void>

  // Notifications
  markNotificationRead: (notificationId: string) => void
  markAllNotificationsRead: () => void
  getUnreadNotificationsCount: () => number

  // Settings
  profileVisibility: 'public' | 'scouts_only' | 'private'
  setProfileVisibility: (visibility: 'public' | 'scouts_only' | 'private') => void

  // Clear
  clearAll: () => void
}

export const usePlayerDashboardStore = create<PlayerDashboardStore>((set, get) => ({
  // Initial State
  linkedPlayer: null,
  teamMatches: [],
  achievements: [],
  scoutNotifications: [],
  isLoading: false,
  error: null,
  profileVisibility: 'scouts_only',

  // Set linked player
  setLinkedPlayer: (player) => {
    set({ linkedPlayer: player })
  },

  // Fetch player data (mock)
  fetchPlayerData: async (_playerId: string) => {
    set({ isLoading: true, error: null })

    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 800))

    // Mock data будет загружаться из mocks/
    // В реальности здесь будет API запрос
    set({ isLoading: false })
  },

  // Fetch team calendar (mock)
  fetchTeamCalendar: async (_teamId: string) => {
    set({ isLoading: true })
    await new Promise((resolve) => setTimeout(resolve, 500))
    // Mock data загружается отдельно
    set({ isLoading: false })
  },

  // Fetch achievements (mock)
  fetchAchievements: async (_playerId: string) => {
    set({ isLoading: true })
    await new Promise((resolve) => setTimeout(resolve, 500))
    set({ isLoading: false })
  },

  // Fetch scout notifications (mock)
  fetchScoutNotifications: async (_playerId: string) => {
    set({ isLoading: true })
    await new Promise((resolve) => setTimeout(resolve, 500))
    set({ isLoading: false })
  },

  // Mark notification as read
  markNotificationRead: (notificationId: string) => {
    set((state) => ({
      scoutNotifications: state.scoutNotifications.map((n) =>
        n.id === notificationId ? { ...n, isRead: true } : n
      ),
    }))
  },

  // Mark all notifications as read
  markAllNotificationsRead: () => {
    set((state) => ({
      scoutNotifications: state.scoutNotifications.map((n) => ({
        ...n,
        isRead: true,
      })),
    }))
  },

  // Get unread notifications count
  getUnreadNotificationsCount: () => {
    return get().scoutNotifications.filter((n) => !n.isRead).length
  },

  // Set profile visibility
  setProfileVisibility: (visibility) => {
    set({ profileVisibility: visibility })
  },

  // Clear all data (on logout)
  clearAll: () => {
    set({
      linkedPlayer: null,
      teamMatches: [],
      achievements: [],
      scoutNotifications: [],
      isLoading: false,
      error: null,
    })
  },
}))
