package transactions

import (
	"path/filepath"
	"testing"

	"github.com/grabielcruz/transportation_back/database"
)

// create pending bill with transaction and closed it
// errors
func TestCloseBillService(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	// account, err := money_accounts.CreateMoneyAccount(money_accounts.GenerateAccountFields())
	// assert.Nil(t, err)
	// person, err := persons.CreatePerson(persons.GeneratePersonFields())
	// assert.Nil(t, err)
	// person_account, err := person_accounts.CreatePersonAccount(person.ID, person_accounts.GeneratePersonAccountFields())
	// assert.Nil(t, err)

	// t.Run("Create pending bill and close it, zero fee", func(t *testing.T) {
	// 	// this bill always will be positive, which means the money will enter the account, never leave it
	// 	billFields := bills.GenerateBillFields(person.ID)
	// 	// force same currency
	// 	billFields.Currency = account.Currency

	// 	pendingBill, err := bills.CreatePendingBill(billFields)
	// 	assert.Nil(t, err)
	// 	closedBill, err := ClosePendingBill(pendingBill.ID, account.ID, person_account.ID, time.Now(), 0)
	// 	assert.Nil(t, err)

	// 	// transaction should be closed by now
	// 	closingTransaction, err := GetTransaction(closedBill.TransactionId)
	// 	assert.Nil(t, err)

	// 	// transaction amount and closed bill amount should be the same
	// 	assert.Equal(t, closingTransaction.Amount, closedBill.Amount)

	// 	// bill should not exist in pending bills
	// 	_, err = bills.GetOnePendingBill(closedBill.ID)
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, errors_handler.DB001, err.Error())

	// 	// bill should exist on closed bills
	// 	recallBill, err := bills.GetOneClosedBill(closedBill.ID)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, recallBill, closedBill)

	// 	// closed bill and transaction should point to the same person
	// 	assert.Equal(t, closingTransaction.PersonId, closedBill.PersonId)

	// 	// transaction person account and person account should point to the same person
	// 	assert.Equal(t, closingTransaction.PersonId, person_account.PersonId)

	// 	// closed bill and account should have the same currency
	// 	assert.Equal(t, closingTransaction.Currency, closedBill.Currency)

	// 	// closed bill and person account should have the same currency
	// 	assert.Equal(t, closingTransaction.Currency, person_account.Currency)

	// 	// check account balance
	// 	updated_account, err := money_accounts.GetOneMoneyAccount(account.ID)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, closingTransaction.Balance, updated_account.Balance)
	// })

	// // complex relations, between tests reset database
	// database.CreateTables(sqlPath)

}
