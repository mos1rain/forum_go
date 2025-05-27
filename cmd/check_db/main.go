package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Подключаемся к базе данных
	db, err := sql.Open("postgres", "postgres://postgres:28072005@localhost:5432/forum?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем существование таблицы users
	var tableExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'users'
		)`).Scan(&tableExists)
	if err != nil {
		log.Fatalf("Failed to check table existence: %v", err)
	}

	if !tableExists {
		log.Fatal("Table 'users' does not exist!")
	}

	// Проверяем структуру таблицы
	rows, err := db.Query(`
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_name = 'users'
		ORDER BY ordinal_position`)
	if err != nil {
		log.Fatalf("Failed to get table structure: %v", err)
	}
	defer rows.Close()

	fmt.Println("Table 'users' structure:")
	fmt.Println("------------------------")
	for rows.Next() {
		var columnName, dataType, isNullable string
		if err := rows.Scan(&columnName, &dataType, &isNullable); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("Column: %s, Type: %s, Nullable: %s\n", columnName, dataType, isNullable)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}

	fmt.Println("\nMigration check completed successfully!")
}
