package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to postgres database to create new database
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = 'seller')").Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		// Terminate existing connections to template1
		_, _ = db.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'template1' AND pid <> pg_backend_pid()")

		// Create database without template
		_, err = db.Exec(`CREATE DATABASE "seller" WITH TEMPLATE template0`)
		if err != nil {
			// If it still fails, just use the simpler command
			_, err = db.Exec(`CREATE DATABASE "seller"`)
			if err != nil {
				log.Printf("Warning: Could not create database: %v", err)
				log.Println("You may need to create it manually using: CREATE DATABASE \"seller\";")
			} else {
				fmt.Println("Database 'seller' created successfully!")
			}
		} else {
			fmt.Println("Database 'seller' created successfully!")
		}
	} else {
		fmt.Println("Database 'seller' already exists.")
	}
}
