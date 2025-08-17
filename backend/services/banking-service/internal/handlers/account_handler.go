package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"microbank/banking-service/internal/services"
)

// AccountHandler handles account-related HTTP requests
type AccountHandler struct {
	accountService     *services.AccountService
	transactionService *services.TransactionService
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService *services.AccountService, transactionService *services.TransactionService) *AccountHandler {
	return &AccountHandler{
		accountService:     accountService,
		transactionService: transactionService,
	}
}

// GetBalance retrieves the current account balance for the authenticated user
func (h *AccountHandler) GetBalance(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "User information not found in context",
			},
		})
		return
	}

	// Parse user ID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Invalid user ID format",
			},
		})
		return
	}

	// Get account balance
	balance, err := h.accountService.GetAccountBalance(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "ACCOUNT_NOT_FOUND",
				"message": "Account not found",
				"details": err.Error(),
			},
		})
		return
	}

	// Return balance
	c.JSON(http.StatusOK, gin.H{
		"message": "Balance retrieved successfully",
		"balance": balance,
		"currency": "USD",
	})
}

// GetTransactions retrieves transaction history for the authenticated user
func (h *AccountHandler) GetTransactions(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "User information not found in context",
			},
		})
		return
	}

	// Parse user ID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Invalid user ID format",
			},
		})
		return
	}

	// Get query parameters for pagination
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get transactions
	transactions, err := h.transactionService.GetTransactionsByUserID(userUUID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "FETCH_TRANSACTIONS_FAILED",
				"message": "Failed to fetch transactions",
				"details": err.Error(),
			},
		})
		return
	}

	// Convert transactions to response format
	var transactionResponses []gin.H
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, gin.H{
			"id":             transaction.ID,
			"type":           transaction.Type,
			"amount":         transaction.Amount,
			"balance_before": transaction.BalanceBefore,
			"balance_after":  transaction.BalanceAfter,
			"description":    transaction.Description,
			"created_at":     transaction.CreatedAt,
		})
	}

	// Return transactions
	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions retrieved successfully",
		"transactions": transactionResponses,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(transactionResponses),
		},
	})
}
