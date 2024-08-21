package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

var Db *sql.DB

func InitDB() {
	var err error
	// Replace DSN with your actual MySQL user, password, and database
	//dsn := "username:password@tcp(127.0.0.1:3306)/dbname"
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBTABLE"))
	//log.Printf("connection string: %v\n", dsn)
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Check if the database is reachable
	if err := Db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
}
