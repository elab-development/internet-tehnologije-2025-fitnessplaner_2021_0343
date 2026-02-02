package main

import (
	"log"
	"net/http"

	"backend/routes"
	"backend/utils"
)

func main() {
	// Initialize database
	if err := utils.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer utils.CloseDB()

	// Setup routes
	handler := routes.SetupRoutes()

	// Start server
	port := ":8080"
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📡 API endpoints available at http://localhost%s/api", port)
	
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
