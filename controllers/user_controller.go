package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"backend/auth"
	"backend/middleware"
	"backend/models"
	"backend/utils"
)

// Register hendluje registraciju korisnika
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validacija cilja
	if req.Goal != "lose_weight" && req.Goal != "hypertrophy" {
		utils.JSONError(w, "Goal must be 'lose_weight' or 'hypertrophy'", http.StatusBadRequest)
		return
	}

	// Provera da li email vec postoji
	var existingID int
	err := utils.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&existingID)
	if err != sql.ErrNoRows {
		if err == nil {
			utils.JSONError(w, "Email already exists", http.StatusConflict)
			return
		}
		log.Printf("❌ Error checking if user exists: %v", err)
		utils.JSONError(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	// hesiranje lozinke/sifre
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.JSONError(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Set default role - ALWAYS use "user" for new registrations
	// Don't trust frontend to send correct role value - force it to be "user"
	role := "user"

	// Log ako je pokusaj da se setuje druga uloga
	if req.Role != "" {
		log.Printf("📝 Received role from request: '%s' (will use 'user' instead)", req.Role)
	}

	// Verifikacija tipa kolone i default vrednositi u bazi za kolonu role/uloga
	var roleType, roleDefault string
	err = utils.DB.QueryRow(
		"SELECT DATA_TYPE, COLUMN_DEFAULT FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'users' AND column_name = 'role'",
	).Scan(&roleType, &roleDefault)
	if err == nil {
		log.Printf("🔍 Role column type: %s, default: %s", roleType, roleDefault)
	}

	// Log informacija o korsniku koji se registruje
	log.Printf("📝 Inserting user: name='%s', email='%s', goal='%s', role='%s', height=%v, weight=%v",
		req.Name, req.Email, req.Goal, role, req.Height, req.Weight)

	// Insertovanje korisnika u bazu
	var result sql.Result
	if req.Height != nil && req.Weight != nil {
		// Pokusaj insertovanje sa visinom i tezinom
		result, err = utils.DB.Exec(
			"INSERT INTO users (name, email, password, goal, role, height, weight) VALUES (?, ?, ?, ?, ?, ?, ?)",
			req.Name, req.Email, hashedPassword, req.Goal, role, req.Height, req.Weight,
		)
		if err != nil {
			log.Printf("⚠️  Insert with height/weight failed: %v, trying without...", err)
			// Fallnack na insertovanje bez visine i tezine
			result, err = utils.DB.Exec(
				"INSERT INTO users (name, email, password, goal, role) VALUES (?, ?, ?, ?, ?)",
				req.Name, req.Email, hashedPassword, req.Goal, role,
			)
		}
	} else {
		// Insertovanje bez visine i tezine
		result, err = utils.DB.Exec(
			"INSERT INTO users (name, email, password, goal, role) VALUES (?, ?, ?, ?, ?)",
			req.Name, req.Email, hashedPassword, req.Goal, role,
		)
		if err != nil {
			// Ako insert sa ulogom ne uspe pokusaj bez uloge
			log.Printf("⚠️  Insert with role failed: %v, trying without role...", err)
			result, err = utils.DB.Exec(
				"INSERT INTO users (name, email, password, goal) VALUES (?, ?, ?, ?)",
				req.Name, req.Email, hashedPassword, req.Goal,
			)
		}
	}

	if err != nil {
		log.Printf("❌ Error creating user: %v", err)
		utils.JSONError(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()

	// Generisi token
	token, err := auth.GenerateToken(int(userID), req.Email)
	if err != nil {
		utils.JSONError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		ID:     int(userID),
		Name:   req.Name,
		Email:  req.Email,
		Goal:   req.Goal,
		Role:   role,
		Height: req.Height,
		Weight: req.Weight,
	}

	response := models.LoginResponse{
		User:  user,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Login handluje prijavu korisnika
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetchovanje korisnika iz baze
	var user models.User
	var password sql.NullString
	var height, weight sql.NullFloat64
	err := utils.DB.QueryRow(
		"SELECT id, name, email, password, goal, role, height, weight FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &password, &user.Goal, &user.Role, &height, &weight)

	// Convertovanjee sql.NullFloat64 u *float64
	if height.Valid {
		user.Height = &height.Float64
	}
	if weight.Valid {
		user.Weight = &weight.Float64
	}

	if password.Valid && password.String != "" {
		user.Password = password.String
	} else {
		log.Printf("⚠️  User %s has NULL or empty password", req.Email)
		utils.JSONError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err == sql.ErrNoRows {
		utils.JSONError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("❌ Error fetching user: %v", err)
		utils.JSONError(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	// Proveri lozinku/sifru
	if !auth.CheckPassword(req.Password, user.Password) {
		utils.JSONError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generisi token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		utils.JSONError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// ne salji sifru u response
	user.Password = ""

	response := models.LoginResponse{
		User:  &user,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProfile vraca profil autentifikovanog korisnika
func GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == 0 {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	var height, weight sql.NullFloat64
	err := utils.DB.QueryRow(
		"SELECT id, name, email, goal, role, height, weight FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Goal, &user.Role, &height, &weight)

	// pretvara sql.NullFloat64 u *float64
	if height.Valid {
		user.Height = &height.Float64
	}
	if weight.Valid {
		user.Weight = &weight.Float64
	}

	if err == sql.ErrNoRows {
		utils.JSONError(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		utils.JSONError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Logout hendluje odjavu korisnika
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
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
