package main

import (
	"fmt"
	"log"
	"os"

	"kyooar/internal/shared/config"
	"kyooar/internal/shared/database"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-query" {
		log.Fatal("Usage: go run main.go -query 'SQL_QUERY'")
	}

	query := os.Args[2]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Execute query
	rows, err := db.Raw(query).Rows()
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("Failed to get columns:", err)
	}

	fmt.Printf("Columns: %v\n", columns)
	fmt.Println("Results:")

	// Prepare interface slice for scan
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Print results
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal("Failed to scan row:", err)
		}

		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			fmt.Printf("%s: %v, ", col, val)
		}
		fmt.Println()
	}
}