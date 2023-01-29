package transactions

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/bills"
	"github.com/grabielcruz/transportation_back/modules/currencies"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/grabielcruz/transportation_back/modules/person_accounts"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/stretchr/testify/assert"
)

// create pending bill with transaction and closed it
// errors
func TestCloseBillService(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	m_account, err := money_accounts.CreateMoneyAccount(money_accounts.GenerateAccountFields())
	assert.Nil(t, err)
	person, err := persons.CreatePerson(persons.GeneratePersonFields())
	assert.Nil(t, err)

	p_account_fields := person_accounts.GeneratePersonAccountFields()

	// froce same currency on m_account and p_account
	p_account_fields.Currency = m_account.Currency
	p_account, err := person_accounts.CreatePersonAccount(person.ID, p_account_fields)
	assert.Nil(t, err)

	money_accounts.SetAccountsBalance(m_account.ID, initial_balance)
	_, err = currencies.CreateCurrency("ABC")
	assert.Nil(t, err)

	t.Run("Create pending bill and close it, zero fee", func(t *testing.T) {
		// this bill could be wheter to pay (negative) or to charge (positive)
		billFields := bills.GenerateBillFields(person.ID)
		// force same currency
		billFields.Currency = m_account.Currency

		pendingBill, err := bills.CreatePendingBill(billFields)
		assert.Nil(t, err)
		closedBill, err := ClosePendingBill(pendingBill.ID, m_account.ID, p_account.ID, time.Now(), 0)
		assert.Nil(t, err)

		// transaction should be closed by now
		closingTransaction, err := GetTransaction(closedBill.TransactionId)
		assert.Nil(t, err)

		// transaction amount and closed bill amount should be the same
		assert.Equal(t, closingTransaction.Amount, closedBill.Amount)

		// bill should not exist in pending bills
		_, err = bills.GetOnePendingBill(closedBill.ID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())

		// bill should exist on closed bills
		recallBill, err := bills.GetOneClosedBill(closedBill.ID)
		assert.Nil(t, err)
		assert.Equal(t, recallBill, closedBill)

		// closed bill and transaction should point to the same person
		assert.Equal(t, closingTransaction.PersonId, closedBill.PersonId)

		// transaction and person account should point to the same person
		assert.Equal(t, closingTransaction.PersonId, p_account.PersonId)

		// closed bill and account should have the same currency
		assert.Equal(t, closingTransaction.Currency, closedBill.Currency)

		// closed bill and person account should have the same currency
		assert.Equal(t, closingTransaction.Currency, p_account.Currency)

		// check account balance
		updated_account, err := money_accounts.GetOneMoneyAccount(m_account.ID)
		assert.Nil(t, err)
		assert.Equal(t, closingTransaction.Balance, updated_account.Balance)
	})

	deleteAllTransactions()
	bills.EmptyBills()
	money_accounts.SetAccountsBalance(m_account.ID, initial_balance)

	t.Run("Create pending bill and close it with fee", func(t *testing.T) {
		// this bill always will be negative, which means the money will always leave the account
		billFields := bills.GenerateBillToPayFields(person.ID)
		// force same currency
		billFields.Currency = m_account.Currency

		pendingBill, err := bills.CreatePendingBill(billFields)
		assert.Nil(t, err)
		closedBill, err := ClosePendingBill(pendingBill.ID, m_account.ID, p_account.ID, time.Now(), 0.003)
		assert.Nil(t, err)

		// transaction should be closed by now
		closingTransaction, err := GetTransaction(closedBill.TransactionId)
		assert.Nil(t, err)

		// transaction amount and closed bill amount should be the same
		assert.Equal(t, closingTransaction.Amount, closedBill.Amount)

		// bill should not exist in pending bills
		_, err = bills.GetOnePendingBill(closedBill.ID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())

		// bill should exist on closed bills
		recallBill, err := bills.GetOneClosedBill(closedBill.ID)
		assert.Nil(t, err)
		assert.Equal(t, recallBill, closedBill)

		// closed bill and transaction should point to the same person
		assert.Equal(t, closingTransaction.PersonId, closedBill.PersonId)

		// transaction and person account should point to the same person
		assert.Equal(t, closingTransaction.PersonId, p_account.PersonId)

		// closed bill and account should have the same currency
		assert.Equal(t, closingTransaction.Currency, closedBill.Currency)

		// closed bill and person account should have the same currency
		assert.Equal(t, closingTransaction.Currency, p_account.Currency)

		// check account balance
		updated_account, err := money_accounts.GetOneMoneyAccount(m_account.ID)
		assert.Nil(t, err)
		assert.Equal(t, closingTransaction.Balance, updated_account.Balance)
		// balance should be less by now, because money was spent
		assert.Less(t, updated_account.Balance, float64(initial_balance))
	})

	deleteAllTransactions()
	bills.EmptyBills()
	money_accounts.SetAccountsBalance(m_account.ID, initial_balance)

	t.Run("Error when closing bill to charge with fee", func(t *testing.T) {
		// this bill will always be positive, which means the money will enter the account, never leave it
		billFields := bills.GenerateBillToChargeFields(person.ID)
		// force same currency
		billFields.Currency = m_account.Currency

		pendingBill, err := bills.CreatePendingBill(billFields)
		assert.Nil(t, err)
		_, err = ClosePendingBill(pendingBill.ID, m_account.ID, p_account.ID, time.Now(), 0.003)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR014, err.Error())

	})

	deleteAllTransactions()
	bills.EmptyBills()
	money_accounts.SetAccountsBalance(m_account.ID, initial_balance)

	t.Run("Error when closing unexisting pending bill", func(t *testing.T) {
		// zero uuid
		_, err = ClosePendingBill(uuid.UUID{}, m_account.ID, p_account.ID, time.Now(), 0.003)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR012, err.Error())

		// rand uuid
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = ClosePendingBill(randId, m_account.ID, p_account.ID, time.Now(), 0.003)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR012, err.Error())
	})

	t.Run("Error when closing pending bill with unexisting money account", func(t *testing.T) {
		billFields := bills.GenerateBillToChargeFields(person.ID)
		// force same currency
		billFields.Currency = m_account.Currency

		pendingBill, err := bills.CreatePendingBill(billFields)
		assert.Nil(t, err)
		// zero uuid
		_, err = ClosePendingBill(pendingBill.ID, uuid.UUID{}, p_account.ID, time.Now(), 0.003)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR013, err.Error())

		// rand uuid
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = ClosePendingBill(pendingBill.ID, randId, p_account.ID, time.Now(), 0.003)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR013, err.Error())
	})

	t.Run("Error when closing pending bill with unexisting person account", func(t *testing.T) {
		billFields := bills.GenerateBillToChargeFields(person.ID)
		// force same currency
		billFields.Currency = m_account.Currency

		pendingBill, err := bills.CreatePendingBill(billFields)
		assert.Nil(t, err)
		// zero uuid
		_, err = ClosePendingBill(pendingBill.ID, m_account.ID, uuid.UUID{}, time.Now(), 0.003)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.PA002, err.Error())

		// rand uuid
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = ClosePendingBill(pendingBill.ID, m_account.ID, randId, time.Now(), 0.003)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.PA002, err.Error())
	})

	t.Run("Error when there is currency missmatch between money account and person account", func(t *testing.T) {
		p_account_fields := person_accounts.GeneratePersonAccountFields()

		// froce different currency on m_account and p_account
		p_account_fields.Currency = "ABC"
		p_account, err := person_accounts.CreatePersonAccount(person.ID, p_account_fields)
		assert.Nil(t, err)
		billFields := bills.GenerateBillToChargeFields(person.ID)
		// force same currency
		billFields.Currency = m_account.Currency

		pendingBill, err := bills.CreatePendingBill(billFields)
		assert.Nil(t, err)

		_, err = ClosePendingBill(pendingBill.ID, m_account.ID, p_account.ID, time.Now(), 0.003)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR011, err.Error())
	})

	t.Run("Error when there is currency missmatch between money account and pending bill", func(t *testing.T) {
		billFields := bills.GenerateBillToChargeFields(person.ID)
		// force different currency
		billFields.Currency = "ABC"

		pendingBill, err := bills.CreatePendingBill(billFields)
		assert.Nil(t, err)

		_, err = ClosePendingBill(pendingBill.ID, m_account.ID, p_account.ID, time.Now(), 0.003)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR011, err.Error())
	})

	t.Run("Error when pending bill and person account have a different owner", func(t *testing.T) {
		person2, err := persons.CreatePerson(persons.GeneratePersonFields())
		assert.Nil(t, err)
		billFields := bills.GenerateBillToChargeFields(person2.ID)
		// force same currency
		billFields.Currency = m_account.Currency
		pendingBill, err := bills.CreatePendingBill(billFields)
		assert.Nil(t, err)

		_, err = ClosePendingBill(pendingBill.ID, m_account.ID, p_account.ID, time.Now(), 0.003)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR010, err.Error())
	})
}
