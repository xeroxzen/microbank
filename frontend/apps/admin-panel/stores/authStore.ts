import { create } from "zustand";
import { persist } from "zustand/middleware";

export interface AdminUser {
  id: string;
  email: string;
  name: string;
  is_blacklisted: boolean;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

interface AuthState {
  user: AdminUser | null;
  accessToken: string | null;
  refreshToken: string | null;
  isLoading: boolean;
  error: string | null;
  hasCheckedAuth: boolean;
}

interface AuthActions {
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  checkAuth: () => Promise<void>;
  refreshAccessToken: () => Promise<void>;
  clearError: () => void;
}

type AuthStore = AuthState & AuthActions;

const CLIENT_SERVICE_URL =
  process.env.NEXT_PUBLIC_CLIENT_SERVICE_URL || "http://localhost:8081";

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      // State
      user: null,
      accessToken: null,
      refreshToken: null,
      isLoading: false,
      error: null,
      hasCheckedAuth: false,

      // Actions
      login: async (email: string, password: string) => {
        set({ isLoading: true, error: null });

        try {
          const response = await fetch(
            `${CLIENT_SERVICE_URL}/api/v1/auth/login`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({ email, password }),
            }
          );

          const data = await response.json();

          if (!response.ok) {
            throw new Error(data.error?.message || "Login failed");
          }

          // Verify the user is an admin
          if (!data.user.is_admin) {
            throw new Error("Access denied. Admin privileges required.");
          }

          set({
            user: data.user,
            accessToken: data.tokens.access_token,
            refreshToken: data.tokens.refresh_token,
            isLoading: false,
            error: null,
            hasCheckedAuth: true,
          });
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : "Login failed",
            isLoading: false,
          });
          throw error;
        }
      },

      logout: () => {
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          error: null,
          hasCheckedAuth: false,
        });
      },

      checkAuth: async () => {
        // Prevent multiple simultaneous calls or if already checked
        if (get().isLoading || get().hasCheckedAuth) {
          return;
        }

        set({ isLoading: true });

        const { accessToken, refreshToken } = get();

        // If user is already authenticated, skip
        if (get().user && accessToken) {
          set({ isLoading: false });
          return;
        }

        // If no tokens at all, we can't proceed
        if (!accessToken && !refreshToken) {
          set({ isLoading: false });
          return;
        }

        // If no access token but we have a refresh token, try to refresh first
        if (!accessToken && refreshToken) {
          try {
            await get().refreshAccessToken();
            const newState = get();
            if (!newState.accessToken) {
              set({ isLoading: false });
              return;
            }
          } catch (error) {
            set({ isLoading: false });
            return;
          }
        }

        // If still no access token, we can't proceed
        if (!get().accessToken) {
          set({ isLoading: false });
          return;
        }

        set({ isLoading: true });

        try {
          const response = await fetch(
            `${CLIENT_SERVICE_URL}/api/v1/auth/validate`,
            {
              headers: {
                Authorization: `Bearer ${get().accessToken}`,
              },
            }
          );

          if (response.ok) {
            const data = await response.json();

            // Verify the user is still an admin
            if (!data.user.is_admin) {
              get().logout();
              return;
            }

            set({ user: data.user, isLoading: false, hasCheckedAuth: true });
          } else {
            // Token is invalid, try to refresh
            await get().refreshAccessToken();
          }
        } catch (error) {
          console.error("Token validation error:", error);
          set({ isLoading: false });
          get().logout();
        }
      },

      refreshAccessToken: async () => {
        const { refreshToken } = get();

        if (!refreshToken) {
          get().logout();
          return;
        }

        try {
          const response = await fetch(
            `${CLIENT_SERVICE_URL}/api/v1/auth/refresh`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({ refresh_token: refreshToken }),
            }
          );

          const data = await response.json();

          if (response.ok) {
            set({
              accessToken: data.tokens.access_token,
              isLoading: false,
              hasCheckedAuth: true,
            });
          } else {
            get().logout();
          }
        } catch (error) {
          console.error("Token refresh error:", error);
          get().logout();
        }
      },

      clearError: () => {
        set({ error: null });
      },
    }),
    {
      name: "admin-auth-storage",
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        hasCheckedAuth: state.hasCheckedAuth,
      }),
    }
  )
);

