package main

import (
    "log"
    "os"

    "github.com/sagorsarker04/Developer-Assignment/internal/database"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Connect to the database
    database.Connect()
    defer database.Close()

    log.Println("App is running on port", os.Getenv("SERVER_PORT"))
}
