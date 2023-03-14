package main

import (
	"fmt"
	"os"
)

func failWith(message string) {
	fmt.Println(message)
	os.Exit(1)
}
