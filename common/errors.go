package common

import (
	"fmt"
	"os"
)

// FailWith exits with an error message controllably
func FailWith(message string) {
	fmt.Println(message)
	os.Exit(1)
}
