package models

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
)

// Transaction represents a banking transaction
type Transaction struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	AccountID     uuid.UUID       `json:"account_id" db:"account_id"`
	UserID        uuid.UUID       `json:"user_id" db:"user_id"`
	Type          TransactionType `json:"type" db:"type"`
	Amount        float64         `json:"amount" db:"amount"`
	BalanceBefore float64         `json:"balance_before" db:"balance_before"`
	BalanceAfter  float64         `json:"balance_after" db:"balance_after"`
	Description   string          `json:"description" db:"description"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}

// TransactionRequest represents the data needed to create a transaction
type TransactionRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"max=255"`
}

// TransactionResponse represents the transaction data sent in responses
type TransactionResponse struct {
	ID            uuid.UUID       `json:"id"`
	AccountID     uuid.UUID       `json:"account_id"`
	UserID        uuid.UUID       `json:"user_id"`
	Type          TransactionType `json:"type"`
	Amount        float64         `json:"amount"`
	BalanceBefore float64         `json:"balance_before"`
	BalanceAfter  float64         `json:"balance_after"`
	Description   string          `json:"description"`
	CreatedAt     time.Time       `json:"created_at"`
}

// ToResponse converts a Transaction to TransactionResponse
func (t *Transaction) ToResponse() TransactionResponse {
	return TransactionResponse{
		ID:            t.ID,
		AccountID:     t.AccountID,
		UserID:        t.UserID,
		Type:          t.Type,
		Amount:        t.Amount,
		BalanceBefore: t.BalanceBefore,
		BalanceAfter:  t.BalanceAfter,
		Description:   t.Description,
		CreatedAt:     t.CreatedAt,
	}
}
