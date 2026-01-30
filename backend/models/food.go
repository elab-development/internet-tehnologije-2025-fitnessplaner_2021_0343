package models

// Food represents a food item
type Food struct {
	ID       int     `json:"id" db:"id"`
	Name     string  `json:"name" db:"name"`
	Barcode  string  `json:"barcode,omitempty" db:"barcode"`
	Calories float64 `json:"calories" db:"calories"`
	Protein  float64 `json:"protein" db:"protein"`
	Carbs    float64 `json:"carbs" db:"carbs"`
	Fat      float64 `json:"fat" db:"fat"`
}

// MealPlan represents a meal plan for a user
type MealPlan struct {
	ID            int     `json:"id" db:"id"`
	UserID        int     `json:"user_id" db:"user_id"`
	Goal          string  `json:"goal" db:"goal"`
	Foods         []Food  `json:"foods"`
	TotalCalories float64 `json:"total_calories"`
	TotalProtein  float64 `json:"total_protein"`
	TotalCarbs    float64 `json:"total_carbs"`
	TotalFat      float64 `json:"total_fat"`
}

// FoodSearchRequest represents a request to search for food
type FoodSearchRequest struct {
	Barcode string `json:"barcode" binding:"required"`
}
