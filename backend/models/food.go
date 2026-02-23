package models

<<<<<<< HEAD
// Food represents a food item
=======
// Food predstavlja jednu namirnicu
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
type Food struct {
	ID       int     `json:"id" db:"id"`
	Name     string  `json:"name" db:"name"`
	Barcode  string  `json:"barcode,omitempty" db:"barcode"`
	Calories float64 `json:"calories" db:"calories"`
	Protein  float64 `json:"protein" db:"protein"`
	Carbs    float64 `json:"carbs" db:"carbs"`
	Fat      float64 `json:"fat" db:"fat"`
}

<<<<<<< HEAD
// MealPlan represents a meal plan for a user
=======
// MealPlan predstavlja meal plan za korisnika
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
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

<<<<<<< HEAD
// FoodSearchRequest represents a request to search for food
=======
// FoodSearchRequest predstavlja zahtev za pretragu hrane
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
type FoodSearchRequest struct {
	Barcode string `json:"barcode" binding:"required"`
}
