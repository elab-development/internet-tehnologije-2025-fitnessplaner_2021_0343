package routes

import (
	"net/http"
	"os"
	"path/filepath"

	"backend/controllers"
	"backend/middleware"
)

// SetupRoutes konfiguriše sve rute
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Javne rute
	mux.HandleFunc("/api/register", controllers.Register)
	mux.HandleFunc("/api/login", controllers.Login)

	// Zaštićene rute - Autentifikacija
	mux.Handle("/api/logout", middleware.Auth(http.HandlerFunc(controllers.Logout)))
	mux.Handle("/api/profile", middleware.Auth(http.HandlerFunc(controllers.GetProfile)))

	// Zaštićene rute - Hrana i Meal Plan
	mux.Handle("/api/food/search", middleware.Auth(http.HandlerFunc(controllers.SearchFood)))
	mux.Handle("/api/meal-plan", middleware.Auth(http.HandlerFunc(controllers.GenerateMealPlan)))

	// Zaštićene rute - Treninzi (GET, POST, PUT, DELETE)
	mux.Handle("/api/workouts", middleware.Auth(http.HandlerFunc(controllers.GetWorkouts)))
	mux.Handle("/api/workouts/create", middleware.Auth(http.HandlerFunc(controllers.CreateWorkout)))
	mux.Handle("/api/workouts/update", middleware.Auth(http.HandlerFunc(controllers.UpdateWorkout)))
	mux.Handle("/api/workouts/delete", middleware.Auth(http.HandlerFunc(controllers.DeleteWorkout)))

	// Zaštićene rute - Napredak (GET, POST, PUT, DELETE)
	mux.Handle("/api/progress", middleware.Auth(http.HandlerFunc(controllers.GetProgress)))
	mux.Handle("/api/progress/create", middleware.Auth(http.HandlerFunc(controllers.CreateProgress)))
	mux.Handle("/api/progress/update", middleware.Auth(http.HandlerFunc(controllers.UpdateProgress)))
	mux.Handle("/api/progress/delete", middleware.Auth(http.HandlerFunc(controllers.DeleteProgress)))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// OpenAPI specifikacija za Swagger UI
	mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		// Učitaj OpenAPI fajl sa diska
		openAPIPath := filepath.Join("docs", "openapi.yaml")
		content, err := os.ReadFile(openAPIPath)
		if err != nil {
			http.Error(w, "Failed to load OpenAPI specification", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write(content)
	})

	// Primena middleware-a
	handler := middleware.CORS(mux)
	handler = middleware.Logging(handler)

	return handler
}
