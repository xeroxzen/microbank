"use client";

import { useAuth } from "@/hooks/useAuth";
import AdminLoginForm from "@/components/auth/AdminLoginForm";
import AdminDashboard from "@/components/admin/AdminDashboard";

export default function AdminPanel() {
  const { user, isLoading, isAuthenticated, logout } = useAuth();

  // Show loading state while checking authentication
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-red-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading admin panel...</p>
        </div>
      </div>
    );
  }

  // Show login form if not authenticated
  if (!isAuthenticated) {
    return <AdminLoginForm />;
  }

  // Show admin dashboard if authenticated
  return <AdminDashboard user={user!} onLogout={logout} />;
}
