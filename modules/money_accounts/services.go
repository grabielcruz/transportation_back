package money_accounts

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
)

func GetMoneyAccounts() []MoneyAccount {
	var moneyAccounts []MoneyAccount
	rows, err := database.DB.Query("SELECT * FROM money_accounts WHERE id <> $1;", uuid.UUID{})
	if err != nil {
		errors_handler.HandleError(err)
		return moneyAccounts
	}
	defer rows.Close()

	for rows.Next() {
		var ma MoneyAccount
		err = rows.Scan(&ma.ID, &ma.Name, &ma.Balance, &ma.Details, &ma.Currency, &ma.CreatedAt, &ma.UpdatedAt)
		if err != nil {
			errors_handler.HandleError(err)
		}
		moneyAccounts = append(moneyAccounts, ma)
	}
	if err != nil {
		errors_handler.HandleError(rows.Err())
	}
	return moneyAccounts
}

func CreateMoneyAccount(fields MoneyAccountFields) (MoneyAccount, error) {
	var nma MoneyAccount
	row := database.DB.QueryRow(
		"INSERT INTO money_accounts (name, details, currency) VALUES ($1, $2, $3) RETURNING *;",
		fields.Name, fields.Details, fields.Currency)
	err := row.Scan(&nma.ID, &nma.Name, &nma.Balance, &nma.Details, &nma.Currency, &nma.CreatedAt, &nma.UpdatedAt)
	if err != nil {
		return nma, errors_handler.MapDBErrors(err)
	}
	return nma, nil
}

func GetOneMoneyAccount(account_id uuid.UUID) (MoneyAccount, error) {
	var ma MoneyAccount
	if account_id == (uuid.UUID{}) {
		return ma, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("SELECT * FROM money_accounts WHERE id = $1;", account_id)
	err := row.Scan(&ma.ID, &ma.Name, &ma.Balance, &ma.Details, &ma.Currency, &ma.CreatedAt, &ma.UpdatedAt)
	if err != nil {
		return ma, errors_handler.MapDBErrors(err)
	}
	return ma, nil
}

func GetAccountsCurrency(account_id uuid.UUID) (string, error) {
	currency := ""
	if account_id == (uuid.UUID{}) {
		return currency, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("SELECT currency FROM money_accounts WHERE id = $1;", account_id)
	err := row.Scan(&currency)
	if err != nil {
		return currency, errors_handler.MapDBErrors(err)
	}
	return currency, nil
}

func UpdateMoneyAccount(account_id uuid.UUID, fields MoneyAccountFields) (MoneyAccount, error) {
	var uma MoneyAccount
	if account_id == (uuid.UUID{}) {
		return uma, fmt.Errorf(errors_handler.DB001)
	}
	// should not update currency
	row := database.DB.QueryRow("UPDATE money_accounts SET name = $1, details = $2, updated_at = $3 WHERE id = $4 RETURNING *;",
		fields.Name, fields.Details, time.Now(), account_id)
	err := row.Scan(&uma.ID, &uma.Name, &uma.Balance, &uma.Details, &uma.Currency, &uma.CreatedAt, &uma.UpdatedAt)
	if err != nil {
		return uma, errors_handler.MapDBErrors(err)
	}
	return uma, nil
}

func getAccountsName(account_id uuid.UUID) (string, error) {
	name := ""
	if account_id == (uuid.UUID{}) {
		return name, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("SELECT name FROM money_accounts WHERE id = $1;", account_id)
	err := row.Scan(&name)
	if err != nil {
		return name, errors_handler.MapDBErrors(err)
	}
	return name, nil
}

func DeleteOneMoneyAccount(account_id uuid.UUID) (common.ID, error) {
	id := common.ID{}
	if account_id == (uuid.UUID{}) {
		return id, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("DELETE FROM money_accounts WHERE id = $1 RETURNING id;", account_id)
	err := row.Scan(&id.ID)
	if err != nil {
		return id, errors_handler.MapDBErrors(err)
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
	if account_id == (uuid.UUID{}) {
		return id, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING id;",
		balance, account_id)
	err := row.Scan(&id.ID)
	if err != nil {
		return id, errors_handler.MapDBErrors(err)
	}
	return id, nil
}

func DeleteAllMoneyAccounts() {
	database.DB.QueryRow("DELETE FROM money_accounts WHERE id <> $1;", uuid.UUID{})
}
