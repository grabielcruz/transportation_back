package person_accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPersonAccountFields(t *testing.T) {
	fields := PersonAccountFields{}
	err := checkPersonAccountFields(fields)
	assert.Equal(t, "Name is required", err.Error())
	fields.Name = "John"
	err = checkPersonAccountFields(fields)
	assert.Equal(t, "Description is required", err.Error())
	fields.Description = "Hola mundo"
	err = checkPersonAccountFields(fields)
	assert.Equal(t, "Currency code should be 3 upper case letters", err.Error())

}
