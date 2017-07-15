package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"

	"github.com/domino14/cool-api/hooked"
)

// A database creation function. On production this shouldn't exist,
// we need it here for initial bootstrapping.
func createDB(user, pass, host, port, dbName string) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/?sslmode=disable", user, pass, host, port)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[DEBUG] Created db object")

	_, err = db.Exec("CREATE DATABASE " + dbName)

	if err != nil {
		// Probably OK, the database already exists.
		log.Printf("[INFO] Error creating database: %s", err)
	}
}

// Create database if it doesn't exist, and load initial fixtures.
func initializeDB() {
	// Env vars from local_config.env
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	createDB(user, pass, host, port, dbName)

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

	// Create tables if they don't exist.
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users(
            sid varchar(24) primary key,
            firstname varchar(40) NOT NULL,
            lastname varchar(40) NOT NULL
        )`)
	if err != nil {
		log.Printf("[INFO] Create table users error: %s", err)
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS stories(
            sid varchar(24) primary key,
            title varchar(128) NOT NULL,
            author varchar(24) REFERENCES users(sid)
        )`)
	if err != nil {
		log.Printf("[INFO] Create table stories error: %s", err)
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS activities(
            sid varchar(24) primary key,
            action varchar(24) NOT NULL,
            date timestamptz NOT NULL,
            actor varchar(24) REFERENCES users(sid),
            user2 varchar(24) REFERENCES users(sid)
        )`)
	if err != nil {
		log.Printf("[INFO] Create table activities error: %s", err)
	}

	// Technically notifications can be computed from the activities,
	// but that is inefficient. Let's instead create a notifications
	// table.
	// In this case, the `notified` is the user whose notification we are getting.
	// The `actor` is the performer of the action, which is not necessarily
	// the same as the notified user.
	// In the case of a follow, it will be mapped to `user2` in the API output.
	// Story is null if this is a follow.
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS notifications(
            id uuid primary key,
            notified varchar(24) REFERENCES users(sid),
            actor varchar(24) REFERENCES users(sid),
            action varchar(24) NOT NULL,
            story varchar(24) REFERENCES stories(sid)
        )`)
	if err != nil {
		log.Printf("[INFO] Create table notifications error: %s", err)
	}

	hooked.LoadFixtures(db) // in loader.go
}

func main() {
	log.Println("Connecting to db...")
	initializeDB()
	hooked.Serve("8086")
}
