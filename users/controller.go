package main

import (
	"database/sql"
	"strings"
)

func AddUser(u User) (User, error) {
	result, err := Db.Exec("INSERT INTO users (FirstName, Surname, Email, PasswordHash) VALUES (?, ?, ?, ?)",
		u.FirstName, u.Surname, u.Email, u.PasswordHash)
	if err != nil {
		logger.Printf("addUser: %v", err)
		return User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		logger.Printf("addUser: %v", err)
		return User{}, err
	}
	u.Id = uint64(id)
	return u, nil
}

func GetUserByID(id uint64) (User, error) {
	var user User

	// Query the database for the user by ID
	row := Db.QueryRow("SELECT id, FirstName, Surname, Email, PasswordHash FROM users WHERE id = ?", id)

	// Scan the result into the user struct
	err := row.Scan(&user.Id, &user.FirstName, &user.Surname, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Printf("getUserByID: no user found with id %v", id)
			return User{}, nil // No user found
		}
		logger.Printf("getUserByID: %v", err)
		return User{}, err // Other errors
	}

	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User

	// Query the database for the user by ID
	row := Db.QueryRow("SELECT id, FirstName, Surname, Email, PasswordHash FROM users WHERE Email = ?", strings.TrimSpace(email))

	// Scan the result into the user struct
	err := row.Scan(&user.Id, &user.FirstName, &user.Surname, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Printf("getUserByID: no user found with email %v", email)
			return User{}, nil // No user found
		}
		logger.Printf("getUserByID: %v", err)
		return User{}, err // Other errors
	}

	return user, nil
}
