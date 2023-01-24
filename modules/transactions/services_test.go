package transactions

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/bills"
	"github.com/grabielcruz/transportation_back/modules/config"
	"github.com/grabielcruz/transportation_back/modules/money_accounts"
	"github.com/grabielcruz/transportation_back/modules/person_accounts"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/grabielcruz/transportation_back/utility"
	"github.com/stretchr/testify/assert"
)

func TestTransactionServices(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	account, err := money_accounts.CreateMoneyAccount(money_accounts.GenerateAccountFields())
	assert.Nil(t, err)
	person, err := persons.CreatePerson(persons.GeneratePersonFields())
	assert.Nil(t, err)

	t.Run("Get transaction response with zero transactions initially", func(t *testing.T) {
		transactions, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Len(t, transactions.Transactions, 0)
		assert.Equal(t, transactions.Count, 0)
		assert.Equal(t, transactions.Offset, config.Offset)
		assert.Equal(t, transactions.Limit, config.Limit)
	})

	// test error when creating transaction without a person

	t.Run("Create one transaction with a person", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		newTransaction, err := CreateTransaction(transactionFields, person.ID, true)
		assert.Nil(t, err)
		updatedAccount, err := money_accounts.GetOneMoneyAccount(newTransaction.AccountId)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.Balance, updatedAccount.Balance)
		assert.Equal(t, newTransaction.AccountId, updatedAccount.ID)
		assert.Equal(t, newTransaction.PersonName, person.Name)
		assert.Equal(t, newTransaction.PersonId, person.ID)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction with a person and a person account", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		// person account
		personAccountFields := person_accounts.GeneratePersonAccountFields()
		// force same currency
		personAccountFields.Currency = account.Currency
		newPersonAccount, err := person_accounts.CreatePersonAccount(person.ID, personAccountFields)
		assert.Nil(t, err)
		//
		transactionFields.PersonAccountId = newPersonAccount.ID
		newTransaction, err := CreateTransaction(transactionFields, person.ID, true)
		assert.Nil(t, err)
		updatedAccount, err := money_accounts.GetOneMoneyAccount(newTransaction.AccountId)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.Balance, updatedAccount.Balance)
		assert.Equal(t, newTransaction.AccountId, updatedAccount.ID)
		assert.Equal(t, newTransaction.PersonName, person.Name)
		assert.Equal(t, newTransaction.PersonId, person.ID)
		// persons account
		assert.Equal(t, newTransaction.PersonAccountId, newPersonAccount.ID)
		assert.Equal(t, newTransaction.PersonAccountName, newPersonAccount.Name)
		assert.Equal(t, newTransaction.PersonAccountDescription, newPersonAccount.Description)
		assert.Equal(t, newTransaction.Currency, newPersonAccount.Currency)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	person_accounts.DeleteAllPersonAccounts()
	deleteAllTransactions()

	t.Run("Error when creating a transaction with an unexisting person account different than zero", func(t *testing.T) {

	})

	t.Run("Error when creating transaction with unexisting account", func(t *testing.T) {
		zeroId := uuid.UUID{}
		transactionFields := GenerateTransactionFields(zeroId)
		_, err := CreateTransaction(transactionFields, zeroId, true)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR001, err.Error())
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when generating negative balance", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.Amount *= -1
		_, err := CreateTransaction(transactionFields, person.ID, true)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR002, err.Error())
		updatedAccount, err := money_accounts.GetOneMoneyAccount(transactionFields.AccountId)
		assert.Nil(t, err)
		// accounts balance should remain unmodified, which means it is equal to zero
		assert.Equal(t, float64(0), updatedAccount.Balance)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction and get it in paginated response", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		newTransaction, err := CreateTransaction(transactionFields, person.ID, true)
		assert.Nil(t, err)
		transactions, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Equal(t, config.Offset, transactions.Offset)
		assert.Equal(t, config.Limit, transactions.Limit)
		assert.Equal(t, 1, transactions.Count)
		assert.Equal(t, newTransaction, transactions.Transactions[0])
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction without fee and get it with single response", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		newTransaction, err := CreateTransaction(transactionFields, person.ID, true)
		assert.Nil(t, err)
		transaction, err := GetTransaction(newTransaction.ID)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.ID, transaction.ID)
		assert.Equal(t, newTransaction.AccountId, transaction.AccountId)
		assert.Equal(t, newTransaction.Amount, transaction.Amount)
		assert.Equal(t, newTransaction.Balance, transaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt, transaction.CreatedAt)
		assert.Equal(t, newTransaction.UpdatedAt, transaction.UpdatedAt)
		assert.Equal(t, newTransaction.Date, transaction.Date)
		assert.Equal(t, newTransaction.Description, transaction.Description)
		assert.Equal(t, newTransaction.PersonId, transaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, transaction.PersonName)
		assert.Equal(t, newTransaction.Currency, account.Currency)
		// fee stuff
		assert.Equal(t, newTransaction.Fee, transaction.Fee)
		assert.Equal(t, newTransaction.AmountWithFee, transaction.AmountWithFee)
		assert.Equal(t, newTransaction.AmountWithFee, utility.RoundToTwoDecimalPlaces(newTransaction.Amount*(1+newTransaction.Fee)))
		assert.Equal(t, transaction.AmountWithFee, utility.RoundToTwoDecimalPlaces(newTransaction.Amount*(1+newTransaction.Fee)))
		// these uuids should be zero
		assert.Equal(t, newTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, newTransaction.RevertBillId, uuid.UUID{})
		assert.Equal(t, transaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, transaction.RevertBillId, uuid.UUID{})
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction with fee and get it with single response", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.Fee = utility.GetRandomFee()
		newTransaction, err := CreateTransaction(transactionFields, person.ID, true)
		assert.Nil(t, err)
		transaction, err := GetTransaction(newTransaction.ID)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.ID, transaction.ID)
		assert.Equal(t, newTransaction.AccountId, transaction.AccountId)
		assert.Equal(t, newTransaction.Amount, transaction.Amount)
		assert.Equal(t, newTransaction.Balance, transaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt, transaction.CreatedAt)
		assert.Equal(t, newTransaction.UpdatedAt, transaction.UpdatedAt)
		assert.Equal(t, newTransaction.Date, transaction.Date)
		assert.Equal(t, newTransaction.Description, transaction.Description)
		assert.Equal(t, newTransaction.PersonId, transaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, transaction.PersonName)
		assert.Equal(t, newTransaction.Currency, account.Currency)
		// fee stuff
		assert.Equal(t, newTransaction.Fee, transaction.Fee)
		assert.Equal(t, newTransaction.AmountWithFee, transaction.AmountWithFee)
		assert.Equal(t, newTransaction.AmountWithFee, utility.RoundToTwoDecimalPlaces(newTransaction.Amount*(1+newTransaction.Fee)))
		assert.Equal(t, transaction.AmountWithFee, utility.RoundToTwoDecimalPlaces(newTransaction.Amount*(1+newTransaction.Fee)))
		// these uuids should be zero
		assert.Equal(t, newTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, newTransaction.RevertBillId, uuid.UUID{})
		assert.Equal(t, transaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, transaction.RevertBillId, uuid.UUID{})
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("It should create transaction with person zero when not blocked", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		newTransaction, err := CreateTransaction(transactionFields, uuid.UUID{}, false)
		assert.Nil(t, err)
		transaction, err := GetTransaction(newTransaction.ID)
		assert.Nil(t, err)
		assert.Equal(t, newTransaction.ID, transaction.ID)
		assert.Equal(t, newTransaction.AccountId, transaction.AccountId)
		assert.Equal(t, newTransaction.Amount, transaction.Amount)
		assert.Equal(t, newTransaction.Balance, transaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt, transaction.CreatedAt)
		assert.Equal(t, newTransaction.UpdatedAt, transaction.UpdatedAt)
		assert.Equal(t, newTransaction.Date, transaction.Date)
		assert.Equal(t, newTransaction.Description, transaction.Description)
		assert.Equal(t, newTransaction.PersonId, transaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, transaction.PersonName)
		assert.Equal(t, newTransaction.Currency, account.Currency)
		// fee stuff
		assert.Equal(t, newTransaction.Fee, transaction.Fee)
		assert.Equal(t, newTransaction.AmountWithFee, transaction.AmountWithFee)
		assert.Equal(t, newTransaction.AmountWithFee, utility.RoundToTwoDecimalPlaces(newTransaction.Amount*(1+newTransaction.Fee)))
		assert.Equal(t, transaction.AmountWithFee, utility.RoundToTwoDecimalPlaces(newTransaction.Amount*(1+newTransaction.Fee)))
		// these uuids should be zero
		assert.Equal(t, newTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, newTransaction.RevertBillId, uuid.UUID{})
		assert.Equal(t, transaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, transaction.RevertBillId, uuid.UUID{})
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when creating transaction without a person when blocked", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.Fee = utility.GetRandomFee()
		_, err := CreateTransaction(transactionFields, uuid.UUID{}, true)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR007, err.Error())
	})

	t.Run("Error when getting non registered transaction", func(t *testing.T) {
		// with zero uuid
		_, err := GetTransaction(uuid.UUID{})
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())

		// with random uuid
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = GetTransaction(randId)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Error when creating transaction with amount zero", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.Amount = float64(0)
		_, err := CreateTransaction(transactionFields, person.ID, true)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR008, err.Error())
	})

	t.Run("Error when creating transaction with negative fee", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.Fee = -0.05
		_, err := CreateTransaction(transactionFields, person.ID, true)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR009, err.Error())
	})

	t.Run("Error when creating transaction with a fee greater than one", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.Fee = 1.05
		_, err := CreateTransaction(transactionFields, person.ID, true)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.TR009, err.Error())
	})

	t.Run("Execute 100 transactions without fee and get accounts balance right", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(100)
		sum := utility.GetSumOfAmounts(amounts)
		for _, v := range amounts {
			personId := person.ID
			transactionFields := GenerateTransactionFields(account.ID)
			transactionFields.Fee = 0
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields, personId, true)
			assert.Nil(t, err)
		}
		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)
		assert.Equal(t, sum, updatedAccount.Balance)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Execute 100 transactions with fee of 5% and get accounts balance right", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(100)
		sum := utility.GetSumOfAmountsWithFee(amounts, 0.05)
		for _, v := range amounts {
			personId := person.ID
			transactionFields := GenerateTransactionFields(account.ID)
			transactionFields.Amount = v
			transactionFields.Fee = 0.05
			_, err := CreateTransaction(transactionFields, personId, true)
			assert.Nil(t, err)
		}
		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)
		assert.Equal(t, sum, updatedAccount.Balance)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Execute 10 transaction and the first transaction in the slice should be the last one executed", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(10)
		for _, v := range amounts {
			personId := person.ID

			transactionFields := GenerateTransactionFields(account.ID)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields, personId, true)
			assert.Nil(t, err)
		}
		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)
		transactions, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Equal(t, transactions.Transactions[0].Balance, updatedAccount.Balance)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Execute 51 transaction and get in last page the initial transaction, and count equal 51", func(t *testing.T) {
		amounts := utility.GetSliceOfAmounts(51)
		for _, v := range amounts {
			personId := person.ID
			transactionFields := GenerateTransactionFields(account.ID)
			transactionFields.Amount = v
			_, err := CreateTransaction(transactionFields, personId, true)
			assert.Nil(t, err)
		}
		transactions, err := GetTransactions(account.ID, config.Limit, 50)
		assert.Nil(t, err)
		assert.Equal(t, transactions.Transactions[0].Balance, transactions.Transactions[0].AmountWithFee)
		assert.Equal(t, 51, transactions.Count)
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction without fee, it creates a pending bill. When deletion, pending bill also is deleted", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.Fee = 0
		newTransaction, err := CreateTransaction(transactionFields, person.ID, true)
		assert.Nil(t, err)
		// pending bill
		newPendingBill, err := bills.GetOneBill(newTransaction.PendingBillId)
		assert.Nil(t, err)
		assert.Equal(t, newPendingBill.ParentTransactionId, newTransaction.ID)
		assert.Equal(t, newPendingBill.Amount, newTransaction.AmountWithFee)
		assert.Equal(t, newPendingBill.Amount, newTransaction.Amount) // fee 0
		assert.Equal(t, newPendingBill.Date, newTransaction.Date)
		assert.Equal(t, newPendingBill.Description, newTransaction.Description)

		// delete
		deletedLastTransaction, err := DeleteLastTransaction()
		assert.Nil(t, err)

		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)

		assert.Equal(t, newTransaction.ID, deletedLastTransaction.ID)
		assert.Equal(t, newTransaction.AccountId, deletedLastTransaction.AccountId)
		assert.Equal(t, newTransaction.Amount, deletedLastTransaction.Amount)
		assert.Equal(t, newTransaction.Balance, deletedLastTransaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt, deletedLastTransaction.CreatedAt)
		assert.Equal(t, newTransaction.UpdatedAt, deletedLastTransaction.UpdatedAt)
		assert.Equal(t, newTransaction.Date, deletedLastTransaction.Date)
		assert.Equal(t, newTransaction.Description, deletedLastTransaction.Description)
		assert.Equal(t, newTransaction.PersonId, deletedLastTransaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, deletedLastTransaction.PersonName)
		assert.Equal(t, float64(0), updatedAccount.Balance)
		assert.Equal(t, newTransaction.Currency, account.Currency)
		// these uuids should be zero
		assert.Equal(t, newTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, newTransaction.RevertBillId, uuid.UUID{})
		// assert.Equal(t, deletedLastTransaction.PendingBillId, uuid.UUID{})
		assert.Equal(t, deletedLastTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, deletedLastTransaction.RevertBillId, uuid.UUID{})

		transactions, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Equal(t, config.Offset, transactions.Offset)
		assert.Equal(t, config.Limit, transactions.Limit)
		assert.Equal(t, 0, transactions.Count)
		assert.Len(t, transactions.Transactions, 0)

		// pending bill also deleted
		_, err = bills.GetOneBill(newTransaction.PendingBillId)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Create one transaction with fee and delete it", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.Fee = utility.GetRandomFee()
		newTransaction, err := CreateTransaction(transactionFields, person.ID, true)
		assert.Nil(t, err)

		// pending bill
		newPendingBill, err := bills.GetOneBill(newTransaction.PendingBillId)
		assert.Nil(t, err)
		assert.Equal(t, newPendingBill.ParentTransactionId, newTransaction.ID)
		assert.LessOrEqual(t, newPendingBill.Amount, newTransaction.AmountWithFee)
		assert.Equal(t, newPendingBill.Amount, newTransaction.Amount)
		assert.Equal(t, newPendingBill.Date, newTransaction.Date)
		assert.Equal(t, newPendingBill.Description, newTransaction.Description)

		deletedLastTransaction, err := DeleteLastTransaction()
		assert.Nil(t, err)
		updatedAccount, err := money_accounts.GetOneMoneyAccount(account.ID)
		assert.Nil(t, err)

		assert.Equal(t, newTransaction.ID, deletedLastTransaction.ID)
		assert.Equal(t, newTransaction.AccountId, deletedLastTransaction.AccountId)
		assert.Equal(t, newTransaction.Amount, deletedLastTransaction.Amount)
		assert.Equal(t, newTransaction.Balance, deletedLastTransaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt, deletedLastTransaction.CreatedAt)
		assert.Equal(t, newTransaction.UpdatedAt, deletedLastTransaction.UpdatedAt)
		assert.Equal(t, newTransaction.Date, deletedLastTransaction.Date)
		assert.Equal(t, newTransaction.Description, deletedLastTransaction.Description)
		assert.Equal(t, newTransaction.PersonId, deletedLastTransaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, deletedLastTransaction.PersonName)
		assert.Equal(t, float64(0), updatedAccount.Balance)
		assert.Equal(t, newTransaction.Currency, account.Currency)
		// these uuids should be zero
		assert.Equal(t, newTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, newTransaction.RevertBillId, uuid.UUID{})
		// assert.Equal(t, deletedLastTransaction.PendingBillId, uuid.UUID{})
		assert.Equal(t, deletedLastTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, deletedLastTransaction.RevertBillId, uuid.UUID{})

		transactions, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Equal(t, config.Offset, transactions.Offset)
		assert.Equal(t, config.Limit, transactions.Limit)
		assert.Equal(t, 0, transactions.Count)
		assert.Len(t, transactions.Transactions, 0)

		// pending bill also deleted
		_, err = bills.GetOneBill(newTransaction.PendingBillId)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	money_accounts.ResetAccountsBalance(account.ID)
	deleteAllTransactions()

	t.Run("Error when deleting last transaction with no transactions", func(t *testing.T) {
		_, err := DeleteLastTransaction()
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Error when deleting pending bill associated with transaction", func(t *testing.T) {
		transactionFields := GenerateTransactionFields(account.ID)
		transactionFields.Fee = utility.GetRandomFee()
		newTransaction, err := CreateTransaction(transactionFields, person.ID, true)
		assert.Nil(t, err)
		// this deletion should be forbidden
		_, err = bills.DeleteBill(newTransaction.PendingBillId)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.BL003, err.Error())
		sameTransaction, err := GetTransaction(newTransaction.ID)
		assert.Nil(t, err)

		assert.Equal(t, newTransaction.ID, sameTransaction.ID)
		assert.Equal(t, newTransaction.AccountId, sameTransaction.AccountId)
		assert.Equal(t, newTransaction.Amount, sameTransaction.Amount)
		assert.Equal(t, newTransaction.Balance, sameTransaction.Balance)
		assert.Equal(t, newTransaction.CreatedAt, sameTransaction.CreatedAt)
		assert.Equal(t, newTransaction.UpdatedAt, sameTransaction.UpdatedAt)
		assert.Equal(t, newTransaction.Date, sameTransaction.Date)
		assert.Equal(t, newTransaction.Description, sameTransaction.Description)
		assert.Equal(t, newTransaction.PersonId, sameTransaction.PersonId)
		assert.Equal(t, newTransaction.PersonName, sameTransaction.PersonName)
		assert.Equal(t, newTransaction.Currency, account.Currency)
		// these uuids should be zero
		assert.Equal(t, newTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, newTransaction.RevertBillId, uuid.UUID{})
		// assert.Equal(t, sameTransaction.PendingBillId, uuid.UUID{})
		assert.Equal(t, sameTransaction.ClosedBillId, uuid.UUID{})
		assert.Equal(t, sameTransaction.RevertBillId, uuid.UUID{})

		transactions, err := GetTransactions(account.ID, config.Limit, config.Offset)
		assert.Nil(t, err)
		assert.Equal(t, config.Offset, transactions.Offset)
		assert.Equal(t, config.Limit, transactions.Limit)
		assert.Equal(t, 1, transactions.Count)
		assert.Len(t, transactions.Transactions, 1)
	})

	// at the end of all transactions services tests
	money_accounts.DeleteAllMoneyAccounts()
	persons.DeleteAllPersons()
}
