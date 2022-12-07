package money_accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckAccountFields(t *testing.T) {
	fields := MoneyAccountFields{}
	err := checkAccountFields(fields)
	assert.NotNil(t, err)
	fields.Name = "John"
	err = checkAccountFields(fields)
	assert.NotNil(t, err)
}
