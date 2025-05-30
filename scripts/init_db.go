package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func main() {
	// Подключение к SQLite
	db, err := sql.Open("sqlite", "./forum.db")
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	// Чтение SQL файла
	sqlFile := filepath.Join("migrations", "001_init.sql")
	sqlBytes, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		fmt.Printf("Failed to read SQL file: %v\n", err)
		os.Exit(1)
	}

	// Выполнение SQL
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		fmt.Printf("Failed to execute SQL: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database initialized successfully!")
}
