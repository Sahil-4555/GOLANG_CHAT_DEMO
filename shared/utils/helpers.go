package utils

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func SliceContains(string string, array ...string) bool {
	for _, n := range array {
		if string == n {
			return true
		}
	}
	return false
}

func GenerateID() string {
	u1 := uuid.Must(uuid.NewV4())
	return fmt.Sprintf("%s", u1)
}
