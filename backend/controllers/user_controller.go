package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"app/auth"
	"app/middleware"
	"app/models"
	"app/utils"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate goal
	if req.Goal != "lose_weight" && req.Goal != "hypertrophy" {
		http.Error(w, "Goal must be 'lose_weight' or 'hypertrophy'", http.StatusBadRequest)
		return
	}

	// Check if user exists
	var existingID int
	err := utils.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&existingID)
	if err != sql.ErrNoRows {
		if err == nil {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "user"
	}
	if role != "admin" && role != "user" && role != "premium" {
		role = "user"
	}

	// Insert user
	result, err := utils.DB.Exec(
		"INSERT INTO users (name, email, password, goal, role) VALUES (?, ?, ?, ?, ?)",
		req.Name, req.Email, hashedPassword, req.Goal, role,
	)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()

	// Generate token
	token, err := auth.GenerateToken(int(userID), req.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		ID:    int(userID),
		Name:  req.Name,
		Email: req.Email,
		Goal:  req.Goal,
		Role:  role,
	}

	response := models.LoginResponse{
		User:  user,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from database
	var user models.User
	err := utils.DB.QueryRow(
		"SELECT id, name, email, password, goal, role FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Goal, &user.Role)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Don't send password
	user.Password = ""

	response := models.LoginResponse{
		User:  &user,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProfile returns the current user's profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	err := utils.DB.QueryRow(
		"SELECT id, name, email, goal, role FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Goal, &user.Role)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Logout handles user logout (client-side token removal, but we can log it)
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// In a stateless JWT system, logout is handled client-side
	// But we can return success to confirm
	response := map[string]string{
		"message": "Logged out successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
