package currencies

import (
	"fmt"
	"regexp"
)

func CheckValidCurrency(currency string) error {
	validCurrency := regexp.MustCompile(`^[A-Z]{3}$`)
	if !validCurrency.MatchString(currency) {
		return fmt.Errorf(`Currency code should be 3 upper case letters`)
	}
	return nil
}
