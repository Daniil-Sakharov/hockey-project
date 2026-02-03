import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface ScoutProfile {
  name: string
  role: 'scout' | 'coach' | 'analyst'
  club?: string
  avatarUrl?: string
}

interface ScoutStore {
  // Профиль скаута
  profile: ScoutProfile

  // Избранные игроки (watchlist)
  watchlist: string[]
  addToWatchlist: (playerId: string) => void
  removeFromWatchlist: (playerId: string) => void
  isInWatchlist: (playerId: string) => boolean
  clearWatchlist: () => void

  // Недавно просмотренные
  recentlyViewed: string[]
  addToRecentlyViewed: (playerId: string) => void
  clearRecentlyViewed: () => void

  // Выбранные для сравнения
  compareList: string[]
  addToCompare: (playerId: string) => void
  removeFromCompare: (playerId: string) => void
  clearCompareList: () => void

  // Обновление профиля
  updateProfile: (profile: Partial<ScoutProfile>) => void
}

const MAX_RECENTLY_VIEWED = 10
const MAX_COMPARE_LIST = 3

export const useScoutStore = create<ScoutStore>()(
  persist(
    (set, get) => ({
      // Дефолтный профиль скаута
      profile: {
        name: 'Скаут',
        role: 'scout',
        club: undefined,
      },

      // Watchlist
      watchlist: [],

      addToWatchlist: (playerId: string) => {
        const { watchlist } = get()
        if (!watchlist.includes(playerId)) {
          set({ watchlist: [...watchlist, playerId] })
        }
      },

      removeFromWatchlist: (playerId: string) => {
        set({ watchlist: get().watchlist.filter((id) => id !== playerId) })
      },

      isInWatchlist: (playerId: string) => {
        return get().watchlist.includes(playerId)
      },

      clearWatchlist: () => {
        set({ watchlist: [] })
      },

      // Recently Viewed
      recentlyViewed: [],

      addToRecentlyViewed: (playerId: string) => {
        const { recentlyViewed } = get()
        // Удаляем если уже есть (чтобы переместить в начало)
        const filtered = recentlyViewed.filter((id) => id !== playerId)
        // Добавляем в начало и ограничиваем размер
        const updated = [playerId, ...filtered].slice(0, MAX_RECENTLY_VIEWED)
        set({ recentlyViewed: updated })
      },

      clearRecentlyViewed: () => {
        set({ recentlyViewed: [] })
      },

      // Compare List
      compareList: [],

      addToCompare: (playerId: string) => {
        const { compareList } = get()
        if (!compareList.includes(playerId) && compareList.length < MAX_COMPARE_LIST) {
          set({ compareList: [...compareList, playerId] })
        }
      },

      removeFromCompare: (playerId: string) => {
        set({ compareList: get().compareList.filter((id) => id !== playerId) })
      },

      clearCompareList: () => {
        set({ compareList: [] })
      },

      // Profile
      updateProfile: (profile: Partial<ScoutProfile>) => {
        set({ profile: { ...get().profile, ...profile } })
      },
    }),
    {
      name: 'scout-storage',
      partialize: (state) => ({
        profile: state.profile,
        watchlist: state.watchlist,
        recentlyViewed: state.recentlyViewed,
      }),
    }
  )
)
