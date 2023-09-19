package main

import (
	"chargeCode/internal/config"
	"chargeCode/internal/database"
	"chargeCode/internal/delivery"
	"chargeCode/internal/repository"
	"chargeCode/internal/usecase"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("If you are loading the .env file within a Docker environment, you can safely ignore this error message: %v", err)
	}

	// Initialize the logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	appConfig, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Error loading app config file: %v", err)
	}

	// Initialize the database connection
	db, err := database.NewDBConnection(appConfig)
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close() // Close the database connection when done

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db, appConfig)
	userUC := usecase.NewUserUseCase(userRepo)

	chargeCodeRepo := repository.NewChargeCodeRepository(db, appConfig)
	chargeCodeUC := usecase.NewChargeCodeUseCase(chargeCodeRepo)

	transactionRepo := repository.NewTransactionRepository(db, appConfig)
	transactionUC := usecase.NewTransactionUseCase(transactionRepo)

	// Pass the UserUseCase instance, not a pointer, to SetupRouter
	router := delivery.SetupRouter(userUC, chargeCodeUC, transactionUC) // Pass userUC, not &userUC

	// Start the server
	logger.Printf("Server started on port %s", appConfig.ApplicationPort)
	router.Run(":" + appConfig.ApplicationPort)
}
