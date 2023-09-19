// config/service.go
package config

import (
	"errors"
	"os"
	"strconv"
)

type AppConfig struct {
	MaxPage              int
	MaxPageSize          int
	ApplicationPort      string
	MysqlUrl             string
	MaxChargeCodeAmount  float64
	MinChargeCodeAmount  float64
	MaxTransactionAmount float64
	MinTransactionAmount float64
}

func LoadConfig() (*AppConfig, error) {
	maxPageStr := os.Getenv("MAX_PAGE")
	if maxPageStr == "" {
		return nil, errors.New("MAX_PAGE environment variable is not set")
	}

	maxPageSizeStr := os.Getenv("MAX_PAGE_SIZE")
	if maxPageSizeStr == "" {
		return nil, errors.New("MAX_PAGE_SIZE environment variable is not set")
	}

	applicationPortStr := os.Getenv("APPLICATION_PORT")
	if applicationPortStr == "" {
		return nil, errors.New("APPLICATION_PORT environment variable is not set")
	}

	mysqlURL := os.Getenv("MYSQL_URL")
	if mysqlURL == "" {
		return nil, errors.New("MYSQL_URL environment variable is not set")
	}

	maxChargeCodeAmountStr := os.Getenv("MAX_CHARGE_CODE_AMOUNT")
	if maxChargeCodeAmountStr == "" {
		return nil, errors.New("MAX_CHARGE_CODE_AMOUNT environment variable is not set")
	}

	minChargeCodeAmountStr := os.Getenv("MIN_CHARGE_CODE_AMOUNT")
	if maxChargeCodeAmountStr == "" {
		return nil, errors.New("MIN_CHARGE_CODE_AMOUNT environment variable is not set")
	}

	maxTransactionAmountStr := os.Getenv("MAX_TRANSACTION_AMOUNT")
	if maxTransactionAmountStr == "" {
		return nil, errors.New("MAX_TRANSACTION_AMOUNT environment variable is not set")
	}

	minTransactionAmountStr := os.Getenv("MIN_TRANSACTION_AMOUNT")
	if minTransactionAmountStr == "" {
		return nil, errors.New("MIN_TRANSACTION_AMOUNT environment variable is not set")
	}

	// Convert environment variables to their respective types
	maxPage, err := strconv.Atoi(maxPageStr)
	if err != nil {
		return nil, err
	}

	if maxPage <= 0 {
		return nil, errors.New("maxPage most bigger than zero")
	}

	maxPageSize, err := strconv.Atoi(maxPageSizeStr)
	if err != nil {
		return nil, err
	}

	if maxPageSize <= 0 {
		return nil, errors.New("maxPageSize most bigger than zero")
	}

	max_charge_code_amount, err := strconv.ParseFloat(maxChargeCodeAmountStr, 64)
	if err != nil {
		return nil, err
	}

	if max_charge_code_amount <= 0 {
		return nil, errors.New("max_charge_code_amount most bigger than zero")
	}

	min_charge_code_amount, err := strconv.ParseFloat(minChargeCodeAmountStr, 64)
	if err != nil {
		return nil, err
	}

	if min_charge_code_amount <= 0 {
		return nil, errors.New("min_charge_code_amount most bigger than zero")
	}

	max_TRANSACTION_AMOUNT, err := strconv.ParseFloat(maxTransactionAmountStr, 64)
	if err != nil {
		return nil, err
	}

	min_TRANSACTION_AMOUNT, err := strconv.ParseFloat(minTransactionAmountStr, 64)
	if err != nil {
		return nil, err
	}

	// if min_TRANSACTION_AMOUNT <= 0 {
	// 	return nil, errors.New("min_TRANSACTION_AMOUNT most bigger than zero")
	// }

	//

	//
	return &AppConfig{
		MaxPage:              maxPage,
		MaxPageSize:          maxPageSize,
		ApplicationPort:      applicationPortStr,
		MysqlUrl:             mysqlURL,
		MaxChargeCodeAmount:  max_charge_code_amount,
		MinChargeCodeAmount:  min_charge_code_amount,
		MaxTransactionAmount: max_TRANSACTION_AMOUNT,
		MinTransactionAmount: min_TRANSACTION_AMOUNT,
	}, nil
}
