package currencies

import (
	"fmt"

	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
)

func GetCurrencies() []string {
	currencies := []string{}
	rows, err := database.DB.Query("SELECT currency FROM currencies WHERE currency <> $1;", "000")
	errors_handler.CheckError(err)
	defer rows.Close()

	for rows.Next() {
		var currency string
		err = rows.Scan(&currency)
		errors_handler.CheckError(err)
		currencies = append(currencies, currency)
	}
	errors_handler.CheckError(rows.Err())
	return currencies
}

func CreateCurrency(newCurrency string) (string, error) {
	createdCurrency := ""
	err := checkValidCurrency(newCurrency)
	if err != nil {
		return createdCurrency, err
	}
	row := database.DB.QueryRow("INSERT INTO currencies (currency) VALUES ($1) RETURNING currency;", newCurrency)
	err = row.Scan(&createdCurrency)
	if err != nil {
		return createdCurrency, mapCurrencyDBError(err)
	}
	return createdCurrency, nil
}

func DeleteCurrency(currency string) (string, error) {
	deletedCurrency := ""
	err := checkValidCurrency(currency)
	if err != nil {
		return deletedCurrency, err
	}
	if currency == "VED" || currency == "USD" {
		return deletedCurrency, fmt.Errorf("Could not delete VED or USD currency")
	}
	row := database.DB.QueryRow("DELETE FROM currencies WHERE currency = $1 RETURNING currency;", currency)
	err = row.Scan(&deletedCurrency)
	if err != nil {
		return deletedCurrency, mapCurrencyDBError(err)
	}
	return deletedCurrency, nil
}

func resetCurrencies() {
	database.DB.QueryRow("DELETE FROM currencies WHERE currency <> $1;", "000")
	database.DB.QueryRow("INSERT INTO currencies (currency) VALUES ('VED'), ('USD');")
}

func mapCurrencyDBError(err error) error {
	switch err.Error() {
	case "pq: duplicate key value violates unique constraint \"currencies_pkey\"":
		return fmt.Errorf("Currency already exists")
	case "pq: update or delete on table \"currencies\" violates foreign key constraint \"money_accounts_currency_fkey\" on table \"money_accounts\"":
		return fmt.Errorf("Currency is being used")
	}
	return err
}
