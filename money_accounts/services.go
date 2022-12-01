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

func CreateMoneyAccount(fields MoneyAccountFields) MoneyAccount {
	var nma MoneyAccount
	row := database.DB.QueryRow(
		"INSERT INTO money_accounts (name, is_cash, currency) VALUES ($1, $2, $3) RETURNING *;",
		fields.Name, fields.IsCash, fields.Currency)
	err := row.Scan(&nma.ID, &nma.Name, &nma.Balance, &nma.IsCash, &nma.Currency, &nma.CreatedAt, &nma.UpdatedAt)
	errors_handler.CheckError(err)
	return nma
}

func GetOneMoneyAccount(acount_id uuid.UUID) (MoneyAccount, error) {
	var ma MoneyAccount
	row := database.DB.QueryRow("SELECT * FROM money_accounts WHERE id=$1;", acount_id)
	err := row.Scan(&ma.ID, &ma.Name, &ma.Balance, &ma.IsCash, &ma.Currency, &ma.CreatedAt, &ma.UpdatedAt)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return ma, err
		}
		errors_handler.CheckError(err)
	}
	return ma, nil
}

func UpdateMoneyAccount(account_id uuid.UUID, fields MoneyAccountFields) (MoneyAccount, error) {
	var uma MoneyAccount
	row := database.DB.QueryRow("UPDATE money_accounts SET name = $1, is_cash = $2, currency = $3 WHERE id = $4 RETURNING *;",
		fields.Name, fields.IsCash, fields.Currency, account_id)
	err := row.Scan(&uma.ID, &uma.Name, &uma.Balance, &uma.IsCash, &uma.Currency, &uma.CreatedAt, &uma.UpdatedAt)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return uma, err
		}
		errors_handler.CheckError(err)
	}
	return uma, nil
}

func UpdatedMoneyAccountsBalance(account_id uuid.UUID, balance MoneyAccountBalance) (MoneyAccount, error) {
	var uma MoneyAccount
	row := database.DB.QueryRow("UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING *;",
		balance.Balance, account_id)
	err := row.Scan(&uma.ID, &uma.Name, &uma.Balance, &uma.IsCash, &uma.Currency, &uma.CreatedAt, &uma.UpdatedAt)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return uma, err
		}
		errors_handler.CheckError(err)
	}
	return uma, nil
}

func DeleteOneMoneyAccount(account_id uuid.UUID) (MoneyAccount, error) {
	var dma MoneyAccount
	row := database.DB.QueryRow("DELETE FROM money_accounts WHERE id=$1 RETURNING *;", account_id)
	err := row.Scan(&dma.ID, &dma.Name, &dma.Balance, &dma.IsCash, &dma.Currency, &dma.CreatedAt, &dma.UpdatedAt)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return dma, err
		}
		errors_handler.CheckError(err)
	}
	return dma, nil
}

func deleteAllMoneyAccounts() {
	database.DB.QueryRow("DELETE FROM money_accounts;")
}
