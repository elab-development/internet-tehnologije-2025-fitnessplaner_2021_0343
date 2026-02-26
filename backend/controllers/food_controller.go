package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"backend/middleware"
	"backend/models"
	"backend/utils"
)

// Open Food Facts API response
type OFFProduct struct {
	Product struct {
		ProductName string `json:"product_name"`
		Nutriments  struct {
			EnergyKcal    float64 `json:"energy-kcal_100g"`
			Proteins      float64 `json:"proteins_100g"`
			Carbohydrates float64 `json:"carbohydrates_100g"`
			Fat           float64 `json:"fat_100g"`
		} `json:"nutriments"`
	} `json:"product"`
	Status int `json:"status"`
}

// SearchFood searches for food by barcode using Open Food Facts API
func SearchFood(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.FoodSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call Open Food Facts API
	url := fmt.Sprintf("https://world.openfoodfacts.org/api/v2/product/%s.json", req.Barcode)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch food data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var product OFFProduct
	if err := json.Unmarshal(body, &product); err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	if product.Status == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	food := models.Food{
		Name:     product.Product.ProductName,
		Barcode:  req.Barcode,
		Calories: product.Product.Nutriments.EnergyKcal,
		Protein:  product.Product.Nutriments.Proteins,
		Carbs:    product.Product.Nutriments.Carbohydrates,
		Fat:      product.Product.Nutriments.Fat,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(food)
}

// GenerateMealPlan generates a meal plan based on user's goal
func GenerateMealPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Proveri da li korisnik postoji u bazi
	var userExists int
	err := utils.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&userExists)
	if err != nil {
		log.Printf("❌ Error checking if user exists: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	if userExists == 0 {
		log.Printf("❌ User with ID %d does not exist in database", userID)
		http.Error(w, "User not found. Please log in again.", http.StatusUnauthorized)
		return
	}

	// Get user's goal
	var goal string
	err = utils.DB.QueryRow("SELECT goal FROM users WHERE id = ?", userID).Scan(&goal)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Example barcodes - in production, this would be more sophisticated
	barcodes := []string{"3274080005003", "3017620425035"}
	var foods []models.Food

	for _, barcode := range barcodes {
		url := fmt.Sprintf("https://world.openfoodfacts.org/api/v2/product/%s.json", barcode)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		var product OFFProduct
		if err := json.Unmarshal(body, &product); err != nil {
			continue
		}

		if product.Status == 0 {
			continue
		}

		food := models.Food{
			Name:     product.Product.ProductName,
			Barcode:  barcode,
			Calories: product.Product.Nutriments.EnergyKcal,
			Protein:  product.Product.Nutriments.Proteins,
			Carbs:    product.Product.Nutriments.Carbohydrates,
			Fat:      product.Product.Nutriments.Fat,
		}
		foods = append(foods, food)
	}

	// Calculate totals
	var totalCalories, totalProtein, totalCarbs, totalFat float64
	for _, food := range foods {
		totalCalories += food.Calories
		totalProtein += food.Protein
		totalCarbs += food.Carbs
		totalFat += food.Fat
	}

	mealPlan := models.MealPlan{
		UserID:        userID,
		Goal:          goal,
		Foods:         foods,
		TotalCalories: totalCalories,
		TotalProtein:  totalProtein,
		TotalCarbs:    totalCarbs,
		TotalFat:      totalFat,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mealPlan)
}

