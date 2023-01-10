package persons

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/stretchr/testify/assert"
)

func TestPersonService(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()

	// zero person should be couned
	t.Run("Get zero persons initially", func(t *testing.T) {
		persons := GetPersons()
		assert.Len(t, persons, 0)
	})

	t.Run("Create one person", func(t *testing.T) {
		personFields := GeneratePersonFields()
		createdPerson, err := CreatePerson(personFields)
		assert.Nil(t, err)
		assert.Equal(t, personFields.Name, createdPerson.Name)
		assert.Equal(t, personFields.Document, createdPerson.Document)
	})

	DeleteAllPersons()

	t.Run("Create two person and get an slice of persons", func(t *testing.T) {
		CreatePerson(GeneratePersonFields())
		CreatePerson(GeneratePersonFields())
		persons := GetPersons()
		assert.Len(t, persons, 2)
	})

	DeleteAllPersons()

	t.Run("Create one person and get it", func(t *testing.T) {
		newPerson, err := CreatePerson(GeneratePersonFields())
		assert.Nil(t, err)
		obtainedPerson, err := GetOnePerson(newPerson.ID)
		assert.Nil(t, err)
		assert.Equal(t, newPerson.ID, obtainedPerson.ID)
	})

	DeleteAllPersons()

	t.Run("Error when getting unexisting person", func(t *testing.T) {
		zeroUUID := uuid.UUID{}
		_, err := GetOnePerson(zeroUUID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("It should create and update one person", func(t *testing.T) {
		createFields := GeneratePersonFields()
		updateFields := GeneratePersonFields()
		newPerson, err := CreatePerson(createFields)
		assert.Nil(t, err)
		updatedPerson, err := UpdatePerson(newPerson.ID, updateFields)
		assert.Nil(t, err)
		assert.Equal(t, newPerson.ID, updatedPerson.ID)
		assert.Equal(t, updateFields.Name, updatedPerson.Name)
		assert.Equal(t, updateFields.Document, updatedPerson.Document)
		assert.Greater(t, updatedPerson.UpdatedAt, newPerson.CreatedAt)
	})

	DeleteAllPersons()

	t.Run("It should genereate error when trying to update unexisting person", func(t *testing.T) {
		zeroUUID := uuid.UUID{}
		zeroFields := PersonFields{}
		_, err := UpdatePerson(zeroUUID, zeroFields)
		assert.NotNil(t, err)
	})

	t.Run("Create a person and delete it", func(t *testing.T) {
		newPerson, err := CreatePerson(GeneratePersonFields())
		assert.Nil(t, err)
		deletedId, err := DeleteOnePerson(newPerson.ID)
		assert.Nil(t, err)
		assert.Equal(t, newPerson.ID, deletedId.ID)
		_, err = GetOnePerson(deletedId.ID)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Error when attempting to delete an unexisting person", func(t *testing.T) {
		zeroUUID := uuid.UUID{}
		_, err := DeleteOnePerson(zeroUUID)
		assert.NotNil(t, err)
	})

	t.Run("Create one person and get its name", func(t *testing.T) {
		newPerson, err := CreatePerson(GeneratePersonFields())
		assert.Nil(t, err)
		name, err := GetPersonsName(newPerson.ID)
		assert.Nil(t, err)
		assert.Equal(t, newPerson.Name, name)
	})

	t.Run("Error when getting unexisting persons name", func(t *testing.T) {
		_, err := GetPersonsName(uuid.UUID{})
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.DB001, err.Error())
	})

	t.Run("Error when creating two persons with the same document", func(t *testing.T) {
		fields1 := GeneratePersonFields()
		fields1.Document = "v7777777"
		fields2 := GeneratePersonFields()
		fields2.Document = "v7777777"

		_, err := CreatePerson(fields1)
		assert.Nil(t, err)

		_, err = CreatePerson(fields2)
		assert.NotNil(t, err)
		assert.Equal(t, errors_handler.PE001, err.Error())
	})

}
