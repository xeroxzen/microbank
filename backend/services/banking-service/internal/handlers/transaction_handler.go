package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"microbank/banking-service/internal/models"
	"microbank/banking-service/internal/services"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	transactionService *services.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// Deposit handles deposit requests
func (h *TransactionHandler) Deposit(c *gin.Context) {
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

	// Bind and validate request body
	var request models.TransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
				"details": err.Error(),
			},
		})
		return
	}

	// Process deposit
	transaction, err := h.transactionService.ProcessDeposit(userUUID, request.Amount, request.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DEPOSIT_FAILED",
				"message": "Failed to process deposit",
				"details": err.Error(),
			},
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Deposit processed successfully",
		"transaction": transaction.ToResponse(),
	})
}

// Withdraw handles withdrawal requests
func (h *TransactionHandler) Withdraw(c *gin.Context) {
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

	// Bind and validate request body
	var request models.TransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
				"details": err.Error(),
			},
		})
		return
	}

	// Process withdrawal
	transaction, err := h.transactionService.ProcessWithdrawal(userUUID, request.Amount, request.Description)
	if err != nil {
		// Check for specific error types
		if err.Error() == "insufficient funds: requested "+fmt.Sprintf("%f", request.Amount)+", available "+fmt.Sprintf("%f", 0.0) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INSUFFICIENT_FUNDS",
					"message": "Insufficient funds for withdrawal",
					"details": gin.H{
						"requested_amount": request.Amount,
					},
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "WITHDRAWAL_FAILED",
				"message": "Failed to process withdrawal",
				"details": err.Error(),
			},
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Withdrawal processed successfully",
		"transaction": transaction.ToResponse(),
	})
}

// GetTransaction retrieves a specific transaction by ID
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	// Get transaction ID from URL parameter
	transactionIDStr := c.Param("id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_TRANSACTION_ID",
				"message": "Invalid transaction ID format",
			},
		})
		return
	}

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

	// Get transaction
	transaction, err := h.transactionService.GetTransactionByID(transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "TRANSACTION_NOT_FOUND",
				"message": "Transaction not found",
				"details": err.Error(),
			},
		})
		return
	}

	// Check if the transaction belongs to the authenticated user
	if transaction.UserID != userUUID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "ACCESS_DENIED",
				"message": "Access denied to this transaction",
			},
		})
		return
	}

	// Return transaction
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction retrieved successfully",
		"transaction": transaction.ToResponse(),
	})
}
