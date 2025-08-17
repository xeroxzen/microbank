"use client";

import { useState } from "react";
import { AdminUser } from "@/stores/authStore";
import AdminNavigation from "@/components/layout/AdminNavigation";
import SystemOverview from "./SystemOverview";
import UserManagement from "./UserManagement";

interface AdminDashboardProps {
  user: AdminUser;
  onLogout: () => void;
}

export default function AdminDashboard({
  user,
  onLogout,
}: AdminDashboardProps) {
  const [activeTab, setActiveTab] = useState("overview");

  return (
    <div className="min-h-screen bg-gray-50">
      <AdminNavigation user={user} onLogout={onLogout} />

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-gray-900">
              Welcome back, {user.name}!
            </h1>
            <p className="mt-2 text-gray-600">
              Manage users, monitor system status, and administer the Microbank
              platform.
            </p>
          </div>

          {/* Tab Navigation */}
          <div className="border-b border-gray-200 mb-8">
            <nav className="-mb-px flex space-x-8">
              {[
                { id: "overview", name: "System Overview", icon: "📊" },
                { id: "users", name: "User Management", icon: "👥" },
              ].map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`py-2 px-1 border-b-2 font-medium text-sm ${
                    activeTab === tab.id
                      ? "border-red-500 text-red-600"
                      : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                  }`}
                >
                  <span className="mr-2">{tab.icon}</span>
                  {tab.name}
                </button>
              ))}
            </nav>
          </div>

          {/* Tab Content */}
          <div className="space-y-6">
            {activeTab === "overview" && <SystemOverview />}
            {activeTab === "users" && <UserManagement />}
          </div>
        </div>
      </main>
    </div>
  );
}

