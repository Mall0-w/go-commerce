package main

type User struct {
	ID           string  `json:"id"`
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

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
