package persons

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/grabielcruz/transportation_back/database"
	"github.com/stretchr/testify/assert"
)

func TestPersonService(t *testing.T) {
	envPath := filepath.Clean("../../.env_test")
	sqlPath := filepath.Clean("../../database/database.sql")
	database.SetupDB(envPath)
	database.CreateTables(sqlPath)
	defer database.CloseConnection()

	// zero person should be couned
	t.Run("Get one person initially", func(t *testing.T) {
		persons := GetPersons()
		assert.Len(t, persons, 1)
	})

	t.Run("Create one person", func(t *testing.T) {
		personFields := GeneratePersonFields()
		createdPerson := CreatePerson(personFields)
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
		newPerson := CreatePerson(GeneratePersonFields())
		obtainedPerson, err := GetOnePerson(newPerson.ID)
		assert.Nil(t, err)
		assert.Equal(t, newPerson.ID, obtainedPerson.ID)
	})

	DeleteAllPersons()

	t.Run("Error when getting unexisting person", func(t *testing.T) {
		zeroUUID := uuid.UUID{}
		_, err := GetOnePerson(zeroUUID)
		assert.NotNil(t, err)
		assert.Equal(t, "sql: no rows in result set", err.Error())
	})

	t.Run("It should create and update one person", func(t *testing.T) {
		createFields := GeneratePersonFields()
		updateFields := GeneratePersonFields()
		newPerson := CreatePerson(createFields)
		updatedPerson, err := UpdatePerson(newPerson.ID, updateFields)
		assert.Nil(t, err)
		assert.Equal(t, newPerson.ID, updatedPerson.ID)
		assert.Equal(t, updateFields.Name, updatedPerson.Name)
		assert.Equal(t, updateFields.Document, updatedPerson.Document)
		assert.NotEqual(t, newPerson.UpdatedAt, updatedPerson.UpdatedAt)
	})

	DeleteAllPersons()

	t.Run("It should genereate error when trying to update unexisting person", func(t *testing.T) {
		zeroUUID := uuid.UUID{}
		zeroFields := PersonFields{}
		_, err := UpdatePerson(zeroUUID, zeroFields)
		assert.NotNil(t, err)
	})

	t.Run("Create a person and delete it", func(t *testing.T) {
		newPerson := CreatePerson(GeneratePersonFields())
		deletedId, err := DeleteOnePerson(newPerson.ID)
		assert.Nil(t, err)
		assert.Equal(t, newPerson.ID, deletedId.ID)
		_, err = GetOnePerson(deletedId.ID)
		assert.NotNil(t, err)
		assert.Equal(t, "sql: no rows in result set", err.Error())
	})

	t.Run("Error when attempting to delete an unexisting person", func(t *testing.T) {
		zeroUUID := uuid.UUID{}
		_, err := DeleteOnePerson(zeroUUID)
		assert.NotNil(t, err)
	})

	t.Run("Create one person and get its name", func(t *testing.T) {
		newPerson := CreatePerson(GeneratePersonFields())
		name, err := GetPersonsName(newPerson.ID)
		assert.Nil(t, err)
		assert.Equal(t, newPerson.Name, name)
	})

	t.Run("Error when getting unexisting persons name", func(t *testing.T) {
		_, err := GetPersonsName(uuid.UUID{})
		assert.NotNil(t, err)
		assert.Equal(t, "sql: no rows in result set", err.Error())
	})

}
