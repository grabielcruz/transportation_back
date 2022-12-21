package money_accounts

import (
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
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
		err = rows.Scan(&ma.ID, &ma.Name, &ma.Balance, &ma.Details, &ma.Currency, &ma.CreatedAt, &ma.UpdatedAt)
		errors_handler.CheckError(err)
		moneyAccounts = append(moneyAccounts, ma)
	}
	errors_handler.CheckError(rows.Err())
	return moneyAccounts
}

func CreateMoneyAccount(fields MoneyAccountFields) MoneyAccount {
	var nma MoneyAccount
	row := database.DB.QueryRow(
		"INSERT INTO money_accounts (name, details, currency) VALUES ($1, $2, $3) RETURNING *;",
		fields.Name, fields.Details, fields.Currency)
	err := row.Scan(&nma.ID, &nma.Name, &nma.Balance, &nma.Details, &nma.Currency, &nma.CreatedAt, &nma.UpdatedAt)
	errors_handler.CheckError(err)
	return nma
}

func GetOneMoneyAccount(acount_id uuid.UUID) (MoneyAccount, error) {
	var ma MoneyAccount
	row := database.DB.QueryRow("SELECT * FROM money_accounts WHERE id = $1;", acount_id)
	err := row.Scan(&ma.ID, &ma.Name, &ma.Balance, &ma.Details, &ma.Currency, &ma.CreatedAt, &ma.UpdatedAt)
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
	row := database.DB.QueryRow("UPDATE money_accounts SET name = $1, details = $2, currency = $3, updated_at = $4 WHERE id = $5 RETURNING *;",
		fields.Name, fields.Details, fields.Currency, time.Now(), account_id)
	err := row.Scan(&uma.ID, &uma.Name, &uma.Balance, &uma.Details, &uma.Currency, &uma.CreatedAt, &uma.UpdatedAt)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return uma, err
		}
		errors_handler.CheckError(err)
	}
	return uma, nil
}

func getAccountsName(acount_id uuid.UUID) (string, error) {
	name := ""
	row := database.DB.QueryRow("SELECT name FROM money_accounts WHERE id = $1;", acount_id)
	err := row.Scan(&name)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return name, err
		}
		errors_handler.CheckError(err)
	}
	return name, nil
}

func DeleteOneMoneyAccount(account_id uuid.UUID) (common.ID, error) {
	id := common.ID{}
	row := database.DB.QueryRow("DELETE FROM money_accounts WHERE id = $1 RETURNING id;", account_id)
	err := row.Scan(&id.ID)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return id, err
		}
		errors_handler.CheckError(err)
	}
	return id, nil
}

// ResetAccountsBalance sets the accounts with the specify id to zero
func ResetAccountsBalance(account_id uuid.UUID) (common.ID, error) {
	newBalance := float64(0)
	id, err := setAccountsBalance(account_id, newBalance)
	return id, err
}

func setAccountsBalance(account_id uuid.UUID, balance float64) (common.ID, error) {
	id := common.ID{}
	row := database.DB.QueryRow("UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING id;",
		balance, account_id)
	err := row.Scan(&id.ID)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return id, err
		}
		errors_handler.CheckError(err)
	}
	return id, nil
}

func DeleteAllMoneyAccounts() {
	database.DB.QueryRow("DELETE FROM money_accounts;")
}
