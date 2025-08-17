package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"microbank/banking-service/internal/models"
)

// AccountRepositoryImpl handles all database operations related to accounts
type AccountRepositoryImpl struct {
	db *PostgresDB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *PostgresDB) AccountRepository {
	return &AccountRepositoryImpl{db: db}
}

// CreateAccount creates a new account for a user
func (r *AccountRepositoryImpl) CreateAccount(userID uuid.UUID) (*models.Account, error) {
	query := `
		INSERT INTO accounts (id, user_id, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, balance, created_at, updated_at`

	now := time.Now()
	account := &models.Account{
		ID:        uuid.New(),
		UserID:    userID,
		Balance:   0.00,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := r.db.QueryRow(
		query,
		account.ID,
		account.UserID,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// GetOrCreateAccount gets an existing account or creates a new one for a user
func (r *AccountRepositoryImpl) GetOrCreateAccount(userID uuid.UUID) (*models.Account, error) {
	// Check if account exists
	exists, err := r.AccountExists(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check account existence: %w", err)
	}

	if exists {
		// Get existing account
		account, err := r.GetAccountByUserID(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing account: %w", err)
		}
		return account, nil
	}

	// Create new account
	account, err := r.CreateAccount(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new account: %w", err)
	}

	return account, nil
}

// GetAccountByUserID retrieves an account by user ID
func (r *AccountRepositoryImpl) GetAccountByUserID(userID uuid.UUID) (*models.Account, error) {
	query := `
		SELECT id, user_id, balance, created_at, updated_at
		FROM accounts WHERE user_id = $1`

	account := &models.Account{}
	err := r.db.QueryRow(query, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found for user")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

// GetAccountByID retrieves an account by its ID
func (r *AccountRepositoryImpl) GetAccountByID(id uuid.UUID) (*models.Account, error) {
	query := `
		SELECT id, user_id, balance, created_at, updated_at
		FROM accounts WHERE id = $1`

	account := &models.Account{}
	err := r.db.QueryRow(query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

// UpdateBalance updates the account balance
func (r *AccountRepositoryImpl) UpdateBalance(accountID uuid.UUID, newBalance float64) error {
	query := `
		UPDATE accounts 
		SET balance = $1, updated_at = $2
		WHERE id = $3`

	result, err := r.db.Exec(query, newBalance, time.Now(), accountID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found for balance update")
	}

	return nil
}

// AccountExists checks if an account exists for a user
func (r *AccountRepositoryImpl) AccountExists(userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM accounts WHERE user_id = $1)`

	var exists bool
	err := r.db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if account exists: %w", err)
	}

	return exists, nil
}

// GetAllAccounts retrieves all accounts (for admin purposes)
func (r *AccountRepositoryImpl) GetAllAccounts() ([]models.Account, error) {
	query := `
		SELECT id, user_id, balance, created_at, updated_at
		FROM accounts
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account row: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over account rows: %w", err)
	}

	return accounts, nil
}
