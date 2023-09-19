package repository

// Import the errors package
import (
	"chargeCode/internal/config"
	"chargeCode/internal/usecase"
	"database/sql"
	"errors"
	"fmt"
)

type ChargeCodeRepository struct {
	// Implement data storage and retrieval methods here
	db     *sql.DB
	config *config.AppConfig
}

func NewChargeCodeRepository(db *sql.DB, config *config.AppConfig) *ChargeCodeRepository {
	// Initialize the database connection

	return &ChargeCodeRepository{db: db, config: config}
}

func (cu *ChargeCodeRepository) GetChargeCodes(page int, pageSize int) ([]*usecase.ChargeCode, error) {

	if page > cu.config.MaxPage {
		return nil, errors.New("page exceeds the maximum allowed limit")
	}

	if pageSize > cu.config.MaxPageSize {
		return nil, errors.New("page size exceeds the maximum allowed limit")
	}

	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}

	// Calculate the OFFSET based on the page number and page size
	offset := (page - 1) * pageSize
	// Query all transactions from the 'transaction' table
	query := `
        SELECT charge_code_id, code, max_uses, current_uses, amount
        FROM charge_code
        LIMIT ? OFFSET ?
    `

	rows, err := cu.db.Query(query, pageSize, offset)
	if err != nil {
		fmt.Printf("Error querying charge codes: %v", err)
		return nil, errors.New("database query error")
	}

	defer rows.Close()

	ChargeCodes := []*usecase.ChargeCode{}

	// Iterate through the result rows
	for rows.Next() {
		var ChargeCodeID int
		var Code string
		var MaxUses int
		var CurrentUses int
		var Amount float64

		if err := rows.Scan(&ChargeCodeID, &Code, &MaxUses, &CurrentUses, &Amount); err != nil {
			fmt.Printf("Error scanning charge code row: %v", err)
			return nil, errors.New("database query error")
		}

		// Adding a new User object to the slice

		newChargeCode := &usecase.ChargeCode{ChargeCodeID: ChargeCodeID, Code: Code, MaxUses: MaxUses, CurrentUses: CurrentUses, Amount: Amount}
		ChargeCodes = append(ChargeCodes, newChargeCode)
		// Process the retrieved data here
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	// Return the list of transactions and no error

	if len(ChargeCodes) > 0 {
		return ChargeCodes, nil // Success case, return charge codes and no error
	}
	return nil, errors.New("charge codes not found")

	// Replace someError with an actual error value
}

func (cu *ChargeCodeRepository) GetChargeCodeByID(id int) (*usecase.ChargeCode, error) {

	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}
	// Query all transactions from the 'transaction' table

	query := `
	SELECT charge_code_id, code, max_uses, current_uses, amount
	FROM charge_code
	WHERE charge_code_id = ?
	`

	var (
		ChargeCodeID int
		Code         string
		MaxUses      int
		CurrentUses  int
		Amount       float64
	)

	err := cu.db.QueryRow(query, id).Scan(&ChargeCodeID, &Code, &MaxUses, &CurrentUses, &Amount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("charge code not found")
		} else {
			fmt.Printf("Error querying charge code by ID: %v", err)
			return nil, errors.New("database query error")
		}
	} else {

		newChargeCode := &usecase.ChargeCode{ChargeCodeID: ChargeCodeID, Code: Code, MaxUses: MaxUses, CurrentUses: CurrentUses, Amount: Amount}

		return newChargeCode, nil // Success case, return user and no error

	}
}

func (cu *ChargeCodeRepository) GetChargeCodeByCode(code string) (*usecase.ChargeCode, error) {
	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}
	// Query all transactions from the 'transaction' table

	query := `
	SELECT charge_code_id, code, max_uses, current_uses, amount
	FROM charge_code
	WHERE code = ?
	`

	var (
		ChargeCodeID int
		Code         string
		MaxUses      int
		CurrentUses  int
		Amount       float64
	)

	err := cu.db.QueryRow(query, code).Scan(&ChargeCodeID, &Code, &MaxUses, &CurrentUses, &Amount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("charge code not found")
		} else {
			fmt.Printf("Error querying charge code by code: %v", err)
			return nil, errors.New("database query error")
		}
	} else {

		newChargeCode := &usecase.ChargeCode{ChargeCodeID: ChargeCodeID, Code: Code, MaxUses: MaxUses, CurrentUses: CurrentUses, Amount: Amount}

		return newChargeCode, nil // Success case, return user and no error

	}
}

func (cu *ChargeCodeRepository) CreateChargeCode(chargeCode *usecase.ChargeCode) (*usecase.ChargeCode, error) {

	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}

	if chargeCode.Amount > cu.config.MaxChargeCodeAmount {
		return nil, errors.New("amount is very big")
	}

	if chargeCode.Amount < cu.config.MinChargeCodeAmount {
		return nil, errors.New("amount is very small")
	}

	// Insert the new charge code into the 'charge_code' table
	_, err := cu.db.Exec(`
	   INSERT INTO charge_code (code, max_uses, amount)
	   VALUES (?, ?, ?)
   `, chargeCode.Code, chargeCode.MaxUses, chargeCode.Amount)
	if err != nil {
		fmt.Printf("Error creating charge code: %v", err)
		return nil, errors.New("database error")
	}

	return chargeCode, nil
}

func (cu *ChargeCodeRepository) DeleteChargeCode(id int) error {

	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return errors.New("internal Server Error")
	}
	// Delete the charge code by ID from the 'charge_code' table
	_, err := cu.db.Exec(`
   DELETE FROM charge_code
   WHERE charge_code_id = ?
`, id)
	if err != nil {
		fmt.Printf("Error deleting charge code: %v", err)
		return errors.New("database error")
	}
	return nil
}

func (cu *ChargeCodeRepository) UpdateChargeCode(chargeCode *usecase.ChargeCode) (*usecase.ChargeCode, error) {
	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}

	// Update the charge code by ID in the 'charge_code' table
	_, err := cu.db.Exec(`
	   UPDATE charge_code
	   SET code = ?, max_uses = ?, amount = ?, current_uses = ?
	   WHERE charge_code_id = ?
   `, chargeCode.Code, chargeCode.MaxUses, chargeCode.Amount, chargeCode.CurrentUses, chargeCode.ChargeCodeID)
	if err != nil {
		fmt.Printf("Error deleting charge code: %v", err)
		return nil, errors.New("database error")
	}

	return chargeCode, nil
}

func (cu *ChargeCodeRepository) GetUserChargeCodes(userId int, page int, pageSize int) ([]*usecase.ChargeCode, error) {

	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}
	if page > cu.config.MaxPage {
		return nil, errors.New("page exceeds the maximum allowed limit")
	}

	if pageSize > cu.config.MaxPageSize {
		return nil, errors.New("page size exceeds the maximum allowed limit")
	}

	// Calculate the offset based on the page and pageSize
	offset := (page - 1) * pageSize

	// Query the charge codes by user ID with pagination from the 'user_charge_code' table
	rows, err := cu.db.Query(`
        SELECT uc.charge_code_id, cc.code, cc.max_uses, cc.current_uses, cc.amount
        FROM user_charge_code uc
        INNER JOIN charge_code cc ON uc.charge_code_id = cc.charge_code_id
        WHERE uc.user_id = ?
        LIMIT ? OFFSET ?
    `, userId, pageSize, offset)

	if err != nil {
		fmt.Printf("Error querying user charge codes: %v", err)
		return nil, errors.New("database error")
	}
	defer rows.Close()

	ChargeCodes := []*usecase.ChargeCode{}

	// Iterate through the result rows
	for rows.Next() {
		var chargeCodeID int
		var code string
		var maxUses, currentUses int
		var amount float64

		if err := rows.Scan(&chargeCodeID, &code, &maxUses, &currentUses, &amount); err != nil {
			fmt.Printf("Error scanning charge code row: %v", err)
			return nil, errors.New("database error")
		}

		// Adding a new User object to the slice

		newChargeCode := &usecase.ChargeCode{ChargeCodeID: chargeCodeID, Code: code, MaxUses: maxUses, CurrentUses: currentUses, Amount: amount}
		ChargeCodes = append(ChargeCodes, newChargeCode)
		// Process the retrieved data here
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error iterating through charge code rows: %v", err)
		return nil, errors.New("database error")
	}

	// Return the list of transactions and no error
	if len(ChargeCodes) > 0 {
		return ChargeCodes, nil // Success case, return charge codes and no error
	}
	return nil, errors.New("charge codes not found")
}
