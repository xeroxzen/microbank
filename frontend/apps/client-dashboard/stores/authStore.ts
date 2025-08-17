import { create } from "zustand";
import { persist } from "zustand/middleware";

export interface User {
  id: string;
  email: string;
  name: string;
  is_blacklisted: boolean;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isLoading: boolean;
  error: string | null;
  hasCheckedAuth: boolean;
}

interface AuthActions {
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, name: string, password: string) => Promise<void>;
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

      register: async (email: string, name: string, password: string) => {
        set({ isLoading: true, error: null });

        try {
          const response = await fetch(
            `${CLIENT_SERVICE_URL}/api/v1/auth/register`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({ email, name, password }),
            }
          );

          const data = await response.json();

          if (!response.ok) {
            throw new Error(data.error?.message || "Registration failed");
          }

          set({
            user: data.user,
            isLoading: false,
            error: null,
            hasCheckedAuth: true,
          });
        } catch (error) {
          set({
            error:
              error instanceof Error ? error.message : "Registration failed",
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
          console.log(
            "checkAuth already in progress or already checked, skipping..."
          );
          return;
        }

        // Set loading state immediately to prevent multiple calls
        set({ isLoading: true });

        const { accessToken, refreshToken } = get();
        console.log("checkAuth called with:", {
          accessToken: !!accessToken,
          refreshToken: !!refreshToken,
        });

        // If user is already authenticated, skip
        if (get().user && accessToken) {
          console.log("User already authenticated, skipping checkAuth");
          set({ isLoading: false });
          return;
        }

        // If no tokens at all, we can't proceed
        if (!accessToken && !refreshToken) {
          console.log("No tokens available, cannot proceed");
          set({ isLoading: false });
          return;
        }

        // If no access token but we have a refresh token, try to refresh first
        if (!accessToken && refreshToken) {
          console.log(
            "No access token but refresh token available, attempting refresh..."
          );
          try {
            await get().refreshAccessToken();
            // After refresh, get the new access token and continue
            const newState = get();
            if (!newState.accessToken) {
              console.log("Refresh failed to provide access token");
              set({ isLoading: false });
              return;
            }
            console.log("Refresh successful, proceeding with validation");
          } catch (error) {
            console.log("Refresh attempt failed:", error);
            set({ isLoading: false });
            return;
          }
        }

        // If still no access token, we can't proceed
        if (!get().accessToken) {
          console.log("No access token available, cannot proceed");
          set({ isLoading: false });
          return;
        }

        console.log("Validating access token...");
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
            console.log("Token validation successful");
            set({ user: data.user, isLoading: false, hasCheckedAuth: true });
          } else {
            console.log("Token validation failed, attempting refresh...");
            // Token is invalid, try to refresh
            await get().refreshAccessToken();
          }
        } catch (error) {
          console.error("Token validation error:", error);
          set({ isLoading: false });
          // If refresh fails, logout
          get().logout();
        }
      },

      refreshAccessToken: async () => {
        const { refreshToken } = get();

        if (!refreshToken) {
          console.log("No refresh token available, logging out");
          get().logout();
          return;
        }

        try {
          console.log("Attempting to refresh access token...");
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
          console.log("Refresh response:", { status: response.status, data });

          if (response.ok) {
            console.log("Token refresh successful, setting new access token");
            set({
              accessToken: data.tokens.access_token,
              isLoading: false,
              hasCheckedAuth: true,
            });
          } else {
            console.log("Token refresh failed:", data.error);
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
      name: "auth-storage",
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        hasCheckedAuth: state.hasCheckedAuth,
      }),
    }
  )
);
