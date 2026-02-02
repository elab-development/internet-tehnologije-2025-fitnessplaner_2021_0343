package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/middleware"
	"backend/models"
	"backend/utils"
)

// ========== WORKOUTS ==========

func GetWorkouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := utils.DB.Query(
		"SELECT id, user_id, name, description, duration, calories_burned, workout_date, created_at, updated_at FROM workouts WHERE user_id = ? ORDER BY workout_date DESC",
		userID,
	)
	if err != nil {
		log.Printf("❌ Error querying workouts: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workouts []models.Workout
	for rows.Next() {
		var workout models.Workout
		var description sql.NullString
		if err := rows.Scan(&workout.ID, &workout.UserID, &workout.Name, &description, &workout.Duration, &workout.CaloriesBurned, &workout.WorkoutDate, &workout.CreatedAt, &workout.UpdatedAt); err != nil {
			log.Printf("❌ Error scanning workout row: %v", err)
			continue
		}
		if description.Valid {
			workout.Description = description.String
		}
		workouts = append(workouts, workout)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Error iterating workout rows: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

func CreateWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.WorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workoutDate, err := time.Parse("2006-01-02", req.WorkoutDate)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	result, err := utils.DB.Exec(
		"INSERT INTO workouts (user_id, name, description, duration, calories_burned, workout_date) VALUES (?, ?, ?, ?, ?, ?)",
		userID, req.Name, req.Description, req.Duration, req.CaloriesBurned, workoutDate,
	)
	if err != nil {
		log.Printf("❌ Error creating workout: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create workout: %v", err), http.StatusInternalServerError)
		return
	}

	workoutID, _ := result.LastInsertId()
	var workout models.Workout
	var description sql.NullString
	err = utils.DB.QueryRow("SELECT id, user_id, name, description, duration, calories_burned, workout_date, created_at, updated_at FROM workouts WHERE id = ?", workoutID).Scan(
		&workout.ID, &workout.UserID, &workout.Name, &description, &workout.Duration, &workout.CaloriesBurned, &workout.WorkoutDate, &workout.CreatedAt, &workout.UpdatedAt,
	)
	if err != nil {
		log.Printf("❌ Error fetching created workout: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch created workout: %v", err), http.StatusInternalServerError)
		return
	}
	if description.Valid {
		workout.Description = description.String
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}

func UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	workoutID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	var ownerID int
	if err := utils.DB.QueryRow("SELECT user_id FROM workouts WHERE id = ?", workoutID).Scan(&ownerID); err == sql.ErrNoRows {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("❌ Error checking workout ownership: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	} else if ownerID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var req models.WorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	workoutDate, err := time.Parse("2006-01-02", req.WorkoutDate)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	_, err = utils.DB.Exec("UPDATE workouts SET name = ?, description = ?, duration = ?, calories_burned = ?, workout_date = ? WHERE id = ?",
		req.Name, req.Description, req.Duration, req.CaloriesBurned, workoutDate, workoutID)
	if err != nil {
		log.Printf("❌ Error updating workout: %v", err)
		http.Error(w, fmt.Sprintf("Failed to update workout: %v", err), http.StatusInternalServerError)
		return
	}

	var workout models.Workout
	var description sql.NullString
	err = utils.DB.QueryRow("SELECT id, user_id, name, description, duration, calories_burned, workout_date, created_at, updated_at FROM workouts WHERE id = ?", workoutID).Scan(
		&workout.ID, &workout.UserID, &workout.Name, &description, &workout.Duration, &workout.CaloriesBurned, &workout.WorkoutDate, &workout.CreatedAt, &workout.UpdatedAt,
	)
	if err != nil {
		log.Printf("❌ Error fetching updated workout: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch updated workout: %v", err), http.StatusInternalServerError)
		return
	}
	if description.Valid {
		workout.Description = description.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

func DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	workoutID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	var ownerID int
	if err := utils.DB.QueryRow("SELECT user_id FROM workouts WHERE id = ?", workoutID).Scan(&ownerID); err == sql.ErrNoRows {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("❌ Error checking workout ownership: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	} else if ownerID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	_, err := utils.DB.Exec("DELETE FROM workouts WHERE id = ?", workoutID)
	if err != nil {
		log.Printf("❌ Error deleting workout: %v", err)
		http.Error(w, fmt.Sprintf("Failed to delete workout: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Workout deleted successfully"})
}

// ========== PROGRESS ==========

func GetProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := utils.DB.Query(
		"SELECT id, user_id, weight, body_fat, muscle_mass, notes, progress_date, created_at, updated_at FROM progress WHERE user_id = ? ORDER BY progress_date DESC",
		userID,
	)
	if err != nil {
		log.Printf("❌ Error querying progress: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var progressList []models.Progress
	for rows.Next() {
		var progress models.Progress
		var bodyFat, muscleMass sql.NullFloat64
		var notes sql.NullString
		if err := rows.Scan(&progress.ID, &progress.UserID, &progress.Weight, &bodyFat, &muscleMass, &notes, &progress.ProgressDate, &progress.CreatedAt, &progress.UpdatedAt); err != nil {
			log.Printf("❌ Error scanning progress row: %v", err)
			continue
		}
		if bodyFat.Valid {
			progress.BodyFat = bodyFat.Float64
		}
		if muscleMass.Valid {
			progress.MuscleMass = muscleMass.Float64
		}
		if notes.Valid {
			progress.Notes = notes.String
		}
		progressList = append(progressList, progress)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Error iterating progress rows: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progressList)
}

func CreateProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	progressDate, err := time.Parse("2006-01-02", req.ProgressDate)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	result, err := utils.DB.Exec(
		"INSERT INTO progress (user_id, weight, body_fat, muscle_mass, notes, progress_date) VALUES (?, ?, ?, ?, ?, ?)",
		userID, req.Weight, req.BodyFat, req.MuscleMass, req.Notes, progressDate,
	)
	if err != nil {
		log.Printf("❌ Error creating progress: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create progress entry: %v", err), http.StatusInternalServerError)
		return
	}

	progressID, _ := result.LastInsertId()
	var progress models.Progress
	var bodyFat, muscleMass sql.NullFloat64
	var notes sql.NullString
	err = utils.DB.QueryRow("SELECT id, user_id, weight, body_fat, muscle_mass, notes, progress_date, created_at, updated_at FROM progress WHERE id = ?", progressID).Scan(
		&progress.ID, &progress.UserID, &progress.Weight, &bodyFat, &muscleMass, &notes, &progress.ProgressDate, &progress.CreatedAt, &progress.UpdatedAt,
	)
	if err != nil {
		log.Printf("❌ Error fetching created progress: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch created progress: %v", err), http.StatusInternalServerError)
		return
	}
	if bodyFat.Valid {
		progress.BodyFat = bodyFat.Float64
	}
	if muscleMass.Valid {
		progress.MuscleMass = muscleMass.Float64
	}
	if notes.Valid {
		progress.Notes = notes.String
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(progress)
}

func UpdateProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	progressID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	var ownerID int
	if err := utils.DB.QueryRow("SELECT user_id FROM progress WHERE id = ?", progressID).Scan(&ownerID); err == sql.ErrNoRows {
		http.Error(w, "Progress entry not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("❌ Error checking progress ownership: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	} else if ownerID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var req models.ProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	progressDate, err := time.Parse("2006-01-02", req.ProgressDate)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	_, err = utils.DB.Exec("UPDATE progress SET weight = ?, body_fat = ?, muscle_mass = ?, notes = ?, progress_date = ? WHERE id = ?",
		req.Weight, req.BodyFat, req.MuscleMass, req.Notes, progressDate, progressID)
	if err != nil {
		log.Printf("❌ Error updating progress: %v", err)
		http.Error(w, fmt.Sprintf("Failed to update progress: %v", err), http.StatusInternalServerError)
		return
	}

	var progress models.Progress
	var bodyFat, muscleMass sql.NullFloat64
	var notes sql.NullString
	err = utils.DB.QueryRow("SELECT id, user_id, weight, body_fat, muscle_mass, notes, progress_date, created_at, updated_at FROM progress WHERE id = ?", progressID).Scan(
		&progress.ID, &progress.UserID, &progress.Weight, &bodyFat, &muscleMass, &notes, &progress.ProgressDate, &progress.CreatedAt, &progress.UpdatedAt,
	)
	if err != nil {
		log.Printf("❌ Error fetching updated progress: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch updated progress: %v", err), http.StatusInternalServerError)
		return
	}
	if bodyFat.Valid {
		progress.BodyFat = bodyFat.Float64
	}
	if muscleMass.Valid {
		progress.MuscleMass = muscleMass.Float64
	}
	if notes.Valid {
		progress.Notes = notes.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

func DeleteProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	progressID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	var ownerID int
	if err := utils.DB.QueryRow("SELECT user_id FROM progress WHERE id = ?", progressID).Scan(&ownerID); err == sql.ErrNoRows {
		http.Error(w, "Progress entry not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("❌ Error checking progress ownership: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	} else if ownerID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	_, err := utils.DB.Exec("DELETE FROM progress WHERE id = ?", progressID)
	if err != nil {
		log.Printf("❌ Error deleting progress: %v", err)
		http.Error(w, fmt.Sprintf("Failed to delete progress: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Progress entry deleted successfully"})
}
