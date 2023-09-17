package main

import (
	"chargeCode/internal/database"
	"chargeCode/internal/delivery"
	"chargeCode/internal/repository"
	"chargeCode/internal/usecase"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// // Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		println("Error loading .env file: %v", err)
	}

	//Initial Database
	// Initialize the database connection
	db, err := database.NewDBConnection()
	if err != nil {
		println(err)
	}
	defer db.Close() // Close the database connection when done

	// Initialize dependencies
	//user dependencies
	userRepo := repository.NewUserRepository(db)
	userUC := usecase.NewUserUseCase(userRepo)
	//chargeCode dependencies
	chargeCodeRepo := repository.NewChargeCodeRepository(db)
	chargeCodeUC := usecase.NewChargeCodeUseCase(chargeCodeRepo)

	//transaction dependencies
	transactionRepo := repository.NewTransactionRepository(db)
	transactionUC := usecase.NewTransactionUseCase(transactionRepo)

	// Pass the UserUseCase instance, not a pointer, to SetupRouter
	router := delivery.SetupRouter(userUC, chargeCodeUC, transactionUC) // Pass userUC, not &userUC

	router.Run(":" + os.Getenv("APPLICATION_PORT"))
}
