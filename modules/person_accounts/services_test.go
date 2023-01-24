package person_accounts

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/grabielcruz/transportation_back/modules/persons"
	"github.com/stretchr/testify/assert"
)

func TestPersonAccountsService(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()
	person, err := persons.CreatePerson(persons.GeneratePersonFields())
	assert.Nil(t, err)

	t.Run("Get zero accounts for a person initially", func(t *testing.T) {
		personAccounts, err := GetPersonAccounts(person.ID)
		assert.Nil(t, err)
		assert.Len(t, personAccounts, 0)
	})

	t.Run("Create a person account", func(t *testing.T) {
		personAccountFields := GeneratePersonAccountFields()
		newAccount, err := CreatePersonAccount(person.ID, personAccountFields)
		assert.Nil(t, err)
		assert.Equal(t, newAccount.PersonId, person.ID)
		assert.Equal(t, newAccount.Name, personAccountFields.Name)
		assert.Equal(t, newAccount.Description, personAccountFields.Description)
		assert.Equal(t, newAccount.Currency, personAccountFields.Currency)
	})

	DeleteAllPersonAccounts()

	t.Run("Create two person accounts and get an slice of person accounts", func(t *testing.T) {
		fields1 := GeneratePersonAccountFields()
		fields2 := GeneratePersonAccountFields()
		account1, err := CreatePersonAccount(person.ID, fields1)
		assert.Nil(t, err)
		account2, err := CreatePersonAccount(person.ID, fields2)
		assert.Nil(t, err)
		personAccounts, err := GetPersonAccounts(person.ID)
		assert.Nil(t, err)
		assert.Len(t, personAccounts, 2)
		assert.Equal(t, account1, personAccounts[0])
		assert.Equal(t, account2, personAccounts[1])
	})

	DeleteAllPersonAccounts()

	t.Run("Create one person account and get it", func(t *testing.T) {
		newAccount, err := CreatePersonAccount(person.ID, GeneratePersonAccountFields())
		assert.Nil(t, err)
		obtainedAccount, err := GetOnePersonAccount(newAccount.ID)
		assert.Nil(t, err)
		assert.Equal(t, newAccount, obtainedAccount)
	})

	DeleteAllPersonAccounts()

	t.Run("Error when getting unexisting person account", func(t *testing.T) {
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = GetOnePersonAccount(randId)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("It should create and update one person account", func(t *testing.T) {
		createFields := GeneratePersonAccountFields()
		updateFields := GeneratePersonAccountFields()
		newAccount, err := CreatePersonAccount(person.ID, createFields)
		assert.Nil(t, err)
		updatedAccount, err := UpdatePersonAccount(newAccount.ID, updateFields)
		assert.Nil(t, err)
		assert.Equal(t, newAccount.ID, updatedAccount.ID)
		assert.Equal(t, updatedAccount.Name, updateFields.Name)
		assert.Equal(t, updatedAccount.Description, updateFields.Description)
		// assert.Equal(t, updatedAccount.Currency, updateFields.Currency)  -> does not update currency
	})

	DeleteAllPersonAccounts()

	t.Run("It should generate and error when trying to update unexisting account", func(t *testing.T) {
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = UpdatePersonAccount(randId, GeneratePersonAccountFields())
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Create one person account and delete it", func(t *testing.T) {
		newAccount, err := CreatePersonAccount(person.ID, GeneratePersonAccountFields())
		assert.Nil(t, err)
		deletedId, err := DeletePersonAccount(newAccount.ID)
		assert.Nil(t, err)
		assert.Equal(t, newAccount.ID, deletedId.ID)
		_, err = GetOnePersonAccount(deletedId.ID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	DeleteAllPersonAccounts()

	t.Run("Error when trying to delete an unexisting person account", func(t *testing.T) {
		randId, err := uuid.NewRandom()
		assert.Nil(t, err)
		_, err = DeletePersonAccount(randId)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})
}
