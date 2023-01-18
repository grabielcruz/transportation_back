package bills

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/modules/currencies"
)

func checkBillFields(fields BillFields) error {
	if fields.PersonId == (uuid.UUID{}) {
		return fmt.Errorf("Person id should be not zero uuid")
	}
	if fields.Description == "" {
		return fmt.Errorf("Description is required")
	}
	if fields.Amount <= 0 {
		return fmt.Errorf("Amount should be greater than zero")
	}
	if err := currencies.CheckValidCurrency(fields.Currency); err != nil {
		return err
	}
	return nil
}
