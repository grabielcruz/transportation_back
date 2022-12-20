package transactions

import (
	"fmt"

	"github.com/grabielcruz/transportation_back/database"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/grabielcruz/transportation_back/utility"
)

func GetTransactions(offset int, limit int) TransationResponse {
	tr := TransationResponse{}
	return tr
}

func CreateTransaction(fields TransactionFields) (Transaction, error) {
	tr := Transaction{}
	var oldBalance float64 = 0
	var updatedBalance float64 = 0
	tx_err := fmt.Errorf("Could not complete transaction")

	tr.PersonName, _ = persons.GetPersonsName(fields.PersonId)
	tr.AccountName, _ = money_accounts.GetAccountsName(fields.AccountId)

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return tr, tx_err
	}
	row := tx.QueryRow(`SELECT balance FROM money_accounts WHERE id = $1;`, fields.AccountId)
	err = row.Scan(&oldBalance)
	if err != nil {
		tx.Rollback()
		return tr, tx_err
	}

	newBalance := oldBalance + utility.RoundToTwoDecimalPlaces(fields.Amount)

	row = tx.QueryRow(`UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING balance;`, newBalance, fields.AccountId)
	err = row.Scan(&updatedBalance)
	if err != nil {
		tx.Rollback()
		return tr, tx_err
	}

	row = tx.QueryRow(`INSERT INTO transactions (account_id, person_id, date, amount, description, balance) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`, fields.AccountId, fields.PersonId, fields.Date, fields.Amount, fields.Description, updatedBalance)
	err = row.Scan(&tr.ID, &tr.AccountId, &tr.PersonId, &tr.Date, &tr.Amount, &tr.Description, &tr.Balance, &tr.CreatedAt, &tr.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return tr, tx_err
	}

	tx.Commit()
	return tr, nil
}
