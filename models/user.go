package models

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Hidden from JSON
	Goal      string    `json:"goal" db:"goal"`  // lose_weight ili hypertrophy
	Role      string    `json:"role" db:"role"`  // admin, user, premium
	Height    *float64  `json:"height,omitempty" db:"height"` // Height in cm
	Weight    *float64  `json:"weight,omitempty" db:"weight"` // Weight in kg
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Goal     string   `json:"goal" binding:"required,oneof=lose_weight hypertrophy"`
	Role     string   `json:"role"` // Optional, defaults to "user"
	Height   *float64 `json:"height,omitempty"` // Optional height in cm
	Weight   *float64 `json:"weight,omitempty"` // Optional weight in kg
}

// LoginRequest represents login data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response with token
type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
