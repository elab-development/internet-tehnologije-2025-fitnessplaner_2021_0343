package models

import "time"

// User predstavlja korisnika u sistemu
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Sakriveno od JSON-a
	Goal      string    `json:"goal" db:"goal"`  // lose_weight ili hypertrophy
	Role      string    `json:"role" db:"role"`  // admin, user, premium
	Height    *float64  `json:"height,omitempty" db:"height"` // Visina u cm
	Weight    *float64  `json:"weight,omitempty" db:"weight"` // Težina u kg
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest predstavlja podatke za registraciju
type RegisterRequest struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Goal     string   `json:"goal" binding:"required,oneof=lose_weight hypertrophy"`
	Role     string   `json:"role"` // Opciono, podrazumevano "user"
	Height   *float64 `json:"height,omitempty"` // Opciona visina u cm
	Weight   *float64 `json:"weight,omitempty"` // Opciona težina u kg
}

// LoginRequest predstavlja podatke za prijavu
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse predstavlja odgovor na prijavu sa tokenom
type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
