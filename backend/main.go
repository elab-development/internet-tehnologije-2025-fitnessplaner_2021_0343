package main

import (
	"log"
	"net/http"

	"backend/routes"
	"backend/utils"
)

func main() {
	// Inicijalizacija baze podataka
	if err := utils.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer utils.CloseDB()

	// Podešavanje ruta
	handler := routes.SetupRoutes()

	// Pokretanje servera
	port := ":8080"
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📡 API endpoints available at http://localhost%s/api", port)
	
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
