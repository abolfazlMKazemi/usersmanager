package delivery

import (
	"chargeCode/internal/usecase"

	// docs "chargeCode/cmd/docs"

	"github.com/gin-gonic/gin"
	// swaggerfiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"

	_ "chargeCode/cmd/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(userUC *usecase.UserUseCase, chargeCodeUC *usecase.ChargeCodeUseCase, transactionUC *usecase.TransactionUseCase) *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	userHandler := NewUserHandler(userUC)
	ChargeCodeHandler := NewChargeCodeHandler(chargeCodeUC)
	transactionandler := NewTransactionHandler(transactionUC)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// // Specify the Swagger JSON file path
	// docs.SwaggerInfo.BasePath = ""

	user := router.Group("/api/v1/user")
	{
		user.GET("/:phoneNumber", userHandler.GetUserByPhoneNumber)
		user.GET("/chargeCode/:chargeCodeId", userHandler.ListOfUsersUseChargeCode)
		user.GET("/balance/:userId", userHandler.GetUserBalance)
		user.PUT("/", userHandler.UpdateUser)
	}

	chargeCode := router.Group("/api/v1/chargeCode")
	{
		chargeCode.POST("/", ChargeCodeHandler.CreateChargeCode)
		chargeCode.GET("", ChargeCodeHandler.GetChargeCodes)
		chargeCode.GET("/:id", ChargeCodeHandler.GetChargeCodeByID)
		chargeCode.GET("/code/:code", ChargeCodeHandler.GetChargeCodeByCode)
		chargeCode.GET("/user/:userId", ChargeCodeHandler.GetUserChargeCodes)
		chargeCode.DELETE("/:id", ChargeCodeHandler.DeleteChargeCodeByID)
		chargeCode.PUT("/", ChargeCodeHandler.UpdateChargeCode)
	}

	transaction := router.Group("/api/v1/transaction")
	{
		transaction.POST("/", transactionandler.CreateTransaction)
		transaction.POST("/charge", transactionandler.CreateChargeTransaction)
		transaction.GET("", transactionandler.GetTransactions)
		transaction.GET(":id", transactionandler.GetTransactionByID)
		transaction.GET("user/:userId", transactionandler.GetUserTransactionsByUserID)
		transaction.GET("user/totalNumber/:userId", transactionandler.GetUserTotalTransaction)

	}

	router.Run(":8080")
	return router
}
