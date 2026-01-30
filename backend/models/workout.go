package models

import "time"

// Workout represents a workout session
type Workout struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Duration    int       `json:"duration" db:"duration"` // in minutes
	CaloriesBurned float64 `json:"calories_burned" db:"calories_burned"`
	WorkoutDate time.Time `json:"workout_date" db:"workout_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// WorkoutRequest represents workout creation/update data
type WorkoutRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Duration    int    `json:"duration" binding:"required,min=1"`
	CaloriesBurned float64 `json:"calories_burned" binding:"required,min=0"`
	WorkoutDate string `json:"workout_date" binding:"required"`
}

