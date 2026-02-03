package models

import "time"

// User predstavlja korisnika u sistemu
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Sakriveno od JSONa
	Goal      string    `json:"goal" db:"goal"`  // lose_weight ili hypertrophy
	Role      string    `json:"role" db:"role"`  // admin, user, premium
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest predstavlja podatke za registraciju
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Goal     string `json:"goal" binding:"required,oneof=lose_weight hypertrophy"`
	Role     string `json:"role"` // Opciono, default "user"
}

// LoginRequest predstavlja podatke za login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse predstavlja odgovor na prijavu koja sadrzi token
type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
}
