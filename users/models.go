package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID           uint64  `json:"id"`
	FirstName    *string `json:"firstName,omitempty"` // Pointer means optional
	Surname      *string `json:"surname,omitempty"`   // Pointer means optional
	Email        string  `json:"email"`
	PasswordHash []byte  `json:"-"` // Omit from JSON responses
}

type UserPost struct {
	FirstName *string `json:"firstName"` // Optional fields
	Surname   *string `json:"surname"`   // Optional fields
	Email     string  `json:"email"`
	Password  string  `json:"password"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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
