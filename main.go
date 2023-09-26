package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	maxConnections = 10
)

func main() {
	// main configs
	dbHost := "localhost"
	dbPort := "5432"
	dbName := "root"
	dbUser := "root"
	dbPassword := "example"

	// build connection
	dbURL := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUser, dbPassword)

	// Create a slice to store connections
	var dbConnections []*sql.DB

	// Init pool
	for i := 0; i < maxConnections; i++ {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatal(err)
		}
		dbConnections = append(dbConnections, db)
	}

	defer func() {
		for _, db := range dbConnections {
			db.Close()
		}
	}()
	fmt.Println("Connection Pool Created!")

	// Usage example
	conn := dbConnections[0]
	_, err := conn.Exec("INSERT INTO my_pool_table (email, phone) VALUES ($1, $2)", "golang@gmail", "123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Insert successfully!")

	// remember to release the connection
	conn = nil
	// After you finish using a connection, you can set the 'conn' variable back to nil.

	// Reuse conn
	for i, db := range dbConnections {
		if db == nil {
			conn = db
			dbConnections[i] = conn // in use again
			break
		}
	}

	fmt.Println("The end")
}
