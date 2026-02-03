import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { UserProfile, UserRole, SubscriptionTier, FeatureKey } from '@/entities/user'
import { hasFeatureAccess } from '@/entities/user'

interface AuthStore {
  // State
  isAuthenticated: boolean
  user: UserProfile | null
  token: string | null
  isLoading: boolean
  error: string | null

  // Auth Actions
  login: (email: string, password: string) => Promise<boolean>
  register: (email: string, password: string) => Promise<boolean>
  logout: () => void
  clearError: () => void

  // Role
  updateRole: (role: UserRole) => void

  // Player Link
  linkPlayer: (playerId: string, fullName: string, birthDate: string) => Promise<boolean>
  unlinkPlayer: () => void

  // Profile
  updateProfile: (data: Partial<UserProfile>) => void
  updateSubscription: (tier: SubscriptionTier) => void

  // Feature Access
  hasFeature: (feature: FeatureKey) => boolean
  getSubscriptionTier: () => SubscriptionTier
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      // Initial State
      isAuthenticated: false,
      user: null,
      token: null,
      isLoading: false,
      error: null,

      // Login (mock implementation)
      login: async (email: string, password: string) => {
        set({ isLoading: true, error: null })

        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 1000))

        // Mock validation
        if (!email || !password) {
          set({ isLoading: false, error: 'Введите email и пароль' })
          return false
        }

        // Mock success - создаём пользователя
        const mockUser: UserProfile = {
          id: 'user-' + Date.now(),
          email,
          role: 'fan',
          subscription: {
            tier: 'free',
            startDate: new Date().toISOString(),
            endDate: null,
            autoRenew: false,
            price: 0,
          },
          createdAt: new Date().toISOString(),
        }

        set({
          isAuthenticated: true,
          user: mockUser,
          token: 'mock-jwt-token-' + Date.now(),
          isLoading: false,
          error: null,
        })

        return true
      },

      // Register (mock implementation)
      register: async (email: string, password: string) => {
        set({ isLoading: true, error: null })

        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 1000))

        // Mock validation
        if (!email || !password) {
          set({ isLoading: false, error: 'Заполните все поля' })
          return false
        }

        if (password.length < 6) {
          set({ isLoading: false, error: 'Пароль должен быть не менее 6 символов' })
          return false
        }

        // Mock success
        const mockUser: UserProfile = {
          id: 'user-' + Date.now(),
          email,
          role: 'fan',
          subscription: {
            tier: 'free',
            startDate: new Date().toISOString(),
            endDate: null,
            autoRenew: false,
            price: 0,
          },
          createdAt: new Date().toISOString(),
        }

        set({
          isAuthenticated: true,
          user: mockUser,
          token: 'mock-jwt-token-' + Date.now(),
          isLoading: false,
          error: null,
        })

        return true
      },

      // Logout
      logout: () => {
        set({
          isAuthenticated: false,
          user: null,
          token: null,
          error: null,
        })
      },

      // Clear error
      clearError: () => {
        set({ error: null })
      },

      // Update role
      updateRole: (role: UserRole) => {
        const { user } = get()
        if (user) {
          set({ user: { ...user, role } })
        }
      },

      // Link player to account (mock)
      linkPlayer: async (playerId: string, fullName: string, birthDate: string) => {
        set({ isLoading: true, error: null })

        // Simulate API call - поиск игрока по ФИО + дата рождения
        await new Promise((resolve) => setTimeout(resolve, 1500))

        const { user } = get()
        if (!user) {
          set({ isLoading: false, error: 'Необходимо авторизоваться' })
          return false
        }

        // Mock: проверяем что "нашли" игрока
        // В реальности здесь будет запрос к API
        if (!fullName || !birthDate) {
          set({ isLoading: false, error: 'Игрок не найден. Проверьте ФИО и дату рождения' })
          return false
        }

        // Успех - привязываем игрока
        set({
          user: {
            ...user,
            linkedPlayerId: playerId,
          },
          isLoading: false,
          error: null,
        })

        return true
      },

      // Unlink player
      unlinkPlayer: () => {
        const { user } = get()
        if (user) {
          set({
            user: {
              ...user,
              linkedPlayerId: undefined,
            },
          })
        }
      },

      // Update profile
      updateProfile: (data: Partial<UserProfile>) => {
        const { user } = get()
        if (user) {
          set({
            user: {
              ...user,
              ...data,
            },
          })
        }
      },

      // Update subscription (mock - для демо переключения уровней)
      updateSubscription: (tier: SubscriptionTier) => {
        const { user } = get()
        if (user) {
          const prices: Record<SubscriptionTier, number> = {
            free: 0,
            pro: 990,
            ultra: 2490,
          }

          set({
            user: {
              ...user,
              subscription: {
                ...user.subscription,
                tier,
                price: prices[tier],
                startDate: new Date().toISOString(),
                endDate: tier === 'free' ? null : new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(),
              },
            },
          })
        }
      },

      // Check feature access
      hasFeature: (feature: FeatureKey) => {
        const { user } = get()
        if (!user) return false
        return hasFeatureAccess(user.subscription.tier, feature)
      },

      // Get current subscription tier
      getSubscriptionTier: () => {
        const { user } = get()
        return user?.subscription.tier || 'free'
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        isAuthenticated: state.isAuthenticated,
        user: state.user,
        token: state.token,
      }),
    }
  )
)
