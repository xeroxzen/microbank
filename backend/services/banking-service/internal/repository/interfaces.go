package repository

import (
	"github.com/google/uuid"
	"microbank/banking-service/internal/models"
)

// AccountRepository defines the interface for account data operations
type AccountRepository interface {
	CreateAccount(userID uuid.UUID) (*models.Account, error)
	GetAccountByUserID(userID uuid.UUID) (*models.Account, error)
	GetAccountByID(id uuid.UUID) (*models.Account, error)
	GetOrCreateAccount(userID uuid.UUID) (*models.Account, error)
	UpdateBalance(accountID uuid.UUID, newBalance float64) error
	AccountExists(userID uuid.UUID) (bool, error)
	GetAllAccounts() ([]models.Account, error)
}

// TransactionRepository defines the interface for transaction operations
type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uuid.UUID) (*models.Transaction, error)
	GetTransactionsByUserID(userID uuid.UUID, limit, offset int) ([]models.Transaction, error)
	GetTransactionsByAccountID(accountID uuid.UUID, limit, offset int) ([]models.Transaction, error)
	GetTransactionCountByUserID(userID uuid.UUID) (int, error)
	GetAllTransactions(limit, offset int) ([]models.Transaction, error)
}
