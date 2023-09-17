package repository

// Import the errors package
import (
	"chargeCode/internal/usecase"
	"database/sql"
	"errors"
	"fmt"
)

type ChargeCodeRepository struct {
	// Implement data storage and retrieval methods here
	db *sql.DB
}

func NewChargeCodeRepository(db *sql.DB) *ChargeCodeRepository {
	// Initialize the database connection

	return &ChargeCodeRepository{db: db}
}

func (cu *ChargeCodeRepository) GetChargeCodes() ([]*usecase.ChargeCode, error) {
	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("Internal Server Error")
	}
	// Query all transactions from the 'transaction' table
	rows, err := cu.db.Query(`
	SELECT charge_code_id, code, max_uses, current_uses, amount
	FROM charge_code
`)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
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
			fmt.Println(err)
			return nil, errors.New(err.Error())
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
		return nil, errors.New("Internal Server Error")
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
			fmt.Println(err)
			return nil, errors.New(err.Error())
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
		return nil, errors.New("Internal Server Error")
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
			fmt.Println(err)
			return nil, errors.New(err.Error())
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
		return nil, errors.New("Internal Server Error")
	}

	// Insert the new charge code into the 'charge_code' table
	_, err := cu.db.Exec(`
	   INSERT INTO charge_code (code, max_uses, amount)
	   VALUES (?, ?, ?)
   `, chargeCode.Code, chargeCode.MaxUses, chargeCode.Amount)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}

	return chargeCode, nil
}

func (cu *ChargeCodeRepository) DeleteChargeCode(id int) error {

	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return errors.New("Internal Server Error")
	}
	// Delete the charge code by ID from the 'charge_code' table
	_, err := cu.db.Exec(`
   DELETE FROM charge_code
   WHERE charge_code_id = ?
`, id)
	if err != nil {
		fmt.Println(err)
		return errors.New(err.Error())
	}
	return nil
}

func (cu *ChargeCodeRepository) UpdateChargeCode(chargeCode *usecase.ChargeCode) (*usecase.ChargeCode, error) {
	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("Internal Server Error")
	}

	// Update the charge code by ID in the 'charge_code' table
	_, err := cu.db.Exec(`
	   UPDATE charge_code
	   SET code = ?, max_uses = ?, amount = ?, current_uses = ?
	   WHERE charge_code_id = ?
   `, chargeCode.Code, chargeCode.MaxUses, chargeCode.Amount, chargeCode.CurrentUses, chargeCode.ChargeCodeID)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}

	return chargeCode, nil
}

func (cu *ChargeCodeRepository) GetUserChargeCodes(userId int) ([]*usecase.ChargeCode, error) {
	// Ensure the database connection is valid
	if err := cu.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("Internal Server Error")
	}
	// Query the charge codes by user ID from the 'user_charge_code' table
	rows, err := cu.db.Query(`
        SELECT uc.charge_code_id, cc.code, cc.max_uses, cc.current_uses, cc.amount
        FROM user_charge_code uc
        INNER JOIN charge_code cc ON uc.charge_code_id = cc.charge_code_id
        WHERE uc.user_id = ?
    `, userId)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
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
			fmt.Println(err)
			return nil, errors.New(err.Error())
		}

		// Adding a new User object to the slice

		newChargeCode := &usecase.ChargeCode{ChargeCodeID: chargeCodeID, Code: code, MaxUses: maxUses, CurrentUses: currentUses, Amount: amount}
		ChargeCodes = append(ChargeCodes, newChargeCode)
		// Process the retrieved data here
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}

	// Return the list of transactions and no error
	if len(ChargeCodes) > 0 {
		return ChargeCodes, nil // Success case, return charge codes and no error
	}
	return nil, errors.New("transaction not found")
}
