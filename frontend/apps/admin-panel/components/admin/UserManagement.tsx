"use client";

import { useState, useEffect } from "react";
import { useAuthStore } from "@/stores/authStore";

interface User {
  id: string;
  email: string;
  name: string;
  is_blacklisted: boolean;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

export default function UserManagement() {
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [showBlacklistModal, setShowBlacklistModal] = useState(false);
  const [blacklistReason, setBlacklistReason] = useState("");

  const { accessToken } = useAuthStore();

  const CLIENT_SERVICE_URL =
    process.env.NEXT_PUBLIC_CLIENT_SERVICE_URL || "http://localhost:8081";

  // Fetch all users
  const fetchUsers = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `${CLIENT_SERVICE_URL}/api/v1/admin/clients`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error("Failed to fetch users");
      }

      const data = await response.json();
      setUsers(data.users || []);
    } catch (error) {
      setError(
        error instanceof Error ? error.message : "Failed to fetch users"
      );
    } finally {
      setIsLoading(false);
    }
  };

  // Blacklist a user
  const blacklistUser = async (userId: string, reason: string) => {
    try {
      const response = await fetch(
        `${CLIENT_SERVICE_URL}/api/v1/admin/clients/${userId}/blacklist`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify({ reason }),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to blacklist user");
      }

      // Update local state
      setUsers(
        users.map((user) =>
          user.id === userId ? { ...user, is_blacklisted: true } : user
        )
      );

      setShowBlacklistModal(false);
      setSelectedUser(null);
      setBlacklistReason("");
    } catch (error) {
      setError(
        error instanceof Error ? error.message : "Failed to blacklist user"
      );
    }
  };

  // Remove user from blacklist
  const removeFromBlacklist = async (userId: string) => {
    try {
      const response = await fetch(
        `${CLIENT_SERVICE_URL}/api/v1/admin/clients/${userId}/blacklist`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error("Failed to remove user from blacklist");
      }

      // Update local state
      setUsers(
        users.map((user) =>
          user.id === userId ? { ...user, is_blacklisted: false } : user
        )
      );
    } catch (error) {
      setError(
        error instanceof Error
          ? error.message
          : "Failed to remove user from blacklist"
      );
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleBlacklistClick = (user: User) => {
    setSelectedUser(user);
    setShowBlacklistModal(true);
  };

  const handleBlacklistSubmit = () => {
    if (selectedUser && blacklistReason.trim()) {
      blacklistUser(selectedUser.id, blacklistReason.trim());
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-red-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-gray-900">User Management</h2>
        <button
          onClick={fetchUsers}
          className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200"
        >
          Refresh
        </button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
          {error}
        </div>
      )}

      <div className="bg-white shadow overflow-hidden sm:rounded-md">
        <ul className="divide-y divide-gray-200">
          {users.map((user) => (
            <li key={user.id} className="px-6 py-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center">
                  <div className="flex-shrink-0 h-10 w-10">
                    <div className="h-10 w-10 rounded-full bg-gray-300 flex items-center justify-center">
                      <span className="text-sm font-medium text-gray-700">
                        {user.name.charAt(0).toUpperCase()}
                      </span>
                    </div>
                  </div>
                  <div className="ml-4">
                    <div className="text-sm font-medium text-gray-900">
                      {user.name}
                    </div>
                    <div className="text-sm text-gray-500">{user.email}</div>
                    <div className="flex items-center space-x-2 mt-1">
                      {user.is_admin && (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                          Admin
                        </span>
                      )}
                      {user.is_blacklisted && (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                          Blacklisted
                        </span>
                      )}
                    </div>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  {user.is_blacklisted ? (
                    <button
                      onClick={() => removeFromBlacklist(user.id)}
                      className="bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded text-sm font-medium transition-colors duration-200"
                    >
                      Remove from Blacklist
                    </button>
                  ) : (
                    <button
                      onClick={() => handleBlacklistClick(user)}
                      className="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded text-sm font-medium transition-colors duration-200"
                    >
                      Blacklist
                    </button>
                  )}
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>

      {/* Blacklist Modal */}
      {showBlacklistModal && selectedUser && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
            <div className="mt-3">
              <h3 className="text-lg font-medium text-gray-900 mb-4">
                Blacklist User
              </h3>
              <p className="text-sm text-gray-500 mb-4">
                Are you sure you want to blacklist{" "}
                <strong>{selectedUser.name}</strong>?
              </p>
              <div className="mb-4">
                <label
                  htmlFor="reason"
                  className="block text-sm font-medium text-gray-700 mb-2"
                >
                  Reason (optional)
                </label>
                <textarea
                  id="reason"
                  value={blacklistReason}
                  onChange={(e) => setBlacklistReason(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-red-500 focus:border-red-500"
                  rows={3}
                  placeholder="Enter reason for blacklisting..."
                />
              </div>
              <div className="flex justify-end space-x-3">
                <button
                  onClick={() => {
                    setShowBlacklistModal(false);
                    setSelectedUser(null);
                    setBlacklistReason("");
                  }}
                  className="bg-gray-300 hover:bg-gray-400 text-gray-700 px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                >
                  Cancel
                </button>
                <button
                  onClick={handleBlacklistSubmit}
                  className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                >
                  Confirm Blacklist
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

