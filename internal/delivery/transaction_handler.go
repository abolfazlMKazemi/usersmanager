// internal/delivery/transaction_handler.go
package delivery

import (
	"chargeCode/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Transaction struct {
	//TransactionID int     `json:"transaction_id"`
	PhoneNumber string  `json:"phoneNumber" binding:"required"`
	Amount      float64 `json:"amount"`
}

type ChargeCodeTransaction struct {
	//TransactionID int `json:"transaction_id"`
	PhoneNumber  string `json:"phoneNumber" binding:"required"`
	ChargeCodeID int    `json:"ChargeCodeID" binding:"required"`
}

type TransactionHandler struct {
	TransactionUseCase *usecase.TransactionUseCase `json:"TransactionUseCase"`
}

func NewTransactionHandler(transactionUC *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{TransactionUseCase: transactionUC}

}

// CreateTransaction godoc
// @Summary Create a new Transaction
// @Description Create a new Transaction using the provided data.
// @Tags Transaction
// @Accept json
// @Produce json
// @Param Transaction body Transaction true "Transaction object to create"
// @Success 200 {object} Transaction
// @Router /api/v1/transaction [post]
func (tH *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction usecase.Transaction

	// Parse the request body into a ChargeCode struct
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// At this point, chargeCode contains the data from the request body
	// You can use it as needed, such as passing it to your use case for creation

	_, err := tH.TransactionUseCase.CreateTransaction(&transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//c.JSON(http.StatusOK, createdTransaction)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreateChargeTransaction godoc
// @Summary Create a new ChargeCodeTransaction
// @Description Create a new ChargeCodeTransaction using the provided data.
// @Tags Transaction
// @Accept json
// @Produce json
// @Param ChargeCodeTransaction body ChargeCodeTransaction true "ChargeCodeTransaction object to create"
// @Success 200 {object} ChargeCodeTransaction
// @Router /api/v1/transaction/charge [post]
func (tH *TransactionHandler) CreateChargeTransaction(c *gin.Context) {
	var chargeCodeTransaction usecase.ChargeCodeTransaction

	// Parse the request body into a ChargeCode struct
	if err := c.ShouldBindJSON(&chargeCodeTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// At this point, chargeCode contains the data from the request body
	// You can use it as needed, such as passing it to your use case for creation

	_, err := tH.TransactionUseCase.CreateChargeTransaction(&chargeCodeTransaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetTransactions godoc
// @Summary Get transactions with pagination
// @Description Get transactions with pagination.
// @Tags Transaction
// @ID get-paginated-transactions
// @Produce json
// @Param page query integer false "Page number most start from 1"
// @Param pageSize query integer false "Number of items per page"
// @Success 200 {object} Transaction
// @Router /api/v1/transaction [get]
func (cH *TransactionHandler) GetTransactions(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "parsing error"})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "parsing error"})
		return
	}
	transactions, err := cH.TransactionUseCase.GetTransactions(page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// GetTransactionByID godoc
// @Summary Get Transaction by ID
// @Description Get a Transaction by their unique ID.
// @Tags Transaction
// @ID get-transaction-by-id
// @Produce json
// @Param id path int true "transaction ID" Example: 123
// @Success 200 {object} Transaction
// @Router /api/v1/transaction/{id} [get]
func (tH *TransactionHandler) GetTransactionByID(c *gin.Context) {
	transactionID, _ := strconv.Atoi(c.Param("id"))
	println(transactionID)
	transaction, err := tH.TransactionUseCase.GetTransactionByID(transactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

// GetUserTransactionsByUserID godoc
// @Summary Get Transactions by user ID
// @Description Get transactions for a user by their unique user ID with pagination.
// @Tags Transaction
// @ID get-transactions-by-user-id
// @Produce json
// @Param userId path int true "User ID" Example: 123
// @Param page query int false "Page number most start from 1"
// @Param pageSize query int false "Number of items per page"
// @Success 200 {object} Transaction
// @Router /api/v1/transaction/user/{userId} [get]
func (tH *TransactionHandler) GetUserTransactionsByUserID(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userId"))
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "parsing error"})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "parsing error"})
		return
	}
	transactions, err := tH.TransactionUseCase.GetUserTransactionsByUserID(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// GetUserTotalTransaction godoc
// @Summary Get Total Transaction by user ID
// @Description Get Total a Transaction by their unique user ID.
// @Tags Transaction
// @ID get-total-transaction-by-user-id
// @Produce json
// @Param userId path int true "user ID" Example: 123
// @Success 200 {object} Transaction
// @Router /api/v1/transaction/user/totalNumber/{userId} [get]
func (tH *TransactionHandler) GetUserTotalTransaction(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userId"))
	num, err := tH.TransactionUseCase.GetUserTotalTransaction(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, num)
}
