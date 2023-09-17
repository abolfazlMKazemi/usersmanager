// internal/repository/user_repository.go
package repository

import (
	"chargeCode/internal/usecase"
	"database/sql"
	"errors" // Import the errors package
	"fmt"
)

type UserRepository struct {
	// Implement data storage and retrieval methods here
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	// Initialize the database connection

	return &UserRepository{db: db}
}

func (ur *UserRepository) GetUserByPhoneNumber(phoneNumber string) (*usecase.User, error) {

	// Ensure the database connection is valid
	if err := ur.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("Internal Server Error")
	}
	// Query to retrieve user by phone number
	query := "SELECT user_id, phoneNumber, balance FROM user WHERE phoneNumber = ?"

	var (
		userID  int
		phone   string
		balance float64
	)

	err := ur.db.QueryRow(query, phoneNumber).Scan(&userID, &phone, &balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		} else {
			fmt.Println(err)
			return nil, errors.New(err.Error())
		}
	} else {
		user := &usecase.User{
			ID:          userID,
			PhoneNumber: phone, // Replace with actual data retrieval logic
			Balance:     balance,
		}

		return user, nil // Success case, return user and no error

	}
}

func (ur *UserRepository) UpdateUser(user *usecase.User) (*usecase.User, error) {

	// Ensure the database connection is valid
	if err := ur.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("Internal Server Error")
	}

	_, err := ur.db.Exec("UPDATE user SET  phoneNumber=?, balance=? WHERE user_id=?", user.PhoneNumber, user.Balance, user.ID)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}
	return user, nil
}

func (ur *UserRepository) ListOfUsersUseChargeCode(chargeCodeId int) ([]*usecase.User, error) {

	// Ensure the database connection is valid
	if err := ur.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("Internal Server Error")
	}

	users := []*usecase.User{}

	// Prepare the SQL query
	query := `
	   SELECT u.*
	   FROM user u
	   JOIN user_charge_code ucc ON u.user_id = ucc.user_id
	   JOIN charge_code cc ON ucc.charge_code_id = cc.charge_code_id
	   WHERE cc.charge_code_id = ?;
   `

	// Execute the query
	rows, err := ur.db.Query(query, chargeCodeId)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var (
			userID      int
			phoneNumber string
			balance     float64 // Assuming balance is a decimal column
		)
		if err := rows.Scan(&userID, &phoneNumber, &balance); err != nil {
			fmt.Println(err)
			return nil, errors.New(err.Error())
		}

		// Adding a new User object to the slice
		newUser := &usecase.User{ID: userID, PhoneNumber: phoneNumber, Balance: balance}
		users = append(users, newUser)
	}

	if len(users) > 0 {
		return users, nil // Success case, return charge codes and no error
	}

	return nil, errors.New("users  not found") // Replace someError with an actual error value
}

func (ur *UserRepository) GetUserBalance(userId int) (float64, error) {

	// Ensure the database connection is valid
	if err := ur.db.Ping(); err != nil {
		fmt.Println(err)
		return 0, errors.New("Internal Server Error")
	}

	query := "SELECT balance FROM user WHERE user_id = ?"
	// Execute the query and retrieve the balance
	var balance float64
	err := ur.db.QueryRow(query, userId).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("User not found")
		} else {
			fmt.Println(err)
			return 0, errors.New(err.Error())
		}
	} else {
		return balance, nil
	}

}
