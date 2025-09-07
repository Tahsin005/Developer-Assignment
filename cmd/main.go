package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/routes"
)

func init() {
	if err := database.ConnDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	database.CreateSystemAdminIfNotExists()
}

func main() {
	// Get config singleton
	cfg := config.GetConfig()

	log.Println("Server initialized")
	defer database.DB.Close()

	router := routes.SetupRoutes()

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		router.ServeHTTP(w, r)
	})

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server is running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, corsHandler))
}
