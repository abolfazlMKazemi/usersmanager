// internal/repository/user_repository.go
package repository

import (
	"chargeCode/internal/config"
	"chargeCode/internal/usecase"
	"database/sql"
	"errors" // Import the errors package
	"fmt"
	"regexp"
)

type UserRepository struct {
	// Implement data storage and retrieval methods here
	db     *sql.DB
	config *config.AppConfig
}

func NewUserRepository(db *sql.DB, config *config.AppConfig) *UserRepository {
	// Initialize the database connection

	return &UserRepository{db: db, config: config}
}

func (ur *UserRepository) GetUserByPhoneNumber(phoneNumber string) (*usecase.User, error) {

	// Ensure the database connection is valid
	if err := ur.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}

	formatPattern := `^09\d{9}$`

	// Compile the regular expression
	re := regexp.MustCompile(formatPattern)
	if !re.MatchString(phoneNumber) {
		return nil, errors.New("phone number is incorrect Most Like 09121114323")
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
			fmt.Println("Database error:", err)
			return nil, errors.New("database query error")
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

	formatPattern := `^09\d{9}$`

	// Compile the regular expression
	re := regexp.MustCompile(formatPattern)
	if !re.MatchString(user.PhoneNumber) {
		return nil, errors.New("phone number is incorrect Most Like 09121114323")
	}

	if user.Balance < 0 {
		return nil, errors.New("invalid user data")
	}

	// Ensure the database connection is valid
	if err := ur.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}

	_, err := ur.db.Exec("UPDATE user SET  phoneNumber=?, balance=? WHERE user_id=?", user.PhoneNumber, user.Balance, user.ID)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("database update error")
	}
	return user, nil
}

func (ur *UserRepository) ListOfUsersUseChargeCode(chargeCodeId int, page int, pageSize int) ([]*usecase.User, error) {

	if page > ur.config.MaxPage {
		return nil, errors.New("page exceeds the maximum allowed limit")
	}

	if pageSize > ur.config.MaxPageSize {
		return nil, errors.New("page size exceeds the maximum allowed limit")
	}

	// Ensure the database connection is valid
	if err := ur.db.Ping(); err != nil {
		fmt.Println(err)
		return nil, errors.New("internal Server Error")
	}

	//check charge code exist
	chargeCodeRepository := NewChargeCodeRepository(ur.db, ur.config)
	_, err := chargeCodeRepository.GetChargeCodeByID(chargeCodeId)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	users := []*usecase.User{}

	// Prepare the SQL query
	// Calculate the OFFSET based on the page number and page size
	offset := (page - 1) * pageSize

	// Prepare the SQL query with pagination
	query := `
			SELECT u.*
			FROM user u
			JOIN user_charge_code ucc ON u.user_id = ucc.user_id
			JOIN charge_code cc ON ucc.charge_code_id = cc.charge_code_id
			WHERE cc.charge_code_id = ?
			LIMIT ? OFFSET ?;
		`

	// Execute the query
	rows, err := ur.db.Query(query, chargeCodeId, pageSize, offset)
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
			return nil, errors.New("database query error")
		}

		// Adding a new User object to the slice
		newUser := &usecase.User{ID: userID, PhoneNumber: phoneNumber, Balance: balance}
		users = append(users, newUser)
	}

	if len(users) > 0 {
		return users, nil // Success case, return charge codes and no error
	}

	return nil, errors.New("no users found for the charge code ID")
}

func (ur *UserRepository) GetUserBalance(userId int) (float64, error) {

	// Ensure the database connection is valid
	if err := ur.db.Ping(); err != nil {
		fmt.Println(err)
		return 0, errors.New("internal Server Error")
	}

	// Check if the user exists based on user ID
	var userExists bool
	query := "SELECT COUNT(*) FROM user WHERE user_id = ?"
	err := ur.db.QueryRow(query, userId).Scan(&userExists)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("database query error")
	}

	if !userExists {
		return 0, errors.New("user not found")
	}

	// Retrieve the user's balance
	query = "SELECT balance FROM user WHERE user_id = ?"
	var balance float64
	err = ur.db.QueryRow(query, userId).Scan(&balance)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("database query error")
	}

	return balance, nil

}
