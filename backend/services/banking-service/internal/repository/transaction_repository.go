package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"microbank/banking-service/internal/models"
)

// TransactionRepositoryImpl handles all database operations related to transactions
type TransactionRepositoryImpl struct {
	db *PostgresDB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *PostgresDB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

// CreateTransaction creates a new transaction record
func (r *TransactionRepositoryImpl) CreateTransaction(transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (id, account_id, user_id, type, amount, balance_before, balance_after, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(
		query,
		transaction.ID,
		transaction.AccountID,
		transaction.UserID,
		transaction.Type,
		transaction.Amount,
		transaction.BalanceBefore,
		transaction.BalanceAfter,
		transaction.Description,
		transaction.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetTransactionByID retrieves a transaction by its ID
func (r *TransactionRepositoryImpl) GetTransactionByID(id uuid.UUID) (*models.Transaction, error) {
	query := `
		SELECT id, account_id, user_id, type, amount, balance_before, balance_after, description, created_at
		FROM transactions WHERE id = $1`

	transaction := &models.Transaction{}
	err := r.db.QueryRow(query, id).Scan(
		&transaction.ID,
		&transaction.AccountID,
		&transaction.UserID,
		&transaction.Type,
		&transaction.Amount,
		&transaction.BalanceBefore,
		&transaction.BalanceAfter,
		&transaction.Description,
		&transaction.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

// GetTransactionsByUserID retrieves all transactions for a specific user
func (r *TransactionRepositoryImpl) GetTransactionsByUserID(userID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	query := `
		SELECT id, account_id, user_id, type, amount, balance_before, balance_after, description, created_at
		FROM transactions 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.BalanceBefore,
			&transaction.BalanceAfter,
			&transaction.Description,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over transaction rows: %w", err)
	}

	return transactions, nil
}

// GetTransactionsByAccountID retrieves all transactions for a specific account
func (r *TransactionRepositoryImpl) GetTransactionsByAccountID(accountID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	query := `
		SELECT id, account_id, user_id, type, amount, balance_before, balance_after, description, created_at
		FROM transactions 
		WHERE account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.BalanceBefore,
			&transaction.BalanceAfter,
			&transaction.Description,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over transaction rows: %w", err)
	}

	return transactions, nil
}

// GetTransactionCountByUserID gets the total count of transactions for a user
func (r *TransactionRepositoryImpl) GetTransactionCountByUserID(userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE user_id = $1`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction count: %w", err)
	}

	return count, nil
}

// GetAllTransactions retrieves all transactions (for admin purposes)
func (r *TransactionRepositoryImpl) GetAllTransactions(limit, offset int) ([]models.Transaction, error) {
	query := `
		SELECT id, account_id, user_id, type, amount, balance_before, balance_after, description, created_at
		FROM transactions 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.BalanceBefore,
			&transaction.BalanceAfter,
			&transaction.Description,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over transaction rows: %w", err)
	}

	return transactions, nil
}
