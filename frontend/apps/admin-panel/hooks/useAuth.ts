import { useEffect } from "react";
import { useAuthStore } from "@/stores/authStore";

export const useAuth = () => {
  const { user, isLoading, error, hasCheckedAuth, checkAuth, logout } =
    useAuthStore();

  useEffect(() => {
    if (!hasCheckedAuth) {
      checkAuth();
    }
  }, [hasCheckedAuth, checkAuth]);

  return {
    user,
    isLoading,
    error,
    hasCheckedAuth,
    isAuthenticated: !!user,
    logout,
  };
};

