"use client";

import { useState, useEffect } from "react";
import { useAuthStore } from "@/stores/authStore";

type TransactionType = "deposit" | "withdrawal";

export default function TransactionForm() {
  const { accessToken } = useAuthStore();
  const [type, setType] = useState<TransactionType>("deposit");
  const [amount, setAmount] = useState("");
  const [description, setDescription] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!amount || parseFloat(amount) <= 0) {
      setMessage({ type: "error", text: "Please enter a valid amount" });
      return;
    }

    if (!description.trim()) {
      setMessage({ type: "error", text: "Please enter a description" });
      return;
    }

    setIsLoading(true);
    setMessage(null);

    try {
      const endpoint = type === "deposit" ? "deposit" : "withdraw";

      // Debug logging
      console.log("Making transaction request:", {
        endpoint,
        accessToken: accessToken ? `${accessToken.substring(0, 20)}...` : null,
        amount: parseFloat(amount),
        description: description.trim(),
      });

      const response = await fetch(
        `http://localhost:8080/api/v1/transactions/${endpoint}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify({
            amount: parseFloat(amount),
            description: description.trim(),
          }),
        }
      );

      const data = await response.json();

      // Debug logging
      console.log("Transaction response:", {
        status: response.status,
        statusText: response.statusText,
        data,
      });

      if (response.ok) {
        setMessage({
          type: "success",
          text: `${type === "deposit" ? "Deposit" : "Withdrawal"} processed successfully! New balance: $${data.transaction.balance_after.toFixed(2)}`,
        });
        setAmount("");
        setDescription("");
      } else {
        setMessage({
          type: "error",
          text: data.error?.message || `Failed to process ${type}`,
        });
      }
    } catch (error) {
      setMessage({
        type: "error",
        text: "Network error occurred",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const clearMessage = () => {
    setMessage(null);
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card">
        <h2 className="text-2xl font-bold text-gray-900 mb-6">
          {type === "deposit" ? "Deposit Money" : "Withdraw Money"}
        </h2>

        {/* Transaction Type Toggle */}
        <div className="mb-6">
          <div className="flex bg-gray-100 rounded-lg p-1">
            <button
              type="button"
              onClick={() => {
                setType("deposit");
                clearMessage();
              }}
              className={`flex-1 py-2 px-4 rounded-md text-sm font-medium transition-colors ${
                type === "deposit"
                  ? "bg-white text-blue-600 shadow-sm"
                  : "text-gray-600 hover:text-gray-900"
              }`}
            >
              💰 Deposit
            </button>
            <button
              type="button"
              onClick={() => {
                setType("withdrawal");
                clearMessage();
              }}
              className={`flex-1 py-2 px-4 rounded-md text-sm font-medium transition-colors ${
                type === "withdrawal"
                  ? "bg-white text-blue-600 shadow-sm"
                  : "text-gray-600 hover:text-gray-900"
              }`}
            >
              💸 Withdraw
            </button>
          </div>
        </div>

        {/* Message Display */}
        {message && (
          <div
            className={`mb-6 p-4 rounded-lg ${
              message.type === "success"
                ? "bg-green-50 border border-green-200 text-green-700"
                : "bg-red-50 border border-red-200 text-red-700"
            }`}
          >
            <div className="flex items-center justify-between">
              <span>{message.text}</span>
              <button
                onClick={clearMessage}
                className="text-gray-400 hover:text-gray-600"
                aria-label="Close message"
                title="Close message"
              >
                <svg
                  className="h-4 w-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Amount Input */}
          <div>
            <label
              htmlFor="amount"
              className="block text-sm font-medium text-gray-700 mb-2"
            >
              Amount ({type === "deposit" ? "to deposit" : "to withdraw"})
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <span className="text-gray-500 sm:text-sm">$</span>
              </div>
              <input
                type="number"
                id="amount"
                name="amount"
                step="0.01"
                min="0.01"
                required
                className="input-field pl-7"
                placeholder="0.00"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
              />
            </div>
            <p className="mt-1 text-sm text-gray-500">
              Enter the amount you want to{" "}
              {type === "deposit" ? "add to" : "remove from"} your account
            </p>
          </div>

          {/* Description Input */}
          <div>
            <label
              htmlFor="description"
              className="block text-sm font-medium text-gray-700 mb-2"
            >
              Description
            </label>
            <input
              type="text"
              id="description"
              name="description"
              required
              className="input-field"
              placeholder={`e.g., ${type === "deposit" ? "Salary deposit, Gift from family" : "ATM withdrawal, Online purchase"}`}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
            <p className="mt-1 text-sm text-gray-500">
              Provide a brief description of this {type}
            </p>
          </div>

          {/* Submit Button */}
          <div className="pt-4">
            <button
              type="submit"
              disabled={isLoading}
              className="btn-primary w-full flex justify-center items-center"
            >
              {isLoading ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  Processing...
                </>
              ) : (
                `${type === "deposit" ? "Deposit" : "Withdraw"} Money`
              )}
            </button>
          </div>
        </form>

        {/* Info Box */}
        <div className="mt-8 p-4 bg-blue-50 border border-blue-200 rounded-lg">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg
                className="h-5 w-5 text-blue-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800">
                {type === "deposit"
                  ? "Deposit Information"
                  : "Withdrawal Information"}
              </h3>
              <div className="mt-2 text-sm text-blue-700">
                <p>
                  {type === "deposit"
                    ? "Deposits are processed immediately and will be reflected in your balance right away."
                    : "Withdrawals are processed immediately. Please ensure you have sufficient funds in your account."}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
