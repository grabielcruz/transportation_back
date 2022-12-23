package transactions

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/grabielcruz/transportation_back/utility"
)

func GetTransactions(account_id uuid.UUID, limit int, offset int) (TransationResponse, error) {
	transactionResponse := TransationResponse{}

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return transactionResponse, fmt.Errorf("Could not begin transaction")
	}

	row := tx.QueryRow("SELECT COUNT(*) FROM transactions WHERE account_id = $1;", account_id)
	err = row.Scan(&transactionResponse.Pagination.Count)
	if err != nil {
		tx.Rollback()
		return transactionResponse, fmt.Errorf("Could not get balance from account")
	}

	rows, err := tx.Query("SELECT * FROM transactions WHERE account_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;", account_id, limit, offset)
	if err != nil {
		tx.Rollback()
		return transactionResponse, fmt.Errorf("Could not get transactions")
	}

	for rows.Next() {
		t := Transaction{}
		err = rows.Scan(&t.ID, &t.AccountId, &t.PersonId, &t.Date, &t.Amount, &t.Description, &t.Balance, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			tx.Rollback()
			return transactionResponse, fmt.Errorf("Could not read transaction")
		}
		t.PersonName, _ = persons.GetPersonsName(t.PersonId)
		transactionResponse.Transactions = append(transactionResponse.Transactions, t)
	}

	tx.Commit()

	transactionResponse.Pagination.Limit = limit
	transactionResponse.Pagination.Offset = offset
	return transactionResponse, nil
}

func CreateTransaction(fields TransactionFields) (Transaction, error) {
	tr := Transaction{}

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

func UpdateLastTransaction(transaction_id uuid.UUID, fields TransactionFields) (Transaction, error) {
	oldT := Transaction{}
	udT := Transaction{}
	updatedBalance := float64(0)

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return udT, fmt.Errorf("Could not begin transaction")
	}

	row := tx.QueryRow("SELECT * FROM transactions ORDER BY created_at DESC LIMIT 1;")
	err = row.Scan(&oldT.ID, &oldT.AccountId, &oldT.PersonId, &oldT.Date, &oldT.Amount, &oldT.Description, &oldT.Balance, &oldT.CreatedAt, &oldT.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return udT, fmt.Errorf("No transaction found in database")
	}

	if oldT.ID != transaction_id {
		tx.Rollback()
		return udT, fmt.Errorf("The transaction requested is not the last transaction")
	}

	newBalance := utility.RoundToTwoDecimalPlaces(oldT.Balance - oldT.Amount + fields.Amount)

	if newBalance < 0 {
		tx.Rollback()
		return udT, fmt.Errorf("Transaction should not generate a negative balance")
	}

	row = tx.QueryRow(`UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING balance;`, newBalance, fields.AccountId)
	err = row.Scan(&updatedBalance)
	if err != nil {
		tx.Rollback()
		return udT, fmt.Errorf("Could not update account's balance")
	}

	if newBalance != updatedBalance {
		tx.Rollback()
		return udT, fmt.Errorf("New balance and updated balance missmatch, oldBalance = %v, newBalance = %v, updatedBalance = %v", oldT.Balance, newBalance, updatedBalance)
	}

	row = tx.QueryRow(`UPDATE transactions SET account_id = $1, person_id = $2, date = $3, amount = $4, description = $5, balance = $6, created_at = $7, updated_at = $8 WHERE id = $9 RETURNING *;`,
		fields.AccountId, fields.PersonId, fields.Date, fields.Amount, fields.Description, updatedBalance, oldT.CreatedAt, time.Now(), transaction_id)
	err = row.Scan(&udT.ID, &udT.AccountId, &udT.PersonId, &udT.Date, &udT.Amount, &udT.Description, &udT.Balance, &udT.CreatedAt, &udT.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return udT, fmt.Errorf("Could not update transaction with the id %v", transaction_id)
	}

	tx.Commit()

	udT.PersonName, _ = persons.GetPersonsName(udT.PersonId)

	return udT, nil
}

func deleteAllTransactions() {
	database.DB.QueryRow("DELETE FROM transactions;")
}
