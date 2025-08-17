"use client";

import { AdminUser } from "@/stores/authStore";

interface AdminNavigationProps {
  user: AdminUser;
  onLogout: () => void;
}

export default function AdminNavigation({
  user,
  onLogout,
}: AdminNavigationProps) {
  return (
    <nav className="bg-red-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="h-8 w-8 bg-white rounded-full flex items-center justify-center">
                <svg
                  className="h-6 w-6 text-red-700"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </div>
            </div>
            <div className="ml-4">
              <h1 className="text-white text-xl font-bold">Microbank Admin</h1>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <div className="text-white text-sm">
              <span className="font-medium">{user.name}</span>
              <span className="ml-2 text-red-200">({user.email})</span>
            </div>
            <button
              onClick={onLogout}
              className="bg-red-600 hover:bg-red-500 text-white px-3 py-2 rounded-md text-sm font-medium transition-colors duration-200"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}

