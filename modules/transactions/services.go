package transactions

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/bills"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/grabielcruz/transportation_back/modules/person_accounts"
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

	row := tx.QueryRow("SELECT COUNT(*) FROM transactions WHERE account_id = $1 AND id <> $2;", account_id, uuid.UUID{})
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
		err = rows.Scan(&t.ID, &t.AccountId, &t.PersonId, &t.PersonAccountId, &t.PersonAccountName, &t.PersonAccountDescription, &t.Date, &t.Amount, &t.Fee, &t.AmountWithFee, &t.Description, &t.Balance, &t.PendingBillId, &t.ClosedBillId, &t.RevertBillId, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			tx.Rollback()
			return transactionResponse, fmt.Errorf(errors_handler.DB005)
		}
		t.PersonName, err = persons.GetPersonsName(t.PersonId)
		if err != nil {
			errors_handler.HandleError(err)
		}
		t.Currency, err = money_accounts.GetAccountsCurrency(account_id)
		if err != nil {
			errors_handler.HandleError(err)
		}
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

// CreateTransaction will throw an error when person_id is zero uuid and
// will always creates a pending bill when the property block_zero_person is set to true,
// otherwise in case of use of zero person uuid,t should register a transaction with zero person uuid,
// and it will not create a new pending bill
func CreateTransaction(fields TransactionFields, person_id uuid.UUID, block_zero_person bool) (Transaction, error) {
	tr := Transaction{}
	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf(errors_handler.DB002)
	}

	tr, err = createTransaction(tx, fields, person_id, block_zero_person)
	if err != nil {
		tx.Rollback()
		return tr, err
	}

	// create pending bill from transaction
	billFields := bills.BillFields{
		PersonId:            tr.PersonId,
		Date:                tr.Date,
		Description:         tr.Description,
		Currency:            tr.Currency,
		Amount:              tr.Amount,
		ParentTransactionId: tr.ID,
		ParentBillCrossId:   uuid.UUID{},
	}

	bill_id := uuid.UUID{}
	row := tx.QueryRow("INSERT INTO pending_bills (person_id, date, description, currency, amount, parent_transaction_id, parent_bill_cross_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;", billFields.PersonId, billFields.Date, billFields.Description, billFields.Currency, billFields.Amount, billFields.ParentTransactionId, billFields.ParentBillCrossId)
	err = row.Scan(&bill_id)
	if err != nil {
		return tr, errors_handler.MapDBErrors(err)
	}

	row = tx.QueryRow("UPDATE transactions SET pending_bill_id = $1 WHERE id = $2 RETURNING pending_bill_id;", bill_id, tr.ID)
	err = row.Scan(&tr.PendingBillId)
	if err != nil {
		return tr, errors_handler.MapDBErrors(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf(errors_handler.DB003)
	}
	return tr, nil
}

// This function is used by two separate handlers
func createTransaction(tx *sql.Tx, fields TransactionFields, person_id uuid.UUID, block_zero_person bool) (Transaction, error) {
	tr := Transaction{}

	personAccount, err := person_accounts.GetOnePersonAccount(fields.PersonAccountId)
	if err != nil {
		// error valid only when the id is not null
		if fields.PersonAccountId != (uuid.UUID{}) {
			return tr, err
		}
	}

	if fields.AccountId == (uuid.UUID{}) {
		return tr, fmt.Errorf(errors_handler.TR001)
	}

	if block_zero_person && person_id == (uuid.UUID{}) {
		return tr, fmt.Errorf(errors_handler.TR007)
	}

	if fields.Amount == float64(0) {
		return tr, fmt.Errorf(errors_handler.TR008)
	}

	if fields.Fee < float64(0) || fields.Fee > float64(1) {
		return tr, fmt.Errorf(errors_handler.TR009)
	}

	transactionCurrency, err := money_accounts.GetAccountsCurrency(fields.AccountId)
	if err != nil {
		errors_handler.HandleError(err)
	}

	// person account specified
	if fields.PersonAccountId != (uuid.UUID{}) {

		// account does not belong to the person
		if personAccount.PersonId != person_id {
			return tr, fmt.Errorf(errors_handler.TR010)
		}
		// currency mismatch
		if personAccount.Currency != transactionCurrency {
			return tr, fmt.Errorf(errors_handler.TR011)
		}
	}

	var oldBalance float64 = 0
	var updatedBalance float64 = 0

	row := tx.QueryRow(`SELECT balance FROM money_accounts WHERE id = $1;`, fields.AccountId)
	err = row.Scan(&oldBalance)
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf(errors_handler.TR001)
	}
	amount := utility.RoundToTwoDecimalPlaces(fields.Amount)
	fee := utility.RoundToTwoDecimalPlaces(fields.Fee)
	amountWithFee := amount * (1 + fee)
	newBalance := utility.RoundToTwoDecimalPlaces(oldBalance + utility.RoundToTwoDecimalPlaces(amountWithFee))
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
		return tr, fmt.Errorf(errors_handler.TR006, oldBalance, newBalance, updatedBalance)
	}

	row = tx.QueryRow(`INSERT INTO transactions (account_id, person_id, person_account_id, person_account_name, person_account_description, date, amount, fee, amount_with_fee, description, balance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;`,
		fields.AccountId, person_id, personAccount.ID, personAccount.Name, personAccount.Description, fields.Date, fields.Amount, fields.Fee, amountWithFee, fields.Description, updatedBalance)
	err = row.Scan(&tr.ID, &tr.AccountId, &tr.PersonId, &tr.PersonAccountId, &tr.PersonAccountName, &tr.PersonAccountDescription, &tr.Date, &tr.Amount, &tr.Fee, &tr.AmountWithFee, &tr.Description, &tr.Balance, &tr.PendingBillId, &tr.ClosedBillId, &tr.RevertBillId, &tr.CreatedAt, &tr.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return tr, fmt.Errorf(errors_handler.DB007)
	}

	tr.PersonName, err = persons.GetPersonsName(tr.PersonId)
	if err != nil {
		errors_handler.HandleError(err)
	}
	tr.Currency, err = money_accounts.GetAccountsCurrency(tr.AccountId)
	if err != nil {
		errors_handler.HandleError(err)
	}

	return tr, nil
}

func GetTransaction(transaction_id uuid.UUID) (Transaction, error) {
	t := Transaction{}
	if transaction_id == (uuid.UUID{}) {
		return t, fmt.Errorf(errors_handler.DB001)
	}
	if transaction_id == (uuid.UUID{}) {
		return t, fmt.Errorf(errors_handler.DB001)
	}
	row := database.DB.QueryRow("SELECT * FROM transactions WHERE id = $1;", transaction_id)
	err := row.Scan(&t.ID, &t.AccountId, &t.PersonId, &t.PersonAccountId, &t.PersonAccountName, &t.PersonAccountDescription, &t.Date, &t.Amount, &t.Fee, &t.AmountWithFee, &t.Description, &t.Balance, &t.PendingBillId, &t.ClosedBillId, &t.RevertBillId, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return t, fmt.Errorf(errors_handler.DB001)
	}
	t.PersonName, err = persons.GetPersonsName(t.PersonId)
	if err != nil {
		errors_handler.HandleError(err)
	}
	t.Currency, err = money_accounts.GetAccountsCurrency(t.AccountId)
	if err != nil {
		errors_handler.HandleError(err)
	}
	return t, nil
}

func DeleteLastTransaction() (Transaction, error) {
	lT := Transaction{} // last transaction
	updatedBalance := float64(0)

	tx, err := database.DB.Begin()
	if err != nil {
		tx.Rollback()
		return lT, fmt.Errorf(errors_handler.DB002)
	}

	row := tx.QueryRow("DELETE FROM transactions WHERE id in (SELECT id FROM transactions WHERE id <> $1 ORDER BY created_at DESC LIMIT 1) RETURNING *;", uuid.UUID{})
	err = row.Scan(&lT.ID, &lT.AccountId, &lT.PersonId, &lT.PersonAccountId, &lT.PersonAccountName, &lT.PersonAccountDescription, &lT.Date, &lT.Amount, &lT.Fee, &lT.AmountWithFee, &lT.Description, &lT.Balance, &lT.PendingBillId, &lT.ClosedBillId, &lT.RevertBillId, &lT.CreatedAt, &lT.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return lT, fmt.Errorf(errors_handler.DB001)
	}

	newBalance := utility.RoundToTwoDecimalPlaces(lT.Balance - lT.AmountWithFee)
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

	lT.PersonName, err = persons.GetPersonsName(lT.PersonId)
	if err != nil {
		errors_handler.HandleError(err)
	}
	lT.Currency, err = money_accounts.GetAccountsCurrency(lT.AccountId)
	if err != nil {
		errors_handler.HandleError(err)
	}

	err = tx.Commit()
	if err != nil {
		return lT, fmt.Errorf(errors_handler.DB003)
	}

	return lT, nil
}

func deleteAllTransactions() {
	database.DB.QueryRow("DELETE FROM transactions WHERE id <> $1;", uuid.UUID{})
	bills.EmptyBills()
}
