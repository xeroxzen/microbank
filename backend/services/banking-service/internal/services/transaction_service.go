package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"microbank/banking-service/internal/models"
	"microbank/banking-service/internal/repository"
)

// TransactionService handles transaction-related business logic
type TransactionService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(transactionRepo repository.TransactionRepository, accountRepo repository.AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

// ProcessDeposit processes a deposit transaction
func (s *TransactionService) ProcessDeposit(userID uuid.UUID, amount float64, description string) (*models.Transaction, error) {
	// Validate amount
	if amount <= 0 {
		return nil, fmt.Errorf("deposit amount must be greater than zero")
	}

	// Get or create account for user
	account, err := s.accountRepo.GetOrCreateAccount(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create account: %w", err)
	}

	// Calculate new balance
	balanceBefore := account.Balance
	balanceAfter := balanceBefore + amount

	// Create transaction record
	transaction := &models.Transaction{
		ID:            uuid.New(),
		AccountID:     account.ID,
		UserID:        userID,
		Type:          models.TransactionTypeDeposit,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   description,
		CreatedAt:     time.Now(),
	}

	// Save transaction to database
	if err := s.transactionRepo.CreateTransaction(transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// Update account balance
	if err := s.accountRepo.UpdateBalance(account.ID, balanceAfter); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	return transaction, nil
}

// ProcessWithdrawal processes a withdrawal transaction
func (s *TransactionService) ProcessWithdrawal(userID uuid.UUID, amount float64, description string) (*models.Transaction, error) {
	// Validate amount
	if amount <= 0 {
		return nil, fmt.Errorf("withdrawal amount must be greater than zero")
	}

	// Get account for user
	account, err := s.accountRepo.GetAccountByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Check if user has sufficient funds
	if account.Balance < amount {
		return nil, fmt.Errorf("insufficient funds: requested %f, available %f", amount, account.Balance)
	}

	// Calculate new balance
	balanceBefore := account.Balance
	balanceAfter := balanceBefore - amount

	// Create transaction record
	transaction := &models.Transaction{
		ID:            uuid.New(),
		AccountID:     account.ID,
		UserID:        userID,
		Type:          models.TransactionTypeWithdrawal,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   description,
		CreatedAt:     time.Now(),
	}

	// Save transaction to database
	if err := s.transactionRepo.CreateTransaction(transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// Update account balance
	if err := s.accountRepo.UpdateBalance(account.ID, balanceAfter); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	return transaction, nil
}

// GetTransactionByID retrieves a specific transaction
func (s *TransactionService) GetTransactionByID(transactionID uuid.UUID) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.GetTransactionByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

// GetTransactionsByUserID retrieves transactions for a specific user
func (s *TransactionService) GetTransactionsByUserID(userID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	// Set default values if not provided
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	transactions, err := s.transactionRepo.GetTransactionsByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

// GetTransactionCountByUserID gets the total count of transactions for a user
func (s *TransactionService) GetTransactionCountByUserID(userID uuid.UUID) (int, error) {
	count, err := s.transactionRepo.GetTransactionCountByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction count: %w", err)
	}

	return count, nil
}

// GetAllTransactions retrieves all transactions (for admin purposes)
func (s *TransactionService) GetAllTransactions(limit, offset int) ([]models.Transaction, error) {
	// Set default values if not provided
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	transactions, err := s.transactionRepo.GetAllTransactions(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
