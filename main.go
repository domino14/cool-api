package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	// "github.com/domino14/cool-api/hooked"
)

// A database creation function. On production this shouldn't exist,
// we need it here for initial bootstrapping.
func createDB(user, pass, host, port, dbName string) error {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/?sslmode=disable", user, pass, host, port)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[DEBUG] Created db object")

	_, err = db.Exec("CREATE DATABASE " + dbName)
	return err
}

// Create database if it doesn't exist, and load initial fixtures.
func initializeDB() {
	// Env vars from local_config.env
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	err := createDB(user, pass, host, port, dbName)
	if err != nil {
		// Probably OK, the database already exists.
		log.Printf("[INFO] Error creating database: %s", err)
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port,
		dbName)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Printf("[DEBUG] Connected to database %s", dbName)
}

func main() {
	log.Println("Connecting to db...")
	initializeDB()
}
