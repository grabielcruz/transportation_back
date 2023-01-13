package currencies

import (
	"fmt"

	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
)

func GetCurrencies() []string {
	currencies := []string{}
	rows, err := database.DB.Query("SELECT currency FROM currencies WHERE currency <> $1;", "000")
	errors_handler.HandleError(err)
	defer rows.Close()

	for rows.Next() {
		var currency string
		err = rows.Scan(&currency)
		errors_handler.HandleError(err)
		currencies = append(currencies, currency)
	}
	errors_handler.HandleError(rows.Err())
	return currencies
}

func CreateCurrency(newCurrency string) (string, error) {
	createdCurrency := ""
	err := CheckValidCurrency(newCurrency)
	if err != nil {
		return createdCurrency, err
	}
	row := database.DB.QueryRow("INSERT INTO currencies (currency) VALUES ($1) RETURNING currency;", newCurrency)
	err = row.Scan(&createdCurrency)
	if err != nil {
		return createdCurrency, errors_handler.MapDBErrors(err)
	}
	return createdCurrency, nil
}

func DeleteCurrency(currency string) (string, error) {
	deletedCurrency := ""
	err := CheckValidCurrency(currency)
	if err != nil {
		return deletedCurrency, err
	}
	if currency == "VED" || currency == "USD" {
		return deletedCurrency, fmt.Errorf("Could not delete VED or USD currency")
	}
	row := database.DB.QueryRow("DELETE FROM currencies WHERE currency = $1 RETURNING currency;", currency)
	err = row.Scan(&deletedCurrency)
	if err != nil {
		return deletedCurrency, errors_handler.MapDBErrors(err)
	}
	return deletedCurrency, nil
}

func resetCurrencies() {
	database.DB.QueryRow("DELETE FROM currencies WHERE currency <> $1;", "000")
	database.DB.QueryRow("INSERT INTO currencies (currency) VALUES ('VED'), ('USD');")
}
