package money_accounts

import "fmt"

func checkAccountFields(fields MoneyAccountFields) error {
	if fields.Name == "" {
		return fmt.Errorf("Name is required")
	}
	if fields.Currency == "" {
		return fmt.Errorf("Currency is required")
	}
	return nil
}
