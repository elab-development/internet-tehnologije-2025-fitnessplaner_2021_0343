package routes

import (
	"net/http"

	"app/controllers"
	"app/middleware"
)

// SetupRoutes configures all routes
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/register", controllers.Register)
	mux.HandleFunc("/api/login", controllers.Login)

	// Protected routes - Auth
	mux.Handle("/api/logout", middleware.Auth(http.HandlerFunc(controllers.Logout)))
	mux.Handle("/api/profile", middleware.Auth(http.HandlerFunc(controllers.GetProfile)))

	// Protected routes - Food & Meal Plan
	mux.Handle("/api/food/search", middleware.Auth(http.HandlerFunc(controllers.SearchFood)))
	mux.Handle("/api/meal-plan", middleware.Auth(http.HandlerFunc(controllers.GenerateMealPlan)))

	// Protected routes - Workouts (GET, POST, PUT, DELETE)
	mux.Handle("/api/workouts", middleware.Auth(http.HandlerFunc(controllers.GetWorkouts)))
	mux.Handle("/api/workouts/create", middleware.Auth(http.HandlerFunc(controllers.CreateWorkout)))
	mux.Handle("/api/workouts/update", middleware.Auth(http.HandlerFunc(controllers.UpdateWorkout)))
	mux.Handle("/api/workouts/delete", middleware.Auth(http.HandlerFunc(controllers.DeleteWorkout)))

	// Protected routes - Progress (GET, POST, PUT, DELETE)
	mux.Handle("/api/progress", middleware.Auth(http.HandlerFunc(controllers.GetProgress)))
	mux.Handle("/api/progress/create", middleware.Auth(http.HandlerFunc(controllers.CreateProgress)))
	mux.Handle("/api/progress/update", middleware.Auth(http.HandlerFunc(controllers.UpdateProgress)))
	mux.Handle("/api/progress/delete", middleware.Auth(http.HandlerFunc(controllers.DeleteProgress)))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware
	handler := middleware.CORS(mux)
	handler = middleware.Logging(handler)

	return handler
}
