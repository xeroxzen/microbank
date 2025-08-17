"use client";

import { useState, useEffect } from "react";
import { useAuthStore } from "@/stores/authStore";

interface Transaction {
  id: string;
  type: "deposit" | "withdrawal";
  amount: number;
  description: string;
  balance_before: number;
  balance_after: number;
  created_at: string;
}

interface TransactionResponse {
  message: string;
  pagination: {
    count: number;
    limit: number;
    offset: number;
  };
  transactions: Transaction[];
}

export default function TransactionHistory() {
  const { accessToken } = useAuthStore();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [filter, setFilter] = useState<"all" | "deposit" | "withdrawal">("all");

  const limit = 10;
  const offset = (currentPage - 1) * limit;

  useEffect(() => {
    fetchTransactions();
  }, [currentPage, filter]);

  const fetchTransactions = async () => {
    if (!accessToken) return;

    try {
      setIsLoading(true);
      setError(null);

      const response = await fetch(
        `http://localhost:8080/api/v1/account/transactions?limit=${limit}&offset=${offset}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        }
      );

      if (response.ok) {
        const data: TransactionResponse = await response.json();
        setTransactions(data.transactions);
        setTotalCount(data.pagination.count);
      } else {
        const errorData = await response.json();
        setError(errorData.error?.message || "Failed to fetch transactions");
      }
    } catch (err) {
      setError("Network error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  const filteredTransactions = transactions.filter((transaction) => {
    if (filter === "all") return true;
    return transaction.type === filter;
  });

  const totalPages = Math.ceil(totalCount / limit);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const formatAmount = (amount: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(amount);
  };

  if (isLoading && transactions.length === 0) {
    return (
      <div className="card">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
          <div className="space-y-4">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-16 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (error && transactions.length === 0) {
    return (
      <div className="card">
        <div className="text-center">
          <div className="text-red-600 mb-4">
            <svg
              className="mx-auto h-12 w-12"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
              />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            Error Loading Transactions
          </h3>
          <p className="text-gray-600 mb-4">{error}</p>
          <button onClick={fetchTransactions} className="btn-primary">
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header and Filters */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">
            Transaction History
          </h2>
          <p className="mt-1 text-sm text-gray-600">
            View all your recent deposits and withdrawals
          </p>
        </div>

        {/* Filter Buttons */}
        <div className="mt-4 sm:mt-0 flex bg-gray-100 rounded-lg p-1">
          {(["all", "deposit", "withdrawal"] as const).map((filterType) => (
            <button
              key={filterType}
              onClick={() => {
                setFilter(filterType);
                setCurrentPage(1);
              }}
              className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
                filter === filterType
                  ? "bg-white text-blue-600 shadow-sm"
                  : "text-gray-600 hover:text-gray-900"
              }`}
            >
              {filterType === "all" && "📋 All"}
              {filterType === "deposit" && "💰 Deposits"}
              {filterType === "withdrawal" && "💸 Withdrawals"}
            </button>
          ))}
        </div>
      </div>

      {/* Transactions List */}
      <div className="card">
        {filteredTransactions.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-gray-400 mb-4">
              <svg
                className="mx-auto h-12 w-12"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
                />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              No transactions found
            </h3>
            <p className="text-gray-600">
              {filter === "all"
                ? "You haven't made any transactions yet."
                : `No ${filter}s found in your history.`}
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {filteredTransactions.map((transaction) => (
              <div
                key={transaction.id}
                className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <div className="flex items-center space-x-4">
                  <div
                    className={`h-10 w-10 rounded-full flex items-center justify-center ${
                      transaction.type === "deposit"
                        ? "bg-green-100 text-green-600"
                        : "bg-red-100 text-red-600"
                    }`}
                  >
                    {transaction.type === "deposit" ? "💰" : "💸"}
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">
                      {transaction.description}
                    </p>
                    <p className="text-sm text-gray-500">
                      {formatDate(transaction.created_at)}
                    </p>
                  </div>
                </div>

                <div className="text-right">
                  <p
                    className={`font-semibold ${
                      transaction.type === "deposit"
                        ? "text-green-600"
                        : "text-red-600"
                    }`}
                  >
                    {transaction.type === "deposit" ? "+" : "-"}{" "}
                    {formatAmount(transaction.amount)}
                  </p>
                  <p className="text-sm text-gray-500">
                    Balance: {formatAmount(transaction.balance_after)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <div className="text-sm text-gray-700">
            Showing {offset + 1} to {Math.min(offset + limit, totalCount)} of{" "}
            {totalCount} results
          </div>

          <div className="flex space-x-2">
            <button
              onClick={() => setCurrentPage(currentPage - 1)}
              disabled={currentPage === 1}
              className="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Previous
            </button>

            <span className="px-3 py-2 text-sm text-gray-700">
              Page {currentPage} of {totalPages}
            </span>

            <button
              onClick={() => setCurrentPage(currentPage + 1)}
              disabled={currentPage === totalPages}
              className="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
