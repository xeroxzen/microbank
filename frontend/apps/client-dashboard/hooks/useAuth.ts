import { useEffect } from "react";
import { useAuthStore } from "@/stores/authStore";

export function useAuth() {
  const { user, isLoading, checkAuth, logout } = useAuthStore();

  useEffect(() => {
    // Only check auth once on mount if we don't have a user
    if (!user) {
      checkAuth();
    }
  }, []); // Only run once on mount

  return {
    user,
    isLoading,
    logout,
    isAuthenticated: !!user,
  };
}
