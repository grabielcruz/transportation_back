package money_accounts

import (
	"github.com/google/uuid"
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

func GetOneMoneyAccount(acount_id uuid.UUID) (MoneyAccount, error) {
	var ma MoneyAccount
	row := database.DB.QueryRow("SELECT * FROM money_accounts WHERE id=$1;", acount_id)
	err := row.Scan(&ma.ID, &ma.Name, &ma.Balance, &ma.IsCash, &ma.Currency, &ma.CreatedAt, &ma.UpdatedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ma, err
		}
		errors_handler.CheckError(err)
	}
	return ma, nil
}

func DeleteOneMoneyAccount(account_id uuid.UUID) error {
	row := database.DB.QueryRow("DELETE FROM money_accounts WHERE id=$1", account_id)
	err := row.Scan()
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return err
		}
		errors_handler.CheckError(err)
	}
	return nil
}

func deleteAllMoneyAccounts() {
	database.DB.QueryRow("DELETE FROM money_accounts;")
}
