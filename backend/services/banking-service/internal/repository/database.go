package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// PostgresDB holds the database connection
type PostgresDB struct {
	*sql.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB() (*PostgresDB, error) {
	// Get database connection parameters from environment
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "banking_service")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Successfully connected to PostgreSQL database")

	// Initialize database schema
	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return &PostgresDB{db}, nil
}

// initSchema creates the necessary database tables if they don't exist
func initSchema(db *sql.DB) error {
	// Create accounts table
	createAccountsTable := `
	CREATE TABLE IF NOT EXISTS accounts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID UNIQUE NOT NULL,
		balance DECIMAL(15,2) DEFAULT 0.00,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Create transactions table
	createTransactionsTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,
		user_id UUID NOT NULL,
		type VARCHAR(20) NOT NULL CHECK (type IN ('deposit', 'withdrawal')),
		amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
		balance_before DECIMAL(15,2) NOT NULL,
		balance_after DECIMAL(15,2) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Create indexes for better performance
	createIndexes := `
	CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
	CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);`

	// Execute schema creation
	queries := []string{createAccountsTable, createTransactionsTable, createIndexes}
	
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute schema query: %w", err)
		}
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
