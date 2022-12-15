package money_accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckAccountFields(t *testing.T) {
	fields := MoneyAccountFields{}
	err := checkAccountFields(fields)
	assert.Equal(t, "Name is required", err.Error())
	fields.Name = "John"
	err = checkAccountFields(fields)
	assert.Equal(t, "Currency is required", err.Error())
}
