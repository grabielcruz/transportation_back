package persons

import "fmt"

func checkPersonFields(fields PersonFields) error {
	if fields.Name == "" {
		return fmt.Errorf("Name is required")
	}
	return nil
}
