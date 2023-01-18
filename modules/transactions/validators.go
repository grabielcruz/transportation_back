package transactions

import "fmt"

func checkTransactionFields(fields TransactionFields) error {
	if fields.Description == "" {
		return fmt.Errorf("Transaction should have a description")
	}
	if fields.Amount == float64(0) {
		return fmt.Errorf("Amount should be greater than zero")
	}
	return nil
}
