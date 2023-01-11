package transactions

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/grabielcruz/transportation_back/utility"
)

func GetTransactions(account_id uuid.UUID, limit int, offset int) (TransationResponse, error) {
	transactionResponse := TransationResponse{}

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return transactionResponse, fmt.Errorf(errors_handler.DB002)
	}

	row := tx.QueryRow("SELECT COUNT(*) FROM transactions WHERE account_id = $1;", account_id)
	err = row.Scan(&transactionResponse.Count)
	if err != nil {
		tx.Rollback()
		return transactionResponse, fmt.Errorf(errors_handler.DB004)
	}

	rows, err := tx.Query("SELECT * FROM transactions WHERE account_id = $1 AND id <> $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4;", account_id, uuid.UUID{}, limit, offset)
	if err != nil {
		tx.Rollback()
		return transactionResponse, fmt.Errorf(errors_handler.DB005)
	}

	for rows.Next() {
		t := Transaction{}
		err = rows.Scan(&t.ID, &t.AccountId, &t.PersonId, &t.Date, &t.Amount, &t.Description, &t.Balance, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			tx.Rollback()
			return transactionResponse, fmt.Errorf(errors_handler.DB006)
		}
		t.PersonName, _ = persons.GetPersonsName(t.PersonId)
		transactionResponse.Transactions = append(transactionResponse.Transactions, t)
	}

	transactionResponse.Limit = limit
	transactionResponse.Offset = offset

	err = tx.Commit()
	if err != nil {
		return transactionResponse, fmt.Errorf(errors_handler.DB003)
	}

	return transactionResponse, nil
}

func CreateTransaction(fields TransactionFields) (Transaction, error) {
	tr := Transaction{}

	if fields.AccountId == (uuid.UUID{}) {
		return tr, fmt.Errorf(errors_handler.TR001)
	}

	var oldBalance float64 = 0
	var updatedBalance float64 = 0

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf(errors_handler.DB002)
	}
	row := tx.QueryRow(`SELECT balance FROM money_accounts WHERE id = $1;`, fields.AccountId)
	err = row.Scan(&oldBalance)
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf(errors_handler.TR001)
	}

	newBalance := utility.RoundToTwoDecimalPlaces(oldBalance + utility.RoundToTwoDecimalPlaces(fields.Amount))
	if newBalance < 0 {
		tx.Rollback()
		return tr, fmt.Errorf(errors_handler.TR002)
	}

	row = tx.QueryRow(`UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING balance;`, newBalance, fields.AccountId)
	err = row.Scan(&updatedBalance)
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf(errors_handler.TR005)
	}

	if newBalance != updatedBalance {
		tx.Rollback()
		return tr, fmt.Errorf("New balance and updated balance missmatch, oldBalance = %v, newBalance = %v, updatedBalance = %v", oldBalance, newBalance, updatedBalance)
	}

	row = tx.QueryRow(`INSERT INTO transactions (account_id, person_id, date, amount, description, balance) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;`, fields.AccountId, fields.PersonId, fields.Date, fields.Amount, fields.Description, updatedBalance)
	err = row.Scan(&tr.ID, &tr.AccountId, &tr.PersonId, &tr.Date, &tr.Amount, &tr.Description, &tr.Balance, &tr.CreatedAt, &tr.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf(errors_handler.DB007)
	}

	tx.Commit()

	tr.PersonName, _ = persons.GetPersonsName(tr.PersonId)

	return tr, nil
}

func GetTransaction(transaction_id uuid.UUID) (Transaction, error) {
	t := Transaction{}
	row := database.DB.QueryRow("SELECT * FROM transactions WHERE id = $1;", transaction_id)
	err := row.Scan(&t.ID, &t.AccountId, &t.PersonId, &t.Date, &t.Amount, &t.Description, &t.Balance, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return t, fmt.Errorf(errors_handler.DB008)
	}
	t.PersonName, _ = persons.GetPersonsName(t.PersonId)
	return t, nil
}

func UpdateLastTransaction(transaction_id uuid.UUID, fields TransactionFields) (Transaction, error) {
	oldT := Transaction{}
	udT := Transaction{}
	updatedBalance := float64(0)

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return udT, fmt.Errorf(errors_handler.DB002)
	}

	row := tx.QueryRow("SELECT * FROM transactions ORDER BY created_at DESC LIMIT 1;")
	err = row.Scan(&oldT.ID, &oldT.AccountId, &oldT.PersonId, &oldT.Date, &oldT.Amount, &oldT.Description, &oldT.Balance, &oldT.CreatedAt, &oldT.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return udT, fmt.Errorf(errors_handler.TR004)
	}

	if oldT.ID != transaction_id {
		tx.Rollback()
		return udT, fmt.Errorf(errors_handler.TR003)
	}

	newBalance := utility.RoundToTwoDecimalPlaces(oldT.Balance - oldT.Amount + fields.Amount)

	if newBalance < 0 {
		tx.Rollback()
		return udT, fmt.Errorf(errors_handler.TR002)
	}

	row = tx.QueryRow(`UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING balance;`, newBalance, fields.AccountId)
	err = row.Scan(&updatedBalance)
	if err != nil {
		tx.Rollback()
		return udT, fmt.Errorf(errors_handler.TR005)
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

	err = tx.Commit()
	if err != nil {
		return udT, fmt.Errorf(errors_handler.DB003)
	}

	udT.PersonName, _ = persons.GetPersonsName(udT.PersonId)

	return udT, nil
}

func DeleteLastTransaction() (Transaction, error) {
	lT := Transaction{} // last transaction
	updatedBalance := float64(0)

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return lT, fmt.Errorf(errors_handler.DB002)
	}

	row := tx.QueryRow("DELETE FROM transactions WHERE id in (SELECT id FROM transactions ORDER BY created_at DESC LIMIT 1) RETURNING *;")
	err = row.Scan(&lT.ID, &lT.AccountId, &lT.PersonId, &lT.Date, &lT.Amount, &lT.Description, &lT.Balance, &lT.CreatedAt, &lT.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return lT, fmt.Errorf(errors_handler.TR004)
	}

	newBalance := utility.RoundToTwoDecimalPlaces(lT.Balance - lT.Amount)
	// This should never happend
	if newBalance < 0 {
		tx.Rollback()
		return lT, fmt.Errorf(errors_handler.TR002)
	}

	row = tx.QueryRow(`UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING balance;`, newBalance, lT.AccountId)
	err = row.Scan(&updatedBalance)
	if err != nil {
		tx.Rollback()
		return lT, fmt.Errorf(errors_handler.TR005)
	}

	err = tx.Commit()
	if err != nil {
		return lT, fmt.Errorf(errors_handler.DB003)
	}

	return lT, nil
}

func deleteAllTransactions() {
	database.DB.QueryRow("DELETE FROM transactions;")
}
