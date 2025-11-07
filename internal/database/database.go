package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the database (PostgreSQL or SQLite)
func InitDB() error {
	var err error
	var dbURL string

	// Check if PostgreSQL URL is provided (for Vercel/production)
	dbURL = os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL") // Vercel uses DATABASE_URL
	}

	if dbURL != "" {
		// Use PostgreSQL
		DB, err = sql.Open("postgres", dbURL)
		if err != nil {
			return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
		}

		// Test connection
		if err = DB.Ping(); err != nil {
			return fmt.Errorf("failed to ping PostgreSQL: %v", err)
		}

		log.Println("Connected to PostgreSQL database")
	} else {
		// Fallback to SQLite for local development
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "gourl.db"
		}

		DB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
		if err != nil {
			return fmt.Errorf("failed to connect to SQLite: %v", err)
		}

		// Enable foreign keys for SQLite
		if _, err := DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
			log.Printf("Warning: Could not enable foreign keys: %v", err)
		}

		log.Println("Connected to SQLite database")
	}

	// Create tables (works for both PostgreSQL and SQLite)
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// createTables creates all necessary tables
func createTables() error {
	// Detect database type
	isPostgres := os.Getenv("DATABASE_URL") != "" || os.Getenv("POSTGRES_URL") != ""

	var createTableSQL string

	if isPostgres {
		// PostgreSQL syntax
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_username ON users(username);
		CREATE INDEX IF NOT EXISTS idx_email ON users(email);
		
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			code VARCHAR(255) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			user_id INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		);
		
		CREATE INDEX IF NOT EXISTS idx_code ON urls(code);
		CREATE INDEX IF NOT EXISTS idx_user_id ON urls(user_id);
		
		CREATE TABLE IF NOT EXISTS clicks (
			id SERIAL PRIMARY KEY,
			url_id INTEGER NOT NULL,
			ip_address VARCHAR(255),
			user_agent TEXT,
			referrer TEXT,
			country VARCHAR(100),
			clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (url_id) REFERENCES urls(id) ON DELETE CASCADE
		);
		
		CREATE INDEX IF NOT EXISTS idx_url_id ON clicks(url_id);
		CREATE INDEX IF NOT EXISTS idx_clicked_at ON clicks(clicked_at);
		`
	} else {
		// SQLite syntax
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_username ON users(username);
		CREATE INDEX IF NOT EXISTS idx_email ON users(email);
		
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			user_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		);
		
		CREATE INDEX IF NOT EXISTS idx_code ON urls(code);
		CREATE INDEX IF NOT EXISTS idx_user_id ON urls(user_id);
		
		CREATE TABLE IF NOT EXISTS clicks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url_id INTEGER NOT NULL,
			ip_address TEXT,
			user_agent TEXT,
			referrer TEXT,
			country TEXT,
			clicked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (url_id) REFERENCES urls(id) ON DELETE CASCADE
		);
		
		CREATE INDEX IF NOT EXISTS idx_url_id ON clicks(url_id);
		CREATE INDEX IF NOT EXISTS idx_clicked_at ON clicks(clicked_at);
		`
	}

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// Migrate existing databases: add country column if it doesn't exist
	return migrateAddCountryColumn(isPostgres)
}

// migrateAddCountryColumn adds the country column to existing clicks table if it doesn't exist
func migrateAddCountryColumn(isPostgres bool) error {
	var alterSQL string
	if isPostgres {
		// PostgreSQL: Check if column exists, if not add it
		alterSQL = `
			DO $$ 
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM information_schema.columns 
					WHERE table_name = 'clicks' AND column_name = 'country'
				) THEN
					ALTER TABLE clicks ADD COLUMN country VARCHAR(100);
				END IF;
			END $$;
		`
	} else {
		// SQLite: Check if column exists, if not add it
		// SQLite doesn't support IF NOT EXISTS for ALTER TABLE ADD COLUMN directly
		// So we'll try to add it and ignore the error if it already exists
		alterSQL = `ALTER TABLE clicks ADD COLUMN country TEXT;`
	}

	_, err := DB.Exec(alterSQL)
	// Ignore error if column already exists
	if err != nil {
		errStr := err.Error()
		if isPostgres || (!strings.Contains(errStr, "duplicate column") && !strings.Contains(errStr, "already exists")) {
			log.Printf("Warning: Could not add country column (might already exist): %v", err)
		}
	}
	return nil
}

// IsPostgres returns true if using PostgreSQL
func IsPostgres() bool {
	return os.Getenv("DATABASE_URL") != "" || os.Getenv("POSTGRES_URL") != ""
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
