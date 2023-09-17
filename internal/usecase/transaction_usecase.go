package usecase

import "time"

type Transaction struct {
	TransactionID int       `json:"transaction_id"`
	PhoneNumber   string    `json:"phoneNumber" binding:"required"`
	Amount        float64   `json:"amount" binding:"required"`
	Timestamp     time.Time `json:"timestamp"`
}

type ChargeCodeTransaction struct {
	TransactionID int       `json:"transaction_id"`
	PhoneNumber   string    `json:"phoneNumber" binding:"required"`
	ChargeCodeID  int       `json:"ChargeCodeID" binding:"required"`
	Timestamp     time.Time `json:"timestamp"`
}

type TransactionRepository interface {
	CreateTransaction(transaction *Transaction) (*Transaction, error)
	CreateChargeTransaction(chargeCodeTransaction *ChargeCodeTransaction) (*ChargeCodeTransaction, error)
	GetTransactions() ([]*Transaction, error)
	GetTransactionByID(id int) (*Transaction, error)
	GetUserTransactionsByUserID(userId int) ([]*Transaction, error)
	GetUserTotalTransaction(userId int) (int, error)
}

type TransactionUseCase struct {
	TransactionRepository TransactionRepository
}

func NewTransactionUseCase(transactionRepo TransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{TransactionRepository: transactionRepo}
}

func (tu *TransactionUseCase) CreateTransaction(transaction *Transaction) (*Transaction, error) {
	return tu.TransactionRepository.CreateTransaction(transaction)
}

func (tu *TransactionUseCase) CreateChargeTransaction(chargeCodeTransaction *ChargeCodeTransaction) (*ChargeCodeTransaction, error) {
	return tu.TransactionRepository.CreateChargeTransaction(chargeCodeTransaction)
}

func (tu *TransactionUseCase) GetTransactions() ([]*Transaction, error) {
	return tu.TransactionRepository.GetTransactions()
}

func (tu *TransactionUseCase) GetTransactionByID(id int) (*Transaction, error) {
	return tu.TransactionRepository.GetTransactionByID(id)
}

func (tu *TransactionUseCase) GetUserTransactionsByUserID(userId int) ([]*Transaction, error) {
	return tu.TransactionRepository.GetUserTransactionsByUserID(userId)
}

func (tu *TransactionUseCase) GetUserTotalTransaction(userId int) (int, error) {
	return tu.TransactionRepository.GetUserTotalTransaction(userId)
}
