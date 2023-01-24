package person_accounts

import (
	"fmt"

	"github.com/grabielcruz/transportation_back/modules/currencies"
)

func checkPersonAccountFields(fields PersonAccountFields) error {
	if fields.Name == "" {
		return fmt.Errorf("Name is required")
	}
	if fields.Description == "" {
		return fmt.Errorf("Description is required")
	}
	err := currencies.CheckValidCurrency(fields.Currency)
	if err != nil {
		return err
	}
	return nil
}
