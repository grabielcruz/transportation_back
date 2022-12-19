package money_accounts

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/common"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/utility"
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

func GetAccountsBalance(account_id uuid.UUID) (float64, error) {
	var balance float64 = 0
	row := database.DB.QueryRow("SELECT balance FROM money_accounts WHERE id = $1;", account_id)
	err := row.Scan(&balance)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return balance, err
		}
		errors_handler.CheckError(err)
	}
	return balance, nil
}

func AddToBalance(account_id uuid.UUID, amount float64) (AccountNameAndBalance, error) {
	anb := AccountNameAndBalance{}
	oldBalance, err := GetAccountsBalance(account_id)
	if err != nil {
		return anb, err
	}
	newBalance := oldBalance + utility.RoundToTwoDecimalPlaces(amount)
	if newBalance < 0 {
		err := fmt.Errorf("New balance can't be a negative number")
		return anb, err
	}
	anb, err = updateMoneyAccountBalance(account_id, newBalance)
	if err != nil {
		return anb, err
	}
	return anb, nil
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

func updateMoneyAccountBalance(account_id uuid.UUID, balance float64) (AccountNameAndBalance, error) {
	var uma AccountNameAndBalance
	row := database.DB.QueryRow("UPDATE money_accounts SET balance = $1, updated_at = $2 WHERE id = $3 RETURNING id, name, balance;",
		balance, time.Now(), account_id)
	err := row.Scan(&uma.ID, &uma.Name, &uma.Balance)
	if err != nil {
		if errors_handler.CheckEmptyRowError(err) {
			return uma, err
		}
		errors_handler.CheckError(err)
	}
	return uma, nil
}

func deleteAllMoneyAccounts() {
	database.DB.QueryRow("DELETE FROM money_accounts;")
}
