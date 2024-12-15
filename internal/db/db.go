package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB // Global variable for database connection

func InitDB(connString string) error {
	var err error
	DB, err = sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping the database: %w", err)
	}

	fmt.Println("Connected to the database successfully!")
	return nil
}
