package transactions

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/grabielcruz/transportation_back/utility"
)

func GetTransactions(account_id uuid.UUID, limit int, offset int) TransationResponse {
	transactionResponse := TransationResponse{}
	rows, err := database.DB.Query("SELECT * FROM transactions WHERE account_id = $1 LIMIT $2 OFFSET $3;", account_id, limit, offset)
	errors_handler.CheckError(err)
	defer rows.Close()

	for rows.Next() {
		t := Transaction{}
		err = rows.Scan(&t.ID, &t.AccountId, &t.PersonId, &t.Date, &t.Amount, &t.Description, &t.Balance, &t.CreatedAt, &t.UpdatedAt)
		errors_handler.CheckError(err)
		t.PersonName, _ = persons.GetPersonsName(t.PersonId)
		transactionResponse.Transactions = append(transactionResponse.Transactions, t)
	}
	transactionResponse.Pagination.Count = len(transactionResponse.Transactions)
	transactionResponse.Pagination.Limit = limit
	transactionResponse.Pagination.Offset = offset
	return transactionResponse
}

func CreateTransaction(fields TransactionFields) (Transaction, error) {
	tr := Transaction{}

	if fields.Description == "" {
		return tr, fmt.Errorf("Transaction should have a description")
	}

	if fields.Amount == float64(0) {
		return tr, fmt.Errorf("Amount should be greater than zero")
	}

	var oldBalance float64 = 0
	var updatedBalance float64 = 0

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf("Could not begin transaction")
	}
	row := tx.QueryRow(`SELECT balance FROM money_accounts WHERE id = $1;`, fields.AccountId)
	err = row.Scan(&oldBalance)
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf("Could not get balance from account")
	}

	newBalance := utility.RoundToTwoDecimalPlaces(oldBalance + utility.RoundToTwoDecimalPlaces(fields.Amount))
	if newBalance < 0 {
		tx.Rollback()
		return tr, fmt.Errorf("Transaction should not generate a negative balance")
	}

	row = tx.QueryRow(`UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING balance;`, newBalance, fields.AccountId)
	err = row.Scan(&updatedBalance)
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf("Could not update account's balance")
	}

	if newBalance != updatedBalance {
		tx.Rollback()
		return tr, fmt.Errorf("New balance and updated balance missmatch, oldBalance = %v, newBalance = %v, updatedBalance = %v", oldBalance, newBalance, updatedBalance)
	}

	row = tx.QueryRow(`INSERT INTO transactions (account_id, person_id, date, amount, description, balance) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`, fields.AccountId, fields.PersonId, fields.Date, fields.Amount, fields.Description, updatedBalance)
	err = row.Scan(&tr.ID, &tr.AccountId, &tr.PersonId, &tr.Date, &tr.Amount, &tr.Description, &tr.Balance, &tr.CreatedAt, &tr.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf("Could not insert a new transaction into transaction table")
	}

	tx.Commit()

	tr.PersonName, _ = persons.GetPersonsName(tr.PersonId)

	return tr, nil
}

func deleteAllTransactions() {
	database.DB.QueryRow("DELETE FROM transactions;")
}
