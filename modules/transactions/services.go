package transactions

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
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
		return transactionResponse, fmt.Errorf(errors_handler.UM001)
	}

	rows, err := tx.Query("SELECT * FROM transactions WHERE account_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;", account_id, limit, offset)
	if err != nil {
		tx.Rollback()
		return transactionResponse, fmt.Errorf(errors_handler.TR007)
	}

	for rows.Next() {
		t := Transaction{}
		err = rows.Scan(&t.ID, &t.AccountId, &t.PersonId, &t.Date, &t.Amount, &t.Description, &t.Balance, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			tx.Rollback()
			return transactionResponse, fmt.Errorf(errors_handler.TR009)
		}
		t.PersonName, _ = persons.GetPersonsName(t.PersonId)
		transactionResponse.Transactions = append(transactionResponse.Transactions, t)
	}

	tx.Commit()

	transactionResponse.Limit = limit
	transactionResponse.Offset = offset
	return transactionResponse, nil
}

func CreateTransaction(fields TransactionFields) (Transaction, error) {
	tr := Transaction{}

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

func DeleteLastTransaction() (TrashedTransaction, error) {
	dT := TrashedTransaction{} // deleted transaction
	lT := Transaction{}        // last transaction
	updatedBalance := float64(0)

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return dT, fmt.Errorf(errors_handler.DB002)
	}

	row := tx.QueryRow("DELETE FROM transactions WHERE id in (SELECT id FROM transactions ORDER BY created_at DESC LIMIT 1) RETURNING *;")
	err = row.Scan(&lT.ID, &lT.AccountId, &lT.PersonId, &lT.Date, &lT.Amount, &lT.Description, &lT.Balance, &lT.CreatedAt, &lT.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return dT, fmt.Errorf(errors_handler.TR004)
	}

	row = tx.QueryRow("INSERT INTO trashed_transactions VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;", lT.ID, lT.AccountId, lT.PersonId, lT.Date, lT.Amount, lT.Description, lT.CreatedAt, lT.UpdatedAt, time.Now())
	err = row.Scan(&dT.ID, &dT.AccountId, &dT.PersonId, &dT.Date, &dT.Amount, &dT.Description, &dT.CreatedAt, &dT.UpdatedAt, &dT.DeletedAt)
	if err != nil {
		tx.Rollback()
		return dT, fmt.Errorf(errors_handler.TR006)
	}

	newBalance := utility.RoundToTwoDecimalPlaces(lT.Balance - lT.Amount)
	// This should never happend
	if newBalance < 0 {
		tx.Rollback()
		return dT, fmt.Errorf(errors_handler.TR002)
	}

	row = tx.QueryRow(`UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING balance;`, newBalance, lT.AccountId)
	err = row.Scan(&updatedBalance)
	if err != nil {
		tx.Rollback()
		return dT, fmt.Errorf(errors_handler.TR005)
	}

	err = tx.Commit()
	if err != nil {
		return dT, fmt.Errorf(errors_handler.DB003)
	}

	return dT, nil
}

func GetTrashedTransactions() ([]TrashedTransaction, error) {
	trashed := []TrashedTransaction{}
	rows, err := database.DB.Query("SELECT * FROM trashed_transactions;")
	if err != nil {
		return trashed, fmt.Errorf(errors_handler.TR008)
	}

	for rows.Next() {
		tt := TrashedTransaction{}
		err = rows.Scan(&tt.ID, &tt.AccountId, &tt.PersonId, &tt.Date, &tt.Amount, &tt.Description, &tt.CreatedAt, &tt.UpdatedAt, &tt.DeletedAt)
		if err != nil {
			return trashed, fmt.Errorf(errors_handler.TR010)
		}
		tt.PersonName, _ = persons.GetPersonsName(tt.PersonId)
		trashed = append(trashed, tt)
	}

	return trashed, nil
}

func RestoreTrashedTransaction(trashed_transaction_id uuid.UUID) (Transaction, error) {
	rt := Transaction{}        // restored transaction
	dT := TrashedTransaction{} // deleted transaction
	updatedBalance := float64(0)

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return rt, fmt.Errorf(errors_handler.DB002)
	}

	row := tx.QueryRow("DELETE FROM trashed_transactions WHERE id = $1 RETURNING *;", trashed_transaction_id)
	err = row.Scan(&dT.ID, &dT.AccountId, &dT.PersonId, &dT.Date, &dT.Amount, &dT.Description, &dT.CreatedAt, &dT.UpdatedAt, &dT.DeletedAt)
	if err != nil {
		tx.Rollback()
		return rt, fmt.Errorf(errors_handler.TR011)
	}

	moneyAccount, err := money_accounts.GetOneMoneyAccount(dT.AccountId)
	if err != nil {
		tx.Rollback()
		return rt, err
	}

	newBalance := utility.RoundToTwoDecimalPlaces(moneyAccount.Balance + dT.Amount)
	if newBalance < 0 {
		tx.Rollback()
		return rt, fmt.Errorf(errors_handler.TR002)
	}

	row = tx.QueryRow("INSERT INTO transactions (id, account_id, person_id, date, amount, description, balance) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *", dT.ID, dT.AccountId, dT.PersonId, dT.Date, dT.Amount, dT.Description, newBalance)
	err = row.Scan(&rt.ID, &rt.AccountId, &rt.PersonId, &rt.Date, &rt.Amount, &rt.Description, &rt.Balance, &rt.CreatedAt, &rt.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return rt, fmt.Errorf(errors_handler.TR012)
	}

	row = tx.QueryRow("UPDATE money_accounts SET balance = $1 WHERE id = $2 RETURNING balance;", newBalance, rt.AccountId)
	err = row.Scan(&updatedBalance)
	if err != nil {
		tx.Rollback()
		return rt, fmt.Errorf(errors_handler.TR005)
	}

	err = tx.Commit()
	if err != nil {
		return rt, fmt.Errorf(errors_handler.DB003)
	}

	return rt, nil
}

func DeleteTrashedTransaction(trashed_transaction_id uuid.UUID) (TrashedTransaction, error) {
	tt := TrashedTransaction{} // trashed_transaction
	row := database.DB.QueryRow("DELETE FROM trashed_transactions WHERE id = $1 RETURNING *;", trashed_transaction_id)
	err := row.Scan(&tt.ID, &tt.AccountId, &tt.PersonId, &tt.Date, &tt.Amount, &tt.Description, &tt.CreatedAt, &tt.UpdatedAt, &tt.DeletedAt)
	if err != nil {
		return tt, fmt.Errorf(errors_handler.TR011)
	}
	return tt, nil
}

func deleteAllTransactions() {
	database.DB.QueryRow("DELETE FROM transactions;")
}

func deleteAllTrashedTransactions() {
	database.DB.QueryRow("DELETE FROM trashed_transactions;")
}

func resetTransactions() {
	deleteAllTransactions()
	deleteAllTrashedTransactions()
}
