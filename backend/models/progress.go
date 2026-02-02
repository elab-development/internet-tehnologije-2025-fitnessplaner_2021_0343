package models

import "time"

// Progress predstavlja pracenje napretka korisnika
type Progress struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Weight    float64   `json:"weight" db:"weight"` // in kg
	BodyFat   float64   `json:"body_fat" db:"body_fat"` // percentage
	MuscleMass float64  `json:"muscle_mass" db:"muscle_mass"` // in kg
	Notes     string    `json:"notes" db:"notes"`
	ProgressDate time.Time `json:"progress_date" db:"progress_date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ProgressRequest predstavlja podatke za kreiranje i azuriranje progresa
type ProgressRequest struct {
	Weight    float64 `json:"weight" binding:"required,min=0"`
	BodyFat   float64 `json:"body_fat" binding:"min=0,max=100"`
	MuscleMass float64 `json:"muscle_mass" binding:"min=0"`
	Notes     string  `json:"notes"`
	ProgressDate string `json:"progress_date" binding:"required"`
}
