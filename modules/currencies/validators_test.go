package currencies

import (
	"testing"

	errors_handler "github.com/grabielcruz/transportation_back/errors"
	"github.com/stretchr/testify/assert"
)

func TestCheckValidCurrency(t *testing.T) {
	goodCurrency := "ABC"
	err := CheckValidCurrency(goodCurrency)
	assert.Nil(t, err)

	longCurrency := "ABCD"
	err = CheckValidCurrency(longCurrency)
	assert.Equal(t, errors_handler.CU002, err.Error())

	lowerCurrency := "wer"
	err = CheckValidCurrency(lowerCurrency)
	assert.Equal(t, errors_handler.CU002, err.Error())

	shortCurrency := "WE"
	err = CheckValidCurrency(shortCurrency)
	assert.Equal(t, errors_handler.CU002, err.Error())

	badCurrency := "!@#"
	err = CheckValidCurrency(badCurrency)
	assert.Equal(t, errors_handler.CU002, err.Error())

	badCurrency2 := "ABCe"
	err = CheckValidCurrency(badCurrency2)
	assert.Equal(t, errors_handler.CU002, err.Error())

	badCurrency3 := "AB C"
	err = CheckValidCurrency(badCurrency3)
	assert.Equal(t, errors_handler.CU002, err.Error())
}
