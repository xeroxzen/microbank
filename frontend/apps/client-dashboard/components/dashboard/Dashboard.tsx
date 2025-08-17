"use client";

import { useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import Navigation from "./Navigation";
import AccountOverview from "./AccountOverview";
import TransactionHistory from "./TransactionHistory";
import TransactionForm from "./TransactionForm";

export default function Dashboard() {
  const { user, logout } = useAuth();
  const [activeTab, setActiveTab] = useState("overview");

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation user={user} onLogout={logout} />

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-gray-900">
              Welcome back, {user.name}!
            </h1>
            <p className="mt-2 text-gray-600">
              Manage your account, view transactions, and stay on top of your
              finances.
            </p>
          </div>

          {/* Tab Navigation */}
          <div className="border-b border-gray-200 mb-8">
            <nav className="-mb-px flex space-x-8">
              {[
                { id: "overview", name: "Account Overview", icon: "📊" },
                { id: "transactions", name: "Transactions", icon: "💳" },
                { id: "transfer", name: "Transfer Money", icon: "💰" },
              ].map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`py-2 px-1 border-b-2 font-medium text-sm ${
                    activeTab === tab.id
                      ? "border-blue-500 text-blue-600"
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
            {activeTab === "overview" && <AccountOverview />}
            {activeTab === "transactions" && <TransactionHistory />}
            {activeTab === "transfer" && <TransactionForm />}
          </div>
        </div>
      </main>
    </div>
  );
}
