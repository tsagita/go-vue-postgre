package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/tsagita/go-vue-postgre/pkg/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	// Get the PORT from the environment variable, default to 3000 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r := handlers.SetupRoutes()
	r.Run(":" + port)
}
