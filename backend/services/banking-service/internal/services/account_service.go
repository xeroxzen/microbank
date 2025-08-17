package services

import (
	"fmt"

	"github.com/google/uuid"
	"microbank/banking-service/internal/models"
	"microbank/banking-service/internal/repository"
)

// AccountService handles account-related business logic
type AccountService struct {
	accountRepo repository.AccountRepository
}

// NewAccountService creates a new account service
func NewAccountService(accountRepo repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

// GetOrCreateAccount gets an existing account or creates a new one for a user
func (s *AccountService) GetOrCreateAccount(userID uuid.UUID) (*models.Account, error) {
	// Check if account exists
	exists, err := s.accountRepo.AccountExists(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check account existence: %w", err)
	}

	if exists {
		// Get existing account
		account, err := s.accountRepo.GetAccountByUserID(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing account: %w", err)
		}
		return account, nil
	}

	// Create new account
	account, err := s.accountRepo.CreateAccount(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new account: %w", err)
	}

	return account, nil
}

// GetAccountByUserID retrieves an account by user ID
func (s *AccountService) GetAccountByUserID(userID uuid.UUID) (*models.Account, error) {
	account, err := s.accountRepo.GetAccountByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

// GetAccountBalance gets the current balance for a user's account
func (s *AccountService) GetAccountBalance(userID uuid.UUID) (float64, error) {
	account, err := s.accountRepo.GetAccountByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get account: %w", err)
	}

	return account.Balance, nil
}

// UpdateAccountBalance updates the account balance
func (s *AccountService) UpdateAccountBalance(accountID uuid.UUID, newBalance float64) error {
	if err := s.accountRepo.UpdateBalance(accountID, newBalance); err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	return nil
}

// GetAllAccounts retrieves all accounts (for admin purposes)
func (s *AccountService) GetAllAccounts() ([]models.Account, error) {
	accounts, err := s.accountRepo.GetAllAccounts()
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}
