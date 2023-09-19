package repository

// Import the errors package
import (
	"chargeCode/internal/config"
	"chargeCode/internal/usecase"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type TransactionRepository struct {
	// Implement data storage and retrieval methods here
	db     *sql.DB
	config *config.AppConfig
}

func NewTransactionRepository(db *sql.DB, config *config.AppConfig) *TransactionRepository {
	// Initialize the database connection
	return &TransactionRepository{db: db, config: config}
}

func (tr *TransactionRepository) CreateTransaction(transaction *usecase.Transaction) (*usecase.Transaction, error) {

	if !(transaction.Amount >= tr.config.MinTransactionAmount && transaction.Amount <= tr.config.MaxTransactionAmount) {
		return nil, errors.New("amount is outside the valid range")
	}

	userRepository := NewUserRepository(tr.db, tr.config)
	currentUser, err := userRepository.GetUserByPhoneNumber(transaction.PhoneNumber)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if transaction.Amount < 0 && (currentUser.Balance+transaction.Amount) < 0 {
		return nil, errors.New("transaction failed. insufficient funds")
	}

	// Insert the new transaction into the 'transaction' table
	_, err = tr.db.Exec("INSERT INTO transaction (user_id, amount) VALUES (?, ?)", currentUser.ID, transaction.Amount)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("database insert error")
	}

	// If the creation is successful, return the created transaction and no error
	return transaction, nil
}

func (tr *TransactionRepository) CreateChargeTransaction(chargeCodeTransaction *usecase.ChargeCodeTransaction) (*usecase.ChargeCodeTransaction, error) {

	formatPattern := `^09\d{9}$`

	// Compile the regular expression
	re := regexp.MustCompile(formatPattern)
	if !re.MatchString(chargeCodeTransaction.PhoneNumber) {
		return nil, errors.New("phone number is incorrect Most Like 09121114323")
	}
	// Ensure the database connection is valid
	if err := tr.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}

	_, err := tr.db.Exec(`
	INSERT INTO user (phoneNumber, balance)
	VALUES (?, ?)
`, chargeCodeTransaction.PhoneNumber, 0)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return nil, errors.New("user with the same phone number already exists")
		}
		fmt.Println(err)
		return nil, errors.New("database insert error")
	}

	userRepository := NewUserRepository(tr.db, tr.config)
	currentUser, err := userRepository.GetUserByPhoneNumber(chargeCodeTransaction.PhoneNumber)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}

	var userAlreadyRedeemed int

	// SQL query to count the rows in user_charge_code
	queryuserAlreadyRedeemed := "SELECT COUNT(*) FROM user_charge_code WHERE user_id = ? AND charge_code_id = ?"
	err = tr.db.QueryRow(queryuserAlreadyRedeemed, currentUser.ID, chargeCodeTransaction.ChargeCodeID).Scan(&userAlreadyRedeemed)

	if err != nil {
		fmt.Println("Error executing queryCheckUseChargeCode:", err)
		return nil, errors.New(err.Error())
	}

	if userAlreadyRedeemed > 0 {
		return nil, errors.New("user has already redeemed this charge_code")
	}

	query := `
	SELECT charge_code_id, max_uses, current_uses, amount
	FROM charge_code
	WHERE charge_code_id = ? AND current_uses < max_uses;
`

	var chargeCodeID, maxUses, currentUses int
	var amount float64

	err = tr.db.QueryRow(query, chargeCodeTransaction.ChargeCodeID).Scan(&chargeCodeID, &maxUses, &currentUses, &amount)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}

	// Call the RedeemChargeCode stored procedure
	_, err = tr.db.Exec("CALL RedeemChargeCode(?, ?)", currentUser.ID, chargeCodeID)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("stored procedure error")
	}
	// If the creation is successful, return the created ChargeCodeTransaction and no error
	return chargeCodeTransaction, nil
}

func (tr *TransactionRepository) GetTransactions(page int, pageSize int) ([]*usecase.Transaction, error) {
	if page > tr.config.MaxPage {
		return nil, errors.New("page exceeds the maximum allowed limit")
	}

	if pageSize > tr.config.MaxPageSize {
		return nil, errors.New("page size exceeds the maximum allowed limit")
	}
	// Ensure the database connection is valid
	if err := tr.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}
	// Calculate the OFFSET based on the page number and page size
	offset := (page - 1) * pageSize

	// Query transactions with pagination
	query := `
		SELECT t.transaction_id, u.phoneNumber, t.amount, t.timestamp
		FROM transaction t
		INNER JOIN user u ON t.user_id = u.user_id
		LIMIT ? OFFSET ?
	`

	rows, err := tr.db.Query(query, pageSize, offset)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("database query error")
	}

	defer rows.Close()

	transactions := []*usecase.Transaction{}

	// Iterate through the result rows
	for rows.Next() {
		var transactionID int
		var phoneNumber string
		var amount float64
		var timestamp string

		if err := rows.Scan(&transactionID, &phoneNumber, &amount, &timestamp); err != nil {
			fmt.Println(err)
			return nil, errors.New("database scan error")
		}

		// Adding a new User object to the slice
		// Parse the string into a time.Time value
		format := "2006-01-02 15:04:05"
		parsedTime, err := time.Parse(format, timestamp)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return nil, errors.New("time parse error")
		}
		newTransaction := &usecase.Transaction{TransactionID: transactionID, PhoneNumber: phoneNumber, Amount: amount, Timestamp: parsedTime}
		transactions = append(transactions, newTransaction)
		// Process the retrieved data here
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, errors.New("database rows error")
	}

	// Return the list of transactions and no error
	if len(transactions) > 0 {
		return transactions, nil // Success case, return charge codes and no error
	}
	return nil, errors.New("transaction not found")

}

func (tr *TransactionRepository) GetTransactionByID(id int) (*usecase.Transaction, error) {

	// Ensure the database connection is valid
	if err := tr.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}
	// Query all transactions from the 'transaction' table

	query := `
    SELECT t.transaction_id, u.phoneNumber, t.amount, t.timestamp
    FROM transaction t
    INNER JOIN user u ON t.user_id = u.user_id
    WHERE t.transaction_id = ?
`

	var (
		transactionID int
		phoneNumber   string
		amount        float64
		timestamp     string
	)

	err := tr.db.QueryRow(query, id).Scan(&transactionID, &phoneNumber, &amount, &timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("transaction not found")
		} else {
			fmt.Println(err)
			return nil, errors.New("database query error")
		}
	} else {
		format := "2006-01-02 15:04:05"
		parsedTime, err := time.Parse(format, timestamp)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return nil, errors.New("time parse error")
		}

		transaction := &usecase.Transaction{
			TransactionID: transactionID,
			PhoneNumber:   phoneNumber,
			Amount:        amount,
			Timestamp:     parsedTime,
		}

		return transaction, nil // Success case, return user and no error

	}
}

func (tr *TransactionRepository) GetUserTransactionsByUserID(id int, page int, pageSize int) ([]*usecase.Transaction, error) {

	if page > tr.config.MaxPage {
		return nil, errors.New("page exceeds the maximum allowed limit")
	}

	if pageSize > tr.config.MaxPageSize {
		return nil, errors.New("page size exceeds the maximum allowed limit")
	}

	// Ensure the database connection is valid
	if err := tr.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}

	// Calculate the OFFSET based on the page number and page size
	offset := (page - 1) * pageSize

	// Query all transactions from the 'transaction' table
	query := `
		SELECT t.transaction_id, u.phoneNumber, t.amount, t.timestamp
		FROM transaction t
		INNER JOIN user u ON t.user_id = u.user_id
		WHERE u.user_id = ?
		LIMIT ? OFFSET ?
	`

	rows, err := tr.db.Query(query, id, pageSize, offset)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("database query error")
	}

	defer rows.Close()

	transactions := []*usecase.Transaction{}

	// Iterate through the result rows
	for rows.Next() {
		var transactionID int
		var phoneNumber string
		var amount float64
		var timestamp string

		if err := rows.Scan(&transactionID, &phoneNumber, &amount, &timestamp); err != nil {
			fmt.Println(err)
			return nil, errors.New("data scan error")
		}

		// Adding a new User object to the slice
		// Parse the string into a time.Time value
		format := "2006-01-02 15:04:05"
		parsedTime, err := time.Parse(format, timestamp)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return nil, errors.New("time parse error")
		}
		newTransaction := &usecase.Transaction{TransactionID: transactionID, PhoneNumber: phoneNumber, Amount: amount, Timestamp: parsedTime}
		transactions = append(transactions, newTransaction)
		// Process the retrieved data here
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, errors.New("row iteration error")
	}

	// Return the list of transactions and no error
	if len(transactions) > 0 {
		return transactions, nil // Success case, return charge codes and no error
	}
	return nil, errors.New("transaction not found")

}

func (tr *TransactionRepository) GetUserTotalTransaction(userId int) (int, error) {

	// Ensure the database connection is valid
	if err := tr.db.Ping(); err != nil {
		fmt.Println(err)
		return 0, errors.New("internal Server Error")
	}

	println(userId)

	// Query the total transaction count for the user
	var totalTransactionCount int
	err := tr.db.QueryRow(`
    SELECT COUNT(*) FROM transaction t
    INNER JOIN user u ON t.user_id = u.user_id
    WHERE u.user_id = ?
`, userId).Scan(&totalTransactionCount)

	println(totalTransactionCount)
	if err != nil {
		fmt.Printf("Error querying total transaction count: %v", err)
		return 0, errors.New("database query error")
	}
	return totalTransactionCount, nil
}
