package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// NewDBConnection initializes a new database connection and returns it.
func NewDBConnection() (*sql.DB, error) {
	//	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	dbUsername := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	// Construct the database URL using environment variables
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUsername, dbPassword, dbHost, dbPort)

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open connections
	db.SetMaxOpenConns(30)

	// Create the database if it doesn't exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS userManager")
	if err != nil {
		db.Close() // Close the connection if database creation fails
		return nil, err
	}

	// Switch to the newly created database
	_, err = db.Exec("USE userManager")
	if err != nil {
		db.Close() // Close the connection if database switching fails
		return nil, err
	}

	// Create tables if they don't exist
	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS user (
			user_id INT PRIMARY KEY AUTO_INCREMENT,
			phoneNumber VARCHAR(20) UNIQUE, -- Add phoneNumber column
			balance DECIMAL(10, 2) DEFAULT 0.00 -- Add balance column with default value
		)`,
		`CREATE TABLE IF NOT EXISTS charge_code (
            charge_code_id INT PRIMARY KEY AUTO_INCREMENT,
            code VARCHAR(255) NOT NULL UNIQUE,
            max_uses INT NOT NULL,
            current_uses INT NOT NULL DEFAULT 0,
            amount DECIMAL(10, 2) NOT NULL
            -- Add other charge code-related columns as needed
        )`,
		`CREATE TABLE IF NOT EXISTS user_charge_code (
            user_charge_code_id INT PRIMARY KEY AUTO_INCREMENT,
            user_id INT NOT NULL,
            charge_code_id INT NOT NULL,
            usage_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES user(user_id),
            FOREIGN KEY (charge_code_id) REFERENCES charge_code(charge_code_id)
        )`,
		`CREATE TABLE IF NOT EXISTS transaction (
            transaction_id INT PRIMARY KEY AUTO_INCREMENT,
            user_id INT NOT NULL,
            amount DECIMAL(10, 2) NOT NULL,
            timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES user(user_id)
        )`,
	}

	for _, query := range createTableQueries {
		_, err = db.Exec(query)
		if err != nil {
			db.Close() // Close the connection if table creation fails
			return nil, err
		}
	}

	// Drop the trigger if it exists (ignore errors if it doesn't exist)
	_, err = db.Exec("DROP TRIGGER IF EXISTS update_user_balance")

	if err != nil {
		db.Close() // Close the connection if trigger creation fails
		return nil, err
	}

	// Create the trigger to update user balance after inserting a transaction
	_, err = db.Exec(`
	  CREATE TRIGGER update_user_balance AFTER INSERT ON transaction
	  FOR EACH ROW
	  BEGIN
		  UPDATE user
		  SET balance = balance + NEW.amount
		  WHERE user_id = NEW.user_id;
	  END;
  `)
	if err != nil {
		db.Close() // Close the connection if trigger creation fails
		return nil, err
	}

	// Drop the RedeemChargeCode procedure if it exists (ignore errors if it doesn't exist)
	_, err = db.Exec("DROP PROCEDURE IF EXISTS RedeemChargeCode")
	if err != nil {
		db.Close() // Close the connection if procedure deletion fails
		return nil, err
	}

	// Create the RedeemChargeCode stored procedure
	_, err = db.Exec(`
	CREATE PROCEDURE RedeemChargeCode(IN in_user_id INT, IN in_charge_code_id INT)
	BEGIN
		DECLARE charge_amount DECIMAL(10, 2);
	
		START TRANSACTION;
	
		-- Update charge_code
		UPDATE charge_code
		SET current_uses = current_uses + 1
		WHERE charge_code_id = in_charge_code_id;
	
		-- Insert into user_charge_code
		INSERT INTO user_charge_code (user_id, charge_code_id)
		VALUES (in_user_id, in_charge_code_id);
	
		-- Get the amount from charge_code
		SELECT amount INTO charge_amount
		FROM charge_code
		WHERE charge_code_id = in_charge_code_id LIMIT 1;
	
		-- Insert into transaction
		INSERT INTO transaction (user_id, amount)
		VALUES (in_user_id, charge_amount);
	
		COMMIT;
	END;
  `)

	if err != nil {
		db.Close() // Close the connection if procedure creation fails
		return nil, err
	}

	return db, nil
}
