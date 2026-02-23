package models

import "time"

<<<<<<< HEAD
// User represents a user in the system
=======
// User predstavlja jednog korisnika
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
<<<<<<< HEAD
	Password  string    `json:"-" db:"password"` // Hidden from JSON
	Goal      string    `json:"goal" db:"goal"`  // lose_weight ili hypertrophy
	Role      string    `json:"role" db:"role"`  // admin, user, premium
	Height    *float64  `json:"height,omitempty" db:"height"` // Height in cm
	Weight    *float64  `json:"weight,omitempty" db:"weight"` // Weight in kg
=======
	Password  string    `json:"-" db:"password"` // sakriveno od jsona
	Goal      string    `json:"goal" db:"goal"`  // lose_weight ili hypertrophy
	Role      string    `json:"role" db:"role"`  // admin, user, premium
	Height    *float64  `json:"height,omitempty" db:"height"` // u cm
	Weight    *float64  `json:"weight,omitempty" db:"weight"` // u kg
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

<<<<<<< HEAD
// RegisterRequest represents registration data
=======
// RegisterRequest predstavlja podatke za registraciju
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
type RegisterRequest struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Goal     string   `json:"goal" binding:"required,oneof=lose_weight hypertrophy"`
<<<<<<< HEAD
	Role     string   `json:"role"` // Optional, defaults to "user"
	Height   *float64 `json:"height,omitempty"` // Optional height in cm
	Weight   *float64 `json:"weight,omitempty"` // Optional weight in kg
}

// LoginRequest represents login data
=======
	Role     string   `json:"role"` // opciono, default je user
	Height   *float64 `json:"height,omitempty"`
	Weight   *float64 `json:"weight,omitempty"` 
}

// LoginRequest predstavlja podatke za login
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

<<<<<<< HEAD
// LoginResponse represents login response with token
=======
// LoginResponse predstavlja login response sa tokenom
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
