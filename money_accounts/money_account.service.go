package money_accounts

import (
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
)

func GetMoneyAccounts() []MoneyAccount {
	var moneyAccounts []MoneyAccount
	rows, err := database.DB.Query("SELECT * FROM money_accounts;")
	errors_handler.CheckError(err)
	defer rows.Close()

	for rows.Next() {
		var ma MoneyAccount
		err = rows.Scan(&ma.ID, &ma.Name, &ma.Balance, &ma.IsCash, &ma.Currency, &ma.CreatedAt, &ma.UpdatedAt)
		errors_handler.CheckError(err)
		moneyAccounts = append(moneyAccounts, ma)
	}
	errors_handler.CheckError(rows.Err())
	return moneyAccounts
}

func CreateMoneyAccount(moneyAccount MoneyAccount) MoneyAccount {
	var nma MoneyAccount
	row := database.DB.QueryRow(
		"INSERT INTO money_accounts (name, balance, is_cash, currency) VALUES ($1, $2, $3, $4) RETURNING *;",
		moneyAccount.Name, moneyAccount.Balance, moneyAccount.IsCash, moneyAccount.Currency)
	err := row.Scan(&nma.ID, &nma.Name, &nma.Balance, &nma.IsCash, &nma.Currency, &nma.CreatedAt, &nma.UpdatedAt)
	errors_handler.CheckError(err)
	return nma
}
